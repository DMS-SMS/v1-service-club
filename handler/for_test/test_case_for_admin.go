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
