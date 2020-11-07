package handler

import (
	"bytes"
	"club/model"
	authproto "club/proto/golang/auth"
	clubproto "club/proto/golang/club"
	consulagent "club/tool/consul/agent"
	"club/tool/mysqlerr"
	"club/tool/random"
	code "club/utils/code/golang"
	topic "club/utils/topic/golang"
	"context"
	"errors"
	"fmt"
	mysqlcode "github.com/VividCortex/mysqlerr"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
	"github.com/micro/go-micro/v2/client"
	microerrors "github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	"gorm.io/gorm"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
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
	callOpts := []client.CallOption{client.WithDialTimeout(time.Second * 2), client.WithRequestTimeout(time.Second * 3), client.WithAddress(selectedNode.Address)}
	respOfReq, err := d.authStudent.GetStudentInformWithUUID(md, authReq, callOpts...)
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

	if d.awsSession != nil {
		spanForS3 := d.tracer.StartSpan("PutObject", opentracing.ChildOf(parentSpan))
		_, err = s3.New(d.awsSession).PutObject(&s3.PutObjectInput{
			Bucket: aws.String(s3Bucket),
			Key:    aws.String(fmt.Sprintf("logos/%s", req.ClubUUID)),
			Body:   bytes.NewReader(req.Logo),
		})
		spanForS3.SetTag("X-Request-Id", reqID).LogFields(log.Error(err))
		spanForS3.Finish()

		if err != nil {
			access.Rollback()
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerMessageFormat, "unable to upload profile to s3, err: " + err.Error())
			return
		}
	}

	spanForDB = d.tracer.StartSpan("ModifyClubInform", opentracing.ChildOf(parentSpan))
	err, rowAffected := access.ModifyClubInform(req.ClubUUID, &model.ClubInform{
		ClubConcept:  model.ClubConcept(req.ClubConcept),
		Introduction: model.Introduction(req.Introduction),
		Link:         model.Link(req.Link),
	})
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Int("RowAffected", int(rowAffected)), log.Error(err))
	spanForDB.Finish()

	switch assertedError := err.(type) {
	case nil:
		break
	case validator.ValidationErrors:
		access.Rollback()
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, "invalid data for club inform model, err: " + assertedError.Error())
		return
	default:
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "ModifyClubInform returns unexpected error, err: " + assertedError.Error())
		return
	}

	if rowAffected == 0 {
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "ModifyClubInform returns 0 row affected")
		return
	}

	access.Commit()
	resp.Status = http.StatusOK
	resp.Message = "success modify club inform"
	return
}

func (d *_default) DeleteClubWithUUID(ctx context.Context, req *clubproto.DeleteClubWithUUIDRequest, resp *clubproto.DeleteClubWithUUIDResponse) (_ error) {
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

	spanForDB = d.tracer.StartSpan("GetCurrentRecruitmentWithClubUUID", opentracing.ChildOf(parentSpan))
	selectedRecruit, err := access.GetCurrentRecruitmentWithClubUUID(req.ClubUUID)
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedRecruit", selectedRecruit), log.Error(err))
	spanForDB.Finish()

	switch err {
	case gorm.ErrRecordNotFound:
		break
	case nil:
		access.Rollback()
		resp.Status = http.StatusConflict
		resp.Code = code.RecruitmentInProgressExist
		resp.Message = fmt.Sprintf(conflictMessageFormat, "there is recruitment which is in progress")
		return
	default:
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "GetCurrentRecruitmentWithClubUUID returns unexpected error, err: " + err.Error())
		return
	}

	spanForDB = d.tracer.StartSpan("DeleteClub", opentracing.ChildOf(parentSpan))
	err, rowsAffected := access.DeleteClub(req.ClubUUID)
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Int("RowAffected", int(rowsAffected)), log.Error(err))
	spanForDB.Finish()

	if err != nil {
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "DeleteClub returns unexpected error, err: " + err.Error())
		return
	}

	if rowsAffected == 0 {
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "DeleteClub returns 0 rows affected")
		return
	}

	spanForDB = d.tracer.StartSpan("DeleteClubInform", opentracing.ChildOf(parentSpan))
	err, rowsAffected = access.DeleteClubInform(req.ClubUUID)
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Int("RowAffected", int(rowsAffected)), log.Error(err))
	spanForDB.Finish()

	if err != nil {
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "DeleteClubInform returns unexpected error, err: " + err.Error())
		return
	}

	if rowsAffected == 0 {
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "DeleteClubInform returns 0 rows affected")
		return
	}

	spanForDB = d.tracer.StartSpan("DeleteAllClubMembers", opentracing.ChildOf(parentSpan))
	err, rowsAffected = access.DeleteAllClubMembers(req.ClubUUID)
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Int("RowAffected", int(rowsAffected)), log.Error(err))
	spanForDB.Finish()

	if err != nil {
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "DeleteAllClubMembers returns unexpected error, err: " + err.Error())
		return
	}

	access.Commit()
	resp.Status = http.StatusOK
	resp.Message = "succeed to delete club"
	return
}

