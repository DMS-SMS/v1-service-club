package test

import (
	"club/model"
	"club/tool/mysqlerr"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func Test_Accessor_CreateClub(t *testing.T) {
	access, err := manager.BeginTx()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		access.Rollback()
		testGroup.Done()
	}()

	tests := []struct {
		UUID, LeaderUUID string
		IsInvalid        bool
		ExpectedError    error
	} {
		{ // success case
			UUID:          "club-123412341234",
			LeaderUUID:    "student-123412341234",
			ExpectedError: nil,
		}, { // uuid invalid case
			UUID:       "club-12341234123",
			LeaderUUID: "student-123412341234",
			IsInvalid:  true,
		}, { // leader uuid invalid case
			UUID:       "club-123412341234",
			LeaderUUID: "student-12341234123",
			IsInvalid:  true,
		}, { // uuid duplicate error
			UUID:       "club-123412341234",
			LeaderUUID: "student-111111111111",
			ExpectedError: mysqlerr.DuplicateEntry(model.ClubInstance.UUID.KeyName(), "club-123412341234"),
		}, { // leader uuid duplicate error
			UUID:       "club-111111111111",
			LeaderUUID: "student-123412341234",
			ExpectedError: mysqlerr.DuplicateEntry(model.ClubInstance.LeaderUUID.KeyName(), "student-123412341234"),
		},
	}

	for _, testCase := range tests {
		_, err := access.CreateClub(&model.Club{
			UUID:       model.UUID(testCase.UUID),
			LeaderUUID: model.LeaderUUID(testCase.LeaderUUID),
		})

		if testCase.IsInvalid {
			_, isInvalid := err.(validator.ValidationErrors)
			assert.Equalf(t, testCase.IsInvalid, isInvalid, "invalid state assertion error (test case: %v)", testCase)
		} else {
			assert.Equalf(t, testCase.ExpectedError, err, "error assertion error (test case: %v)", testCase)
		}
	}
}