package test

import (
	"club/model"
	authproto "club/proto/golang/auth"
	clubproto "club/proto/golang/club"
	topic "club/utils/topic/golang"
	"context"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/stretchr/testify/mock"
	"log"
	"strconv"
	"strings"
	"time"
)

type AddClubMemberCase struct {
	UUID                  string
	ClubUUID, StudentUUID string
	XRequestID            string
	SpanContextString     string
	ExpectedMethods       map[Method]Returns
	ExpectedStatus        uint32
	ExpectedCode          int32
}

func (test *AddClubMemberCase) ChangeEmptyValueToValidValue() {
	if test.XRequestID == EmptyString        { test.XRequestID = validXRequestID }
	if test.SpanContextString == EmptyString { test.SpanContextString = validSpanContextString }
}

func (test *AddClubMemberCase) ChangeEmptyReplaceValueToEmptyValue() {
	if test.XRequestID == EmptyReplaceValueForString        { test.XRequestID = "" }
	if test.SpanContextString == EmptyReplaceValueForString { test.SpanContextString = "" }
}

func (test *AddClubMemberCase) OnExpectMethodsTo(mock *mock.Mock) {
	for method, returns := range test.ExpectedMethods {
		test.onMethod(mock, method, returns)
	}
}

func (test *AddClubMemberCase) onMethod(mock *mock.Mock, method Method, returns Returns) {
	switch method {
	case "GetClubWithClubUUID":
		mock.On(string(method), test.ClubUUID).Return(returns...)
	case "GetNextServiceNode":
		mock.On(string(method), topic.AuthServiceName).Return(returns...)
	case "GetStudentInformWithUUID":
		mock.On(string(method), &authproto.GetStudentInformWithUUIDRequest{
			UUID:        test.UUID,
			StudentUUID: test.StudentUUID,
		}).Return(returns...)
	case "CreateClubMember":
		const indexForClubMember = 0
		const indexForError = 1
		if _, ok := returns[indexForClubMember].(*model.ClubMember); ok && returns[indexForError] == nil {
			memberForResp := test.getClubMember()
			memberForResp.Model = createGormModelOnCurrentTime()
			returns[indexForClubMember] = memberForResp
		}
		mock.On(string(method), test.getClubMember()).Return(returns...)
	case "BeginTx":
		mock.On(string(method)).Return(returns...)
	case "Commit":
		mock.On(string(method)).Return(returns...)
	case "Rollback":
		mock.On(string(method)).Return(returns...)
	default:
		log.Fatalf("this method cannot be registered, method name: %s", method)
	}
}

func (test *AddClubMemberCase) getClubMember() *model.ClubMember {
	return &model.ClubMember{
		ClubUUID:    model.ClubUUID(test.ClubUUID),
		StudentUUID: model.StudentUUID(test.StudentUUID),
	}
}

func (test *AddClubMemberCase) SetRequestContextOf(req *clubproto.AddClubMemberRequest) {
	req.UUID = test.UUID
	req.ClubUUID = test.ClubUUID
	req.StudentUUID = test.StudentUUID
}

func (test *AddClubMemberCase) GetMetadataContext() (ctx context.Context) {
	ctx = context.Background()
	ctx = metadata.Set(ctx, "X-Request-Id", test.XRequestID)
	ctx = metadata.Set(ctx, "Span-Context", test.SpanContextString)
	return
}

type DeleteClubMemberCase struct {
	UUID                  string
	ClubUUID, StudentUUID string
	XRequestID            string
	SpanContextString     string
	ExpectedMethods       map[Method]Returns
	ExpectedStatus        uint32
	ExpectedCode          int32
}

func (test *DeleteClubMemberCase) ChangeEmptyValueToValidValue() {
	if test.XRequestID == EmptyString        { test.XRequestID = validXRequestID }
	if test.SpanContextString == EmptyString { test.SpanContextString = validSpanContextString }
}

func (test *DeleteClubMemberCase) ChangeEmptyReplaceValueToEmptyValue() {
	if test.XRequestID == EmptyReplaceValueForString        { test.XRequestID = "" }
	if test.SpanContextString == EmptyReplaceValueForString { test.SpanContextString = "" }
}

