package test

import (
	"club/model"
	"log"
	"testing"
)

func Test_Accessor_DeleteXXX(t *testing.T) {
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
		},
	} {
		if _, err := access.CreateClubMember(member); err != nil {
			log.Fatal(err)
		}
	}

	for _, recruitment := range []*model.ClubRecruitment{
		{
			UUID:           "recruitment-111111111111",
			ClubUUID:       "club-111111111111",
			RecruitConcept: "첫 번째 공채",
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

	_, _ = access.DeleteClub("club-111111111111") // nil, 1
	_, _ = access.DeleteClubInform("club-111111111111") // nil, 1
	_, _ = access.DeleteClubMember("club-111111111111", "student-111111111111") // nil, 1
	_, _ = access.DeleteRecruitment("recruitment-111111111111") // nil, 1
	_, _ = access.DeleteAllRecruitMember("recruitment-111111111111") // nil, 5

	_, _ = access.DeleteClub("club-222222222222") // nil, 0
	_, _ = access.DeleteClubInform("club-222222222222") // nil, 0
	_, _ = access.DeleteClubMember("club-222222222222", "student-222222222222") // nil, 0
	_, _ = access.DeleteRecruitment("recruitment-222222222222") // nil, 0
	_, _ = access.DeleteAllRecruitMember("recruitment-222222222222") // nil, 0
}