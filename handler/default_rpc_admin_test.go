package handler

import (
	test "club/handler/for_test"
	"club/model"
	authproto "club/proto/golang/auth"
	clubproto "club/proto/golang/club"
	consulagent "club/tool/consul/agent"
	"club/tool/mysqlerr"
	code "club/utils/code/golang"
	"errors"
	mysqlcode "github.com/VividCortex/mysqlerr"
	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
	microerrors "github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
				"GetStudentInformsWithUUIDs": {&authproto.GetStudentInformsWithUUIDsResponse{
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
			ExpectedStatus:   http.StatusCreated,
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
		}, { // invalid request
			ExpectedMethods: map[test.Method]test.Returns{
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformsWithUUIDs": {&authproto.GetStudentInformsWithUUIDsResponse{
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
				"CreateClub":          {&model.Club{}, (validator.ValidationErrors)(nil)},
				"Rollback":            {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusProxyAuthRequired,
		}, { // invalid request (floor -> in 1~5)
			Floor: "100",
			ExpectedMethods: map[test.Method]test.Returns{
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformsWithUUIDs": {&authproto.GetStudentInformsWithUUIDsResponse{
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
		}, { // invalid request
			ExpectedMethods: map[test.Method]test.Returns{
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformsWithUUIDs": {&authproto.GetStudentInformsWithUUIDsResponse{
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
				"CreateClubMembers":   {&model.ClubMember{}, (validator.ValidationErrors)(nil)},
				"Rollback":            {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusProxyAuthRequired,
		}, { // invalid request (logo not exist)
			Logo: []byte(test.EmptyReplaceValueForString),
			ExpectedMethods: map[test.Method]test.Returns{
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformsWithUUIDs": {&authproto.GetStudentInformsWithUUIDsResponse{
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
				"Rollback":            {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusProxyAuthRequired,
		}, { // member uuid arr not include leader uuid
			LeaderUUID:     "student-111111111111",
			MemberUUIDs:    []string{"student-222222222222"},
			ExpectedStatus: http.StatusConflict,
			ExpectedCode:   code.MemberUUIDsNotIncludeLeaderUUID,
		}, { // member uuid arr include not exist uuid
			LeaderUUID:  "student-111111111111",
			MemberUUIDs: []string{"student-111111111111", "student-222222222222"},
			ExpectedMethods: map[test.Method]test.Returns{
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformsWithUUIDs": {&authproto.GetStudentInformsWithUUIDsResponse{
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
			ExpectedCode:   code.MemberUUIDsIncludeNoExistUUID,
		}, { // member uuid arr include not exist uuid
			LeaderUUID:  "student-111111111111",
			MemberUUIDs: []string{"student-111111111111", "student-222222222222"},
			ExpectedMethods: map[test.Method]test.Returns{
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformsWithUUIDs": {&authproto.GetStudentInformsWithUUIDsResponse{
					Status:  http.StatusConflict,
					Code:    code.AdminIDNoExist, // unexpected error code
					Message: "unexpected error code",
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
			ExpectedStatus: http.StatusInternalServerError,
		}, { // GetStudentInformsWithUUIDs return not 200 or 407
			LeaderUUID:  "student-111111111111",
			MemberUUIDs: []string{"student-111111111111"},
			ExpectedMethods: map[test.Method]test.Returns{
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformsWithUUIDs": {&authproto.GetStudentInformsWithUUIDsResponse{
					Status:  http.StatusNetworkAuthenticationRequired,
					Message: "network authentication required",
				}, nil},
			},
			ExpectedStatus: http.StatusNetworkAuthenticationRequired,
		}, { // GetStudentInformsWithUUIDs return time out error
			LeaderUUID:  "student-111111111111",
			MemberUUIDs: []string{"student-111111111111"},
			ExpectedMethods: map[test.Method]test.Returns{
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformsWithUUIDs": {&authproto.GetStudentInformsWithUUIDsResponse{}, &microerrors.Error{
					Code:   http.StatusRequestTimeout,
					Detail: "request time out",
				}},
			},
			ExpectedStatus: http.StatusRequestTimeout,
		}, { // GetStudentInformsWithUUIDs return unexpected error 1
			LeaderUUID:  "student-111111111111",
			MemberUUIDs: []string{"student-111111111111"},
			ExpectedMethods: map[test.Method]test.Returns{
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformsWithUUIDs": {&authproto.GetStudentInformsWithUUIDsResponse{}, &microerrors.Error{
					Code:   http.StatusNetworkAuthenticationRequired,
					Detail: "what is this error?",
				}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // GetStudentInformsWithUUIDs return unexpected error 2
			LeaderUUID:  "student-111111111111",
			MemberUUIDs: []string{"student-111111111111"},
			ExpectedMethods: map[test.Method]test.Returns{
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformsWithUUIDs": {&authproto.GetStudentInformsWithUUIDsResponse{}, errors.New("unexpected error")},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // GetNextServiceNode return any error
			LeaderUUID:  "student-111111111111",
			MemberUUIDs: []string{"student-111111111111"},
			ExpectedMethods: map[test.Method]test.Returns{
				"GetNextServiceNode": {&registry.Node{}, errors.New("I don't know what error is")},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // GetNextServiceNode return ErrAvailableNodeNotFound
			LeaderUUID:  "student-111111111111",
			MemberUUIDs: []string{"student-111111111111"},
			ExpectedMethods: map[test.Method]test.Returns{
				"GetNextServiceNode": {&registry.Node{}, consulagent.ErrAvailableNodeNotFound},
			},
			ExpectedStatus: http.StatusServiceUnavailable,
		}, { // GetClubWithClubUUID unexpected error
			LeaderUUID:  "student-111111111111",
			MemberUUIDs: []string{"student-111111111111"},
			ExpectedMethods: map[test.Method]test.Returns{
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformsWithUUIDs": {&authproto.GetStudentInformsWithUUIDsResponse{
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
				"GetClubWithClubUUID": {&model.Club{}, errors.New("unexpected error")},
				"Rollback":            {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // leader uuid duplicate error
			LeaderUUID:  "student-111111111111",
			MemberUUIDs: []string{"student-111111111111"},
			ExpectedMethods: map[test.Method]test.Returns{
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformsWithUUIDs": {&authproto.GetStudentInformsWithUUIDsResponse{
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
				"CreateClub":          {&model.Club{}, mysqlerr.DuplicateEntry(model.ClubInstance.LeaderUUID.KeyName(), "student-111111111111")},
				"Rollback":            {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusConflict,
			ExpectedCode:   code.ClubLeaderAlreadyExist,
		}, { // CreateClub return unexpected duplicate error
			LeaderUUID:  "student-111111111111",
			MemberUUIDs: []string{"student-111111111111"},
			ExpectedMethods: map[test.Method]test.Returns{
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformsWithUUIDs": {&authproto.GetStudentInformsWithUUIDsResponse{
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
				"CreateClub":          {&model.Club{}, mysqlerr.DuplicateEntry(model.ClubInstance.UUID.KeyName(), "club-111111111111")},
				"Rollback":            {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		},  { // CreateClub return unexpected type of error
			LeaderUUID:  "student-111111111111",
			MemberUUIDs: []string{"student-111111111111"},
			ExpectedMethods: map[test.Method]test.Returns{
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformsWithUUIDs": {&authproto.GetStudentInformsWithUUIDsResponse{
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
				"CreateClub":          {&model.Club{}, errors.New("unexpected type of error")},
				"Rollback":            {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateClub return invalid message in duplicate error
			ExpectedMethods: map[test.Method]test.Returns{
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformsWithUUIDs": {&authproto.GetStudentInformsWithUUIDsResponse{
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
				"CreateClub":          {&model.Club{}, &mysql.MySQLError{Number: mysqlcode.ER_DUP_ENTRY, Message: "Invalid Message"}},
				"Rollback":            {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateClub return unexpected error code
			ExpectedMethods: map[test.Method]test.Returns{
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformsWithUUIDs": {&authproto.GetStudentInformsWithUUIDsResponse{
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
				"CreateClub":          {&model.Club{}, &mysql.MySQLError{Number: mysqlcode.ER_BAD_NULL_ERROR, Message: "Unexpected Err Code"}},
				"Rollback":            {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // Club Name Duplicate error
			Name: "DMS",
			ExpectedMethods: map[test.Method]test.Returns{
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformsWithUUIDs": {&authproto.GetStudentInformsWithUUIDsResponse{
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
				"GetStudentInformsWithUUIDs": {&authproto.GetStudentInformsWithUUIDsResponse{
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
				"CreateClubInform":    {&model.ClubInform{}, mysqlerr.DuplicateEntry(model.ClubInformInstance.Location.KeyName(), "2-2반 교실")},
				"Rollback":            {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusConflict,
			ExpectedCode:   code.ClubLocationDuplicate,
		}, { // CreateClubInform return unexpected duplicate error
			Location: "2-2반 교실",
			ExpectedMethods: map[test.Method]test.Returns{
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformsWithUUIDs": {&authproto.GetStudentInformsWithUUIDsResponse{
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
				"CreateClubInform":    {&model.ClubInform{}, mysqlerr.DuplicateEntry(model.ClubInformInstance.Floor.KeyName(), "3")},
				"Rollback":            {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateClubInform return unexpected type of error
			Location: "2-2반 교실",
			ExpectedMethods: map[test.Method]test.Returns{
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformsWithUUIDs": {&authproto.GetStudentInformsWithUUIDsResponse{
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
				"CreateClubInform":    {&model.ClubInform{}, errors.New("unexpected type of error")},
				"Rollback":            {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateClubInform returns invalid message in duplicate error
			ExpectedMethods: map[test.Method]test.Returns{
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformsWithUUIDs": {&authproto.GetStudentInformsWithUUIDsResponse{
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
				"CreateClubInform":    {&model.ClubInform{}, &mysql.MySQLError{Number: mysqlcode.ER_DUP_ENTRY, Message: "Invalid Message"}},
				"Rollback":            {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateClubInform returns unexpected mysql error
			ExpectedMethods: map[test.Method]test.Returns{
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformsWithUUIDs": {&authproto.GetStudentInformsWithUUIDsResponse{
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
				"CreateClubInform":    {&model.ClubInform{}, &mysql.MySQLError{Number: mysqlcode.ER_BAD_NULL_ERROR, Message: "Unexpected error code"}},
				"Rollback":            {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // Club Member duplicate error
			LeaderUUID:  "student-111111111111",
			MemberUUIDs: []string{"student-111111111111", "student-111111111111"},
			ExpectedMethods: map[test.Method]test.Returns{
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformsWithUUIDs": {&authproto.GetStudentInformsWithUUIDsResponse{
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
		}, { // CreateClubMembers returns unexpected duplicate error
			ExpectedMethods: map[test.Method]test.Returns{
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformsWithUUIDs": {&authproto.GetStudentInformsWithUUIDsResponse{
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
				"CreateClubMembers":   {[]*model.ClubMember{}, mysqlerr.DuplicateEntry(model.ClubMemberInstance.ClubUUID.KeyName(), "club-111111111111")},
				"Rollback":            {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateClubMembers returns unexpected type of error
			ExpectedMethods: map[test.Method]test.Returns{
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformsWithUUIDs": {&authproto.GetStudentInformsWithUUIDsResponse{
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
				"CreateClubMembers":   {[]*model.ClubMember{}, errors.New("unexpected error type")},
				"Rollback":            {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateClubMembers returns invalid message in duplicate error
			ExpectedMethods: map[test.Method]test.Returns{
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformsWithUUIDs": {&authproto.GetStudentInformsWithUUIDsResponse{
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
				"CreateClubMembers":   {[]*model.ClubMember{}, &mysql.MySQLError{Number: mysqlcode.ER_DUP_ENTRY, Message: "Invalid Message"}},
				"Rollback":            {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateClubMembers returns unexpected mysql error
			ExpectedMethods: map[test.Method]test.Returns{
				"GetNextServiceNode": {&registry.Node{
					Id:      "DMS.SMS.v1.service.auth-6b37b034-5f0b-4c9f-a03a-decbcb3799ef",
					Address: "127.0.0.1:10101",
				}, nil},
				"GetStudentInformsWithUUIDs": {&authproto.GetStudentInformsWithUUIDsResponse{
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
				"CreateClubMembers":   {[]*model.ClubMember{}, &mysql.MySQLError{Number: mysqlcode.ER_BAD_NULL_ERROR, Message: "Unexpected Error Code"}},
				"Rollback":            {&gorm.DB{}},
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

		req := new(clubproto.CreateNewClubRequest)
		testCase.SetRequestContextOf(req)
		ctx := testCase.GetMetadataContext()

		resp := new(clubproto.CreateNewClubResponse)
		_ = handler.CreateNewClub(ctx, req, resp)

		testCase.Logo = nil
		assert.Equalf(t, int(testCase.ExpectedStatus), int(resp.Status), "status assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Equalf(t, testCase.ExpectedCode, resp.Code, "code assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Regexpf(t, testCase.ExpectedClubUUID, resp.ClubUUID, "club uuid assertion error (test case: %v, message: %s)", testCase, resp.Message)

		newMock.AssertExpectations(t)
	}
}
