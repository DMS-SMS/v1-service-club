package test

import (
	"club/model"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func Test_Accessor_GetClubWithClubUUID(t *testing.T) {
	access, err := manager.BeginTx()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		access.Rollback()
		testGroup.Done()
	}()

	for _, club := range []*model.Club{
		{
			UUID:       "club-111111111111",
			LeaderUUID: "student-111111111111",
		},
	} {
		if _, err := access.CreateClub(club); err != nil {
			log.Fatal(err, club)
		}
	}

	tests := []struct {
		ClubUUID     string
		ExpectResult *model.Club
		ExpectError  error
	} {
		{
			ClubUUID: "club-111111111111",
			ExpectResult: &model.Club{
				UUID:       "club-111111111111",
				LeaderUUID: "student-111111111111",
			},
			ExpectError: nil,
		}, {
			ClubUUID:     "club-222222222222",
			ExpectResult: &model.Club{},
			ExpectError:  gorm.ErrRecordNotFound,
		},
	}

	for _, test := range tests {
		result, err := access.GetClubWithClubUUID(test.ClubUUID)

		assert.Equalf(t, test.ExpectError, err, "error assertion error (test case: %v)", test)
		assert.Equalf(t, test.ExpectResult, result.ExceptGormModel(), "result club assertion error (test case: %v)", test)
	}
}

func Test_Accessor_GetClubWithLeaderUUID(t *testing.T) {
	access, err := manager.BeginTx()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		access.Rollback()
		testGroup.Done()
	}()

	for _, club := range []*model.Club{
		{
			UUID:       "club-111111111111",
			LeaderUUID: "student-111111111111",
		},
	} {
		if _, err := access.CreateClub(club); err != nil {
			log.Fatal(err, club)
		}
	}

	tests := []struct {
		LeaderUUID   string
		ExpectResult *model.Club
		ExpectError  error
	}{
		{
			LeaderUUID: "student-111111111111",
			ExpectResult: &model.Club{
				UUID:       "club-111111111111",
				LeaderUUID: "student-111111111111",
			},
			ExpectError: nil,
		}, {
			LeaderUUID:   "student-222222222222",
			ExpectResult: &model.Club{},
			ExpectError:  gorm.ErrRecordNotFound,
		},
	}

	for _, test := range tests {
		result, err := access.GetClubWithLeaderUUID(test.LeaderUUID)

		assert.Equalf(t, test.ExpectError, err, "error assertion error (test case: %v)", test)
		assert.Equalf(t, test.ExpectResult, result.ExceptGormModel(), "result club assertion error (test case: %v)", test)
	}
}

func Test_Accessor_GetCurrentRecruitmentWithClubUUID(t *testing.T) {
	access, err := manager.BeginTx()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		access.Rollback()
		testGroup.Done()
	}()

	for _, club := range []*model.Club{
		{
			UUID:       "club-111111111111",
			LeaderUUID: "student-111111111111",
		}, {
			UUID:       "club-222222222222",
			LeaderUUID: "student-222222222222",
		}, {
			UUID:       "club-333333333333",
			LeaderUUID: "student-333333333333",
		},
	} {
		if _, err := access.CreateClub(club); err != nil {
			log.Fatal(err, club)
		}
	}

	now := time.Now()
	startTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	endTime := startTime.Add(time.Hour * 24 * 7)

	for _, recruitment := range []*model.ClubRecruitment{
		{ // 종료된 채용
			UUID:           "recruitment-111111111111",
			ClubUUID:       "club-111111111111",
			RecruitConcept: "첫 번째 공채",
			StartPeriod:    model.StartPeriod(time.Date(2020, time.Month(9), 17, 0, 0, 0, 0, time.UTC)),
			EndPeriod:      model.EndPeriod(time.Date(2020, time.Month(9), 24, 0, 0, 0, 0, time.UTC)),
		}, { // 현재 진행중인 채용
			UUID:           "recruitment-222222222222",
			ClubUUID:       "club-111111111111",
			RecruitConcept: "두 번째 공채",
			StartPeriod:    model.StartPeriod(startTime),
			EndPeriod:      model.EndPeriod(endTime),
		}, { // 종료된 채용
			UUID:           "recruitment-333333333333",
			ClubUUID:       "club-222222222222",
			RecruitConcept: "첫 번째 공채",
			StartPeriod:    model.StartPeriod(time.Date(2020, time.Month(9), 17, 0, 0, 0, 0, time.UTC)),
			EndPeriod:      model.EndPeriod(time.Date(2020, time.Month(9), 24, 0, 0, 0, 0, time.UTC)),
		}, { // 상시 채용
			UUID:           "recruitment-444444444444",
			ClubUUID:       "club-222222222222",
			RecruitConcept: "두 번째 상시 채용",
		},
	} {
		if _, err := access.CreateRecruitment(recruitment); err != nil {
			log.Fatal(err, recruitment)
		}
	}

	tests := []struct {
		ClubUUID   string
		ExpectResult *model.ClubRecruitment
		ExpectError  error
	}{
		{
			ClubUUID: "club-111111111111",
			ExpectResult: &model.ClubRecruitment{
				UUID:           "recruitment-222222222222",
				ClubUUID:       "club-111111111111",
				RecruitConcept: "두 번째 공채",
				StartPeriod:    model.StartPeriod(startTime),
				EndPeriod:      model.EndPeriod(endTime),
			},
			ExpectError: nil,
		}, {
			ClubUUID: "club-222222222222",
			ExpectResult: &model.ClubRecruitment{
				UUID:           "recruitment-444444444444",
				ClubUUID:       "club-222222222222",
				RecruitConcept: "두 번째 상시 채용",
			},
			ExpectError: nil,
		}, {
			ClubUUID:     "club-333333333333",
			ExpectError:  gorm.ErrRecordNotFound,
			ExpectResult: &model.ClubRecruitment{},
		},
	}

	for _, test := range tests {
		result, err := access.GetCurrentRecruitmentWithClubUUID(test.ClubUUID)

		assert.Equalf(t, test.ExpectError, err, "error assertion error (test case: %v)", test)
		assert.Equalf(t, test.ExpectResult, result.ExceptGormModel(), "result club assertion error (test case: %v)", test)
	}
}

