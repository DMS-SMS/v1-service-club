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

func Test_Accessor_GetCurrentRecruitmentsSortByCreateTime(t *testing.T) {
	access, err := manager.BeginTx()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		access.Rollback()
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
		}, {
			UUID:       "club-444444444444",
			LeaderUUID: "student-444444444444",
		}, {
			UUID:       "club-555555555555",
			LeaderUUID: "student-555555555555",
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
			Name:     "PPP",
			Field:    "SW 개발",
			Location: "2-3반 교실",
			Floor:    "3",
			LogoURI:  "logo.com/club-333333333333",
		}, {
			ClubUUID: "club-444444444444",
			Name:     "PMS",
			Field:    "SW 개발",
			Location: "2-4반 교실",
			Floor:    "3",
			LogoURI:  "logo.com/club-444444444444",
		}, {
			ClubUUID: "club-555555555555",
			Name:     "MSMS",
			Field:    "SW 개발",
			Location: "2-5반 교실",
			Floor:    "3",
			LogoURI:  "logo.com/club-555555555555",
		},
	} {
		if _, err := access.CreateClubInform(inform); err != nil {
			log.Fatal(err, inform)
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
		}, { // 상시 채용
			UUID:           "recruitment-555555555555",
			ClubUUID:       "club-333333333333",
			RecruitConcept: "첫 번째 상시 채용",
		}, { // 상시 채용 (모집 삭제 예정)
			UUID:           "recruitment-666666666666",
			ClubUUID:       "club-444444444444",
			RecruitConcept: "첫 번째 상시 채용",
		}, { // 상시 채용 (동아리 삭제 예정)
			UUID:           "recruitment-777777777777",
			ClubUUID:       "club-555555555555",
			RecruitConcept: "첫 번째 상시 채용",
		},
	} {
		if _, err := access.CreateRecruitment(recruitment); err != nil {
			log.Fatal(err, recruitment)
		}
		time.Sleep(time.Millisecond * 500)
	}

	if err, _ := access.DeleteClubInform("club-555555555555"); err != nil {
		log.Fatal(err)
	}

	if err, _ := access.DeleteRecruitment("recruitment-666666666666"); err != nil {
		log.Fatal(err)
	}

	tests := []struct {
		Offset, Limit int
		Field, Name   string
		ExpectResults []*model.ClubRecruitment
		ExpectError   error
	} {
		{
			Offset: 0,
			Limit: 10,
			ExpectResults: []*model.ClubRecruitment{
				{
					UUID:           "recruitment-555555555555",
					ClubUUID:       "club-333333333333",
					RecruitConcept: "첫 번째 상시 채용",
				}, {
					UUID:           "recruitment-444444444444",
					ClubUUID:       "club-222222222222",
					RecruitConcept: "두 번째 상시 채용",
				}, {
					UUID:           "recruitment-222222222222",
					ClubUUID:       "club-111111111111",
					RecruitConcept: "두 번째 공채",
					StartPeriod:    model.StartPeriod(startTime),
					EndPeriod:      model.EndPeriod(endTime),
				},
			},
			ExpectError: nil,
		}, {
			Offset:      10,
			Limit:       10,
			ExpectError: gorm.ErrRecordNotFound,
		}, {
			Offset: 0,
			Limit:  10,
			Name:   "MS",
			ExpectResults: []*model.ClubRecruitment{
				{
					UUID:           "recruitment-444444444444",
					ClubUUID:       "club-222222222222",
					RecruitConcept: "두 번째 상시 채용",
				}, {
					UUID:           "recruitment-222222222222",
					ClubUUID:       "club-111111111111",
					RecruitConcept: "두 번째 공채",
					StartPeriod:    model.StartPeriod(startTime),
					EndPeriod:      model.EndPeriod(endTime),
				},
			},
			ExpectError: nil,
		}, {
			Offset: 0,
			Limit:  3,
			Field:  "SW",
			ExpectResults: []*model.ClubRecruitment{
				{
					UUID:           "recruitment-555555555555",
					ClubUUID:       "club-333333333333",
					RecruitConcept: "첫 번째 상시 채용",
				}, {
					UUID:           "recruitment-444444444444",
					ClubUUID:       "club-222222222222",
					RecruitConcept: "두 번째 상시 채용",
				}, {
					UUID:           "recruitment-222222222222",
					ClubUUID:       "club-111111111111",
					RecruitConcept: "두 번째 공채",
					StartPeriod:    model.StartPeriod(startTime),
					EndPeriod:      model.EndPeriod(endTime),
				},
			},
			ExpectError: nil,
		},
	}

	for _, test := range tests {
		recruitments, err := access.GetCurrentRecruitmentsSortByCreateTime(test.Offset, test.Limit, test.Field, test.Name)

		var exceptedRecruitments []*model.ClubRecruitment
		for _, recruitment := range recruitments {
			exceptedRecruitments = append(exceptedRecruitments, recruitment.ExceptGormModel())
		}

		assert.Equalf(t, test.ExpectError, err, "error assertion error (test case: %v)", test)
		assert.Equalf(t, test.ExpectResults, exceptedRecruitments, "result recruitments assertion error (test case: %v)", test)
	}
}

