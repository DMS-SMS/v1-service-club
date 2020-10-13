package handler

import (
	"club/model"
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
		if queryErr != nil {
			err = queryErr
			break
		}
		informsForResp[index].LeaderUUID = string(selectedClub.LeaderUUID)
		selectedClubs[index] = selectedClub
	}
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedClubs", selectedClubs), log.Error(err))
	spanForDB.Finish()

	if err != nil {
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "GetClubsWithClubUUIDs return errors, err: " + err.Error())
		return
	}

	spanForDB = d.tracer.StartSpan("GetClubMembersListWithClubUUIDs", opentracing.ChildOf(parentSpan))
	selectedMembersList := make([][]*model.ClubMember, len(informsForResp))
	for index, informForResp := range informsForResp {
		selectedMembers, queryErr := access.GetClubMembersWithClubUUID(informForResp.ClubUUID)
		if queryErr == gorm.ErrRecordNotFound {
			err = queryErr
		} else if queryErr != nil {
			err = queryErr
			break
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
		if queryErr == gorm.ErrRecordNotFound {
			err = queryErr
		} else if queryErr != nil {
			err = queryErr
			break
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

	if err != nil && err != gorm.ErrRecordNotFound {
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

func (d *_default) GetClubInformsWithUUIDs(ctx context.Context, req *clubproto.GetClubInformsWithUUIDsRequest, resp *clubproto.GetClubInformsWithUUIDsResponse) (_ error) {
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
	informsForResp := make([]*clubproto.ClubInform, len(req.ClubUUIDs))

	var err error
	selectedClubs := make([]*model.Club, len(req.ClubUUIDs))
	spanForDB := d.tracer.StartSpan("GetClubsWithClubUUIDs", opentracing.ChildOf(parentSpan))
	for index, clubUUID := range req.ClubUUIDs {
		selectedClub, queryErr := access.GetClubWithClubUUID(clubUUID)
		if queryErr != nil {
			err = queryErr
			break
		}
		informsForResp[index] = new(clubproto.ClubInform)
		informsForResp[index].ClubUUID = string(selectedClub.UUID)
		informsForResp[index].LeaderUUID = string(selectedClub.LeaderUUID)
		selectedClubs[index] = selectedClub
	}
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedClubs", selectedClubs), log.Error(err))
	spanForDB.Finish()

	switch err {
	case nil:
		break
	case gorm.ErrRecordNotFound:
		access.Rollback()
		resp.Status = http.StatusNotFound
		resp.Message = fmt.Sprintf(notFoundMessageFormat, "club uuid list include not exist uuid")
		return
	default:
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "GetClubWithClubUUID returns unexpected error, err: " + err.Error())
		return
	}

	selectedInforms := make([]*model.ClubInform, len(req.ClubUUIDs))
	spanForDB = d.tracer.StartSpan("GetClubInformsWithClubUUIDs", opentracing.ChildOf(parentSpan))
	for index, clubUUID := range req.ClubUUIDs {
		selectedInform, queryErr := access.GetClubInformWithClubUUID(clubUUID)
		if queryErr != nil {
			err = queryErr
			break
		}
		informsForResp[index].Name = string(selectedInform.Name)
		informsForResp[index].ClubConcept = string(selectedInform.ClubConcept)
		informsForResp[index].Introduction = string(selectedInform.Introduction)
		informsForResp[index].Floor = string(selectedInform.Floor)
		informsForResp[index].Location = string(selectedInform.Location)
		informsForResp[index].Field = string(selectedInform.Field)
		informsForResp[index].Link = string(selectedInform.Link)
		informsForResp[index].LogoURI = string(selectedInform.LogoURI)
		selectedInforms[index] = selectedInform
	}
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedInforms", selectedInforms), log.Error(err))
	spanForDB.Finish()

	if err != nil {
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "GetClubInformsWithClubUUIDs returns unexepcted error, err: " + err.Error())
		return
	}

	selectedMembersList := make([][]*model.ClubMember, len(req.ClubUUIDs))
	spanForDB = d.tracer.StartSpan("GetClubMembersListWithClubUUIDs", opentracing.ChildOf(parentSpan))
	for index, clubUUID := range req.ClubUUIDs {
		selectedMembers, queryErr := access.GetClubMembersWithClubUUID(clubUUID)
		if queryErr == gorm.ErrRecordNotFound {
			err = queryErr
		} else if queryErr != nil {
			err = queryErr
			break
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
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "GetClubMembersListWithClubUUIDs returns unexepcted error, err: " + err.Error())
		return
	}

	access.Commit()
	resp.Status = http.StatusOK
	resp.Informs = informsForResp
	resp.Message = fmt.Sprintf("get club informs success")
	return
}

func (d *_default) GetRecruitmentInformWithUUID(ctx context.Context, req *clubproto.GetRecruitmentInformWithUUIDRequest, resp *clubproto.GetRecruitmentInformWithUUIDResponse) (_ error) {
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

	spanForDB := d.tracer.StartSpan("GetRecruitmentWithRecruitmentUUID", opentracing.ChildOf(parentSpan))
	selectedRecruitment, err := access.GetRecruitmentWithRecruitmentUUID(req.RecruitmentUUID)
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedRecruitment", selectedRecruitment), log.Error(err))
	spanForDB.Finish()

	switch err {
	case nil:
		break
	case gorm.ErrRecordNotFound:
		access.Rollback()
		resp.Status = http.StatusNotFound
		resp.Message = fmt.Sprintf(notFoundMessageFormat, "that recruitment uuid is not exist")
		return
	default:
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "GetRecruitmentWithRecruitmentUUID returns unexpected error, err: " + err.Error())
		return
	}

	spanForDB = d.tracer.StartSpan("GetClubMembersWithClubUUID", opentracing.ChildOf(parentSpan))
	selectedMembers, err := access.GetRecruitMembersWithRecruitmentUUID(req.RecruitmentUUID)
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedMembers", selectedMembers), log.Error(err))
	spanForDB.Finish()

	if err != nil && err != gorm.ErrRecordNotFound {
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "GetClubMembersWithClubUUID returns unexpected error, err: " + err.Error())
		return
	}

	membersForResp := make([]*clubproto.RecruitMember, len(selectedMembers))
	for index, selectedMember := range selectedMembers {
		memberForResp := &clubproto.RecruitMember{
			Grade:  string(selectedMember.Grade),
			Field:  string(selectedMember.Field),
			Number: string(selectedMember.Number),
		}
		membersForResp[index] = memberForResp
	}

	access.Commit()
	resp.Status = http.StatusOK
	resp.RecruitmentUUID = string(selectedRecruitment.UUID)
	resp.ClubUUID = string(selectedRecruitment.ClubUUID)
	resp.RecruitConcept = string(selectedRecruitment.RecruitConcept)
	resp.RecruitMembers = membersForResp
	startTime, _ := selectedRecruitment.StartPeriod.Value()
	if timeString, ok := startTime.(string); ok {
		resp.StartPeriod = timeString
	}
	endTime, _ := selectedRecruitment.EndPeriod.Value()
	if timeString, ok := endTime.(string); ok {
		resp.EndPeriod = timeString
	}
	resp.Message = "get club inform success"
	return
}

