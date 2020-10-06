package test

import (
	"club/model"
	authproto "club/proto/golang/auth"
	clubproto "club/proto/golang/club"
	"context"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/stretchr/testify/mock"
	"log"
)

type Method string
type Returns []interface{}

type CreateNewClubCase struct {
	UUID, LeaderUUID  string
	Name, Field       string
	MemberUUIDs       []string
	Floor, Location   string
	Logo              []byte
	ClubUUID          string
	XRequestID        string
	SpanContextString string
	ExpectedMethods   map[Method]Returns
	ExpectedStatus    uint32
	ExpectedCode      int32
	ExpectedClubUUID  string
}

func (test *CreateNewClubCase) ChangeEmptyValueToValidValue() {
	if test.UUID == EmptyString              { test.UUID = validAdminUUID }
	if test.Name == EmptyString              { test.Name = validClubName }
	if test.LeaderUUID == EmptyString        { test.LeaderUUID = validLeaderUUID }
	if len(test.MemberUUIDs) == EmptyInt     { test.MemberUUIDs = validMemberUUIDs }
	if test.Field == EmptyString             { test.Field = validField }
	if test.Location == EmptyString          { test.Location = validLocation }
	if test.Floor == EmptyString             { test.Floor = validFloor }
	if string(test.Logo) == EmptyString      { test.Logo = validImageByteArr }
	if test.ClubUUID == EmptyString          { test.ClubUUID = validClubUUID }
	if test.XRequestID == EmptyString        { test.XRequestID = validXRequestID }
	if test.SpanContextString == EmptyString { test.SpanContextString = validSpanContextString }
}

func (test *CreateNewClubCase) ChangeEmptyReplaceValueToEmptyValue() {
	if test.UUID == EmptyReplaceValueForString                     { test.UUID = "" }
	if test.Name == EmptyReplaceValueForString                     { test.Name = "" }
	if test.LeaderUUID == EmptyReplaceValueForString               { test.LeaderUUID = "" }
	if len(test.MemberUUIDs) == emptyReplaceValueForMemberUUIDsLen { test.MemberUUIDs = []string{} }
	if test.Field == EmptyReplaceValueForString                    { test.Field = "" }
	if test.Location == EmptyReplaceValueForString                 { test.Location = "" }
	if test.Floor == EmptyReplaceValueForString                    { test.Floor = "" }
	// logo empty case 테스트 필요
	if string(test.Logo) == EmptyReplaceValueForString             { test.Logo = []byte{} }
	if test.ClubUUID == EmptyReplaceValueForString                 { test.ClubUUID = "" }
	if test.XRequestID == EmptyReplaceValueForString               { test.XRequestID = "" }
	if test.SpanContextString == EmptyReplaceValueForString        { test.SpanContextString = "" }
}

func (test *CreateNewClubCase) OnExpectMethodsTo(mock *mock.Mock) {
	for method, returns := range test.ExpectedMethods {
		test.onMethod(mock, method, returns)
	}
}

func (test *CreateNewClubCase) onMethod(mock *mock.Mock, method Method, returns Returns) {
	switch method {
	case "CreateClub":
		const indexClubModel = 0
		const indexError = 1
		if _, ok := returns[indexClubModel].(*model.Club); ok && returns[indexError] == nil {
			modelToReturn := test.getClubModel()
			modelToReturn.Model = createGormModelOnCurrentTime()
			returns[indexClubModel] = modelToReturn
		}
		mock.On(string(method), test.getClubModel()).Return(returns...)

	case "CreateClubInform":
		const indexClubInformModel = 0
		const indexError = 1
		if _, ok := returns[indexClubInformModel].(*model.ClubInform); ok && returns[indexError] == nil {
			modelToReturn := test.getClubInformModel()
			modelToReturn.Model = createGormModelOnCurrentTime()
			returns[indexClubInformModel] = modelToReturn
		}
		mock.On(string(method), test.getClubInformModel()).Return(returns...)

	case "CreateClubMembers":
		const indexClubMemberModel = 0
		const indexError = 1
		for index := range test.MemberUUIDs {
			if _, ok := returns[indexClubMemberModel].(*model.ClubMember); ok && returns[indexError] == nil {
				modelToReturn := test.getClubMemberModelWithIndex(index)
				modelToReturn.Model = createGormModelOnCurrentTime()
				returns[indexClubMemberModel] = modelToReturn
			}
			mock.On(string(method), test.getClubMemberModelWithIndex(index)).Return(returns...)
		}

	case "GetClubWithClubUUID":
		mock.On(string(method), test.ClubUUID).Return(returns...)

	case "GetStudentInformsWithUUIDs":
		requestCtx := context.Background()
		requestCtx = metadata.Set(requestCtx, "X-Request-Id", test.XRequestID)
		requestCtx = metadata.Set(requestCtx, "Span-Context", "123412341234:123412341234:123412341234")
		mock.On(string(method), requestCtx, &authproto.GetStudentInformsWithUUIDsRequest{
			UUID:         test.UUID,
			StudentUUIDs: test.MemberUUIDs,
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

func (test *CreateNewClubCase) getClubModel() *model.Club {
	return &model.Club{
		UUID:       model.UUID(test.UUID),
		LeaderUUID: model.LeaderUUID(test.LeaderUUID),
	}
}

func (test *CreateNewClubCase) getClubInformModel() *model.ClubInform {
	return &model.ClubInform{
		ClubUUID:     model.ClubUUID(test.ClubUUID),
		Name:         model.Name(test.Name),
		Field:        model.Field(test.Field),
		Location:     model.Location(test.Location),
		Floor:        model.Floor(test.Floor),
	}
}

func (test *CreateNewClubCase) getClubMemberModelWithIndex(index int) *model.ClubMember {
	return &model.ClubMember{
		ClubUUID:    model.ClubUUID(test.ClubUUID),
		StudentUUID: model.StudentUUID(test.MemberUUIDs[index]),
		Club:        nil,
	}
}

func (test *CreateNewClubCase) SetRequestContextOf(req *clubproto.CreateNewClubRequest) {
	req.UUID = test.UUID
	req.Name = test.Name
	req.LeaderUUID = test.LeaderUUID
	req.MemberUUIDs = test.MemberUUIDs
	req.Floor = test.Floor
	req.LeaderUUID = test.LeaderUUID
	req.Floor = test.Floor
	req.Logo = test.Logo
}

func (test *CreateNewClubCase) GetMetadataContext() (ctx context.Context) {
	ctx = context.Background()
	ctx = metadata.Set(ctx, "X-Request-Id", test.XRequestID)
	ctx = metadata.Set(ctx, "Span-Context", test.SpanContextString)
	ctx = metadata.Set(ctx, "ClubUUID", test.ClubUUID)
	return
}
