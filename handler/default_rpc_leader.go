package handler

import (
	"bytes"
	"club/model"
	authproto "club/proto/golang/auth"
	clubproto "club/proto/golang/club"
	consulagent "club/tool/consul/agent"
	"club/tool/mysqlerr"
	code "club/utils/code/golang"
	topic "club/utils/topic/golang"
	"context"
	"errors"
	"fmt"
	mysqlcode "github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
	microerrors "github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	"gorm.io/gorm"
	"net/http"
)

func (d *_default) AddClubMember(ctx context.Context, req *clubproto.AddClubMemberRequest, resp *clubproto.AddClubMemberResponse) (_ error) {
	ctx, proxyAuthenticated, reason := d.getContextFromMetadata(ctx)
	if !proxyAuthenticated {
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, reason)
		return
	}

	switch true {
	case adminUUIDRegex.MatchString(req.UUID):
		break
	case studentUUIDRegex.MatchString(req.UUID):
		break
	default:
		resp.Status = http.StatusForbidden
		resp.Code = code.ForbiddenNotStudentOrAdminUUID
		resp.Message = fmt.Sprintf(forbiddenMessageFormat, "you are not student or admin")
		return
	}

	reqID := ctx.Value("X-Request-Id").(string)
	parentSpan := ctx.Value("Span-Context").(jaeger.SpanContext)

	access := d.accessManage.BeginTx()

	spanForDB := d.tracer.StartSpan("GetClubWithClubUUID", opentracing.ChildOf(parentSpan))
	selectedClub, err := access.GetClubWithClubUUID(req.ClubUUID)
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedClub", selectedClub), log.Error(err))
	spanForDB.Finish()

	switch err {
	case nil:
		break
	case gorm.ErrRecordNotFound:
		access.Rollback()
		resp.Status = http.StatusNotFound
		resp.Code = code.NotFoundClubNoExist
		resp.Message = fmt.Sprintf(notFoundMessageFormat, "club with that uuid not exist")
		return
	default:
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "GetClubWithClubUUID returns unexpected error, err: " + err.Error())
		return
	}

	if !adminUUIDRegex.MatchString(req.UUID) && req.UUID != string(selectedClub.LeaderUUID) {
		access.Rollback()
		resp.Status = http.StatusForbidden
		resp.Code = code.ForbiddenNotClubLeader
		resp.Message = fmt.Sprintf(forbiddenMessageFormat, "you're not admin and not club leader")
		return
	}

	spanForConsul := d.tracer.StartSpan("GetNextServiceNode", opentracing.ChildOf(parentSpan))
	selectedNode, err := d.consulAgent.GetNextServiceNode(topic.AuthServiceName)
	spanForConsul.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedNode", selectedNode), log.Error(err))
	spanForConsul.Finish()

	switch err {
	case nil:
		break
	case consulagent.ErrAvailableNodeNotFound:
		access.Rollback()
		resp.Status = http.StatusServiceUnavailable
		resp.Message = fmt.Sprintf(serviceUnavailableMessageFormat, "there is no available server, service name: " + topic.AuthServiceName)
		return
	default:
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "unable to query in consul agent, err: " + err.Error())
		return
	}

	spanForReq := d.tracer.StartSpan("GetStudentInformWithUUID", opentracing.ChildOf(parentSpan))
	md := metadata.Set(context.Background(), "X-Request-Id", reqID)
	md = metadata.Set(md, "Span-Context", spanForReq.Context().(jaeger.SpanContext).String())
	authReq := &authproto.GetStudentInformWithUUIDRequest{
		UUID:        req.UUID,
		StudentUUID: req.StudentUUID,
	}
	respOfReq, err := d.authStudent.GetStudentInformWithUUID(md, authReq)
	spanForReq.SetTag("X-Request-Id", reqID).LogFields(log.Object("request", authReq), log.Object("response", respOfReq), log.Error(err))
	spanForReq.Finish()

	switch assertedError := err.(type) {
	case nil:
		break
	case *microerrors.Error:
		switch assertedError.Code {
		case http.StatusRequestTimeout:
			access.Rollback()
			resp.Status = http.StatusRequestTimeout
			resp.Message = fmt.Sprintf(requestTimeoutMessageFormat, assertedError.Detail)
			return
		default:
			access.Rollback()
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerMessageFormat, assertedError.Detail)
			return
		}
	default:
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, assertedError.Error())
		return
	}

	switch respOfReq.Status {
	case http.StatusOK:
		break
	case http.StatusNotFound:
		access.Rollback()
		resp.Status = http.StatusNotFound
		resp.Code = code.NotFoundStudentNoExist
		resp.Message = fmt.Sprintf(notFoundMessageFormat, "student with that uuid not eixst")
		return
	default:
		access.Rollback()
		resp.Status = respOfReq.Status
		resp.Message = fmt.Sprintf("GetStudentInformWithUUID unexpected status returned, message: %s", respOfReq.Message)
		return
	}

	spanForDB = d.tracer.StartSpan("CreateClubMember", opentracing.ChildOf(parentSpan))
	createdMember, err := access.CreateClubMember(&model.ClubMember{
		ClubUUID:    model.ClubUUID(req.ClubUUID),
		StudentUUID: model.StudentUUID(req.StudentUUID),
	})
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("CreatedMember", createdMember), log.Error(err))
	spanForDB.Finish()

	switch assertedError := err.(type) {
	case nil:
		break
	case *mysql.MySQLError:
		access.Rollback()
		switch assertedError.Number {
		case mysqlcode.ER_DUP_ENTRY:
			key, entry, err := mysqlerr.ParseDuplicateEntryErrorFrom(assertedError)
			if err != nil {
				resp.Status = http.StatusInternalServerError
				resp.Message = fmt.Sprintf(internalServerMessageFormat, "unable to parse MySQL duplicate error, err: " + err.Error())
				return
			}
			switch key {
			case model.ClubMemberInstance.StudentUUID.KeyName():
				resp.Status = http.StatusConflict
				resp.Code = code.ClubMemberAlreadyExist
				resp.Message = fmt.Sprintf(conflictMessageFormat, "alreay exists as member, entry: " + entry)
				return
			default:
				resp.Status = http.StatusInternalServerError
				resp.Message = fmt.Sprintf(internalServerMessageFormat, "unexpected duplicate entry, key: " + key)
				return
			}
		default:
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerMessageFormat, "unexpected CreateClubMember MySQL error code, err: " + assertedError.Error())
			return
		}
	default:
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "unexpected type of CreateClubMember errors, err: " + assertedError.Error())
		return
	}

	access.Commit()
	resp.Status = http.StatusCreated
	resp.Message = "success to create new club member"
	return
}