func (d *_default) GetRecruitmentUUIDWithClubUUID(ctx context.Context, req *clubproto.GetRecruitmentUUIDWithClubUUIDRequest, resp *clubproto.GetRecruitmentUUIDWithClubUUIDResponse) (_ error) {
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

	spanForDB = d.tracer.StartSpan("GetCurrentRecruitmentWithClubUUID", opentracing.ChildOf(parentSpan))
	selectedRecruitment, err := access.GetCurrentRecruitmentWithClubUUID(req.ClubUUID)
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedRecruitment", selectedRecruitment), log.Error(err))
	spanForDB.Finish()

	switch err {
	case nil:
		break
	case gorm.ErrRecordNotFound:
		access.Rollback()
		resp.Status = http.StatusConflict
		resp.Code = code.ThereIsNoCurrentRecruitment
		resp.Message = fmt.Sprintf(conflictMessageFormat, "no recruitment is currently in progress with that club uuid")
		return
	default:
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "GetCurrentRecruitmentWithClubUUID returns unexpected error, err: " + err.Error())
		return
	}

	access.Commit()
	resp.Status = http.StatusOK
	resp.RecruitmentUUID = string(selectedRecruitment.UUID)
	resp.Message = "get recruitment in progress success"
	return
}

func (d *_default) GetRecruitmentUUIDsWithClubUUIDs(ctx context.Context, req *clubproto.GetRecruitmentUUIDsWithClubUUIDsRequest, resp *clubproto.GetRecruitmentUUIDsWithClubUUIDsResponse) (_ error) {
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

	var err error
	selectedClubs := make([]*model.Club, len(req.ClubUUIDs))
	spanForDB := d.tracer.StartSpan("GetClubsWithClubUUIDs", opentracing.ChildOf(parentSpan))
	for index, clubUUID := range req.ClubUUIDs {
		selectedClub, queryErr := access.GetClubWithClubUUID(clubUUID)
		if queryErr != nil {
			err = queryErr
			break
		}
		selectedClubs[index] = selectedClub
	}
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedClubs", selectedClubs), log.Error(err))
	spanForDB.Finish()

	switch err {
	case nil:
		break
	case gorm.ErrRecordNotFound:
		access.Rollback()
		resp.Status = http.StatusNotFound
		resp.Message = fmt.Sprintf(notFoundMessageFormat, "club uuid list include not exist club")
		return
	default:
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "GetClubsWithClubUUIDs returns unexpected error, err: " + err.Error())
		return
	}

	recruitmentUUIDsForResp := make([]string, len(req.ClubUUIDs))
	selectedRecruits := make([]*model.ClubRecruitment, len(req.ClubUUIDs))
	spanForDB = d.tracer.StartSpan("GetCurrentRecruitmentsWithClubUUIDs", opentracing.ChildOf(parentSpan))
	for index, clubUUID := range req.ClubUUIDs {
		selectedRecruit, queryErr := access.GetCurrentRecruitmentWithClubUUID(clubUUID)
		if queryErr == gorm.ErrRecordNotFound {
			err = queryErr
		} else if queryErr != nil {
			err = queryErr
			break
		}
		selectedRecruits[index] = selectedRecruit
		if selectedRecruit == nil {
			selectedRecruit = &model.ClubRecruitment{}
		}
		recruitmentUUIDsForResp[index] = string(selectedRecruit.UUID)
	}
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedRecruits", selectedRecruits), log.Error(err))
	spanForDB.Finish()

	if err != nil && err != gorm.ErrRecordNotFound {
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "GetCurrentRecruitmentsWithClubUUIDs returns unexpected error, err:" + err.Error())
		return
	}

	access.Commit()
	resp.Status = http.StatusOK
	resp.RecruitmentUUIDs = recruitmentUUIDsForResp
	resp.Message = "get recruitment list in progress success"
	return
}

