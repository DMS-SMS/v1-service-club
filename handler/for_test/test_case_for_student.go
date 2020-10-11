package test

import (
	"club/model"
	clubproto "club/proto/golang/club"
	"context"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"log"
)

type GetClubsSortByUpdateTimeCase struct {
	UUID string
	Field, Name string
	Start, Count uint32
	XRequestID        string
	SpanContextString string
	ExpectedMethods   map[Method]Returns
	ExpectedStatus    uint32
	ExpectedCode      int32
	ExpectClubInforms []*clubproto.ClubInform
}

func (test *GetClubsSortByUpdateTimeCase) ChangeEmptyValueToValidValue() {
	if test.UUID == EmptyString              { test.UUID = validStudentUUID }
	if test.XRequestID == EmptyString        { test.XRequestID = validXRequestID }
	if test.SpanContextString == EmptyString { test.SpanContextString = validSpanContextString }
}

func (test *GetClubsSortByUpdateTimeCase) ChangeEmptyReplaceValueToEmptyValue() {
	if test.UUID == EmptyReplaceValueForString                     { test.UUID = "" }
	if test.XRequestID == EmptyReplaceValueForString               { test.XRequestID = "" }
	if test.SpanContextString == EmptyReplaceValueForString        { test.SpanContextString = "" }
}


func (test *GetClubsSortByUpdateTimeCase) OnExpectMethodsTo(mock *mock.Mock) {
	for method, returns := range test.ExpectedMethods {
		test.onMethod(mock, method, returns)
	}
}

