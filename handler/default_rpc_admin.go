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
	"fmt"
	mysqlcode "github.com/VividCortex/mysqlerr"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
	microerrors "github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	"gorm.io/gorm"
	"net/http"
)

func (d *_default) CreateNewClub(ctx context.Context, req *clubproto.CreateNewClubRequest, resp *clubproto.CreateNewClubResponse) (_ error) {
	ctx, proxyAuthenticated, reason := d.getContextFromMetadata(ctx)
	if !proxyAuthenticated {
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, reason)
		return
	}

	if !adminUUIDRegex.MatchString(req.UUID) {
		resp.Status = http.StatusForbidden
		resp.Message = fmt.Sprintf(forbiddenMessageFormat, "you are not admin")
		return
	}

	reqID := ctx.Value("X-Request-Id").(string)
	parentSpan := ctx.Value("Span-Context").(jaeger.SpanContext)

	if !contains(req.MemberUUIDs, req.LeaderUUID) {
		resp.Status = http.StatusConflict
		resp.Code = code.MemberUUIDsNotIncludeLeaderUUID
		resp.Message = fmt.Sprintf(conflictMessageFormat, "there is'nt leader uuid in member uuid list")
		return
	}

	spanForConsul := d.tracer.StartSpan("GetNextServiceNode", opentracing.ChildOf(parentSpan))
	selectedNode, err := d.consulAgent.GetNextServiceNode(topic.AuthServiceName)
	spanForConsul.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedNode", selectedNode), log.Error(err))
	spanForConsul.Finish()

	// Handling Response Error of Consul Query
	switch err {
	case nil:
		break
	case consulagent.ErrAvailableNodeNotFound:
		resp.Status = http.StatusServiceUnavailable
		resp.Message = fmt.Sprintf(serviceUnavailableMessageFormat, "there is no available server, name: " + topic.AuthServiceName)
		return
	default:
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "unable to query in consul agent, err: " + err.Error())
		return
	}

	spanForReq := d.tracer.StartSpan("GetStudentInformsWithUUIDs", opentracing.ChildOf(parentSpan))
	md := metadata.Set(context.Background(), "X-Request-Id", reqID)
	md = metadata.Set(md, "Span-Context", spanForReq.Context().(jaeger.SpanContext).String())
	authReq := &authproto.GetStudentInformsWithUUIDsRequest{
		UUID:         req.UUID,
		StudentUUIDs: req.MemberUUIDs,
	}
	respOfReq, err := d.authStudent.GetStudentInformsWithUUIDs(md, authReq)
	spanForReq.SetTag("X-Request-Id", reqID).LogFields(log.Object("request", authReq), log.Object("response", respOfReq), log.Error(err))
	spanForReq.Finish()

	// Handling Response Error of Request
	switch assertedError := err.(type) {
	case nil:
		break
	case *microerrors.Error:
		switch assertedError.Code {
		case http.StatusRequestTimeout:
			resp.Status = http.StatusRequestTimeout
			resp.Message = fmt.Sprintf(requestTimeoutMessageFormat, assertedError.Detail)
			return
		default:
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerMessageFormat, assertedError.Detail)
			return
		}
	default:
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, assertedError.Error())
		return
	}

	// Handling Response Code of Request
	switch respOfReq.Status {
	case http.StatusOK:
		break
	case http.StatusConflict:
		switch respOfReq.Code {
		case code.StudentUUIDsContainNoExistUUID:
			resp.Status = http.StatusConflict
			resp.Code = code.MemberUUIDsIncludeNoExistUUID
			resp.Message = fmt.Sprintf(conflictMessageFormat, "there is not exist uuid in member uuid list")
			return
		default:
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerMessageFormat, fmt.Sprintf("auth service returns unexpected code, code: %d", respOfReq.Code))
			return
		}
	default:
		resp.Status = respOfReq.Status
		resp.Message = fmt.Sprintf("unexpected status returned, message: %s", respOfReq.Message)
		return
	}

	access := d.accessManage.BeginTx()

	cUUID, ok := ctx.Value("ClubUUID").(string)
	if !ok || cUUID == "" {
		cUUID = fmt.Sprintf("club-%s", random.StringConsistOfIntWithLength(12))
	}

	for {
		spanForDB := d.tracer.StartSpan("GetClubWithClubUUID", opentracing.ChildOf(parentSpan))
		selectedClub, err := access.GetClubWithClubUUID(cUUID)
		spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedClub", selectedClub), log.Error(err))
		spanForDB.Finish()
		if err == gorm.ErrRecordNotFound {
			break
		}
		if err != nil {
			access.Rollback()
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerMessageFormat, "unexpected error in GetClubWithClubUUID, err: " + err.Error())
			return
		}
		cUUID = fmt.Sprintf("club-%s", random.StringConsistOfIntWithLength(12))
		continue
	}

	spanForDB := d.tracer.StartSpan("CreateClub", opentracing.ChildOf(parentSpan))
	createdClub, err := access.CreateClub(&model.Club{
		UUID:       model.UUID(cUUID),
		LeaderUUID: model.LeaderUUID(req.LeaderUUID),
	})
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("CreatedClub", createdClub), log.Error(err))
	spanForDB.Finish()

	// Handling Error of CreateClub Method
	switch assertedError := err.(type) {
	case nil:
		break
	case validator.ValidationErrors:
		access.Rollback()
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, "invalid data for club model, err: " + err.Error())
		return
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
			case model.ClubInstance.LeaderUUID.KeyName():
				resp.Status = http.StatusConflict
				resp.Code = code.ClubLeaderAlreadyExist
				resp.Message = fmt.Sprintf(conflictMessageFormat, "club with that leader uuid is already exist, entry: " + entry)
				return
			default:
				resp.Status = http.StatusInternalServerError
				resp.Message = fmt.Sprintf(internalServerMessageFormat, "unexpected duplicate entry, key: " + key)
				return
			}
		default:
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerMessageFormat, "unexpected CreateClub MySQL error code, err: " + assertedError.Error())
			return
		}
	default:
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "unexpected type of CreateClub errors, err: " + assertedError.Error())
		return
	}

	if string(req.Logo) == "" {
		access.Rollback()
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, "Logo attribute cannot be null")
		return
	}

	if d.awsSession != nil {
		spanForS3 := d.tracer.StartSpan("PutObject", opentracing.ChildOf(parentSpan))
		_, err = s3.New(d.awsSession).PutObject(&s3.PutObjectInput{
			Bucket: aws.String("dms-sms"),
			Key:    aws.String(fmt.Sprintf("logos/%s", string(createdClub.UUID))),
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

	spanForDB = d.tracer.StartSpan("CreateClubInform", opentracing.ChildOf(parentSpan))
	createdInform, err := access.CreateClubInform(&model.ClubInform{
		ClubUUID: model.ClubUUID(cUUID),
		Name:     model.Name(req.Name),
		Field:    model.Field(req.Field),
		Location: model.Location(req.Location),
		Floor:    model.Floor(req.Floor),
		LogoURI:  model.LogoURI(fmt.Sprintf("logos/%s", string(createdClub.UUID))),
	})
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("CreatedInform", createdInform), log.Error(err))
	spanForDB.Finish()

	// Handling Error of CreateClubInform Method
	switch assertedError := err.(type) {
	case nil:
		break
	case validator.ValidationErrors:
		access.Rollback()
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, "invalid data for club inform model, err: " + err.Error())
		return
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
			case model.ClubInformInstance.Name.KeyName():
				resp.Status = http.StatusConflict
				resp.Code = code.ClubNameDuplicate
				resp.Message = fmt.Sprintf(conflictMessageFormat, "that club name is alreay exist, entry: " + entry)
				return
			case model.ClubInformInstance.Location.KeyName():
				resp.Status = http.StatusConflict
				resp.Code = code.ClubLocationDuplicate
				resp.Message = fmt.Sprintf(conflictMessageFormat, "that club location is alreay exist, entry: " + entry)
				return
			default:
				resp.Status = http.StatusInternalServerError
				resp.Message = fmt.Sprintf(internalServerMessageFormat, "unexpected duplicate entry, key: " + key)
				return
			}
		default:
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerMessageFormat, "unexpected CreateClubInform MySQL error code, err: " + assertedError.Error())
			return
		}
	default:
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "unexpected type of CreateClubInform errors, err: " + assertedError.Error())
		return
	}
}
