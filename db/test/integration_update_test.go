package test

import (
	"club/db/access/errors"
	"club/model"
	"club/tool/mysqlerr"
	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func Test_Accessor_ChangeClubLeader(t *testing.T) {
	access := manager.BeginTx()
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
		},
	} {
		if _, err := access.CreateClub(club); err != nil {
			log.Fatal(err, club)
		}
	}

	tests := []struct {
		ClubUUID          string
		NewLeaderUUID     string
		IsInvalid         bool
		ExpectError       error
		ExpectRowAffected int64
	} {
		{ // success case
			ClubUUID:          "club-111111111111",
			NewLeaderUUID:     "student-333333333333",
			ExpectError:       nil,
			ExpectRowAffected: 1,
		}, { // leader uuid duplicate
			ClubUUID:          "club-111111111111",
			NewLeaderUUID:     "student-222222222222",
			ExpectError:       mysqlerr.DuplicateEntry(model.ClubInstance.LeaderUUID.KeyName(), "student-222222222222"),
			ExpectRowAffected: 0,
		}, { // no exist club uuid
			ClubUUID:          "club-333333333333",
			NewLeaderUUID:     "student-444444444444",
			ExpectRowAffected: 0,
		},  { // invalid leader uuid
			ClubUUID:          "club-444444444444",
			NewLeaderUUID:     "student-4321",
			IsInvalid:         true,
			ExpectRowAffected: 0,
		},
	}

	for _, test := range tests {
		err, rowAffected := access.ChangeClubLeader(test.ClubUUID, test.NewLeaderUUID)

		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			err = mysqlerr.ExceptReferenceInformFrom(mysqlErr)
		}

		if test.IsInvalid {
			_, isInvalid := err.(validator.ValidationErrors)
			assert.Equalf(t, test.IsInvalid, isInvalid, "invalid state assertion error (test case: %v)", test)
		} else {
			assert.Equalf(t, test.ExpectError, err, "error assertion error (test case: %v)", test)
		}
		assert.Equalf(t, test.ExpectRowAffected, rowAffected, "row affected assertion error (test case: %v)", test)
	}
}

func Test_Accessor_ModifyClubInform(t *testing.T) {
	access := manager.BeginTx()
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
		},
	} {
		if _, err := access.CreateClubInform(inform); err != nil {
			log.Fatal(err, inform)
		}
	}

	tests := []struct {
		ClubUUID       string
		RevisionInform *model.ClubInform
		IsInvalid      bool
		ExpectError    error
		ExpectRows     int64
	} {
		{ // success case
			ClubUUID: "club-111111111111",
			RevisionInform: &model.ClubInform{
				ClubConcept:  "DMS 개발 및 운영",
				Introduction: "우리 DMS 동아리 완전 좋아요~",
				Link:         "facebook.com/DSM-DMS",
			},
			ExpectError: nil,
			ExpectRows:  1,
		}, { // name duplicate error
			ClubUUID: "club-111111111111",
			RevisionInform: &model.ClubInform{
				Name: "SMS",
			},
			ExpectError: mysqlerr.DuplicateEntry(model.ClubInformInstance.Name.KeyName(), "SMS"),
			ExpectRows:  0,
		}, { // name duplicate error
			ClubUUID: "club-111111111111",
			RevisionInform: &model.ClubInform{
				Location: "2-2반 교실",
			},
			ExpectError: mysqlerr.DuplicateEntry(model.ClubInformInstance.Location.KeyName(), "2-2반 교실"),
			ExpectRows:  0,
		}, { // floor invalid
			ClubUUID: "club-111111111111",
			RevisionInform: &model.ClubInform{
				Floor: "7",
			},
			IsInvalid:  true,
			ExpectRows: 0,
		}, { // name invalid
			ClubUUID: "club-111111111111",
			RevisionInform: &model.ClubInform{
				Name: "이거 30글자 넘음 30글자 넘으면 유효성 검사 부분에서 오류가 나야만함 그래야만함 나겠죠?",
			},
			IsInvalid:  true,
			ExpectRows: 0,
		}, {
			ClubUUID: "club-111111111111",
			RevisionInform: &model.ClubInform{
				ClubUUID: "club-123412341234",
			},
			ExpectError: errors.ClubUUIDCannotBeChanged,
			ExpectRows:  0,
		},
	}

	for _, test := range tests {
		err, rowAffected := access.ModifyClubInform(test.ClubUUID, test.RevisionInform)

		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			err = mysqlerr.ExceptReferenceInformFrom(mysqlErr)
		}

		if test.IsInvalid {
			_, isInvalid := err.(validator.ValidationErrors)
			assert.Equalf(t, test.IsInvalid, isInvalid, "invalid state assertion error (test case: %v)", test)
		} else {
			assert.Equalf(t, test.ExpectError, err, "error assertion error (test case: %v)", test)
		}
		assert.Equalf(t, test.ExpectRows, rowAffected, "row affected assertion error (test case: %v)", test)
	}

	confirmTests := []struct {
		ClubUUID   string
		ExpectResult *model.ClubInform
		ExpectError  error
	}{
		{
			ClubUUID: "club-111111111111",
			ExpectResult: &model.ClubInform{
				ClubUUID:     "club-111111111111",
				Name:         "DMS",
				Field:        "SW 개발",
				Location:     "2-1반 교실",
				Floor:        "3",
				LogoURI:      "logo.com/club-111111111111",
				ClubConcept:  "DMS 개발 및 운영",
				Introduction: "우리 DMS 동아리 완전 좋아요~",
				Link:         "facebook.com/DSM-DMS",
			},
			ExpectError: nil,
		},
	}

	for _, test := range confirmTests {
		result, err := access.GetClubInformWithClubUUID(test.ClubUUID)

		assert.Equalf(t, test.ExpectError, err, "error assertion error (test case: %v)", test)
		assert.Equalf(t, test.ExpectResult, result.ExceptGormModel(), "result club inform assertion error (test case: %v)", test)
	}
}