func (d *_default) DeleteClubMember(ctx context.Context, req *clubproto.DeleteClubMemberRequest, resp *clubproto.DeleteClubMemberResponse) (_ error) {
	ctx, proxyAuthenticated, reason := d.getContextFromMetadata(ctx)
	if !proxyAuthenticated {
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, reason)
		return
	}

	switch true {
	case adminUUIDRegex.MatchString(req.UUID):
		break
	case studentUUIDRegex.MatchString(req.UUID):
		break
	default:
		resp.Status = http.StatusForbidden
		resp.Code = code.ForbiddenNotStudentOrAdminUUID
		resp.Message = fmt.Sprintf(forbiddenMessageFormat, "you are not student or admin")
		return
	}

	reqID := ctx.Value("X-Request-Id").(string)
	parentSpan := ctx.Value("Span-Context").(jaeger.SpanContext)

	access := d.accessManage.BeginTx()

	spanForDB := d.tracer.StartSpan("GetClubWithClubUUID", opentracing.ChildOf(parentSpan))
	selectedClub, err := access.GetClubWithClubUUID(req.ClubUUID)
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedClub", selectedClub), log.Error(err))
	spanForDB.Finish()

	switch err {
	case nil:
		break
	case gorm.ErrRecordNotFound:
		access.Rollback()
		resp.Status = http.StatusNotFound
		resp.Code = code.NotFoundClubNoExist
		resp.Message = fmt.Sprintf(notFoundMessageFormat, "club with that uuid not exist")
		return
	default:
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "GetClubWithClubUUID returns unexpected error, err: " + err.Error())
		return
	}

	if !adminUUIDRegex.MatchString(req.UUID) && req.UUID != string(selectedClub.LeaderUUID) {
		access.Rollback()
		resp.Status = http.StatusForbidden
		resp.Code = code.ForbiddenNotClubLeader
		resp.Message = fmt.Sprintf(forbiddenMessageFormat, "you're not admin and not club leader")
		return
	}

	spanForDB = d.tracer.StartSpan("DeleteClubMember", opentracing.ChildOf(parentSpan))
	err, rowAffected := access.DeleteClubMember(req.ClubUUID, req.StudentUUID)
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Int("RowAffected", int(rowAffected)), log.Error(err))
	spanForDB.Finish()

	if err != nil {
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "DeleteClubMember returns unexpected error, err: " + err.Error())
		return
	}

	if rowAffected == 0 {
		access.Rollback()
		resp.Status = http.StatusNotFound
		resp.Code = code.NotFoundClubMemberNoExist
		resp.Message = fmt.Sprintf("club member with that student uuid not exist")
		return
	}

	access.Commit()
	resp.Status = http.StatusOK
	resp.Message = "success delete club member"
	return
}