func (test *DeleteClubMemberCase) OnExpectMethodsTo(mock *mock.Mock) {
	for method, returns := range test.ExpectedMethods {
		test.onMethod(mock, method, returns)
	}
}

func (test *DeleteClubMemberCase) onMethod(mock *mock.Mock, method Method, returns Returns) {
	switch method {
	case "GetClubWithClubUUID":
		mock.On(string(method), test.ClubUUID).Return(returns...)
	case "DeleteClubMember":
		mock.On(string(method), test.ClubUUID, test.StudentUUID).Return(returns...)
	case "BeginTx":
		mock.On(string(method)).Return(returns...)
	case "Commit":
		mock.On(string(method)).Return(returns...)
	case "Rollback":
		mock.On(string(method)).Return(returns...)
	default:
		log.Fatalf("this method cannot be registered, method name: %s", method)
	}
}

func (test *DeleteClubMemberCase) SetRequestContextOf(req *clubproto.DeleteClubMemberRequest) {
	req.UUID = test.UUID
	req.ClubUUID = test.ClubUUID
	req.StudentUUID = test.StudentUUID
}

func (test *DeleteClubMemberCase) GetMetadataContext() (ctx context.Context) {
	ctx = context.Background()
	ctx = metadata.Set(ctx, "X-Request-Id", test.XRequestID)
	ctx = metadata.Set(ctx, "Span-Context", test.SpanContextString)
	return
}

type ChangeClubLeaderCase struct {
	UUID, ClubUUID    string
	NewLeaderUUID     string
	XRequestID        string
	SpanContextString string
	ExpectedMethods   map[Method]Returns
	ExpectedStatus    uint32
	ExpectedCode      int32
}

func (test *ChangeClubLeaderCase) ChangeEmptyValueToValidValue() {
	if test.XRequestID == EmptyString        { test.XRequestID = validXRequestID }
	if test.SpanContextString == EmptyString { test.SpanContextString = validSpanContextString }
}

func (test *ChangeClubLeaderCase) ChangeEmptyReplaceValueToEmptyValue() {
	if test.XRequestID == EmptyReplaceValueForString        { test.XRequestID = "" }
	if test.SpanContextString == EmptyReplaceValueForString { test.SpanContextString = "" }
}

func (test *ChangeClubLeaderCase) OnExpectMethodsTo(mock *mock.Mock) {
	for method, returns := range test.ExpectedMethods {
		test.onMethod(mock, method, returns)
	}
}

func (test *ChangeClubLeaderCase) onMethod(mock *mock.Mock, method Method, returns Returns) {
	switch method {
	case "GetClubWithClubUUID":
		mock.On(string(method), test.ClubUUID).Return(returns...)
	case "GetClubMembersWithClubUUID":
		mock.On(string(method), test.ClubUUID).Return(returns...)
	case "ChangeClubLeader":
		mock.On(string(method), test.ClubUUID, test.NewLeaderUUID).Return(returns...)
	case "BeginTx":
		mock.On(string(method)).Return(returns...)
	case "Commit":
		mock.On(string(method)).Return(returns...)
	case "Rollback":
		mock.On(string(method)).Return(returns...)
	default:
		log.Fatalf("this method cannot be registered, method name: %s", method)
	}
}

func (test *ChangeClubLeaderCase) SetRequestContextOf(req *clubproto.ChangeClubLeaderRequest) {
	req.UUID = test.UUID
	req.ClubUUID = test.ClubUUID
	req.NewLeaderUUID = test.NewLeaderUUID
}

func (test *ChangeClubLeaderCase) GetMetadataContext() (ctx context.Context) {
	ctx = context.Background()
	ctx = metadata.Set(ctx, "X-Request-Id", test.XRequestID)
	ctx = metadata.Set(ctx, "Span-Context", test.SpanContextString)
	return
}