func Test_Accessor_GetClubInformWithClubUUID(t *testing.T) {
	access, err := manager.BeginTx()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		access.Rollback()
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

	for _, inform := range []*model.ClubInform{
		{
			ClubUUID: "club-111111111111",
			Name:     "DMS",
			Field:    "SW 개발",
			Location: "2-1반 교실",
			Floor:    "3",
			LogoURI:  "logo.com/club-111111111111",
		},
	} {
		if _, err := access.CreateClubInform(inform); err != nil {
			log.Fatal(err, inform)
		}
	}

	tests := []struct {
		ClubUUID   string
		ExpectResult *model.ClubInform
		ExpectError  error
	}{
		{
			ClubUUID: "club-111111111111",
			ExpectResult: &model.ClubInform{
				ClubUUID: "club-111111111111",
				Name:     "DMS",
				Field:    "SW 개발",
				Location: "2-1반 교실",
				Floor:    "3",
				LogoURI:  "logo.com/club-111111111111",
			},
			ExpectError: nil,
		}, {
			ClubUUID:     "club-222222222222",
			ExpectResult: &model.ClubInform{},
			ExpectError:  gorm.ErrRecordNotFound,
		},
	}

	for _, test := range tests {
		result, err := access.GetClubInformWithClubUUID(test.ClubUUID)

		assert.Equalf(t, test.ExpectError, err, "error assertion error (test case: %v)", test)
		assert.Equalf(t, test.ExpectResult, result.ExceptGormModel(), "result club inform assertion error (test case: %v)", test)
	}
}

func Test_Accessor_GetRecruitmentWithRecruitmentUUID(t *testing.T) {
	access, err := manager.BeginTx()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		access.Rollback()
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

	startTime := time.Date(2020, time.Month(9), 17, 0, 0, 0, 0, time.Local)
	endTime := time.Date(2020, time.Month(9), 24, 0, 0, 0, 0, time.Local)

	for _, recruitment := range []*model.ClubRecruitment{
		{ // 종료된 채용
			UUID:           "recruitment-111111111111",
			ClubUUID:       "club-111111111111",
			RecruitConcept: "첫 번째 공채",
			StartPeriod:    model.StartPeriod(startTime),
			EndPeriod:      model.EndPeriod(endTime),
		}, { // 현재 진행중인 채용
			UUID:           "recruitment-222222222222",
			ClubUUID:       "club-111111111111",
			RecruitConcept: "두 번째 공채",
		},
	} {
		if _, err := access.CreateRecruitment(recruitment); err != nil {
			log.Fatal(err, recruitment)
		}
	}

	tests := []struct {
		RecruitmentUUID string
		ExpectResult    *model.ClubRecruitment
		ExpectError     error
	} {
		{
			RecruitmentUUID: "recruitment-111111111111",
			ExpectResult: &model.ClubRecruitment{
				UUID:           "recruitment-111111111111",
				ClubUUID:       "club-111111111111",
				RecruitConcept: "첫 번째 공채",
				StartPeriod:    model.StartPeriod(startTime),
				EndPeriod:      model.EndPeriod(endTime),
			},
			ExpectError: nil,
		}, {
			RecruitmentUUID: "club-333333333333",
			ExpectResult:    &model.ClubRecruitment{},
			ExpectError:     gorm.ErrRecordNotFound,
		},
	}

	for _, test := range tests {
		result, err := access.GetRecruitmentWithRecruitmentUUID(test.RecruitmentUUID)

		assert.Equalf(t, test.ExpectError, err, "error assertion error (test case: %v)", test)
		assert.Equalf(t, test.ExpectResult, result.ExceptGormModel(), "result recruitment assertion error (test case: %v)", test)
	}
}