func (d *_default) RegisterRecruitment(ctx context.Context, req *clubproto.RegisterRecruitmentRequest, resp *clubproto.RegisterRecruitmentResponse) (_ error) {
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

	spanForDB = d.tracer.StartSpan("GetCurrentRecruitmentWithClubUUID", opentracing.ChildOf(parentSpan))
	selectedRecruit, err := access.GetCurrentRecruitmentWithClubUUID(req.ClubUUID)
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedRecruit", selectedRecruit), log.Error(err))
	spanForDB.Finish()

	switch err {
	case gorm.ErrRecordNotFound:
		break
	case nil:
		access.Rollback()
		resp.Status = http.StatusConflict
		resp.Code = code.RecruitmentInProgressAlreadyExist
		resp.Message = fmt.Sprintf(conflictMessageFormat, "recruitment in progress is already exists")
		return
	default:
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "GetCurrentRecruitmentWithClubUUID returns unexpected error, err: " + err.Error())
		return
	}

	rUUID, ok := ctx.Value("RecruitmentUUID").(string)
	if !ok || rUUID == "" {
		rUUID = fmt.Sprintf("recruitment-%s", random.StringConsistOfIntWithLength(12))
	}

	for {
		spanForDB := d.tracer.StartSpan("GetRecruitmentWithRecruitmentUUID", opentracing.ChildOf(parentSpan))
		selectedClub, err := access.GetRecruitmentWithRecruitmentUUID(rUUID)
		spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedClub", selectedClub), log.Error(err))
		spanForDB.Finish()
		if err == gorm.ErrRecordNotFound {
			break
		}
		if err != nil {
			access.Rollback()
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerMessageFormat, "unexpected error in GetRecruitmentWithRecruitmentUUID, err: " + err.Error())
			return
		}
		rUUID = fmt.Sprintf("recruitment-%s", random.StringConsistOfIntWithLength(12))
		continue
	}

	now := time.Now()
	startTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	endTime := time.Time{}
	if req.EndPeriod != "" {
		endTimeSplice := strings.Split(req.EndPeriod, "-")
		if len(endTimeSplice) != 3 {
			access.Rollback()
			resp.Status = http.StatusProxyAuthRequired
			resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, "invalid EndPeriod value")
			return
		}

		err = nil
		const indexForYear = 0
		const indexForMonth = 1
		const indexForDay = 2
		year, convertErr := strconv.Atoi(endTimeSplice[indexForYear])
		if len(endTimeSplice[indexForYear]) != 4 || convertErr != nil { err = errors.New("year invalid") }
		month, convertErr := strconv.Atoi(endTimeSplice[indexForMonth])
		if len(endTimeSplice[indexForMonth]) != 2 || convertErr != nil { err = errors.New("month invalid") }
		day, convertErr := strconv.Atoi(endTimeSplice[indexForDay])
		if len(endTimeSplice[indexForDay]) != 2 || convertErr != nil { err = errors.New("day invalid") }

		if err != nil {
			access.Rollback()
			resp.Status = http.StatusProxyAuthRequired
			resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, "invalid EndPeriod value")
			return
		}

		endTime = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	}

	if !reflect.DeepEqual(endTime, time.Time{}) && endTime.Add(time.Hour).Sub(startTime).Milliseconds() < 0 {
		access.Rollback()
		resp.Status = http.StatusConflict
		resp.Code = code.EndPeriodOlderThanNow
		resp.Message = fmt.Sprintf(conflictMessageFormat, "end period is older than now")
		return
	}

	spanForDB = d.tracer.StartSpan("CreateRecruitment", opentracing.ChildOf(parentSpan))
	createdRecruitment, err := access.CreateRecruitment(&model.ClubRecruitment{
		UUID:           model.UUID(rUUID),
		ClubUUID:       model.ClubUUID(req.ClubUUID),
		RecruitConcept: model.RecruitConcept(req.RecruitConcept),
		StartPeriod:    model.StartPeriod(startTime),
		EndPeriod:      model.EndPeriod(endTime),
	})
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("CreatedRecruitment", createdRecruitment), log.Error(err))
	spanForDB.Finish()

	switch err.(type) {
	case nil:
		break
	case validator.ValidationErrors:
		access.Rollback()
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, "invalid data for club recruit model")
		return
	default:
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "CreateRecruitment returns unexpected error, err: " + err.Error())
		return
	}

	if len(req.RecruitMembers) == 0 {
		access.Rollback()
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, "recruit member list is empty")
		return
	}

	createdRecruitMembers := make([]*model.RecruitMember, len(req.RecruitMembers))
	spanForDB = d.tracer.StartSpan("CreateRecruitMembers", opentracing.ChildOf(parentSpan))
	for index, recruitMember := range req.RecruitMembers {
		createdRecruitMember, commandErr := access.CreateRecruitMember(&model.RecruitMember{
			RecruitmentUUID: model.RecruitmentUUID(string(createdRecruitment.UUID)),
			Grade:           model.Grade(recruitMember.Grade),
			Field:           model.Field(recruitMember.Field),
			Number:          model.Number(recruitMember.Number),
		})
		if commandErr != nil {
			err = commandErr
			break
		}
		createdRecruitMembers[index] = createdRecruitMember
	}
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("CreatedRecruitMembers", createdRecruitMembers), log.Error(err))
	spanForDB.Finish()

	switch err.(type) {
	case nil:
		break
	case validator.ValidationErrors:
		access.Rollback()
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, "invalid data for club recruit model")
		return
	default:
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "CreateRecruitment returns unexpected error, err: " + err.Error())
		return
	}

	access.Commit()
	resp.Status = http.StatusCreated
	resp.RecruitmentUUID = string(createdRecruitment.UUID)
	resp.Message = "succeed to register club recruitment"
	return
}

