package handler

import (
	clubproto "club/proto/golang/club"
	code "club/utils/code/golang"
	"context"
	"fmt"
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
		resp.Code = code.NotFoundClubNotExists
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