func Test_Accessor_GetClubMembersWithClubUUID(t *testing.T) {
	access, err := manager.BeginTx()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		access.Rollback()
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

	for _, member := range []*model.ClubMember{
		{
			ClubUUID:    "club-111111111111",
			StudentUUID: "student-111111111111",
		}, {
			ClubUUID:    "club-111111111111",
			StudentUUID: "student-222222222222",
		}, {
			ClubUUID:    "club-111111111111",
			StudentUUID: "student-333333333333",
		}, {
			ClubUUID:    "club-222222222222",
			StudentUUID: "student-333333333333",
		},
	} {
		if _, err := access.CreateClubMember(member); err != nil {
			log.Fatal(err)
		}
	}

	if err, _ := access.DeleteClubMember("club-111111111111", "student-333333333333"); err != nil {
		log.Fatal(err)
	}

	tests := []struct {
		ClubUUID      string
		ExpectResults []*model.ClubMember
		ExpectError   error
	} {
		{
			ClubUUID: "club-111111111111",
			ExpectResults: []*model.ClubMember{
				{
					ClubUUID:    "club-111111111111",
					StudentUUID: "student-111111111111",
				}, {
					ClubUUID:    "club-111111111111",
					StudentUUID: "student-222222222222",
				},
			},
			ExpectError: nil,
		}, {
			ClubUUID:      "club-333333333333",
			ExpectError:   gorm.ErrRecordNotFound,
		}, {
			ClubUUID:      "club-444444444444",
			ExpectError:   gorm.ErrRecordNotFound,
		},
	}

	for _, test := range tests {
		resultMembers, err := access.GetClubMembersWithClubUUID(test.ClubUUID)

		var exceptedResult []*model.ClubMember
		for _, member := range resultMembers {
			exceptedResult = append(exceptedResult, member.ExceptGormModel())
		}

		assert.Equalf(t, test.ExpectError, err, "error assertion error (test case: %v)", test)
		assert.Equalf(t, test.ExpectResults, exceptedResult, "result members assertion error (test case: %v)", test)
	}
}