func (d *_default) ModifyRecruitment(ctx context.Context, req *clubproto.ModifyRecruitmentRequest, resp *clubproto.ModifyRecruitmentResponse) (_ error) {
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

	spanForDB := d.tracer.StartSpan("GetCurrentRecruitmentWithRecruitmentUUID", opentracing.ChildOf(parentSpan))
	selectedRecruit, err := access.GetCurrentRecruitmentWithRecruitmentUUID(req.RecruitmentUUID)
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedRecruit", selectedRecruit), log.Error(err))
	spanForDB.Finish()

	switch err {
	case nil:
		break
	case gorm.ErrRecordNotFound:
		access.Rollback()
		resp.Status = http.StatusNotFound
		resp.Code = code.NotFoundCurrentRecruitmentNoExist
		resp.Message = fmt.Sprintf(notFoundMessageFormat, "recruitment which is in progress not exists")
		return
	default:
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "GetCurrentRecruitmentWithRecruitmentUUID returns unexpected error, err: " + err.Error())
		return
	}

	spanForDB = d.tracer.StartSpan("GetClubWithClubUUID", opentracing.ChildOf(parentSpan))
	selectedClub, err := access.GetClubWithClubUUID(string(selectedRecruit.ClubUUID))
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedClub", selectedClub), log.Error(err))
	spanForDB.Finish()

	if err != nil {
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

	spanForDB = d.tracer.StartSpan("ModifyRecruitment", opentracing.ChildOf(parentSpan))
	err, rowAffected := access.ModifyRecruitment(string(selectedRecruit.UUID), &model.ClubRecruitment{
		RecruitConcept: model.RecruitConcept(req.RecruitConcept),
	})
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Int("RowAffected", int(rowAffected)), log.Error(err))
	spanForDB.Finish()

	if err != nil {
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "ModifyRecruitment returns unexpected error, err: " + err.Error())
		return
	}

	if rowAffected == 0 {
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "ModifyRecruitment returns 0 rows affected")
		return
	}

	if len(req.RecruitMembers) == 0 {
		access.Commit()
		resp.Status = http.StatusOK
		resp.Message = "succeed to modify club recruitment"
		return
	}

	spanForDB = d.tracer.StartSpan("DeleteAllRecruitMember", opentracing.ChildOf(parentSpan))
	err, rowAffected = access.DeleteAllRecruitMember(string(selectedRecruit.UUID))
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Int("RowAffected", int(rowAffected)), log.Error(err))
	spanForDB.Finish()

	if err != nil {
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "DeleteAllRecruitMember returns unexpected error, err: " + err.Error())
		return
	}

	createdRecruitMembers := make([]*model.RecruitMember, len(req.RecruitMembers))
	spanForDB = d.tracer.StartSpan("CreateRecruitMembers", opentracing.ChildOf(parentSpan))
	for index, recruitMember := range req.RecruitMembers {
		createdRecruitMember, commandErr := access.CreateRecruitMember(&model.RecruitMember{
			RecruitmentUUID: model.RecruitmentUUID(string(selectedRecruit.UUID)),
			Grade:           model.Grade(recruitMember.Grade),
			Field:           model.Field(recruitMember.Field),
			Number:          model.Number(recruitMember.Number),
		})
		if commandErr != nil {
			err = commandErr
			break
		}
		createdRecruitMembers[index] = createdRecruitMember
	}
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("CreatedRecruitMembers", createdRecruitMembers), log.Error(err))
	spanForDB.Finish()

	switch err.(type) {
	case nil:
		break
	case validator.ValidationErrors:
		access.Rollback()
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, "invalid data for club recruit model")
		return
	default:
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "CreateRecruitMembers returns unexpected error, err: " + err.Error())
		return
	}

	access.Commit()
	resp.Status = http.StatusOK
	resp.Message = "succeed to modify club recruitment"
	return
}

