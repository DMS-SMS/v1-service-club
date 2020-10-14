package test

import (
	"club/model"
	authproto "club/proto/golang/auth"
	topic "club/utils/topic/golang"
	"github.com/stretchr/testify/mock"
	"log"
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
func (test *AddClubMemberCase) OnExpectMethodsTo(mock *mock.Mock) {
	for method, returns := range test.ExpectedMethods {
		test.onMethod(mock, method, returns)
	}
}

func (test *AddClubMemberCase) onMethod(mock *mock.Mock, method Method, returns Returns) {
	switch method {
	case "GetClubWithClubUUID":
		mock.On(string(method), test.UUID).Return(returns...)
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

