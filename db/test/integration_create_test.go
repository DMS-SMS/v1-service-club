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

func Test_Accessor_CreateClubInform(t *testing.T) {
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
			UUID:       "club-123412341234",
			LeaderUUID: "student-123412341234",
		}, {
			UUID:       "club-432143214321",
			LeaderUUID: "student-432143214321",
		}, {
			UUID:       "club-111111111111",
			LeaderUUID: "student-111111111111",
		},
	} {
		if _, err := access.CreateClub(club); err != nil {
			log.Fatal(err, club)
		}
	}

	tests := []struct {
		ClubUUID, Name  string
		ClubConcept     string
		Introduction    string
		Field, Location string
		Floor           int64
		Link, LogoURI   string
		IsInvalid       bool
		ExpectedError   error
	} {
		{ // success case 1
			ClubUUID:      "club-123412341234",
			Name:          "DMS",
			ClubConcept:   "DMS, SMS 서비스 개발 및 운영",
			Introduction:  "어서 오세용~ 박진홍이 서식중이에용~",
			Field:         "SW 개발",
			Location:      "2-1반 교실",
			Floor:         3,
			Link:          "link.com",
			LogoURI:       "logo.com",
			ExpectedError: nil,
		}, { // success case 2
			ClubUUID:      "club-432143214321",
			Name:          "SMS",
			Field:         "SW 개발",
			Location:      "2-2반 교실",
			Floor:         3,
			LogoURI:       "logo.com",
			ExpectedError: nil,
		}, { // floor invalid
			ClubUUID:  "club-111111111111",
			Name:      "DSM",
			Field:     "SW 개발",
			Location:  "2-3반 교실",
			Floor:     6, // invalid floor
			LogoURI:   "logo.com",
			IsInvalid: true,
		}, { // club uuid duplicate error
			ClubUUID:      "club-123412341234",
			Name:          "DSM",
			Field:         "SW 개발",
			Location:      "2-3반 교실",
			Floor:         3,
			LogoURI:       "logo.com",
			ExpectedError: mysqlerr.DuplicateEntry(model.ClubInformInstance.ClubUUID.KeyName(), "club-123412341234"),
		}, { // name uuid duplicate error
			ClubUUID:      "club-111111111111",
			Name:          "DMS",
			Field:         "SW 개발",
			Location:      "2-3반 교실",
			Floor:         3,
			LogoURI:       "logo.com",
			ExpectedError: mysqlerr.DuplicateEntry(model.ClubInformInstance.Name.KeyName(), "DMS"),
		}, { // location uuid duplicate error
			ClubUUID:      "club-111111111111",
			Name:          "DSM",
			Field:         "SW 개발",
			Location:      "2-2반 교실",
			Floor:         3,
			LogoURI:       "logo.com",
			ExpectedError: mysqlerr.DuplicateEntry(model.ClubInformInstance.Location.KeyName(), "2-2반 교실"),
		}, { // not exist club uuid
			ClubUUID:      "club-222222222222",
			Name:          "DSM",
			Field:         "SW 개발",
			Location:      "2-3반 교실",
			Floor:         3,
			LogoURI:       "logo.com",
			ExpectedError: clubInformClubUUIDFKConstraintFailError,
		},
	}

	for _, test := range tests {
		_, err := access.CreateClubInform(&model.ClubInform{
			ClubUUID:     model.ClubUUID(test.ClubUUID),
			Name:         model.Name(test.Name),
			ClubConcept:  model.ClubConcept(test.ClubConcept),
			Introduction: model.Introduction(test.Introduction),
			Field:        model.Field(test.Field),
			Location:     model.Location(test.Location),
			Floor:        model.Floor(test.Floor),
			Link:         model.Link(test.Link),
			LogoURI:      model.LogoURI(test.LogoURI),
		})

		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			err = mysqlerr.ExceptReferenceInformFrom(mysqlErr)
		}

		if test.IsInvalid {
			_, isInvalid := err.(validator.ValidationErrors)
			assert.Equalf(t, test.IsInvalid, isInvalid, "invalid state assertion error (test case: %v)", test)
		} else {
			assert.Equalf(t, test.ExpectedError, err, "error assertion error (test case: %v)", test)
		}
	}
}

func Test_Accessor_CreateClubMember(t *testing.T) {
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
			UUID:       "club-123412341234",
			LeaderUUID: "student-123412341234",
		}, {
			UUID:       "club-432143214321",
			LeaderUUID: "student-432143214321",
		},
	} {
		if _, err := access.CreateClub(club); err != nil {
			log.Fatal(err, club)
		}
	}

	tests := []struct {
		ClubUUID      string
		StudentUUID   string
		IsInvalid     bool
		ExpectedError error
	} {
		{ // success case
			ClubUUID:      "club-123412341234",
			StudentUUID:   "student-123412341234",
			ExpectedError: nil,
		}, { // success case
			ClubUUID:      "club-123412341234",
			StudentUUID:   "student-111111111111",
			ExpectedError: nil,
		}, { // student uuid duplicate -> error X
			ClubUUID:      "club-432143214321",
			StudentUUID:   "student-111111111111",
			ExpectedError: nil,
		}, { // club & student uuid duplicate error
			ClubUUID:      "club-123412341234",
			StudentUUID:   "student-111111111111",
			ExpectedError: mysqlerr.DuplicateEntry(model.ClubMemberInstance.StudentUUID.KeyName(), "club-123412341234.student-111111111111"),
		}, { // validate error
			ClubUUID:    "club-12341234123",
			StudentUUID: "student-111111111111",
			IsInvalid:   true,
		}, { // validate error
			ClubUUID:    "club-123412341234",
			StudentUUID: "student-11111111111",
			IsInvalid:   true,
		}, { // no exist club uuid error
			ClubUUID:      "club-111111111111",
			StudentUUID:   "student-111111111111",
			ExpectedError: clubMemberClubUUIDFKConstraintFailError,
		},
	}

	for _, test := range tests {
		_, err := access.CreateClubMember(&model.ClubMember{
			ClubUUID:    model.ClubUUID(test.ClubUUID),
			StudentUUID: model.StudentUUID(test.StudentUUID),
		})

		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			err = mysqlerr.ExceptReferenceInformFrom(mysqlErr)
		}

		if test.IsInvalid {
			_, isInvalid := err.(validator.ValidationErrors)
			assert.Equalf(t, test.IsInvalid, isInvalid, "invalid state assertion error (test case: %v)", test)
		} else {
			assert.Equalf(t, test.ExpectedError, err, "error assertion error (test case: %v)", test)
		}
	}
}