type ModifyClubInformCase struct {
	UUID, ClubUUID    string
	ClubConcept       string
	Introduction      string
	Link              string
	Logo              []byte
	XRequestID        string
	SpanContextString string
	ExpectedMethods   map[Method]Returns
	ExpectedStatus    uint32
	ExpectedCode      int32
}

func (test *ModifyClubInformCase) ChangeEmptyValueToValidValue() {
	if test.UUID == EmptyString              { test.UUID = validLeaderUUID }
	if test.ClubConcept == EmptyString       { test.ClubConcept = validClubConcept }
	if test.Introduction == EmptyString      { test.Introduction = validIntroduction }
	if test.Link == EmptyString              { test.Link = validLink }
	if string(test.Logo) == EmptyString      { test.Logo = validImageByteArr }
	if test.XRequestID == EmptyString        { test.XRequestID = validXRequestID }
	if test.SpanContextString == EmptyString { test.SpanContextString = validSpanContextString }
}

func (test *ModifyClubInformCase) ChangeEmptyReplaceValueToEmptyValue() {
	if test.UUID == EmptyReplaceValueForString              { test.UUID = "" }
	if test.ClubConcept == EmptyReplaceValueForString       { test.ClubConcept = "" }
	if test.Introduction == EmptyReplaceValueForString      { test.Introduction = "" }
	if test.Link == EmptyReplaceValueForString              { test.Link = "" }
	if string(test.Logo) == EmptyReplaceValueForString      { test.Logo = []byte{} }
	if test.XRequestID == EmptyReplaceValueForString        { test.XRequestID = "" }
	if test.SpanContextString == EmptyReplaceValueForString { test.SpanContextString = "" }
}

func (test *ModifyClubInformCase) OnExpectMethodsTo(mock *mock.Mock) {
	for method, returns := range test.ExpectedMethods {
		test.onMethod(mock, method, returns)
	}
}

func (test *ModifyClubInformCase) onMethod(mock *mock.Mock, method Method, returns Returns) {
	switch method {
	case "GetClubWithClubUUID":
		mock.On(string(method), test.ClubUUID).Return(returns...)
	case "ModifyClubInform":
		mock.On(string(method), test.ClubUUID, &model.ClubInform{
			ClubConcept:  model.ClubConcept(test.ClubConcept),
			Introduction: model.Introduction(test.Introduction),
			Link:         model.Link(test.Link),
		}).Return(returns...)
	case "BeginTx":
		mock.On(string(method)).Return(returns...)
	case "Commit":
		mock.On(string(method)).Return(returns...)
	case "Rollback":
		mock.On(string(method)).Return(returns...)
	default:
		log.Fatalf("this method cannot be registered, method name: %s", method)
	}
}

func (test *ModifyClubInformCase) SetRequestContextOf(req *clubproto.ModifyClubInformRequest) {
	req.UUID = test.UUID
	req.ClubUUID = test.ClubUUID
	req.ClubConcept = test.ClubConcept
	req.Introduction = test.Introduction
	req.Link = test.Link
	req.Logo = test.Logo
}

func (test *ModifyClubInformCase) GetMetadataContext() (ctx context.Context) {
	ctx = context.Background()
	ctx = metadata.Set(ctx, "X-Request-Id", test.XRequestID)
	ctx = metadata.Set(ctx, "Span-Context", test.SpanContextString)
	return
}

type DeleteClubWithUUIDCase struct {
	UUID, ClubUUID    string
	XRequestID        string
	SpanContextString string
	ExpectedMethods   map[Method]Returns
	ExpectedStatus    uint32
	ExpectedCode      int32
}

func (test *DeleteClubWithUUIDCase) ChangeEmptyValueToValidValue() {
	if test.XRequestID == EmptyString        { test.XRequestID = validXRequestID }
	if test.SpanContextString == EmptyString { test.SpanContextString = validSpanContextString }
}

func (test *DeleteClubWithUUIDCase) ChangeEmptyReplaceValueToEmptyValue() {
	if test.XRequestID == EmptyReplaceValueForString        { test.XRequestID = "" }
	if test.SpanContextString == EmptyReplaceValueForString { test.SpanContextString = "" }
}