func Test_Accessor_GetRecruitMembersWithRecruitmentUUID(t *testing.T) {
	access, err := manager.BeginTx()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		access.Rollback()
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

	startTime := time.Date(2020, time.Month(9), 17, 0, 0, 0, 0, time.Local)
	endTime := time.Date(2020, time.Month(9), 24, 0, 0, 0, 0, time.Local)

	for _, recruitment := range []*model.ClubRecruitment{
		{ // 종료된 채용
			UUID:           "recruitment-111111111111",
			ClubUUID:       "club-111111111111",
			RecruitConcept: "첫 번째 공채",
			StartPeriod:    model.StartPeriod(startTime),
			EndPeriod:      model.EndPeriod(endTime),
		}, { // 현재 진행중인 채용
			UUID:           "recruitment-222222222222",
			ClubUUID:       "club-111111111111",
			RecruitConcept: "두 번째 공채",
		},
	} {
		if _, err := access.CreateRecruitment(recruitment); err != nil {
			log.Fatal(err, recruitment)
		}
	}

	for _, member := range []*model.RecruitMember{
		{
			RecruitmentUUID: "recruitment-111111111111",
			Grade:           "2",
			Field:           "서버 개발자",
			Number:          "2",
		}, {
			RecruitmentUUID: "recruitment-111111111111",
			Grade:           "2",
			Field:           "웹 프론트 개발자",
			Number:          "2",
		}, {
			RecruitmentUUID: "recruitment-111111111111",
			Grade:           "2",
			Field:           "안드로이드 개발자",
			Number:          "2",
		}, {
			RecruitmentUUID: "recruitment-111111111111",
			Grade:           "2",
			Field:           "iOS 개발자",
			Number:          "2",
		}, {
			RecruitmentUUID: "recruitment-111111111111",
			Grade:           "2",
			Field:           "다 사랑해",
			Number:          "8",
		},
	} {
		if _, err := access.CreateRecruitMember(member); err != nil {
			log.Fatal(err)
		}
	}

	tests := []struct {
		RecruitmentUUID string
		ExpectResults   []*model.RecruitMember
		ExpectError     error
	} {
		{
			RecruitmentUUID: "recruitment-111111111111",
			ExpectResults: []*model.RecruitMember{
				{
					RecruitmentUUID: "recruitment-111111111111",
					Grade:           "2",
					Field:           "서버 개발자",
					Number:          "2",
				}, {
					RecruitmentUUID: "recruitment-111111111111",
					Grade:           "2",
					Field:           "웹 프론트 개발자",
					Number:          "2",
				}, {
					RecruitmentUUID: "recruitment-111111111111",
					Grade:           "2",
					Field:           "안드로이드 개발자",
					Number:          "2",
				}, {
					RecruitmentUUID: "recruitment-111111111111",
					Grade:           "2",
					Field:           "iOS 개발자",
					Number:          "2",
				}, {
					RecruitmentUUID: "recruitment-111111111111",
					Grade:           "2",
					Field:           "다 사랑해",
					Number:          "8",
				},
			},
			ExpectError: nil,
		}, {
			RecruitmentUUID: "recruitment-222222222222",
			ExpectError:     gorm.ErrRecordNotFound,
		}, {
			RecruitmentUUID: "recruitment-333333333333",
			ExpectError:     gorm.ErrRecordNotFound,
		},
	}

	for _, test := range tests {
		resultMembers, err := access.GetRecruitMembersWithRecruitmentUUID(test.RecruitmentUUID)

		var exceptedResult []*model.RecruitMember
		for _, member := range resultMembers {
			exceptedResult = append(exceptedResult, member.ExceptGormModel())
		}

		assert.Equalf(t, test.ExpectError, err, "error assertion error (test case: %v)", test)
		assert.Equalf(t, test.ExpectResults, exceptedResult, "result members assertion error (test case: %v)", test)
	}
}

func Test_Accessor_GetAllClubInforms(t *testing.T) {
	access, err := manager.BeginTx()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		access.Rollback()
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
		}, {
			UUID:       "club-444444444444",
			LeaderUUID: "student-444444444444",
		}, {
			UUID:       "club-555555555555",
			LeaderUUID: "student-555555555555",
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
		}, {
			ClubUUID: "club-444444444444",
			Name:     "MSMS",
			Field:    "SW 개발",
			Location: "2-4반 교실",
			Floor:    "3",
			LogoURI:  "logo.com/club-444444444444",
		},
	} {
		if _, err := access.CreateClubInform(inform); err != nil {
			log.Fatal(err, inform)
		}
	}

	if err, _ := access.DeleteClub("club-444444444444"); err != nil {
		log.Fatal(err)
	}
	if err, _ := access.DeleteClubInform("club-444444444444"); err != nil {
		log.Fatal(err)
	}

	tests := []struct {
		ExpectResults []*model.ClubInform
		ExpectError   error
	} {
		{
			ExpectResults: []*model.ClubInform{
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
			},
			ExpectError: nil,
		},
	}

	for _, test := range tests {
		resultMembers, err := access.GetAllClubInforms()

		var exceptedResult []*model.ClubInform
		for _, member := range resultMembers {
			exceptedResult = append(exceptedResult, member.ExceptGormModel())
		}

		assert.Equalf(t, test.ExpectError, err, "error assertion error (test case: %v)", test)
		assert.Equalf(t, test.ExpectResults, exceptedResult, "result informs assertion error (test case: %v)", test)
	}
}
