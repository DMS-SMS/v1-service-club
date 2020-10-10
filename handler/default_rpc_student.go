package handler

import (
	clubproto "club/proto/golang/club"
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	"gorm.io/gorm"
	"net/http"
)

const defaultCountValue = 10

func (d *_default) GetClubsSortByUpdateTime(ctx context.Context, req *clubproto.GetClubsSortByUpdateTimeRequest, resp *clubproto.GetClubsSortByUpdateTimeResponse) (_ error) {
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
		resp.Message = fmt.Sprintf(forbiddenMessageFormat, "you are not student or admin")
		return
	}

	reqID := ctx.Value("X-Request-Id").(string)
	parentSpan := ctx.Value("Span-Context").(jaeger.SpanContext)

	access := d.accessManage.BeginTx()

	if req.Count == 0 { req.Count = defaultCountValue }
	spanForDB := d.tracer.StartSpan("GetClubInformsSortByUpdateTime", opentracing.ChildOf(parentSpan))
	selectedInforms, err := access.GetClubInformsSortByUpdateTime(int(req.Start), int(req.Count), req.Field, req.Name)
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedInforms", selectedInforms), log.Error(err))
	spanForDB.Finish()

	switch err {
	case nil:
		break
	case gorm.ErrRecordNotFound:
		access.Commit()
		resp.Status = http.StatusOK
		resp.Message = "get clubs success (result not exist)"
		return
	default:
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "GetClubInformsSortByUpdateTime returns unexpected error, error: " + err.Error())
		return
	}

	informs := make([]*clubproto.ClubInform, len(selectedInforms))
	for index, selectedInform := range selectedInforms {
		informs[index] = &clubproto.ClubInform{
			ClubUUID:     string(selectedInform.ClubUUID),
			Name:         string(selectedInform.Name),
			ClubConcept:  string(selectedInform.ClubConcept),
			Introduction: string(selectedInform.Introduction),
			Field:        string(selectedInform.Field),
			Location:     string(selectedInform.Location),
			Floor:        string(selectedInform.Floor),
			Link:         string(selectedInform.Link),
			LogoURI:      string(selectedInform.LogoURI),
		}
	}

	clubUUIDs := make([]string, len(informs))
	for index, inform := range informs {
		clubUUIDs[index] = inform.ClubUUID
	}

	spanForDB = d.tracer.StartSpan("GetClubsWithClubUUIDs", opentracing.ChildOf(parentSpan))
	selectedClubs, err := access.GetClubsWithClubUUIDs(clubUUIDs)
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedClubs", selectedClubs), log.Error(err))
	spanForDB.Finish()

	if err != nil {
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "GetClubsWithClubUUIDs return errors, err: " + err.Error())
		return
	}

	if len(selectedClubs) != len(informs) {
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "GetClubsWithClubUUIDs return abnormal length slice")
		return
	}

	for index, club := range selectedClubs {
		informs[index].LeaderUUID = string(club.LeaderUUID)
	}

	access.Commit()
	resp.Status = http.StatusOK
	resp.Informs = informs
	resp.Message = "get clubs success"

	return
}
