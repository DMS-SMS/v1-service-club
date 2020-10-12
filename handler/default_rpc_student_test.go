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
	"time"
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
				"GetClubMembersWithClubUUIDs": {[][]*model.ClubMember{{
					{
						ClubUUID:    "club-222222222222",
						StudentUUID: "student-222222222222",
					}, {
						ClubUUID:    "club-222222222222",
						StudentUUID: "student-222222222223",
					},
				}, {
					{
						ClubUUID:    "club-333333333333",
						StudentUUID: "student-333333333333",
					},
				}, {
					{
						ClubUUID:    "club-111111111111",
						StudentUUID: "student-111111111111",
					}, {
						ClubUUID:    "club-111111111111",
						StudentUUID: "student-111111111112",
					}, {
						ClubUUID:    "club-111111111111",
						StudentUUID: "student-111111111113",
					},
				}}, nil},
				"Commit": {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusOK,
			ExpectClubInforms: []*clubproto.ClubInform{{
				ClubUUID:     "club-222222222222",
				LeaderUUID:   "student-222222222222",
				MemberUUIDs:  []string{"student-222222222222", "student-222222222223"},
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
				MemberUUIDs: []string{"student-333333333333"},
				Name:       "PMS",
				Field:      "SW 개발",
				Location:   "2-3반 교실",
				Floor:      "3",
				LogoURI:    "logo.com/club-333333333333",
			}, {
				ClubUUID:   "club-111111111111",
				LeaderUUID: "student-111111111111",
				MemberUUIDs:  []string{"student-111111111111", "student-111111111112", "student-111111111113"},
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
				"GetClubMembersWithClubUUIDs": {[][]*model.ClubMember{{
					{
						ClubUUID:    "club-222222222222",
						StudentUUID: "student-222222222222",
					}, {
						ClubUUID:    "club-222222222222",
						StudentUUID: "student-222222222223",
					},
				}, {
					{
						ClubUUID:    "club-333333333333",
						StudentUUID: "student-333333333333",
					},
				}, {
					{
						ClubUUID:    "club-111111111111",
						StudentUUID: "student-111111111111",
					}, {
						ClubUUID:    "club-111111111111",
						StudentUUID: "student-111111111112",
					}, {
						ClubUUID:    "club-111111111111",
						StudentUUID: "student-111111111113",
					},
				}}, nil},
				"Commit": {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusOK,
			ExpectClubInforms: []*clubproto.ClubInform{{
				ClubUUID:     "club-222222222222",
				LeaderUUID:   "student-222222222222",
				MemberUUIDs:  []string{"student-222222222222", "student-222222222223"},
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
				MemberUUIDs: []string{"student-333333333333"},
				Name:       "PMS",
				Field:      "SW 개발",
				Location:   "2-3반 교실",
				Floor:      "3",
				LogoURI:    "logo.com/club-333333333333",
			}, {
				ClubUUID:   "club-111111111111",
				LeaderUUID: "student-111111111111",
				MemberUUIDs:  []string{"student-111111111111", "student-111111111112", "student-111111111113"},
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
				"GetClubMembersWithClubUUIDs": {[][]*model.ClubMember{{
					{
						ClubUUID:    "club-222222222222",
						StudentUUID: "student-222222222222",
					}, {
						ClubUUID:    "club-222222222222",
						StudentUUID: "student-222222222223",
					},
				}}, nil},
				"Commit": {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusOK,
			ExpectClubInforms: []*clubproto.ClubInform{{
				ClubUUID:     "club-222222222222",
				LeaderUUID:   "student-222222222222",
				MemberUUIDs:  []string{"student-222222222222", "student-222222222223"},
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
				"GetClubsWithClubUUIDs": {[]*model.Club{{}}, errors.New("some error occurs")},
				"Rollback":              {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // GetClubMembersWithClubUUID returns NOT FOUND error
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
					UUID:       "club-111111111111",
					LeaderUUID: "student-111111111111",
				}}, nil},
				"GetClubMembersWithClubUUIDs": {[][]*model.ClubMember{{}, {}}, gorm.ErrRecordNotFound},
				"Commit":                      {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusOK,
			ExpectClubInforms: []*clubproto.ClubInform{{
				ClubUUID:     "club-222222222222",
				LeaderUUID:   "student-222222222222",
				MemberUUIDs:  []string{},
				Name:         "SMS",
				ClubConcept:  "DMS의 소속부서 SMS 입니다!",
				Introduction: "School Management System 서비스를 개발 및 운영합니다",
				Link:         "facebook.com/DMS-SMS",
				Field:        "SW 개발",
				Location:     "2-2반 교실",
				Floor:        "3",
				LogoURI:      "logo.com/club-222222222222",
			}, {
				ClubUUID:    "club-111111111111",
				LeaderUUID:  "student-111111111111",
				MemberUUIDs: []string{},
				Name:        "DMS",
				Field:       "SW 개발",
				Location:    "2-1반 교실",
				Floor:       "3",
				LogoURI:     "logo.com/club-111111111111",
			}},
		}, { // GetClubMembersWithClubUUID returns unexpected error
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
					UUID:       "club-111111111111",
					LeaderUUID: "student-111111111111",
				}}, nil},
				"GetClubMembersWithClubUUIDs": {[][]*model.ClubMember{{}, {}}, errors.New("unexpected error")},
				"Rollback":                    {&gorm.DB{}},
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
		assert.Equalf(t, testCase.ExpectClubInforms, resp.Informs, "club informs assertion error (test case: %v, message: %s)", testCase, resp.Message)

		newMock.AssertExpectations(t)
	}
}

func Test_Default_GetRecruitmentsSortByCreateTime(t *testing.T) {
	now := time.Now()
	startTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	endTime := startTime.Add(time.Hour * 24 * 7)

	tests := []*test.GetRecruitmentsSortByCreateTimeCase{
		{ // success case
			Start: 0,
			Count: 10,
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetCurrentRecruitmentsSortByCreateTime": {[]*model.ClubRecruitment{{
					UUID:           "recruitment-555555555555",
					ClubUUID:       "club-333333333333",
					RecruitConcept: "첫 번째 상시 채용",
				}, {
					UUID:           "recruitment-222222222222",
					ClubUUID:       "club-111111111111",
					RecruitConcept: "두 번째 공채",
					StartPeriod:    model.StartPeriod(startTime),
					EndPeriod:      model.EndPeriod(endTime),
				}}, nil},
				"GetRecruitMembersWithRecruitmentUUIDs": {[][]*model.RecruitMember{
					{
						{
							RecruitmentUUID: "recruitment-555555555555",
							Grade:           "1",
							Field:           "서버 개발자",
							Number:          "1",
						}, {
							RecruitmentUUID: "recruitment-555555555555",
							Grade:           "1",
							Field:           "웹 프론트 개발자",
							Number:          "1",
						},
					}, {
						{
							RecruitmentUUID: "recruitment-222222222222",
							Grade:           "1",
							Field:           "모든 분야",
							Number:          "4",
						},
					},
				}, nil},
				"Commit": {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusOK,
			ExpectRecruitments: []*clubproto.RecruitmentInform{{
				RecruitmentUUID: "recruitment-555555555555",
				ClubUUID:        "club-333333333333",
				RecruitConcept:  "첫 번째 상시 채용",
				RecruitMembers: []*clubproto.RecruitMember{{
					Grade:  "1",
					Field:  "서버 개발자",
					Number: "1",
				}, {
					Grade:  "1",
					Field:  "웹 프론트 개발자",
					Number: "1",
				}},
			}, {
				RecruitmentUUID: "recruitment-222222222222",
				ClubUUID:        "club-111111111111",
				RecruitConcept:  "두 번째 공채",
				RecruitMembers: []*clubproto.RecruitMember{{
					Grade:  "1",
					Field:  "모든 분야",
					Number: "4",
				}},
				StartPeriod: startTime.Format("2006-01-02"),
				EndPeriod:   endTime.Format("2006-01-02"),
			}},
		}, { // Start, Count Set X -> default(0, 10)
			UUID:  "admin-111111111111",
			Start: 0,
			Count: 10,
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetCurrentRecruitmentsSortByCreateTime": {[]*model.ClubRecruitment{{
					UUID:           "recruitment-555555555555",
					ClubUUID:       "club-333333333333",
					RecruitConcept: "첫 번째 상시 채용",
				}, {
					UUID:           "recruitment-222222222222",
					ClubUUID:       "club-111111111111",
					RecruitConcept: "두 번째 공채",
					StartPeriod:    model.StartPeriod(startTime),
					EndPeriod:      model.EndPeriod(endTime),
				}}, nil},
				"GetRecruitMembersWithRecruitmentUUIDs": {[][]*model.RecruitMember{
					{
						{
							RecruitmentUUID: "recruitment-555555555555",
							Grade:           "1",
							Field:           "서버 개발자",
							Number:          "1",
						}, {
							RecruitmentUUID: "recruitment-555555555555",
							Grade:           "1",
							Field:           "웹 프론트 개발자",
							Number:          "1",
						},
					}, {
						{
							RecruitmentUUID: "recruitment-222222222222",
							Grade:           "1",
							Field:           "모든 분야",
							Number:          "4",
						},
					},
				}, nil},
				"Commit": {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusOK,
			ExpectRecruitments: []*clubproto.RecruitmentInform{{
				RecruitmentUUID: "recruitment-555555555555",
				ClubUUID:        "club-333333333333",
				RecruitConcept:  "첫 번째 상시 채용",
				RecruitMembers: []*clubproto.RecruitMember{{
					Grade:  "1",
					Field:  "서버 개발자",
					Number: "1",
				}, {
					Grade:  "1",
					Field:  "웹 프론트 개발자",
					Number: "1",
				}},
			}, {
				RecruitmentUUID: "recruitment-222222222222",
				ClubUUID:        "club-111111111111",
				RecruitConcept:  "두 번째 공채",
				RecruitMembers: []*clubproto.RecruitMember{{
					Grade:  "1",
					Field:  "모든 분야",
					Number: "4",
				}},
				StartPeriod: startTime.Format("2006-01-02"),
				EndPeriod:   endTime.Format("2006-01-02"),
			}},
		}, { // Set Name and Field
			Name:  "SMS",
			Field: "SW",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetCurrentRecruitmentsSortByCreateTime": {[]*model.ClubRecruitment{{
					UUID:           "recruitment-555555555555",
					ClubUUID:       "club-333333333333",
					RecruitConcept: "첫 번째 상시 채용",
				}}, nil},
				"GetRecruitMembersWithRecruitmentUUIDs": {[][]*model.RecruitMember{
					{
						{
							RecruitmentUUID: "recruitment-555555555555",
							Grade:           "1",
							Field:           "서버 개발자",
							Number:          "1",
						}, {
							RecruitmentUUID: "recruitment-555555555555",
							Grade:           "1",
							Field:           "웹 프론트 개발자",
							Number:          "1",
						},
					},
				}, nil},
				"Commit": {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusOK,
			ExpectRecruitments: []*clubproto.RecruitmentInform{{
				RecruitmentUUID: "recruitment-555555555555",
				ClubUUID:        "club-333333333333",
				RecruitConcept:  "첫 번째 상시 채용",
				RecruitMembers: []*clubproto.RecruitMember{{
					Grade:  "1",
					Field:  "서버 개발자",
					Number: "1",
				}, {
					Grade:  "1",
					Field:  "웹 프론트 개발자",
					Number: "1",
				}},
			}},
		},  { // no exist X-Request-ID -> Proxy Authorization Required
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
		}, { // GetCurrentRecruitmentsSortByCreateTime record not found
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                                {},
				"GetCurrentRecruitmentsSortByCreateTime": {[]*model.ClubRecruitment{}, gorm.ErrRecordNotFound},
				"Commit":                                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusOK,
		}, { // GetCurrentRecruitmentsSortByCreateTime returns unexpected error
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                                {},
				"GetCurrentRecruitmentsSortByCreateTime": {[]*model.ClubRecruitment{}, errors.New("db connect fail")},
				"Rollback":                               {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // GetRecruitMembersWithRecruitmentUUIDs returns RecordNotFound error
			Name:  "SMS",
			Field: "SW",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetCurrentRecruitmentsSortByCreateTime": {[]*model.ClubRecruitment{{
					UUID:           "recruitment-555555555555",
					ClubUUID:       "club-333333333333",
					RecruitConcept: "첫 번째 상시 채용",
				}}, nil},
				"GetRecruitMembersWithRecruitmentUUIDs": {[][]*model.RecruitMember{{}}, gorm.ErrRecordNotFound},
				"Commit":                                {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusOK,
			ExpectRecruitments: []*clubproto.RecruitmentInform{{
				RecruitmentUUID: "recruitment-555555555555",
				ClubUUID:        "club-333333333333",
				RecruitConcept:  "첫 번째 상시 채용",
				RecruitMembers: []*clubproto.RecruitMember{},
			}},
		}, { // GetRecruitMembersWithRecruitmentUUIDs return unexpected error
			Name:  "SMS",
			Field: "SW",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetCurrentRecruitmentsSortByCreateTime": {[]*model.ClubRecruitment{{
					UUID:           "recruitment-555555555555",
					ClubUUID:       "club-333333333333",
					RecruitConcept: "첫 번째 상시 채용",
				}}, nil},
				"GetRecruitMembersWithRecruitmentUUIDs": {[][]*model.RecruitMember{{}}, errors.New("unexpected error")},
				"Rollback":                              {&gorm.DB{}},
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

		req := new(clubproto.GetRecruitmentsSortByCreateTimeRequest)
		testCase.SetRequestContextOf(req)
		ctx := testCase.GetMetadataContext()

		resp := new(clubproto.GetRecruitmentsSortByCreateTimeResponse)
		_ = handler.GetRecruitmentsSortByCreateTime(ctx, req, resp)

		assert.Equalf(t, int(testCase.ExpectedStatus), int(resp.Status), "status assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Equalf(t, testCase.ExpectedCode, resp.Code, "code assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Equalf(t, testCase.ExpectRecruitments, resp.Recruitments, "recruitments assertion error (test case: %v, message: %s)", testCase, resp.Message)

		newMock.AssertExpectations(t)
	}
}