func (d *_default) ChangeClubLeader(ctx context.Context, req *clubproto.ChangeClubLeaderRequest, resp *clubproto.ChangeClubLeaderResponse) (_ error) {
	ctx, proxyAuthenticated, reason := d.getContextFromMetadata(ctx)
	if !proxyAuthenticated {
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, reason)
		return
	}

	switch true {
	case adminUUIDRegex.MatchString(req.UUID):
		break
	case studentUUIDRegex.MatchString(req.UUID):
		break
	default:
		resp.Status = http.StatusForbidden
		resp.Code = code.ForbiddenNotStudentOrAdminUUID
		resp.Message = fmt.Sprintf(forbiddenMessageFormat, "you are not student or admin")
		return
	}

	reqID := ctx.Value("X-Request-Id").(string)
	parentSpan := ctx.Value("Span-Context").(jaeger.SpanContext)

	access := d.accessManage.BeginTx()

	spanForDB := d.tracer.StartSpan("GetClubWithClubUUID", opentracing.ChildOf(parentSpan))
	selectedClub, err := access.GetClubWithClubUUID(req.ClubUUID)
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedClub", selectedClub), log.Error(err))
	spanForDB.Finish()

	switch err {
	case nil:
		break
	case gorm.ErrRecordNotFound:
		access.Rollback()
		resp.Status = http.StatusNotFound
		resp.Code = code.NotFoundClubNoExist
		resp.Message = fmt.Sprintf(notFoundMessageFormat, "club with that uuid not exist")
		return
	default:
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "GetClubWithClubUUID returns unexpected error, err: " + err.Error())
		return
	}

	if !adminUUIDRegex.MatchString(req.UUID) && req.UUID != string(selectedClub.LeaderUUID) {
		access.Rollback()
		resp.Status = http.StatusForbidden
		resp.Code = code.ForbiddenNotClubLeader
		resp.Message = fmt.Sprintf(forbiddenMessageFormat, "you're not admin and not club leader")
		return
	}

	if req.NewLeaderUUID == string(selectedClub.LeaderUUID) {
		access.Rollback()
		resp.Status = http.StatusConflict
		resp.Code = code.AlreadyClubLeader
		resp.Message = fmt.Sprintf("that student is already club leader")
		return
	}

	spanForDB = d.tracer.StartSpan("GetClubMembersWithClubUUID", opentracing.ChildOf(parentSpan))
	selectedMembers, err := access.GetClubMembersWithClubUUID(req.ClubUUID)
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedMembers", selectedMembers), log.Error(err))
	spanForDB.Finish()

	if err != nil && err != gorm.ErrRecordNotFound {
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "GetClubMembersWithClubUUID returns unexpected error, err: " + err.Error())
		return
	}

	memberUUIDs := make([]string, len(selectedMembers))
	for index, selectedMember := range selectedMembers {
		memberUUIDs[index] = string(selectedMember.StudentUUID)
	}

	if !contains(memberUUIDs, req.NewLeaderUUID) {
		access.Rollback()
		resp.Status = http.StatusNotFound
		resp.Code = code.NotFoundClubMemberNoExist
		resp.Message = fmt.Sprintf(notFoundMessageFormat, "member to be club leader is not exists")
		return
	}

	spanForDB = d.tracer.StartSpan("ChangeClubLeader", opentracing.ChildOf(parentSpan))
	err, rowAffected := access.ChangeClubLeader(req.ClubUUID, req.NewLeaderUUID)
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Int("RowAffected", int(rowAffected)), log.Error(err))
	spanForDB.Finish()

	switch assertedError := err.(type) {
	case nil:
		break
	case *mysql.MySQLError:
		switch assertedError.Number {
		case mysqlcode.ER_DUP_ENTRY:
			key, entry, err := mysqlerr.ParseDuplicateEntryErrorFrom(assertedError)
			if err != nil {
				err = errors.New("unable to parse ChangeClubLeader duplicate error, err: " + err.Error())
				break
			}
			switch key {
			case model.ClubInstance.LeaderUUID.KeyName():
				access.Rollback()
				resp.Status = http.StatusConflict
				resp.Code = code.ClubLeaderDuplicateForChange
				resp.Message = fmt.Sprintf(conflictMessageFormat, "that leader uuid is already other club's leader, entry: " + entry)
				return
			default:
				err = errors.New("unexpected duplicate entry, key: " + key)
			}
		default:
			err = errors.New("unexpected ChangeClubLeader MySQL error code, err: " + assertedError.Error())
		}
	default:
		err = errors.New("unexpected type of ChangeClubLeader error, err: " + assertedError.Error())
	}

	if err != nil {
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, err.Error())
		return
	}

	if rowAffected == 0 {
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "ChangeClubLeader returns 0 row affected")
		return
	}

	access.Commit()
	resp.Status = http.StatusOK
	resp.Message = "success change club leader"
	return
}

