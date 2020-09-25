package test

import (
	"club/model"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
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
