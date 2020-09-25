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