func (d *_default) ModifyClubInform(ctx context.Context, req *clubproto.ModifyClubInformRequest, resp *clubproto.ModifyClubInformResponse) (_ error) {
	ctx, proxyAuthenticated, reason := d.getContextFromMetadata(ctx)
	if !proxyAuthenticated {
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, reason)
		return
	}

	switch true {
	case adminUUIDRegex.MatchString(req.UUID):
		break
	case studentUUIDRegex.MatchString(req.UUID):
		break
	default:
		resp.Status = http.StatusForbidden
		resp.Code = code.ForbiddenNotStudentOrAdminUUID
		resp.Message = fmt.Sprintf(forbiddenMessageFormat, "you are not student or admin")
		return
	}

	reqID := ctx.Value("X-Request-Id").(string)
	parentSpan := ctx.Value("Span-Context").(jaeger.SpanContext)

	access := d.accessManage.BeginTx()

	spanForDB := d.tracer.StartSpan("GetClubWithClubUUID", opentracing.ChildOf(parentSpan))
	selectedClub, err := access.GetClubWithClubUUID(req.ClubUUID)
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedClub", selectedClub), log.Error(err))
	spanForDB.Finish()

	switch err {
	case nil:
		break
	case gorm.ErrRecordNotFound:
		access.Rollback()
		resp.Status = http.StatusNotFound
		resp.Code = code.NotFoundClubNoExist
		resp.Message = fmt.Sprintf(notFoundMessageFormat, "club with that uuid not exist")
		return
	default:
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "GetClubWithClubUUID returns unexpected error, err: " + err.Error())
		return
	}

	if !adminUUIDRegex.MatchString(req.UUID) && req.UUID != string(selectedClub.LeaderUUID) {
		access.Rollback()
		resp.Status = http.StatusForbidden
		resp.Code = code.ForbiddenNotClubLeader
		resp.Message = fmt.Sprintf(forbiddenMessageFormat, "you're not admin and not club leader")
		return
	}
}
