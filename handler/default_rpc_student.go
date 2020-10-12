package handler

import (
	"club/model"
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

	informsForResp := make([]*clubproto.ClubInform, len(selectedInforms))
	for index, selectedInform := range selectedInforms {
		informsForResp[index] = &clubproto.ClubInform{
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

	spanForDB = d.tracer.StartSpan("GetClubsWithClubUUIDs", opentracing.ChildOf(parentSpan))
	selectedClubs := make([]*model.Club, len(informsForResp))
	for index, informForResp := range informsForResp {
		selectedClub, queryErr := access.GetClubWithClubUUID(informForResp.ClubUUID)
		selectedClubs[index] = selectedClub
		if queryErr != nil {
			err = queryErr
		}
	}
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedClubs", selectedClubs), log.Error(err))
	spanForDB.Finish()

	if err != nil {
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "GetClubsWithClubUUIDs return errors, err: " + err.Error())
		return
	}

	for index, club := range selectedClubs {
		informsForResp[index].LeaderUUID = string(club.LeaderUUID)
	}

	spanForDB = d.tracer.StartSpan("GetClubMembersWithClubUUID", opentracing.ChildOf(parentSpan))
	selectedMembersList := make([][]*model.ClubMember, len(informsForResp))
	for index, informForResp := range informsForResp {
		selectedMembers, queryErr := access.GetClubMembersWithClubUUID(informForResp.ClubUUID)
		if queryErr != nil {
			err = queryErr
		}
		membersForResp := make([]string, len(selectedMembers))
		for index, selectedMember := range selectedMembers {
			membersForResp[index] = string(selectedMember.StudentUUID)
		}
		informsForResp[index].MemberUUIDs = membersForResp
		selectedMembersList[index] = selectedMembers
	}
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedMembersList", selectedMembersList), log.Error(err))
	spanForDB.Finish()

	if err != nil && err != gorm.ErrRecordNotFound {
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "GetClubMembersWithClubUUID returns unexpected error, err: " + err.Error())
		return
	}

	access.Commit()
	resp.Status = http.StatusOK
	resp.Informs = informsForResp
	resp.Message = "get clubs success"

	return
}

func (d *_default) GetRecruitmentsSortByCreateTime(ctx context.Context, req *clubproto.GetRecruitmentsSortByCreateTimeRequest, resp *clubproto.GetRecruitmentsSortByCreateTimeResponse) (_ error) {
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
	spanForDB := d.tracer.StartSpan("GetCurrentRecruitmentsSortByCreateTime", opentracing.ChildOf(parentSpan))
	selectedRecruits, err := access.GetCurrentRecruitmentsSortByCreateTime(int(req.Start), int(req.Count), req.Field, req.Name)
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedRecruits", selectedRecruits), log.Error(err))
	spanForDB.Finish()

	switch err {
	case nil:
		break
	case gorm.ErrRecordNotFound:
		access.Commit()
		resp.Status = http.StatusOK
		resp.Message = "get recruitments success (result not exist)"
		return
	default:
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "GetCurrentRecruitmentsSortByCreateTime returns unexpected error, error: " + err.Error())
		return
	}

	recruitmentsForResp := make([]*clubproto.RecruitmentInform, len(selectedRecruits))
	for index, selectedRecruit := range selectedRecruits {
		recruit := &clubproto.RecruitmentInform{
			RecruitmentUUID: string(selectedRecruit.UUID),
			ClubUUID:        string(selectedRecruit.ClubUUID),
			RecruitConcept:  string(selectedRecruit.RecruitConcept),
		}
		startTime, _ := selectedRecruit.StartPeriod.Value()
		if timeString, ok := startTime.(string); ok {
			recruit.StartPeriod = timeString
		}
		endTime, _ := selectedRecruit.EndPeriod.Value()
		if timeString, ok := endTime.(string); ok {
			recruit.EndPeriod = timeString
		}
		recruitmentsForResp[index] = recruit
	}

	spanForDB = d.tracer.StartSpan("GetRecruitMembersListWithRecruitmentUUIDs", opentracing.ChildOf(parentSpan))
	selectedMembersList := make([][]*model.RecruitMember, len(recruitmentsForResp))
	for index, recruitmentForResp := range recruitmentsForResp {
		selectedMembers, queryErr := access.GetRecruitMembersWithRecruitmentUUID(recruitmentForResp.RecruitmentUUID)
		if queryErr != nil {
			err = queryErr
		}
		membersForResp := make([]*clubproto.RecruitMember, len(selectedMembers))
		for index, selectedMember := range selectedMembers {
			membersForResp[index] = &clubproto.RecruitMember{
				Grade:  string(selectedMember.Grade),
				Field:  string(selectedMember.Field),
				Number: string(selectedMember.Number),
			}
		}
		recruitmentsForResp[index].RecruitMembers = membersForResp
		selectedMembersList[index] = selectedMembers
	}
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedMembersList", selectedMembersList), log.Error(err))
	spanForDB.Finish()

	if err != nil && err != gorm.ErrRecordNotFound {
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "GetRecruitMembersWithRecruitmentUUID returns unexpected error, err: " + err.Error())
		return
	}

	access.Commit()
	resp.Status = http.StatusOK
	resp.Message = fmt.Sprintf("get recruitments success (len: %d)", len(recruitmentsForResp))
	resp.Recruitments = recruitmentsForResp
	return
}

func (d *_default) GetClubInformWithUUID(ctx context.Context, req *clubproto.GetClubInformWithUUIDRequest, resp *clubproto.GetClubInformWithUUIDResponse) (_ error) {
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
		resp.Message = fmt.Sprintf(notFoundMessageFormat, "that club uuid is not exist")
		return
	default:
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "GetClubWithClubUUID returns unexpected error, err: " + err.Error())
		return
	}

	spanForDB = d.tracer.StartSpan("GetClubInformWithClubUUID", opentracing.ChildOf(parentSpan))
	selectedInform, err := access.GetClubInformWithClubUUID(string(selectedClub.UUID))
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedInform", selectedInform), log.Error(err))
	spanForDB.Finish()

	if err != nil {
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "GetClubInformWithClubUUID returns unexpected error, err: " + err.Error())
		return
	}

	spanForDB = d.tracer.StartSpan("GetClubMembersWithClubUUID", opentracing.ChildOf(parentSpan))
	selectedMembers, err := access.GetClubMembersWithClubUUID(string(selectedClub.UUID))
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedMembers", selectedMembers), log.Error(err))
	spanForDB.Finish()

	if err != nil {
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "GetClubMembersWithClubUUID returns unexpected error, err: " + err.Error())
		return
	}

	access.Commit()
	membersForResp := make([]string, len(selectedMembers))
	for index, selectedMember := range selectedMembers {
		membersForResp[index] = string(selectedMember.StudentUUID)
	}
	resp.Status = http.StatusOK
	resp.ClubUUID = string(selectedClub.UUID)
	resp.LeaderUUID = string(selectedClub.LeaderUUID)
	resp.MemberUUIDs = membersForResp
	resp.Name = string(selectedInform.Name)
	resp.ClubConcept = string(selectedInform.ClubConcept)
	resp.Introduction = string(selectedInform.Introduction)
	resp.Floor = string(selectedInform.Floor)
	resp.Location = string(selectedInform.Location)
	resp.Field = string(selectedInform.Field)
	resp.Link = string(selectedInform.Link)
	resp.LogoURI = string(selectedInform.LogoURI)
	resp.Message = "get club inform success"

	return
}