func (d *_default) DeleteRecruitmentWithUUID(ctx context.Context, req *clubproto.DeleteRecruitmentWithUUIDRequest, resp *clubproto.DeleteRecruitmentWithUUIDResponse) (_ error) {
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

	spanForDB := d.tracer.StartSpan("GetCurrentRecruitmentWithRecruitmentUUID", opentracing.ChildOf(parentSpan))
	selectedRecruit, err := access.GetCurrentRecruitmentWithRecruitmentUUID(req.RecruitmentUUID)
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedRecruit", selectedRecruit), log.Error(err))
	spanForDB.Finish()

	switch err {
	case nil:
		break
	case gorm.ErrRecordNotFound:
		access.Rollback()
		resp.Status = http.StatusNotFound
		resp.Code = code.NotFoundCurrentRecruitmentNoExist
		resp.Message = fmt.Sprintf(notFoundMessageFormat, "recruitment which is in progress not exists")
		return
	default:
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "GetCurrentRecruitmentWithRecruitmentUUID returns unexpected error, err: " + err.Error())
		return
	}

	spanForDB = d.tracer.StartSpan("GetClubWithClubUUID", opentracing.ChildOf(parentSpan))
	selectedClub, err := access.GetClubWithClubUUID(string(selectedRecruit.ClubUUID))
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedClub", selectedClub), log.Error(err))
	spanForDB.Finish()

	if err != nil {
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

	spanForDB = d.tracer.StartSpan("DeleteRecruitment", opentracing.ChildOf(parentSpan))
	err, rowsAffected := access.DeleteRecruitment(string(selectedRecruit.UUID))
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Int("RowsAffected", int(rowsAffected)), log.Error(err))
	spanForDB.Finish()

	if err != nil {
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "DeleteRecruitment returns unexpected error, err: " + err.Error())
		return
	}

	if rowsAffected == 0 {
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "DeleteRecruitment returns 0 rows affected")
		return
	}

	spanForDB = d.tracer.StartSpan("DeleteAllRecruitMember", opentracing.ChildOf(parentSpan))
	err, rowsAffected = access.DeleteAllRecruitMember(string(selectedRecruit.UUID))
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Int("RowsAffected", int(rowsAffected)), log.Error(err))
	spanForDB.Finish()

	if err != nil {
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "DeleteAllRecruitMember returns unexpected error, err: " + err.Error())
		return
	}

	access.Commit()
	resp.Status = http.StatusOK
	resp.Message = "succeed to delete club recruitment"
	return
}
