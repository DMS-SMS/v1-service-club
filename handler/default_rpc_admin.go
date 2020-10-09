package handler

import (
	clubproto "club/proto/golang/club"
	consulagent "club/tool/consul/agent"
	code "club/utils/code/golang"
	topic "club/utils/topic/golang"
	"context"
	"fmt"
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
}