func (test *DeleteClubWithUUIDCase) OnExpectMethodsTo(mock *mock.Mock) {
	for method, returns := range test.ExpectedMethods {
		test.onMethod(mock, method, returns)
	}
}

func (test *DeleteClubWithUUIDCase) onMethod(mock *mock.Mock, method Method, returns Returns) {
	switch method {
	case "GetClubWithClubUUID":
		mock.On(string(method), test.ClubUUID).Return(returns...)
	case "GetCurrentRecruitmentWithClubUUID":
		mock.On(string(method), test.ClubUUID).Return(returns...)
	case "DeleteClub":
		mock.On(string(method), test.ClubUUID).Return(returns...)
	case "DeleteClubInform":
		mock.On(string(method), test.ClubUUID).Return(returns...)
	case "DeleteAllClubMembers":
		mock.On(string(method), test.ClubUUID).Return(returns...)
	case "BeginTx":
		mock.On(string(method)).Return(returns...)
	case "Commit":
		mock.On(string(method)).Return(returns...)
	case "Rollback":
		mock.On(string(method)).Return(returns...)
	default:
		log.Fatalf("this method cannot be registered, method name: %s", method)
	}
}

func (test *DeleteClubWithUUIDCase) SetRequestContextOf(req *clubproto.DeleteClubWithUUIDRequest) {
	req.UUID = test.UUID
	req.ClubUUID = test.ClubUUID
}

func (test *DeleteClubWithUUIDCase) GetMetadataContext() (ctx context.Context) {
	ctx = context.Background()
	ctx = metadata.Set(ctx, "X-Request-Id", test.XRequestID)
	ctx = metadata.Set(ctx, "Span-Context", test.SpanContextString)
	return
}

type RegisterRecruitmentCase struct {
	UUID, ClubUUID     string
	RecruitmentUUID    string
	RecruitmentConcept string
	EndPeriod          string
	RecruitMembers     []*clubproto.RecruitMember
	XRequestID         string
	SpanContextString  string
	ExpectedMethods    map[Method]Returns
	ExpectedStatus     uint32
	ExpectedCode       int32
}

func (test *RegisterRecruitmentCase) ChangeEmptyValueToValidValue() {
	if test.UUID == EmptyString               { test.UUID = validStudentUUID }
	if test.ClubUUID == EmptyString           { test.ClubUUID = validClubUUID }
	if test.RecruitmentUUID == EmptyString    { test.RecruitmentUUID = validRecruitmentUUID }
	if test.RecruitmentConcept == EmptyString { test.UUID = validRecruitConcept }
	if test.EndPeriod == EmptyString          { test.EndPeriod = validEndPeriod }
	if len(test.RecruitMembers) == EmptyInt   { test.RecruitMembers = validRecruitMembers }
	if test.XRequestID == EmptyString         { test.XRequestID = validXRequestID }
	if test.SpanContextString == EmptyString  { test.SpanContextString = validSpanContextString }
}

func (test *RegisterRecruitmentCase) ChangeEmptyReplaceValueToEmptyValue() {
	if test.UUID == EmptyReplaceValueForString               { test.UUID = "" }
	if test.ClubUUID == EmptyReplaceValueForString           { test.ClubUUID = "" }
	if test.RecruitmentUUID == EmptyReplaceValueForString    { test.RecruitmentUUID = "" }
	if test.RecruitmentConcept == EmptyReplaceValueForString { test.UUID = "" }
	if test.EndPeriod == EmptyReplaceValueForString          { test.EndPeriod = "" }
	if len(test.RecruitMembers) == emptyReplaceValueForRecruitMembersLen {
		test.RecruitMembers = []*clubproto.RecruitMember{}
	}
	if test.XRequestID == EmptyReplaceValueForString         { test.XRequestID = "" }
	if test.SpanContextString == EmptyReplaceValueForString  { test.SpanContextString = "" }
}

func (test *RegisterRecruitmentCase) OnExpectMethodsTo(mock *mock.Mock) {
	for method, returns := range test.ExpectedMethods {
		test.onMethod(mock, method, returns)
	}
}

