package handler

import (
	test "club/handler/for_test"
	"club/model"
	authproto "club/proto/golang/auth"
	consulagent "club/tool/consul/agent"
	"club/tool/mysqlerr"
	code "club/utils/code/golang"
	"errors"
	"github.com/go-playground/validator/v10"
	microerrors "github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/registry"
	"gorm.io/gorm"
	"net/http"
	"testing"
)

func Test_default_CreateNewClub(t *testing.T) {
	const clubUUIDRegexString = "^club-\\d{12}"

	tests := []test.CreateNewClubCase{
		{ // success case
			LeaderUUID:  "student-111111111111",
			MemberUUIDs: []string{"student-111111111111"},
			ExpectedMethods: map[test.Method]test.Returns{
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformsWithUUID": {&authproto.GetStudentInformsWithUUIDsResponse{
					Status:  http.StatusOK,
					Message: "success!",
					StudentInforms: []*authproto.StudentInform{{
						StudentUUID:   "student-111111111111",
						Grade:         2,
						Group:         2,
						StudentNumber: 7,
						Name:          "박진홍",
						PhoneNumber:   "01088378347",
						ImageURI:      "profiles/student-111111111111",
					}},
				}, nil},
				"BeginTx":             {},
				"GetClubWithClubUUID": {&model.Club{}, gorm.ErrRecordNotFound},
				"CreateClub":          {&model.Club{}, nil},
				"CreateClubInform":    {&model.ClubInform{}, nil},
				"CreateClubMembers":   {[]*model.ClubMember{}, nil},
				"Commit":              {&gorm.DB{}},
			},
			ExpectedStatus:   http.StatusOK,
			ExpectedClubUUID: clubUUIDRegexString,
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
		}, { // not admin uuid
			UUID:            "parent-111111111111",
			ExpectedMethods: map[test.Method]test.Returns{},
			ExpectedStatus:  http.StatusForbidden,
		},  { // invalid request (floor -> in 1~5)
			Floor: "100",
			ExpectedMethods: map[test.Method]test.Returns{
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformsWithUUID": {&authproto.GetStudentInformsWithUUIDsResponse{
					Status:  http.StatusOK,
					Message: "success!",
					StudentInforms: []*authproto.StudentInform{{
						StudentUUID:   "student-111111111111",
						Grade:         2,
						Group:         2,
						StudentNumber: 7,
						Name:          "박진홍",
						PhoneNumber:   "01088378347",
						ImageURI:      "profiles/student-111111111111",
					}},
				}, nil},
				"BeginTx":             {},
				"GetClubWithClubUUID": {&model.Club{}, gorm.ErrRecordNotFound},
				"CreateClub":          {&model.Club{}, nil},
				"CreateClubInform":    {&model.ClubInform{}, (validator.ValidationErrors)(nil)},
				"Rollback":            {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusProxyAuthRequired,
		}, { // member uuid arr not include leader uuid
			LeaderUUID:     "student-111111111111",
			MemberUUIDs:    []string{"student-222222222222"},
			ExpectedStatus: http.StatusOK,
			ExpectedCode:   code.MemberUUIDsNotIncludeLeaderUUID,
		}, { // member uuid arr include not exist uuid
			LeaderUUID:  "student-111111111111",
			MemberUUIDs: []string{"student-111111111111", "student-222222222222"},
			ExpectedMethods: map[test.Method]test.Returns{
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformsWithUUID": {&authproto.GetStudentInformsWithUUIDsResponse{
					Status:  http.StatusConflict,
					Code:    code.StudentUUIDsContainNoExistUUID,
					Message: "student uuid array contain no exist uuid",
					StudentInforms: []*authproto.StudentInform{{
						StudentUUID:   "student-111111111111",
						Grade:         2,
						Group:         2,
						StudentNumber: 7,
						Name:          "박진홍",
						PhoneNumber:   "01088378347",
						ImageURI:      "profiles/student-111111111111",
					}, {}},
				}, nil},
			},
			ExpectedStatus: http.StatusConflict,
			ExpectedCode:   code.MemeberUUIDsIncludeNoExistUUID,
		}, { // GetStudentInformsWithUUID return not 200 or 407
			LeaderUUID:  "student-111111111111",
			MemberUUIDs: []string{"student-111111111111"},
			ExpectedMethods: map[test.Method]test.Returns{
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformsWithUUID": {&authproto.GetStudentInformsWithUUIDsResponse{
					Status:  http.StatusInternalServerError,
					Message: "internal server error",
				}, nil},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // GetStudentInformsWithUUID return any error
			LeaderUUID:  "student-111111111111",
			MemberUUIDs: []string{"student-111111111111"},
			ExpectedMethods: map[test.Method]test.Returns{
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformsWithUUID": {&authproto.GetStudentInformsWithUUIDsResponse{}, errors.New("I don't know what error is")},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // GetStudentInformsWithUUID return micro error
			LeaderUUID:  "student-111111111111",
			MemberUUIDs: []string{"student-111111111111"},
			ExpectedMethods: map[test.Method]test.Returns{
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformsWithUUID": {&authproto.GetStudentInformsWithUUIDsResponse{}, microerrors.Error{
					Code: http.StatusNetworkAuthenticationRequired,
				}},
			},
			ExpectedStatus: http.StatusNetworkAuthenticationRequired,
		},{ // GetNextServiceNode return any error
			LeaderUUID:  "student-111111111111",
			MemberUUIDs: []string{"student-111111111111"},
			ExpectedMethods: map[test.Method]test.Returns{
				"GetNextServiceNode": {nil, errors.New("I don't know what error is")},
			},
			ExpectedStatus: http.StatusInternalServerError,
		},  { // GetNextServiceNode return ErrAvailableNodeNotFound
			LeaderUUID:  "student-111111111111",
			MemberUUIDs: []string{"student-111111111111"},
			ExpectedMethods: map[test.Method]test.Returns{
				"GetNextServiceNode": {nil, consulagent.ErrAvailableNodeNotFound},
			},
			ExpectedStatus: http.StatusServiceUnavailable,
		}, { // leader uuid duplicate error
			LeaderUUID:  "student-111111111111",
			MemberUUIDs: []string{"student-111111111111"},
			ExpectedMethods: map[test.Method]test.Returns{
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformsWithUUID": {&authproto.GetStudentInformsWithUUIDsResponse{
					Status:  http.StatusOK,
					Message: "success!",
					StudentInforms: []*authproto.StudentInform{{
						StudentUUID:   "student-111111111111",
						Grade:         2,
						Group:         2,
						StudentNumber: 7,
						Name:          "박진홍",
						PhoneNumber:   "01088378347",
						ImageURI:      "profiles/student-111111111111",
					}},
				}, nil},
				"BeginTx":             {},
				"GetClubWithClubUUID": {&model.Club{}, mysqlerr.DuplicateEntry(model.ClubInstance.LeaderUUID.KeyName(), "student-111111111111")},
				"Rollback":            {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusConflict,
			ExpectedCode:   code.LeaderUUIDDuplicate,
		}, { // Club Name Duplicate error
			Name: "DMS",
			ExpectedMethods: map[test.Method]test.Returns{
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformsWithUUID": {&authproto.GetStudentInformsWithUUIDsResponse{
					Status:  http.StatusOK,
					Message: "success!",
					StudentInforms: []*authproto.StudentInform{{
						StudentUUID:   "student-111111111111",
						Grade:         2,
						Group:         2,
						StudentNumber: 7,
						Name:          "박진홍",
						PhoneNumber:   "01088378347",
						ImageURI:      "profiles/student-111111111111",
					}},
				}, nil},
				"BeginTx":             {},
				"GetClubWithClubUUID": {&model.Club{}, gorm.ErrRecordNotFound},
				"CreateClub":          {&model.Club{}, nil},
				"CreateClubInform":    {&model.ClubInform{}, mysqlerr.DuplicateEntry(model.ClubInformInstance.Name.KeyName(), "DMS")},
				"Rollback":            {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusConflict,
			ExpectedCode:   code.ClubNameDuplicate,
		}, { // Club Location Duplicate error
			Location: "2-2반 교실",
			ExpectedMethods: map[test.Method]test.Returns{
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformsWithUUID": {&authproto.GetStudentInformsWithUUIDsResponse{
					Status:  http.StatusOK,
					Message: "success!",
					StudentInforms: []*authproto.StudentInform{{
						StudentUUID:   "student-111111111111",
						Grade:         2,
						Group:         2,
						StudentNumber: 7,
						Name:          "박진홍",
						PhoneNumber:   "01088378347",
						ImageURI:      "profiles/student-111111111111",
					}},
				}, nil},
				"BeginTx":             {},
				"GetClubWithClubUUID": {&model.Club{}, gorm.ErrRecordNotFound},
				"CreateClub":          {&model.Club{}, nil},
				"CreateClubInform":    {&model.ClubInform{}, mysqlerr.DuplicateEntry(model.ClubInformInstance.Name.KeyName(), "DMS")},
				"Rollback":            {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusConflict,
			ExpectedCode:   code.ClubLocationDuplicate,
		}, { // Club Location Duplicate error
			LeaderUUID:  "student-111111111111",
			MemberUUIDs: []string{"student-111111111111", "student-111111111111"},
			ExpectedMethods: map[test.Method]test.Returns{
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformsWithUUID": {&authproto.GetStudentInformsWithUUIDsResponse{
					Status:  http.StatusOK,
					Message: "success!",
					StudentInforms: []*authproto.StudentInform{{
						StudentUUID:   "student-111111111111",
						Grade:         2,
						Group:         2,
						StudentNumber: 7,
						Name:          "박진홍",
						PhoneNumber:   "01088378347",
						ImageURI:      "profiles/student-111111111111",
					}, {
						StudentUUID:   "student-111111111111",
						Grade:         2,
						Group:         2,
						StudentNumber: 7,
						Name:          "박진홍",
						PhoneNumber:   "01088378347",
						ImageURI:      "profiles/student-111111111111",
					}},
				}, nil},
				"BeginTx":             {},
				"GetClubWithClubUUID": {&model.Club{}, gorm.ErrRecordNotFound},
				"CreateClub":          {&model.Club{}, nil},
				"CreateClubInform":    {&model.ClubInform{}, nil},
				"CreateClubMembers":   {[]*model.ClubMember{}, mysqlerr.DuplicateEntry(model.ClubMemberInstance.StudentUUID.KeyName(), "student-111111111111")},
				"Rollback":            {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusConflict,
			ExpectedCode:   code.ClubMemberDuplicate,
		},
	}
}
