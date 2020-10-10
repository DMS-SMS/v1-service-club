package handler

import (
	test "club/handler/for_test"
	"club/model"
	clubproto "club/proto/golang/club"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"net/http"
	"testing"
)

func Test_Default_GetClubsSortByUpdateTime(t *testing.T) {
	tests := []test.GetClubsSortByUpdateTimeCase{
		{ // success case
			Start: 0,
			Count: 10,
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetClubInformsSortByUpdateTime": {[]*model.ClubInform{{
					ClubUUID:     "club-222222222222",
					Name:         "SMS",
					ClubConcept:  "DMS의 소속부서 SMS 입니다!",
					Introduction: "School Management System 서비스를 개발 및 운영합니다",
					Link:         "facebook.com/DMS-SMS",
					Field:        "SW 개발",
					Location:     "2-2반 교실",
					Floor:        "3",
					LogoURI:      "logo.com/club-222222222222",
				}, {
					ClubUUID: "club-333333333333",
					Name:     "PMS",
					Field:    "SW 개발",
					Location: "2-3반 교실",
					Floor:    "3",
					LogoURI:  "logo.com/club-333333333333",
				}, {
					ClubUUID: "club-111111111111",
					Name:     "DMS",
					Field:    "SW 개발",
					Location: "2-1반 교실",
					Floor:    "3",
					LogoURI:  "logo.com/club-111111111111",
				}}, nil},
				"GetClubsWithClubUUIDs": {[]*model.Club{{
					UUID:       "club-222222222222",
					LeaderUUID: "student-222222222222",
				}, {
					UUID:       "club-333333333333",
					LeaderUUID: "student-333333333333",
				}, {
					UUID:       "club-111111111111",
					LeaderUUID: "student-111111111111",
				}}, nil},
				"Commit": {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusOK,
			ExpectClubInforms: []*clubproto.ClubInform{{
				ClubUUID:     "club-222222222222",
				LeaderUUID:   "student-222222222222",
				Name:         "SMS",
				ClubConcept:  "DMS의 소속부서 SMS 입니다!",
				Introduction: "School Management System 서비스를 개발 및 운영합니다",
				Link:         "facebook.com/DMS-SMS",
				Field:        "SW 개발",
				Location:     "2-2반 교실",
				Floor:        "3",
				LogoURI:      "logo.com/club-222222222222",
			}, {
				ClubUUID:   "club-333333333333",
				LeaderUUID: "student-333333333333",
				Name:       "PMS",
				Field:      "SW 개발",
				Location:   "2-3반 교실",
				Floor:      "3",
				LogoURI:    "logo.com/club-333333333333",
			}, {
				ClubUUID:   "club-111111111111",
				LeaderUUID: "student-111111111111",
				Name:       "DMS",
				Field:      "SW 개발",
				Location:   "2-1반 교실",
				Floor:      "3",
				LogoURI:    "logo.com/club-111111111111",
			}},
		}, { // Start, Count Set X -> default(0, 10)
			UUID: "admin-111111111111",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetClubInformsSortByUpdateTime": {[]*model.ClubInform{{
					ClubUUID:     "club-222222222222",
					Name:         "SMS",
					ClubConcept:  "DMS의 소속부서 SMS 입니다!",
					Introduction: "School Management System 서비스를 개발 및 운영합니다",
					Link:         "facebook.com/DMS-SMS",
					Field:        "SW 개발",
					Location:     "2-2반 교실",
					Floor:        "3",
					LogoURI:      "logo.com/club-222222222222",
				}, {
					ClubUUID: "club-333333333333",
					Name:     "PMS",
					Field:    "SW 개발",
					Location: "2-3반 교실",
					Floor:    "3",
					LogoURI:  "logo.com/club-333333333333",
				}, {
					ClubUUID: "club-111111111111",
					Name:     "DMS",
					Field:    "SW 개발",
					Location: "2-1반 교실",
					Floor:    "3",
					LogoURI:  "logo.com/club-111111111111",
				}}, nil},
				"GetClubsWithClubUUIDs": {[]*model.Club{{
					UUID:       "club-222222222222",
					LeaderUUID: "student-222222222222",
				}, {
					UUID:       "club-333333333333",
					LeaderUUID: "student-333333333333",
				}, {
					UUID:       "club-111111111111",
					LeaderUUID: "student-111111111111",
				}}, nil},
				"Commit": {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusOK,
			ExpectClubInforms: []*clubproto.ClubInform{{
				ClubUUID:     "club-222222222222",
				LeaderUUID:   "student-222222222222",
				Name:         "SMS",
				ClubConcept:  "DMS의 소속부서 SMS 입니다!",
				Introduction: "School Management System 서비스를 개발 및 운영합니다",
				Link:         "facebook.com/DMS-SMS",
				Field:        "SW 개발",
				Location:     "2-2반 교실",
				Floor:        "3",
				LogoURI:      "logo.com/club-222222222222",
			}, {
				ClubUUID:   "club-333333333333",
				LeaderUUID: "student-333333333333",
				Name:       "PMS",
				Field:      "SW 개발",
				Location:   "2-3반 교실",
				Floor:      "3",
				LogoURI:    "logo.com/club-333333333333",
			}, {
				ClubUUID:   "club-111111111111",
				LeaderUUID: "student-111111111111",
				Name:       "DMS",
				Field:      "SW 개발",
				Location:   "2-1반 교실",
				Floor:      "3",
				LogoURI:    "logo.com/club-111111111111",
			}},
		}, { // Set Name and Field
			Name:  "SMS",
			Field: "SW",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetClubInformsSortByUpdateTime": {[]*model.ClubInform{{
					ClubUUID:     "club-222222222222",
					Name:         "SMS",
					ClubConcept:  "DMS의 소속부서 SMS 입니다!",
					Introduction: "School Management System 서비스를 개발 및 운영합니다",
					Link:         "facebook.com/DMS-SMS",
					Field:        "SW 개발",
					Location:     "2-2반 교실",
					Floor:        "3",
					LogoURI:      "logo.com/club-222222222222",
				}}, nil},
				"GetClubsWithClubUUIDs": {[]*model.Club{{
					UUID:       "club-222222222222",
					LeaderUUID: "student-222222222222",
				}}, nil},
				"Commit": {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusOK,
			ExpectClubInforms: []*clubproto.ClubInform{{
				ClubUUID:     "club-222222222222",
				LeaderUUID:   "student-222222222222",
				Name:         "SMS",
				ClubConcept:  "DMS의 소속부서 SMS 입니다!",
				Introduction: "School Management System 서비스를 개발 및 운영합니다",
				Link:         "facebook.com/DMS-SMS",
				Field:        "SW 개발",
				Location:     "2-2반 교실",
				Floor:        "3",
				LogoURI:      "logo.com/club-222222222222",
			}},
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
		}, { // forbidden (not student)
			UUID:           "parent-111111111112",
			ExpectedStatus: http.StatusForbidden,
		}, { // GetClubInformsSortByUpdateTime record not found
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                        {},
				"GetClubInformsSortByUpdateTime": {[]*model.ClubInform{}, gorm.ErrRecordNotFound},
				"Commit":                         {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusOK,
		}, { // GetClubInformsSortByUpdateTime record not found
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                        {},
				"GetClubInformsSortByUpdateTime": {[]*model.ClubInform{}, errors.New("db connect fail")},
				"Rollback":                       {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // GetClubsWithClubUUIDs returns unexpected error
			Name:  "SMS",
			Field: "SW",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetClubInformsSortByUpdateTime": {[]*model.ClubInform{{
					ClubUUID:     "club-222222222222",
					Name:         "SMS",
					ClubConcept:  "DMS의 소속부서 SMS 입니다!",
					Introduction: "School Management System 서비스를 개발 및 운영합니다",
					Link:         "facebook.com/DMS-SMS",
					Field:        "SW 개발",
					Location:     "2-2반 교실",
					Floor:        "3",
					LogoURI:      "logo.com/club-222222222222",
				}}, nil},
				"GetClubsWithClubUUIDs": {[]*model.Club{}, errors.New("some error occurs")},
				"Rollback":              {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // GetClubsWithClubUUIDs return RecordNotFound error
			Name:  "SMS",
			Field: "SW",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetClubInformsSortByUpdateTime": {[]*model.ClubInform{{
					ClubUUID:     "club-222222222222",
					Name:         "SMS",
					ClubConcept:  "DMS의 소속부서 SMS 입니다!",
					Introduction: "School Management System 서비스를 개발 및 운영합니다",
					Link:         "facebook.com/DMS-SMS",
					Field:        "SW 개발",
					Location:     "2-2반 교실",
					Floor:        "3",
					LogoURI:      "logo.com/club-222222222222",
				}}, nil},
				"GetClubsWithClubUUIDs": {[]*model.Club{}, gorm.ErrRecordNotFound},
				"Rollback":              {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // GetClubsWithClubUUIDs return abnormal length array
			Name:  "SMS",
			Field: "SW",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetClubInformsSortByUpdateTime": {[]*model.ClubInform{{
					ClubUUID:     "club-222222222222",
					Name:         "SMS",
					ClubConcept:  "DMS의 소속부서 SMS 입니다!",
					Introduction: "School Management System 서비스를 개발 및 운영합니다",
					Link:         "facebook.com/DMS-SMS",
					Field:        "SW 개발",
					Location:     "2-2반 교실",
					Floor:        "3",
					LogoURI:      "logo.com/club-222222222222",
				}}, nil},
				"GetClubsWithClubUUIDs": {[]*model.Club{}, nil},
				"Rollback":              {&gorm.DB{}},
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

		req := new(clubproto.GetClubsSortByUpdateTimeRequest)
		testCase.SetRequestContextOf(req)
		ctx := testCase.GetMetadataContext()

		resp := new(clubproto.GetClubsSortByUpdateTimeResponse)
		_ = handler.GetClubsSortByUpdateTime(ctx, req, resp)

		assert.Equalf(t, int(testCase.ExpectedStatus), int(resp.Status), "status assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Equalf(t, testCase.ExpectedCode, resp.Code, "code assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Equalf(t, testCase.ExpectClubInforms, resp.Clubs, "club informs assertion error (test case: %v, message: %s)", testCase, resp.Message)

		newMock.AssertExpectations(t)
	}
}