func (test *RegisterRecruitmentCase) onMethod(mock *mock.Mock, method Method, returns Returns) {
	switch method {
	case "GetRecruitmentWithRecruitmentUUID":
		mock.On(string(method), test.RecruitmentUUID).Return(returns...)
	case "GetClubWithClubUUID":
		mock.On(string(method), test.ClubUUID).Return(returns...)
	case "GetCurrentRecruitmentWithClubUUID":
		mock.On(string(method), test.ClubUUID).Return(returns...)
	case "CreateRecruitment":
		const indexForRecruitment = 0
		const indexForError = 1
		if _, ok := returns[indexForRecruitment].(*model.ClubRecruitment); ok && returns[indexForError] == nil {
			recruitment := test.getClubRecruitment()
			recruitment.Model = createGormModelOnCurrentTime()
			returns[indexForRecruitment] = recruitment
		}
		mock.On(string(method), test.getClubRecruitment()).Return(returns...)
	case "CreateRecruitMembers":
		const indexForRecruitMembers = 0
		const indexForError = 1
		for index := range test.RecruitMembers {
			member := test.getRecruitMemberWithIndex(index)
			memberForResp := member
			memberForResp.Model = createGormModelOnCurrentTime()
			mock.On("CreateRecruitMember", member).Return(memberForResp, returns[indexForError])
			if returns[indexForError] != nil {
				break
			}
		}
	case "BeginTx":
		mock.On(string(method)).Return(returns...)
	case "Commit":
		mock.On(string(method)).Return(returns...)
	case "Rollback":
		mock.On(string(method)).Return(returns...)
	default:
		log.Fatalf("this method cannot be registered, method name: %s", method)
	}
}

func (test *RegisterRecruitmentCase) getClubRecruitment() *model.ClubRecruitment {
	recruitment := &model.ClubRecruitment{
		UUID:           model.UUID(test.RecruitmentUUID),
		ClubUUID:       model.ClubUUID(test.ClubUUID),
		RecruitConcept: model.RecruitConcept(test.RecruitmentConcept),
	}
	recruitment.StartPeriod = model.StartPeriod(time.Now())
	endTimeSplice := strings.Split(test.EndPeriod, "-")
	if len(endTimeSplice) == 3 {
		const indexForYear = 0
		const indexForMonth = 1
		const indexForDay = 2
		year, _ := strconv.Atoi(endTimeSplice[indexForYear])
		month, _ := strconv.Atoi(endTimeSplice[indexForMonth])
		day, _ := strconv.Atoi(endTimeSplice[indexForDay])
		recruitment.EndPeriod = model.EndPeriod(time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local))
	} else {
		recruitment.EndPeriod = model.EndPeriod(time.Time{})
	}
	return recruitment
}

func (test *RegisterRecruitmentCase) getRecruitMemberWithIndex(index int) *model.RecruitMember {
	return &model.RecruitMember{
		RecruitmentUUID: model.RecruitmentUUID(test.RecruitmentUUID),
		Grade:           model.Grade(test.RecruitMembers[index].Grade),
		Field:           model.Field(test.RecruitMembers[index].Field),
		Number:          model.Number(test.RecruitMembers[index].Number),
	}
}

func (test *RegisterRecruitmentCase) SetRequestContextOf(req *clubproto.RegisterRecruitmentRequest) {
	req.UUID = test.UUID
	req.ClubUUID = test.ClubUUID
	req.RecruitConcept = test.RecruitmentConcept
	req.RecruitMembers = test.RecruitMembers
	req.EndPeriod = test.EndPeriod
}

func (test *RegisterRecruitmentCase) GetMetadataContext() (ctx context.Context) {
	ctx = context.Background()
	ctx = metadata.Set(ctx, "X-Request-Id", test.XRequestID)
	ctx = metadata.Set(ctx, "Span-Context", test.SpanContextString)
	ctx = metadata.Set(ctx, "RecruitmentUUID", test.RecruitmentUUID)
	return
}