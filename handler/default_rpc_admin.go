package handler

import (
	clubproto "club/proto/golang/club"
	code "club/utils/code/golang"
	"context"
	"fmt"
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

}