func Test_Accessor_GetClubInformsSortByUpdateTime(t *testing.T) {
	access, err := manager.BeginTx()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		access.Rollback()
		testGroup.Done()
	}()

	for _, club := range []*model.Club{
		{
			UUID:       "club-111111111111",
			LeaderUUID: "student-111111111111",
		}, {
			UUID:       "club-222222222222",
			LeaderUUID: "student-222222222222",
		}, {
			UUID:       "club-333333333333",
			LeaderUUID: "student-333333333333",
		},
	} {
		if _, err := access.CreateClub(club); err != nil {
			log.Fatal(err, club)
		}
	}

	for _, inform := range []*model.ClubInform{
		{
			ClubUUID: "club-111111111111",
			Name:     "DMS",
			Field:    "SW 개발",
			Location: "2-1반 교실",
			Floor:    "3",
			LogoURI:  "logo.com/club-111111111111",
		}, {
			ClubUUID: "club-222222222222",
			Name:     "SMS",
			Field:    "SW 개발",
			Location: "2-2반 교실",
			Floor:    "3",
			LogoURI:  "logo.com/club-222222222222",
		}, {
			ClubUUID: "club-333333333333",
			Name:     "PMS",
			Field:    "SW 개발",
			Location: "2-3반 교실",
			Floor:    "3",
			LogoURI:  "logo.com/club-333333333333",
		},
	} {
		if _, err := access.CreateClubInform(inform); err != nil {
			log.Fatal(err, inform)
		}
		time.Sleep(time.Millisecond * 500)
	}

	time.Sleep(time.Millisecond * 500)
	if err, _ := access.ModifyClubInform("club-222222222222", &model.ClubInform{
		ClubConcept:  "DMS의 소속부서 SMS 입니다!",
		Introduction: "School Management System 서비스를 개발 및 운영합니다",
		Link:         "facebook.com/DMS-SMS",
	}); err != nil {
		log.Fatal(err)
	}

	tests := []struct {
		Offset, Limit int
		Field, Name   string
		ExpectResults []*model.ClubInform
		ExpectError   error
	} {
		{
			Offset: 0,
			Limit:  10,
			ExpectResults: []*model.ClubInform{
				{
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
				},
			},
			ExpectError: nil,
		}, {
			Offset:        10,
			Limit:         10,
			ExpectError:   gorm.ErrRecordNotFound,
		}, {
			Offset: 0,
			Limit:  10,
			Field:  "SW",
			ExpectResults: []*model.ClubInform{
				{
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
				},
			},
			ExpectError: nil,
		}, {
			Offset: 0,
			Limit:  10,
			Name:  "D",
			ExpectResults: []*model.ClubInform{
				{
					ClubUUID: "club-111111111111",
					Name:     "DMS",
					Field:    "SW 개발",
					Location: "2-1반 교실",
					Floor:    "3",
					LogoURI:  "logo.com/club-111111111111",
				},
			},
			ExpectError: nil,
		},
	}

	for _, test := range tests {
		informs, err := access.GetClubInformsSortByUpdateTime(test.Offset, test.Limit, test.Field, test.Name)

		var exceptedInforms []*model.ClubInform
		for _, inform := range informs {
			exceptedInforms = append(exceptedInforms, inform.ExceptGormModel())
		}

		assert.Equalf(t, test.ExpectError, err, "error assertion error (test case: %v)", test)
		assert.Equalf(t, test.ExpectResults, exceptedInforms, "result informs assertion error (test case: %v)", test)
	}
}
