package handler

import (
	test "club/handler/for_test"
	"club/model"
	authproto "club/proto/golang/auth"
	clubproto "club/proto/golang/club"
	consulagent "club/tool/consul/agent"
	"club/tool/mysqlerr"
	code "club/utils/code/golang"
	mysqlcode "github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
	microerrors "github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/registry"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"net/http"
	"testing"
)

func Test_Default_AddClubMember(t *testing.T) {
	tests := []test.AddClubMemberCase{
		{ // success case (student uuid)
			UUID:        "student-111111111111",
			ClubUUID:    "club-111111111111",
			StudentUUID: "student-222222222222",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetClubWithClubUUID": {&model.Club{
					UUID:       "club-111111111111",
					LeaderUUID: "student-111111111111",
				}, nil},
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformWithUUID": {&authproto.GetStudentInformWithUUIDResponse{
					Status:        http.StatusOK,
					Message:       "get student inform success",
					Grade:         2,
					Group:         2,
					StudentNumber: 7,
					Name:          "박진홍",
					PhoneNumber:   "01088378347",
					ImageURI:      "profiles/student-111111111111",
				}, nil},
				"CreateClubMember": {&model.ClubMember{
					ClubUUID:    "club-111111111111",
					StudentUUID: "student-222222222222",
				}, nil},
				"Commit": {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusCreated,
		}, { // success case (admin uuid)
			UUID:        "admin-111111111111",
			ClubUUID:    "club-111111111111",
			StudentUUID: "student-222222222222",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetClubWithClubUUID": {&model.Club{
					UUID:       "club-111111111111",
					LeaderUUID: "student-111111111111",
				}, nil},
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformWithUUID": {&authproto.GetStudentInformWithUUIDResponse{
					Status:        http.StatusOK,
					Message:       "get student inform success",
					Grade:         2,
					Group:         2,
					StudentNumber: 7,
					Name:          "박진홍",
					PhoneNumber:   "01088378347",
					ImageURI:      "profiles/student-111111111111",
				}, nil},
				"CreateClubMember": {&model.ClubMember{
					ClubUUID:    "club-111111111111",
					StudentUUID: "student-222222222222",
				}, nil},
				"Commit": {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusCreated,
		}, { // no exist X-Request-ID -> Proxy Authorization Required
			XRequestID:      test.EmptyReplaceValueForString,
			ExpectedMethods: map[test.Method]test.Returns{},
			ExpectedStatus:  http.StatusProxyAuthRequired,
		}, { // invalid X-Request-ID -> Proxy Authorization Required
			XRequestID:      "InvalidXRequestID",
			ExpectedMethods: map[test.Method]test.Returns{},
			ExpectedStatus:  http.StatusProxyAuthRequired,
		}, { // no exist Span-Context -> Proxy Authorization Required
			SpanContextString: test.EmptyReplaceValueForString,
			ExpectedMethods:   map[test.Method]test.Returns{},
			ExpectedStatus:    http.StatusProxyAuthRequired,
		}, { // invalid Span-Context -> Proxy Authorization Required
			SpanContextString: "InvalidSpanContext",
			ExpectedMethods:   map[test.Method]test.Returns{},
			ExpectedStatus:    http.StatusProxyAuthRequired,
		}, { // not student or admin uuid
			UUID:            "parent-111111111111",
			ExpectedMethods: map[test.Method]test.Returns{},
			ExpectedStatus:  http.StatusForbidden,
			ExpectedCode:    code.ForbiddenNotStudentOrAdminUUID,
		}, { // not club leader
			UUID:        "student-222222222222",
			ClubUUID:    "club-111111111111",
			StudentUUID: "student-333333333333",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetClubWithClubUUID": {&model.Club{
					UUID:       "club-111111111111",
					LeaderUUID: "student-111111111111",
				}, nil},
				"Rollback": {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusForbidden,
			ExpectedCode:   code.ForbiddenNotClubLeader,
		}, { // GetClubWithClubUUID returns not exists
			UUID:        "student-111111111111",
			ClubUUID:    "club-111111111111",
			StudentUUID: "student-222222222222",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":             {},
				"GetClubWithClubUUID": {&model.Club{}, gorm.ErrRecordNotFound},
				"Rollback":            {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusNotFound,
			ExpectedCode:   code.NotFoundClubNoExist,
		}, { // GetClubWithClubUUID returns unexpected error
			UUID:        "student-111111111111",
			ClubUUID:    "club-111111111111",
			StudentUUID: "student-222222222222",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":             {},
				"GetClubWithClubUUID": {&model.Club{}, errors.New("unexpected error")},
				"Rollback":            {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // GetNextServiceNode return any error
			UUID:        "student-111111111111",
			ClubUUID:    "club-111111111111",
			StudentUUID: "student-222222222222",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetClubWithClubUUID": {&model.Club{
					UUID:       "club-111111111111",
					LeaderUUID: "student-111111111111",
				}, nil},
				"GetNextServiceNode": {&registry.Node{}, errors.New("I don't know what error is")},
				"Rollback":           {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // GetNextServiceNode return ErrAvailableNodeNotFound
			UUID:        "student-111111111111",
			ClubUUID:    "club-111111111111",
			StudentUUID: "student-222222222222",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetClubWithClubUUID": {&model.Club{
					UUID:       "club-111111111111",
					LeaderUUID: "student-111111111111",
				}, nil},
				"GetNextServiceNode": {&registry.Node{}, consulagent.ErrAvailableNodeNotFound},
				"Rollback":           {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusServiceUnavailable,
		}, { // GetStudentInformWithUUID response 404
			UUID:        "student-111111111111",
			ClubUUID:    "club-111111111111",
			StudentUUID: "student-222222222222",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetClubWithClubUUID": {&model.Club{
					UUID:       "club-111111111111",
					LeaderUUID: "student-111111111111",
				}, nil},
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformWithUUID": {&authproto.GetStudentInformWithUUIDResponse{
					Status:  http.StatusNotFound,
					Message: "student uuid not exist",
				}, nil},
				"Rollback": {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusNotFound,
			ExpectedCode:   code.NotFoundStudentNoExist,
		}, { // GetStudentInformWithUUID response not 200 or 404
			UUID:        "student-111111111111",
			ClubUUID:    "club-111111111111",
			StudentUUID: "student-222222222222",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetClubWithClubUUID": {&model.Club{
					UUID:       "club-111111111111",
					LeaderUUID: "student-111111111111",
				}, nil},
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformWithUUID": {&authproto.GetStudentInformWithUUIDResponse{
					Status:  http.StatusNetworkAuthenticationRequired,
					Message: "I don't know about this error",
				}, nil},
				"Rollback": {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusNetworkAuthenticationRequired,
		}, { // GetStudentInformWithUUID response not 200 or 404
			UUID:        "student-111111111111",
			ClubUUID:    "club-111111111111",
			StudentUUID: "student-222222222222",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetClubWithClubUUID": {&model.Club{
					UUID:       "club-111111111111",
					LeaderUUID: "student-111111111111",
				}, nil},
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformWithUUID": {&authproto.GetStudentInformWithUUIDResponse{
					Status:  http.StatusNetworkAuthenticationRequired,
					Message: "I don't know about this error",
				}, nil},
				"Rollback": {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusNetworkAuthenticationRequired,
		}, { // GetStudentInformWithUUID response timeout error
			UUID:        "student-111111111111",
			ClubUUID:    "club-111111111111",
			StudentUUID: "student-222222222222",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetClubWithClubUUID": {&model.Club{
					UUID:       "club-111111111111",
					LeaderUUID: "student-111111111111",
				}, nil},
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformWithUUID": {&authproto.GetStudentInformWithUUIDResponse{}, &microerrors.Error{
					Code:   http.StatusRequestTimeout,
					Detail: "request time out",
				}},
				"Rollback": {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusRequestTimeout,
		}, { // GetStudentInformWithUUID response unexpected error code
			UUID:        "student-111111111111",
			ClubUUID:    "club-111111111111",
			StudentUUID: "student-222222222222",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetClubWithClubUUID": {&model.Club{
					UUID:       "club-111111111111",
					LeaderUUID: "student-111111111111",
				}, nil},
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformWithUUID": {&authproto.GetStudentInformWithUUIDResponse{}, &microerrors.Error{
					Code:   http.StatusNetworkAuthenticationRequired,
					Detail: "I don't know about this error",
				}},
				"Rollback": {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // GetStudentInformWithUUID response unexpected type of error
			UUID:        "student-111111111111",
			ClubUUID:    "club-111111111111",
			StudentUUID: "student-222222222222",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetClubWithClubUUID": {&model.Club{
					UUID:       "club-111111111111",
					LeaderUUID: "student-111111111111",
				}, nil},
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformWithUUID": {&authproto.GetStudentInformWithUUIDResponse{}, errors.New("unexpected error")},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateClubMember returns duplicate error
			UUID:        "student-111111111111",
			ClubUUID:    "club-111111111111",
			StudentUUID: "student-222222222222",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetClubWithClubUUID": {&model.Club{
					UUID:       "club-111111111111",
					LeaderUUID: "student-111111111111",
				}, nil},
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformWithUUID": {&authproto.GetStudentInformWithUUIDResponse{
					Status:        http.StatusOK,
					Message:       "get student inform success",
					Grade:         2,
					Group:         2,
					StudentNumber: 7,
					Name:          "박진홍",
					PhoneNumber:   "01088378347",
					ImageURI:      "profiles/student-111111111111",
				}, nil},
				"CreateClubMember": {&model.ClubMember{}, mysqlerr.DuplicateEntry(model.ClubMemberInstance.StudentUUID.KeyName(), "student-222222222222")},
				"Rollback":         {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusConflict,
			ExpectedCode:   code.ClubMemberAlreadyExist,
		}, { // CreateClubMember returns unexpected duplicate entry error
			UUID:        "student-111111111111",
			ClubUUID:    "club-111111111111",
			StudentUUID: "student-222222222222",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetClubWithClubUUID": {&model.Club{
					UUID:       "club-111111111111",
					LeaderUUID: "student-111111111111",
				}, nil},
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformWithUUID": {&authproto.GetStudentInformWithUUIDResponse{
					Status:        http.StatusOK,
					Message:       "get student inform success",
					Grade:         2,
					Group:         2,
					StudentNumber: 7,
					Name:          "박진홍",
					PhoneNumber:   "01088378347",
					ImageURI:      "profiles/student-111111111111",
				}, nil},
				"CreateClubMember": {&model.ClubMember{}, mysqlerr.DuplicateEntry(model.ClubMemberInstance.ClubUUID.KeyName(), "club-111111111111")},
				"Rollback":         {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateClubMember returns invalid message in duplicate entry error
			UUID:        "student-111111111111",
			ClubUUID:    "club-111111111111",
			StudentUUID: "student-222222222222",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetClubWithClubUUID": {&model.Club{
					UUID:       "club-111111111111",
					LeaderUUID: "student-111111111111",
				}, nil},
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformWithUUID": {&authproto.GetStudentInformWithUUIDResponse{
					Status:        http.StatusOK,
					Message:       "get student inform success",
					Grade:         2,
					Group:         2,
					StudentNumber: 7,
					Name:          "박진홍",
					PhoneNumber:   "01088378347",
					ImageURI:      "profiles/student-111111111111",
				}, nil},
				"CreateClubMember": {&model.ClubMember{}, &mysql.MySQLError{
					Number:  mysqlcode.ER_DUP_ENTRY,
					Message: "invalid message",
				}},
				"Rollback": {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateClubMember returns unexpected MySQL error number
			UUID:        "student-111111111111",
			ClubUUID:    "club-111111111111",
			StudentUUID: "student-222222222222",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetClubWithClubUUID": {&model.Club{
					UUID:       "club-111111111111",
					LeaderUUID: "student-111111111111",
				}, nil},
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformWithUUID": {&authproto.GetStudentInformWithUUIDResponse{
					Status:        http.StatusOK,
					Message:       "get student inform success",
					Grade:         2,
					Group:         2,
					StudentNumber: 7,
					Name:          "박진홍",
					PhoneNumber:   "01088378347",
					ImageURI:      "profiles/student-111111111111",
				}, nil},
				"CreateClubMember": {&model.ClubMember{}, &mysql.MySQLError{
					Number:  mysqlcode.ER_BAD_NULL_ERROR,
					Message: "unexpected number",
				}},
				"Rollback": {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateClubMember returns unexpected type of error
			UUID:        "student-111111111111",
			ClubUUID:    "club-111111111111",
			StudentUUID: "student-222222222222",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetClubWithClubUUID": {&model.Club{
					UUID:       "club-111111111111",
					LeaderUUID: "student-111111111111",
				}, nil},
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformWithUUID": {&authproto.GetStudentInformWithUUIDResponse{
					Status:        http.StatusOK,
					Message:       "get student inform success",
					Grade:         2,
					Group:         2,
					StudentNumber: 7,
					Name:          "박진홍",
					PhoneNumber:   "01088378347",
					ImageURI:      "profiles/student-111111111111",
				}, nil},
				"CreateClubMember": {&model.ClubMember{}, errors.New("unexpected error")},
				"Rollback":         {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		},
	}

	for _, testCase := range tests {
		newMock := &mock.Mock{}
		handler := newDefaultMockHandler(newMock)

		testCase.ChangeEmptyValueToValidValue()
		testCase.ChangeEmptyReplaceValueToEmptyValue()
		testCase.OnExpectMethodsTo(newMock)

		req := new(clubproto.AddClubMemberRequest)
		testCase.SetRequestContextOf(req)
		ctx := testCase.GetMetadataContext()

		resp := new(clubproto.AddClubMemberResponse)
		_ = handler.AddClubMember(ctx, req, resp)

		assert.Equalf(t, int(testCase.ExpectedStatus), int(resp.Status), "status assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Equalf(t, testCase.ExpectedCode, resp.Code, "code assertion error (test case: %v, message: %s)", testCase, resp.Message)

		newMock.AssertExpectations(t)
	}
}
