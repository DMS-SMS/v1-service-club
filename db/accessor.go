package db

import (
	"club/model"
)

type Accessor interface {
	CreateClub(club *model.Club) (resultClub *model.Club, err error)
	CreateClubInform(inform *model.ClubInform) (resultInform *model.ClubInform, err error)
	CreateClubMember(clubMember *model.ClubMember) (resultMember *model.ClubMember, err error)
	CreateRecruitment(recruit *model.ClubRecruitment) (resultRecruit *model.ClubRecruitment, err error)
	CreateRecruitMember(recruitMember *model.RecruitMember) (resultMember *model.RecruitMember, err error)
}