func (d *_default) GetAllClubFields(ctx context.Context, req *clubproto.GetAllClubFieldsRequest, resp *clubproto.GetAllClubFieldsResponse) (_ error) {
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

	spanForDB := d.tracer.StartSpan("GetAllClubInforms", opentracing.ChildOf(parentSpan))
	allInforms, err := access.GetAllClubInforms()
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("AllInforms", allInforms), log.Error(err))
	spanForDB.Finish()

	if err != nil && err != gorm.ErrRecordNotFound {
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "GetAllClubInforms returns unexpected error, err: " + err.Error())
		return
	}

	fieldsForResp := make([]string, 0, 5)
	for _, inform := range allInforms {
		if contains(fieldsForResp, string(inform.Field)) {
			continue
		}
		fieldsForResp = append(fieldsForResp, string(inform.Field))
	}

	access.Commit()
	resp.Status = http.StatusOK
	resp.Fields = fieldsForResp
	resp.Message = fmt.Sprintf("get all club field success")
	return
}

func (d *_default) GetTotalCountOfClubs(ctx context.Context, req *clubproto.GetTotalCountOfClubsRequest, resp *clubproto.GetTotalCountOfClubsResponse) (_ error) {
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

	spanForDB := d.tracer.StartSpan("GetAllClubInforms", opentracing.ChildOf(parentSpan))
	allInforms, err := access.GetAllClubInforms()
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("AllInforms", allInforms), log.Error(err))
	spanForDB.Finish()

	if err != nil && err != gorm.ErrRecordNotFound {
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "GetAllClubInforms returns unexpected error, err: " + err.Error())
		return
	}

	access.Commit()
	resp.Status = http.StatusOK
	resp.Count = int64(len(allInforms))
	resp.Message = fmt.Sprintf("get total count of club success")
	return
}

func (d *_default) GetTotalCountOfCurrentRecruitments(ctx context.Context, req *clubproto.GetTotalCountOfCurrentRecruitmentsRequest, resp *clubproto.GetTotalCountOfCurrentRecruitmentsResponse) (_ error) {
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

	spanForDB := d.tracer.StartSpan("GetAllCurrentRecruitments", opentracing.ChildOf(parentSpan))
	allRecruitments, err := access.GetAllCurrentRecruitments()
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("AllRecruitments", allRecruitments), log.Error(err))
	spanForDB.Finish()

	if err != nil && err != gorm.ErrRecordNotFound {
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerMessageFormat, "GetAllCurrentRecruitments returns unexpected error, err: " + err.Error())
		return
	}

	access.Commit()
	resp.Status = http.StatusOK
	resp.Count = int64(len(allRecruitments))
	resp.Message = fmt.Sprintf("get total count of current recruitment success")
	return
}