func (test *GetClubsSortByUpdateTimeCase) onMethod(mock *mock.Mock, method Method, returns Returns) {
	switch method {
	case "GetClubInformsSortByUpdateTime":
		const defaultCountValue = 10
		var count = int(test.Count)
		if count == 0 {
			count = defaultCountValue
		}
		mock.On(string(method), int(test.Start), count, test.Field, test.Name).Return(returns...)
	case "GetClubsWithClubUUIDs":
		const indexForClubInforms = 0
		const indexForClubs = 0
		const indexForError = 1
		informs := test.ExpectedMethods["GetClubInformsSortByUpdateTime"][indexForClubInforms].([]*model.ClubInform)
		for index, inform := range informs {
			mock.On("GetClubWithClubUUID", string(inform.ClubUUID)).Return(returns[indexForClubs].([]*model.Club)[index], returns[indexForError])
			if returns[indexForError] != nil && returns[indexForError] != gorm.ErrRecordNotFound {
				break
			}
		}
	case "GetClubMembersWithClubUUIDs":
		const indexForClubInforms = 0
		const indexForCLubMembers = 0
		const indexForError = 1
		informs := test.ExpectedMethods["GetClubInformsSortByUpdateTime"][indexForClubInforms].([]*model.ClubInform)
		for index, inform := range informs {
			mock.On("GetClubMembersWithClubUUID", string(inform.ClubUUID)).Return(returns[indexForCLubMembers].([][]*model.ClubMember)[index], returns[indexForError])
			if returns[indexForError] != nil && returns[indexForError] != gorm.ErrRecordNotFound {
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

func (test *GetClubsSortByUpdateTimeCase) SetRequestContextOf(req *clubproto.GetClubsSortByUpdateTimeRequest) {
	req.UUID = test.UUID
	req.Start = test.Start
	req.Count = test.Count
	req.Field = test.Field
	req.Name = test.Name

}

func (test *GetClubsSortByUpdateTimeCase) GetMetadataContext() (ctx context.Context) {
	ctx = context.Background()
	ctx = metadata.Set(ctx, "X-Request-Id", test.XRequestID)
	ctx = metadata.Set(ctx, "Span-Context", test.SpanContextString)
	return
}

type GetRecruitmentsSortByCreateTimeCase struct {
	UUID               string
	Field, Name        string
	Start, Count       uint32
	XRequestID         string
	SpanContextString  string
	ExpectedMethods    map[Method]Returns
	ExpectedStatus     uint32
	ExpectedCode       int32
	ExpectRecruitments []*clubproto.RecruitmentInform
}

func (test *GetRecruitmentsSortByCreateTimeCase) ChangeEmptyValueToValidValue() {
	if test.UUID == EmptyString              { test.UUID = validStudentUUID }
	if test.XRequestID == EmptyString        { test.XRequestID = validXRequestID }
	if test.SpanContextString == EmptyString { test.SpanContextString = validSpanContextString }
}

func (test *GetRecruitmentsSortByCreateTimeCase) ChangeEmptyReplaceValueToEmptyValue() {
	if test.UUID == EmptyReplaceValueForString                     { test.UUID = "" }
	if test.XRequestID == EmptyReplaceValueForString               { test.XRequestID = "" }
	if test.SpanContextString == EmptyReplaceValueForString        { test.SpanContextString = "" }
}

func (test *GetRecruitmentsSortByCreateTimeCase) OnExpectMethodsTo(mock *mock.Mock) {
	for method, returns := range test.ExpectedMethods {
		test.onMethod(mock, method, returns)
	}
}

func (test *GetRecruitmentsSortByCreateTimeCase) onMethod(mock *mock.Mock, method Method, returns Returns) {
	switch method {
	case "GetCurrentRecruitmentsSortByCreateTime":
		const defaultCountValue = 10
		var count = int(test.Count)
		if count == 0 {
			count = defaultCountValue
		}
		mock.On(string(method), int(test.Start), count, test.Field, test.Name).Return(returns...)
	case "GetRecruitMembersWithRecruitmentUUIDs":
		const indexForRecruitments = 0
		const indexForRecruitMembersList = 0
		const indexForError = 1
		recruitments := test.ExpectedMethods["GetCurrentRecruitmentsSortByCreateTime"][indexForRecruitments].([]*model.ClubRecruitment)
		for index, recruitment := range recruitments {
			mock.On("GetRecruitMembersWithRecruitmentUUID", string(recruitment.UUID)).Return(returns[indexForRecruitMembersList].([][]*model.RecruitMember)[index], returns[indexForError])
			if returns[indexForError] != nil && returns[indexForError] != gorm.ErrRecordNotFound {
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

func (test *GetRecruitmentsSortByCreateTimeCase) SetRequestContextOf(req *clubproto.GetRecruitmentsSortByCreateTimeRequest) {
	req.UUID = test.UUID
	req.Start = test.Start
	req.Count = test.Count
	req.Field = test.Field
	req.Name = test.Name

}

func (test *GetRecruitmentsSortByCreateTimeCase) GetMetadataContext() (ctx context.Context) {
	ctx = context.Background()
	ctx = metadata.Set(ctx, "X-Request-Id", test.XRequestID)
	ctx = metadata.Set(ctx, "Span-Context", test.SpanContextString)
	return
}

type GetClubInformWithUUIDCase struct {
	UUID              string
	ClubUUID          string
	XRequestID        string
	SpanContextString string
	ExpectedMethods   map[Method]Returns
	ExpectedStatus    uint32
	ExpectedCode      int32
	ExpectInform      *clubproto.ClubInform
}

func (test *GetClubInformWithUUIDCase) ChangeEmptyValueToValidValue() {
	if test.UUID == EmptyString              { test.UUID = validStudentUUID }
	if test.ClubUUID == EmptyString          { test.ClubUUID = validClubUUID }
	if test.XRequestID == EmptyString        { test.XRequestID = validXRequestID }
	if test.SpanContextString == EmptyString { test.SpanContextString = validSpanContextString }
}

func (test *GetClubInformWithUUIDCase) ChangeEmptyReplaceValueToEmptyValue() {
	if test.UUID == EmptyReplaceValueForString              { test.UUID = "" }
	if test.ClubUUID == EmptyReplaceValueForString          { test.ClubUUID = "" }
	if test.XRequestID == EmptyReplaceValueForString        { test.XRequestID = "" }
	if test.SpanContextString == EmptyReplaceValueForString { test.SpanContextString = "" }
}

func (test *GetClubInformWithUUIDCase) OnExpectMethodsTo(mock *mock.Mock) {
	for method, returns := range test.ExpectedMethods {
		test.onMethod(mock, method, returns)
	}
}

func (test *GetClubInformWithUUIDCase) onMethod(mock *mock.Mock, method Method, returns Returns) {
	switch method {
	case "GetClubWithClubUUID":
		mock.On(string(method), test.ClubUUID).Return(returns...)
	case "GetClubInformWithClubUUID":
		mock.On(string(method), test.ClubUUID).Return(returns...)
	case "GetClubMembersWithClubUUID":
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

func (test *GetClubInformWithUUIDCase) SetRequestContextOf(req *clubproto.GetClubInformWithUUIDRequest) {
	req.UUID = test.UUID
	req.ClubUUID = test.ClubUUID
}

func (test *GetClubInformWithUUIDCase) GetMetadataContext() (ctx context.Context) {
	ctx = context.Background()
	ctx = metadata.Set(ctx, "X-Request-Id", test.XRequestID)
	ctx = metadata.Set(ctx, "Span-Context", test.SpanContextString)
	return
}
