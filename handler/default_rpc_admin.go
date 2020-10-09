package handler

import (
	authproto "club/proto/golang/auth"
	clubproto "club/proto/golang/club"
	consulagent "club/tool/consul/agent"
	code "club/utils/code/golang"
	topic "club/utils/topic/golang"
	"context"
	"fmt"
	microerrors "github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
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
}