func Test_Default_GetClubInformWithUUID(t *testing.T) {
	tests := []test.GetClubInformWithUUIDCase{
		{ // success case
			UUID:     "student-222222222222",
			ClubUUID: "club-222222222222",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetClubWithClubUUID": {&model.Club{
					UUID:       "club-222222222222",
					LeaderUUID: "student-222222222222",
				}, nil},
				"GetClubInformWithClubUUID": {&model.ClubInform{
					ClubUUID:     "club-222222222222",
					Name:         "SMS",
					ClubConcept:  "DMS의 소속부서 SMS 입니다!",
					Introduction: "School Management System 서비스를 개발 및 운영합니다",
					Link:         "facebook.com/DMS-SMS",
					Field:        "SW 개발",
					Location:     "2-2반 교실",
					Floor:        "3",
					LogoURI:      "logo.com/club-222222222222",
				}, nil},
				"GetClubMembersWithClubUUID": {[]*model.ClubMember{
					{
						ClubUUID:    "club-222222222222",
						StudentUUID: "student-222222222222",
					}, {
						ClubUUID:    "club-222222222222",
						StudentUUID: "student-222222222223",
					},
				}, nil},
				"Commit": {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusOK,
			ExpectInform: &clubproto.ClubInform{
				ClubUUID:     "club-222222222222",
				LeaderUUID:   "student-222222222222",
				MemberUUIDs:  []string{"student-222222222222", "student-222222222223"},
				Name:         "SMS",
				ClubConcept:  "DMS의 소속부서 SMS 입니다!",
				Introduction: "School Management System 서비스를 개발 및 운영합니다",
				Link:         "facebook.com/DMS-SMS",
				Field:        "SW 개발",
				Location:     "2-2반 교실",
				Floor:        "3",
				LogoURI:      "logo.com/club-222222222222",
			},
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
		}, { // club uuid not exist
			UUID:     "admin-222222222222",
			ClubUUID: "club-333333333333",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetClubWithClubUUID": {&model.Club{}, gorm.ErrRecordNotFound},
				"Rollback":            {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusNotFound,
		}, { // GetClubWithClubUUID returns unexpected error
			UUID:     "admin-222222222222",
			ClubUUID: "club-333333333333",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":             {},
				"GetClubWithClubUUID": {&model.Club{}, errors.New("unexpected error")},
				"Rollback":            {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // GetClubInformWithClubUUID returns unexpected error
			UUID:     "student-222222222222",
			ClubUUID: "club-222222222222",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetClubWithClubUUID": {&model.Club{
					UUID:       "club-222222222222",
					LeaderUUID: "student-222222222222",
				}, nil},
				"GetClubInformWithClubUUID": {&model.ClubInform{}, errors.New("unexpected error")},
				"Rollback":                  {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // GetClubMembersWithClubUUID returns unexpected error
			UUID:     "student-222222222222",
			ClubUUID: "club-222222222222",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetClubWithClubUUID": {&model.Club{
					UUID:       "club-222222222222",
					LeaderUUID: "student-222222222222",
				}, nil},
				"GetClubInformWithClubUUID": {&model.ClubInform{
					ClubUUID:     "club-222222222222",
					Name:         "SMS",
					ClubConcept:  "DMS의 소속부서 SMS 입니다!",
					Introduction: "School Management System 서비스를 개발 및 운영합니다",
					Link:         "facebook.com/DMS-SMS",
					Field:        "SW 개발",
					Location:     "2-2반 교실",
					Floor:        "3",
					LogoURI:      "logo.com/club-222222222222",
				}, nil},
				"GetClubMembersWithClubUUID": {[]*model.ClubMember{}, errors.New("unexpected error")},
				"Rollback":                   {&gorm.DB{}},
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

		req := new(clubproto.GetClubInformWithUUIDRequest)
		testCase.SetRequestContextOf(req)
		ctx := testCase.GetMetadataContext()

		resp := new(clubproto.GetClubInformWithUUIDResponse)
		_ = handler.GetClubInformWithUUID(ctx, req, resp)
		var respInform *clubproto.ClubInform
		if testCase.ExpectedStatus == http.StatusOK {
			respInform = &clubproto.ClubInform{
				ClubUUID:     resp.ClubUUID,
				LeaderUUID:   resp.LeaderUUID,
				MemberUUIDs:  resp.MemberUUIDs,
				Name:         resp.Name,
				ClubConcept:  resp.ClubConcept,
				Introduction: resp.Introduction,
				Floor:        resp.Floor,
				Location:     resp.Location,
				Field:        resp.Field,
				Link:         resp.Link,
				LogoURI:      resp.LogoURI,
			}
		}
		assert.Equalf(t, int(testCase.ExpectedStatus), int(resp.Status), "status assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Equalf(t, testCase.ExpectInform, respInform, "club informs assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Equalf(t, testCase.ExpectedCode, resp.Code, "code assertion error (test case: %v, message: %s)", testCase, resp.Message)

		newMock.AssertExpectations(t)
	}
}
/*
func Test_Default_GetClubInformsWithUUIDs(t *testing.T) {
	tests := []test.GetClubInformsWithUUIDsCase{
		{ // success case
			UUID:     "student-111111111111",
			ClubUUIDs: []string{"club-111111111111", "club-222222222222", "club-111111111111"},
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetClubWithClubUUIDs": {[]*model.Club{{
					UUID:       "club-111111111111",
					LeaderUUID: "student-111111111111",
				}, {
					UUID:       "club-222222222222",
					LeaderUUID: "student-222222222222",
				}, {
					UUID:       "club-111111111111",
					LeaderUUID: "student-111111111111",
				}}, nil},
				"GetClubInformWithClubUUIDs": {[]*model.ClubInform{{
					ClubUUID: "club-111111111111",
					Name:     "DMS",
					Field:    "SW 개발",
					Location: "2-1반 교실",
					Floor:    "3",
					LogoURI:  "logo.com/club-111111111111",
				}, {
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
					ClubUUID: "club-111111111111",
					Name:     "DMS",
					Field:    "SW 개발",
					Location: "2-1반 교실",
					Floor:    "3",
					LogoURI:  "logo.com/club-111111111111",
				}}, nil},
				"GetClubMembersWithClubUUIDs": {[][]*model.ClubMember{{
					{
						ClubUUID:    "club-111111111111",
						StudentUUID: "student-111111111111",
					}, {
						ClubUUID:    "club-111111111111",
						StudentUUID: "student-111111111112",
					}, {
						ClubUUID:    "club-111111111111",
						StudentUUID: "student-111111111113",
					},
				}, {
					{
						ClubUUID:    "club-222222222222",
						StudentUUID: "student-222222222222",
					}, {
						ClubUUID:    "club-222222222222",
						StudentUUID: "student-222222222223",
					},
				}, {
					{
						ClubUUID:    "club-111111111111",
						StudentUUID: "student-111111111111",
					}, {
						ClubUUID:    "club-111111111111",
						StudentUUID: "student-111111111112",
					}, {
						ClubUUID:    "club-111111111111",
						StudentUUID: "student-111111111113",
					},
				}}, nil},
				"Commit": {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusOK,
			ExpectInforms: []*clubproto.ClubInform{{
				ClubUUID:   "club-111111111111",
				LeaderUUID: "student-111111111111",
				MemberUUIDs:  []string{"student-111111111111", "student-111111111112", "student-111111111113"},
				Name:       "DMS",
				Field:      "SW 개발",
				Location:   "2-1반 교실",
				Floor:      "3",
				LogoURI:    "logo.com/club-111111111111",
			}, {
				ClubUUID:     "club-222222222222",
				LeaderUUID:   "student-222222222222",
				MemberUUIDs:  []string{"student-222222222222", "student-222222222223"},
				Name:         "SMS",
				ClubConcept:  "DMS의 소속부서 SMS 입니다!",
				Introduction: "School Management System 서비스를 개발 및 운영합니다",
				Link:         "facebook.com/DMS-SMS",
				Field:        "SW 개발",
				Location:     "2-2반 교실",
				Floor:        "3",
				LogoURI:      "logo.com/club-222222222222",
			}, {
				ClubUUID:   "club-111111111111",
				LeaderUUID: "student-111111111111",
				MemberUUIDs:  []string{"student-111111111111", "student-111111111112", "student-111111111113"},
				Name:       "DMS",
				Field:      "SW 개발",
				Location:   "2-1반 교실",
				Floor:      "3",
				LogoURI:    "logo.com/club-111111111111",
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
		}, { // not student or admin uuid
			UUID:            "parent-111111111111",
			ExpectedMethods: map[test.Method]test.Returns{},
			ExpectedStatus:  http.StatusForbidden,
		}, { // club uuids include not exist uuid
			UUID:      "admin-222222222222",
			ClubUUIDs: []string{"club-444444444444", "club-333333333333"},
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":              {},
				"GetClubWithClubUUIDs": {[]*model.Club{{}}, gorm.ErrRecordNotFound},
				"Rollback":             {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusNotFound,
		}, { // GetClubWithClubUUID returns unexpected error
			UUID:      "admin-222222222222",
			ClubUUIDs: []string{"club-333333333333"},
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":              {},
				"GetClubWithClubUUIDs": {[]*model.Club{{}}, errors.New("unexpected error")},
				"Rollback":             {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // GetClubInformWithClubUUID returns unexpected error
			UUID:     "student-222222222222",
			ClubUUIDs: []string{"club-222222222222"},
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetClubWithClubUUIDs": {[]*model.Club{{
					UUID:       "club-222222222222",
					LeaderUUID: "student-222222222222",
				}}, nil},
				"GetClubInformWithClubUUIDs": {[]*model.ClubInform{{}}, errors.New("unexpected error")},
				"Rollback":                   {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // GetClubInformWithClubUUID returns not found error
			UUID:     "student-222222222222",
			ClubUUIDs: []string{"club-222222222222"},
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetClubWithClubUUIDs": {[]*model.Club{{
					UUID:       "club-222222222222",
					LeaderUUID: "student-222222222222",
				}}, nil},
				"GetClubInformWithClubUUIDs": {[]*model.ClubInform{{}}, gorm.ErrRecordNotFound},
				"Rollback":                   {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // GetClubMembersWithClubUUID returns unexpected error
			UUID:     "student-222222222222",
			ClubUUIDs: []string{"club-222222222222"},
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetClubWithClubUUIDs": {[]*model.Club{{
					UUID:       "club-222222222222",
					LeaderUUID: "student-222222222222",
				}}, nil},
				"GetClubInformWithClubUUIDs": {[]*model.ClubInform{{
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
				"GetClubMembersWithClubUUIDs": {[][]*model.ClubMember{{}}, errors.New("unexpected error")},
				"Rollback":                    {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // GetClubMembersWithClubUUID returns not found error
			UUID:     "student-222222222222",
			ClubUUIDs: []string{"club-222222222222"},
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetClubWithClubUUIDs": {[]*model.Club{{
					UUID:       "club-222222222222",
					LeaderUUID: "student-222222222222",
				}}, nil},
				"GetClubInformWithClubUUIDs": {[]*model.ClubInform{{
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
				"GetClubMembersWithClubUUIDs": {[][]*model.ClubMember{{}}, gorm.ErrRecordNotFound},
				"Commit":                      {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusOK,
			ExpectInforms: []*clubproto.ClubInform{{
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
		},
	}

	for _, testCase := range tests {
		newMock := &mock.Mock{}
		handler := newDefaultMockHandler(newMock)

		testCase.ChangeEmptyValueToValidValue()
		testCase.ChangeEmptyReplaceValueToEmptyValue()
		testCase.OnExpectMethodsTo(newMock)

		req := new(clubproto.GetClubInformsWithUUIDsRequest)
		testCase.SetRequestContextOf(req)
		ctx := testCase.GetMetadataContext()

		resp := new(clubproto.GetClubInformsWithUUIDsResponse)
		_ = handler.GetClubInformsWithUUIDs(ctx, req, resp)
		assert.Equalf(t, int(testCase.ExpectedStatus), int(resp.Status), "status assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Equalf(t, testCase.ExpectInforms, resp.Informs, "club informs assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Equalf(t, testCase.ExpectedCode, resp.Code, "code assertion error (test case: %v, message: %s)", testCase, resp.Message)

		newMock.AssertExpectations(t)
	}
}
*/