func Test_Accessor_ModifyRecruitment(t *testing.T) {
	access := manager.BeginTx()
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

	now := time.Now()
	startTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	endTime := startTime.Add(time.Hour * 24 * 7)

	for _, recruitment := range []*model.ClubRecruitment{
		{
			UUID:           "recruitment-111111111111",
			ClubUUID:       "club-111111111111",
			RecruitConcept: "기간제 공채",
			StartPeriod:    model.StartPeriod(startTime),
			EndPeriod:      model.EndPeriod(endTime),
		},
	} {
		if _, err := access.CreateRecruitment(recruitment); err != nil {
			log.Fatal(err, recruitment)
		}
	}

	tests := []struct {
		RecruitmentUUID string
		RevisionRecruit *model.ClubRecruitment
		IsInvalid       bool
		ExpectError     error
		ExpectRows      int64
	} {
		{ // success case
			RecruitmentUUID: "recruitment-111111111111",
			RevisionRecruit: &model.ClubRecruitment{
				RecruitConcept: "응 상시 채용으로 바꿔",
				StartPeriod:    model.StartPeriod(model.ClubRecruitmentInstance.StartPeriod.NullReplaceValue()),
				EndPeriod:      model.EndPeriod(model.ClubRecruitmentInstance.EndPeriod.NullReplaceValue()),
			},
			ExpectError: nil,
			ExpectRows:  1,
		}, { // floor invalid
			RecruitmentUUID: "recruitment-111111111111",
			RevisionRecruit: &model.ClubRecruitment{
				RecruitConcept: "이것도 40자가 넘도록 만들어야되는데 곧 있으면 40자가 될 것 같다 만약 이게 40자가 된다면 유효성 오류가 발생해야한다",
			},
			IsInvalid:  true,
			ExpectRows: 0,
		}, {
			RecruitmentUUID: "recruitment-111111111111",
			RevisionRecruit: &model.ClubRecruitment{
				ClubUUID: "club-123412341234",
			},
			ExpectError: errors.ClubUUIDCannotBeChanged,
			ExpectRows:  0,
		}, {
			RecruitmentUUID: "recruitment-111111111111",
			RevisionRecruit: &model.ClubRecruitment{
				UUID: "recruitment-123412341234",
			},
			ExpectError: errors.RecruitmentUUIDCannotBeChanged,
			ExpectRows:  0,
		},
	}

	for _, test := range tests {
		err, rowAffected := access.ModifyRecruitment(test.RecruitmentUUID, test.RevisionRecruit)

		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			err = mysqlerr.ExceptReferenceInformFrom(mysqlErr)
		}

		if test.IsInvalid {
			_, isInvalid := err.(validator.ValidationErrors)
			assert.Equalf(t, test.IsInvalid, isInvalid, "invalid state assertion error (test case: %v)", test)
		} else {
			assert.Equalf(t, test.ExpectError, err, "error assertion error (test case: %v)", test)
		}
		assert.Equalf(t, test.ExpectRows, rowAffected, "row affected assertion error (test case: %v)", test)
	}

	confirmTests := []struct {
		RecruitmentUUID string
		ExpectResult    *model.ClubRecruitment
		ExpectError     error
	} {
		{
			RecruitmentUUID: "recruitment-111111111111",
			ExpectResult: &model.ClubRecruitment{
				UUID:           "recruitment-111111111111",
				ClubUUID:       "club-111111111111",
				RecruitConcept: "응 상시 채용으로 바꿔",
			},
			ExpectError: nil,
		},
	}

	for _, test := range confirmTests {
		result, err := access.GetRecruitmentWithRecruitmentUUID(test.RecruitmentUUID)

		assert.Equalf(t, test.ExpectError, err, "error assertion error (test case: %v)", test)
		assert.Equalf(t, test.ExpectResult, result.ExceptGormModel(), "result recruitment assertion error (test case: %v)", test)
	}
}
