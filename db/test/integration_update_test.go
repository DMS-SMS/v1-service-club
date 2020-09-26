package test

import (
	"club/model"
	"club/tool/mysqlerr"
	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func Test_Accessor_ChangeClubLeader(t *testing.T) {
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
