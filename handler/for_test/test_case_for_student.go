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
			if returns[indexForError] != nil {
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

type GetClubInformsWithUUIDsCase struct {
	UUID              string
	ClubUUIDs         []string
	XRequestID        string
	SpanContextString string
	ExpectedMethods   map[Method]Returns
	ExpectedStatus    uint32
	ExpectedCode      int32
	ExpectInforms     []*clubproto.ClubInform
}

func (test *GetClubInformsWithUUIDsCase) ChangeEmptyValueToValidValue() {
	if test.XRequestID == EmptyString        { test.XRequestID = validXRequestID }
	if test.SpanContextString == EmptyString { test.SpanContextString = validSpanContextString }
}

func (test *GetClubInformsWithUUIDsCase) ChangeEmptyReplaceValueToEmptyValue() {
	if test.XRequestID == EmptyReplaceValueForString        { test.XRequestID = "" }
	if test.SpanContextString == EmptyReplaceValueForString { test.SpanContextString = "" }
}

func (test *GetClubInformsWithUUIDsCase) OnExpectMethodsTo(mock *mock.Mock) {
	for method, returns := range test.ExpectedMethods {
		test.onMethod(mock, method, returns)
	}
}

func (test *GetClubInformsWithUUIDsCase) onMethod(mock *mock.Mock, method Method, returns Returns) {
	switch method {
	case "GetClubWithClubUUIDs":
		const indexForClubs = 0
		const indexForError = 1
		for index, clubUUID := range test.ClubUUIDs {
			mock.On("GetClubWithClubUUID", clubUUID).Return(returns[indexForClubs].([]*model.Club)[index], returns[indexForError])
			if returns[indexForError] != nil {
				break
			}
		}
	case "GetClubInformWithClubUUIDs":
		const indexForClubInforms = 0
		const indexForError = 1
		for index, clubUUID := range test.ClubUUIDs {
			mock.On("GetClubInformWithClubUUID", clubUUID).Return(returns[indexForClubInforms].([]*model.ClubInform)[index], returns[indexForError])
			if returns[indexForError] != nil {
				break
			}
		}
	case "GetClubMembersWithClubUUIDs":
		const indexForClubMembersList = 0
		const indexForError = 1
		for index, clubUUID := range test.ClubUUIDs {
			mock.On("GetClubMembersWithClubUUID", clubUUID).Return(returns[indexForClubMembersList].([][]*model.ClubMember)[index], returns[indexForError])
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

func (test *GetClubInformsWithUUIDsCase) SetRequestContextOf(req *clubproto.GetClubInformsWithUUIDsRequest) {
	req.UUID = test.UUID
	req.ClubUUIDs = test.ClubUUIDs
}

func (test *GetClubInformsWithUUIDsCase) GetMetadataContext() (ctx context.Context) {
	ctx = context.Background()
	ctx = metadata.Set(ctx, "X-Request-Id", test.XRequestID)
	ctx = metadata.Set(ctx, "Span-Context", test.SpanContextString)
	return
}

type GetRecruitmentInformWithUUIDCase struct {
	UUID              string
	RecruitmentUUID   string
	XRequestID        string
	SpanContextString string
	ExpectedMethods   map[Method]Returns
	ExpectedStatus    uint32
	ExpectedCode      int32
	ExpectedRecruit   *clubproto.RecruitmentInform
}

func (test *GetRecruitmentInformWithUUIDCase) ChangeEmptyValueToValidValue() {
	if test.UUID == EmptyString              { test.UUID = validStudentUUID }
	if test.RecruitmentUUID == EmptyString   { test.RecruitmentUUID = validRecruitmentUUID }
	if test.XRequestID == EmptyString        { test.XRequestID = validXRequestID }
	if test.SpanContextString == EmptyString { test.SpanContextString = validSpanContextString }
}

func (test *GetRecruitmentInformWithUUIDCase) ChangeEmptyReplaceValueToEmptyValue() {
	if test.UUID == EmptyReplaceValueForString              { test.UUID = "" }
	if test.RecruitmentUUID == EmptyReplaceValueForString   { test.RecruitmentUUID = "" }
	if test.XRequestID == EmptyReplaceValueForString        { test.XRequestID = "" }
	if test.SpanContextString == EmptyReplaceValueForString { test.SpanContextString = "" }
}

func (test *GetRecruitmentInformWithUUIDCase) OnExpectMethodsTo(mock *mock.Mock) {
	for method, returns := range test.ExpectedMethods {
		test.onMethod(mock, method, returns)
	}
}

func (test *GetRecruitmentInformWithUUIDCase) onMethod(mock *mock.Mock, method Method, returns Returns) {
	switch method {
	case "GetRecruitmentWithRecruitmentUUID":
		mock.On(string(method), test.RecruitmentUUID).Return(returns...)
	case "GetRecruitMembersWithRecruitmentUUID":
		mock.On(string(method), test.RecruitmentUUID).Return(returns...)
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

func (test *GetRecruitmentInformWithUUIDCase) SetRequestContextOf(req *clubproto.GetRecruitmentInformWithUUIDRequest) {
	req.UUID = test.UUID
	req.RecruitmentUUID = test.RecruitmentUUID
}

func (test *GetRecruitmentInformWithUUIDCase) GetMetadataContext() (ctx context.Context) {
	ctx = context.Background()
	ctx = metadata.Set(ctx, "X-Request-Id", test.XRequestID)
	ctx = metadata.Set(ctx, "Span-Context", test.SpanContextString)
	return
}

type GetRecruitmentUUIDWithClubUUIDCase struct {
	UUID                    string
	ClubUUID                string
	XRequestID              string
	SpanContextString       string
	ExpectedMethods         map[Method]Returns
	ExpectedStatus          uint32
	ExpectedCode            int32
	ExpectedRecruitmentUUID string
}

func (test *GetRecruitmentUUIDWithClubUUIDCase) ChangeEmptyValueToValidValue() {
	if test.UUID == EmptyString              { test.UUID = validStudentUUID }
	if test.ClubUUID == EmptyString          { test.ClubUUID = validClubUUID }
	if test.XRequestID == EmptyString        { test.XRequestID = validXRequestID }
	if test.SpanContextString == EmptyString { test.SpanContextString = validSpanContextString }
}

func (test *GetRecruitmentUUIDWithClubUUIDCase) ChangeEmptyReplaceValueToEmptyValue() {
	if test.UUID == EmptyReplaceValueForString              { test.UUID = "" }
	if test.ClubUUID == EmptyReplaceValueForString          { test.ClubUUID = "" }
	if test.XRequestID == EmptyReplaceValueForString        { test.XRequestID = "" }
	if test.SpanContextString == EmptyReplaceValueForString { test.SpanContextString = "" }
}

func (test *GetRecruitmentUUIDWithClubUUIDCase) OnExpectMethodsTo(mock *mock.Mock) {
	for method, returns := range test.ExpectedMethods {
		test.onMethod(mock, method, returns)
	}
}

func (test *GetRecruitmentUUIDWithClubUUIDCase) onMethod(mock *mock.Mock, method Method, returns Returns) {
	switch method {
	case "GetClubWithClubUUID":
		mock.On(string(method), test.ClubUUID).Return(returns...)
	case "GetCurrentRecruitmentWithClubUUID":
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

func (test *GetRecruitmentUUIDWithClubUUIDCase) SetRequestContextOf(req *clubproto.GetRecruitmentUUIDWithClubUUIDRequest) {
	req.UUID = test.UUID
	req.ClubUUID = test.ClubUUID
}

func (test *GetRecruitmentUUIDWithClubUUIDCase) GetMetadataContext() (ctx context.Context) {
	ctx = context.Background()
	ctx = metadata.Set(ctx, "X-Request-Id", test.XRequestID)
	ctx = metadata.Set(ctx, "Span-Context", test.SpanContextString)
	return
}

type GetRecruitmentUUIDsWithClubUUIDsCase struct {
	UUID                     string
	ClubUUIDs                []string
	XRequestID               string
	SpanContextString        string
	ExpectedMethods          map[Method]Returns
	ExpectedStatus           uint32
	ExpectedCode             int32
	ExpectedRecruitmentUUIDs []string
}

func (test *GetRecruitmentUUIDsWithClubUUIDsCase) ChangeEmptyValueToValidValue() {
	if test.XRequestID == EmptyString        { test.XRequestID = validXRequestID }
	if test.SpanContextString == EmptyString { test.SpanContextString = validSpanContextString }
}

func (test *GetRecruitmentUUIDsWithClubUUIDsCase) ChangeEmptyReplaceValueToEmptyValue() {
	if test.XRequestID == EmptyReplaceValueForString        { test.XRequestID = "" }
	if test.SpanContextString == EmptyReplaceValueForString { test.SpanContextString = "" }
}

func (test *GetRecruitmentUUIDsWithClubUUIDsCase) OnExpectMethodsTo(mock *mock.Mock) {
	for method, returns := range test.ExpectedMethods {
		test.onMethod(mock, method, returns)
	}
}

func (test *GetRecruitmentUUIDsWithClubUUIDsCase) onMethod(mock *mock.Mock, method Method, returns Returns) {
	switch method {
	case "GetClubsWithClubUUIDs":
		const indexForClubs = 0
		const indexForError = 1
		for index, clubUUID := range test.ClubUUIDs {
			mock.On("GetClubWithClubUUID", clubUUID).Return(returns[indexForClubs].([]*model.Club)[index], returns[indexForError])
			if returns[indexForError] != nil {
				break
			}
		}
	case "GetCurrentRecruitmentsWithClubUUIDs":
		const indexForRecruitments = 0
		const indexForError = 1
		for index, clubUUID := range test.ClubUUIDs {
			mock.On("GetCurrentRecruitmentWithClubUUID", clubUUID).Return(returns[indexForRecruitments].([]*model.ClubRecruitment)[index], returns[indexForError])
			if !(returns[indexForError] == nil || returns[indexForError] == gorm.ErrRecordNotFound) {
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

func (test *GetRecruitmentUUIDsWithClubUUIDsCase) SetRequestContextOf(req *clubproto.GetRecruitmentUUIDsWithClubUUIDsRequest) {
	req.UUID = test.UUID
	req.ClubUUIDs = test.ClubUUIDs
}

func (test *GetRecruitmentUUIDsWithClubUUIDsCase) GetMetadataContext() (ctx context.Context) {
	ctx = context.Background()
	ctx = metadata.Set(ctx, "X-Request-Id", test.XRequestID)
	ctx = metadata.Set(ctx, "Span-Context", test.SpanContextString)
	return
}

type GetAllClubFieldsCase struct {
	UUID              string
	XRequestID        string
	SpanContextString string
	ExpectedMethods   map[Method]Returns
	ExpectedStatus    uint32
	ExpectedCode      int32
	ExpectedFields    []string
}

func (test *GetAllClubFieldsCase) ChangeEmptyValueToValidValue() {
	if test.XRequestID == EmptyString        { test.XRequestID = validXRequestID }
	if test.SpanContextString == EmptyString { test.SpanContextString = validSpanContextString }
}

func (test *GetAllClubFieldsCase) ChangeEmptyReplaceValueToEmptyValue() {
	if test.XRequestID == EmptyReplaceValueForString        { test.XRequestID = "" }
	if test.SpanContextString == EmptyReplaceValueForString { test.SpanContextString = "" }
}

func (test *GetAllClubFieldsCase) OnExpectMethodsTo(mock *mock.Mock) {
	for method, returns := range test.ExpectedMethods {
		test.onMethod(mock, method, returns)
	}
}

func (test *GetAllClubFieldsCase) onMethod(mock *mock.Mock, method Method, returns Returns) {
	switch method {
	case "GetAllClubInforms":
		mock.On(string(method)).Return(returns...)
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

func (test *GetAllClubFieldsCase) SetRequestContextOf(req *clubproto.GetAllClubFieldsRequest) {
	req.UUID = test.UUID
}

func (test *GetAllClubFieldsCase) GetMetadataContext() (ctx context.Context) {
	ctx = context.Background()
	ctx = metadata.Set(ctx, "X-Request-Id", test.XRequestID)
	ctx = metadata.Set(ctx, "Span-Context", test.SpanContextString)
	return
}

type GetTotalCountOfClubsCase struct {
	UUID              string
	XRequestID        string
	SpanContextString string
	ExpectedMethods   map[Method]Returns
	ExpectedStatus    uint32
	ExpectedCode      int32
	ExpectedCount     int64
}

func (test *GetTotalCountOfClubsCase) ChangeEmptyValueToValidValue() {
	if test.XRequestID == EmptyString        { test.XRequestID = validXRequestID }
	if test.SpanContextString == EmptyString { test.SpanContextString = validSpanContextString }
}

func (test *GetTotalCountOfClubsCase) ChangeEmptyReplaceValueToEmptyValue() {
	if test.XRequestID == EmptyReplaceValueForString        { test.XRequestID = "" }
	if test.SpanContextString == EmptyReplaceValueForString { test.SpanContextString = "" }
}

func (test *GetTotalCountOfClubsCase) OnExpectMethodsTo(mock *mock.Mock) {
	for method, returns := range test.ExpectedMethods {
		test.onMethod(mock, method, returns)
	}
}

func (test *GetTotalCountOfClubsCase) onMethod(mock *mock.Mock, method Method, returns Returns) {
	switch method {
	case "GetAllClubInforms":
		mock.On(string(method)).Return(returns...)
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

func (test *GetTotalCountOfClubsCase) SetRequestContextOf(req *clubproto.GetTotalCountOfClubsRequest) {
	req.UUID = test.UUID
}

func (test *GetTotalCountOfClubsCase) GetMetadataContext() (ctx context.Context) {
	ctx = context.Background()
	ctx = metadata.Set(ctx, "X-Request-Id", test.XRequestID)
	ctx = metadata.Set(ctx, "Span-Context", test.SpanContextString)
	return
}

type GetTotalCountOfCurrentRecruitmentsCase struct {
	UUID              string
	XRequestID        string
	SpanContextString string
	ExpectedMethods   map[Method]Returns
	ExpectedStatus    uint32
	ExpectedCode      int32
	ExpectedCount     int64
}

func (test *GetTotalCountOfCurrentRecruitmentsCase) ChangeEmptyValueToValidValue() {
	if test.XRequestID == EmptyString        { test.XRequestID = validXRequestID }
	if test.SpanContextString == EmptyString { test.SpanContextString = validSpanContextString }
}

func (test *GetTotalCountOfCurrentRecruitmentsCase) ChangeEmptyReplaceValueToEmptyValue() {
	if test.XRequestID == EmptyReplaceValueForString        { test.XRequestID = "" }
	if test.SpanContextString == EmptyReplaceValueForString { test.SpanContextString = "" }
}

func (test *GetTotalCountOfCurrentRecruitmentsCase) OnExpectMethodsTo(mock *mock.Mock) {
	for method, returns := range test.ExpectedMethods {
		test.onMethod(mock, method, returns)
	}
}

func (test *GetTotalCountOfCurrentRecruitmentsCase) onMethod(mock *mock.Mock, method Method, returns Returns) {
	switch method {
	case "GetAllCurrentRecruitments":
		mock.On(string(method)).Return(returns...)
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

func (test *GetTotalCountOfCurrentRecruitmentsCase) SetRequestContextOf(req *clubproto.GetTotalCountOfCurrentRecruitmentsRequest) {
	req.UUID = test.UUID
}

func (test *GetTotalCountOfCurrentRecruitmentsCase) GetMetadataContext() (ctx context.Context) {
	ctx = context.Background()
	ctx = metadata.Set(ctx, "X-Request-Id", test.XRequestID)
	ctx = metadata.Set(ctx, "Span-Context", test.SpanContextString)
	return
}

type GetClubUUIDWithLeaderUUIDCase struct {
	UUID, LeaderUUID  string
	XRequestID        string
	SpanContextString string
	ExpectedMethods   map[Method]Returns
	ExpectedStatus    uint32
	ExpectedCode      int32
	ExpectedClubUUID  string
}

func (test *GetClubUUIDWithLeaderUUIDCase) ChangeEmptyValueToValidValue() {
	if test.XRequestID == EmptyString        { test.XRequestID = validXRequestID }
	if test.SpanContextString == EmptyString { test.SpanContextString = validSpanContextString }
}

func (test *GetClubUUIDWithLeaderUUIDCase) ChangeEmptyReplaceValueToEmptyValue() {
	if test.XRequestID == EmptyReplaceValueForString        { test.XRequestID = "" }
	if test.SpanContextString == EmptyReplaceValueForString { test.SpanContextString = "" }
}

func (test *GetClubUUIDWithLeaderUUIDCase) OnExpectMethodsTo(mock *mock.Mock) {
	for method, returns := range test.ExpectedMethods {
		test.onMethod(mock, method, returns)
	}
}

func (test *GetClubUUIDWithLeaderUUIDCase) onMethod(mock *mock.Mock, method Method, returns Returns) {
	switch method {
	case "GetClubWithLeaderUUID":
		mock.On(string(method), test.LeaderUUID).Return(returns...)
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

func (test *GetClubUUIDWithLeaderUUIDCase) SetRequestContextOf(req *clubproto.GetClubUUIDWithLeaderUUIDRequest) {
	req.UUID = test.UUID
	req.LeaderUUID = test.LeaderUUID
}

func (test *GetClubUUIDWithLeaderUUIDCase) GetMetadataContext() (ctx context.Context) {
	ctx = context.Background()
	ctx = metadata.Set(ctx, "X-Request-Id", test.XRequestID)
	ctx = metadata.Set(ctx, "Span-Context", test.SpanContextString)
	return
}
