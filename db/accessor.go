package db

import (
	"club/model"
	"github.com/jinzhu/gorm"
)

type Accessor interface {
	CreateClub(club *model.Club) (resultClub *model.Club, err error)
	CreateClubInform(inform *model.ClubInform) (resultInform *model.ClubInform, err error)
	CreateClubMember(clubMember *model.ClubMember) (resultMember *model.ClubMember, err error)
	CreateRecruitment(recruit *model.ClubRecruitment) (resultRecruit *model.ClubRecruitment, err error)
	CreateRecruitMember(recruitMember *model.RecruitMember) (resultMember *model.RecruitMember, err error)

	GetClubWithClubUUID(clubUUID string) (*model.Club, error)
	GetClubWithLeaderUUID(leaderUUID string) (*model.Club, error)
	GetRecruitmentWithClubUUID(clubUUID string) (*model.ClubRecruitment, error)
	GetClubInformsSortByUpdateTime(offset, limit int, field, name string) ([]*model.ClubInform, error)
	GetRecruitmentsSortByCreateTime(offset, limit int, field, name string) ([]*model.ClubRecruitment, error)
	GetClubInformWithClubUUID(clubUUID string) (*model.ClubInform, error)
	GetRecruitmentWithRecruitmentUUID(recruitUUID string) (*model.ClubRecruitment, error)
	GetClubMembersWithClubUUID(clubUUID string) ([]*model.ClubMember, error)
	GetRecruitMembersWithRecruitmentUUID(recruitUUID string) ([]*model.RecruitMember, error)
	GetAllClubInforms() ([]*model.ClubInform, error)
	GetAllRecruitments() ([]*model.ClubRecruitment, error)

	ChangeClubLeader(clubUUID, newLeaderUUID string) error
	ModifyClubInform(clubUUID string, revisionInform *model.ClubInform) error
	ModifyRecruitment(recruitUUID string, revisionRecruit *model.ClubRecruitment) error

	DeleteClub(clubUUID string) error
	DeleteClubInform(clubUUID string) error
	DeleteClubMember(clubUUID, studentUUID string) error
	DeleteRecruitment(recruitUUID string) error
	DeleteAllRecruitMember(recruitUUID string) error

	BeginTx()
	Commit() *gorm.DB
	Rollback() *gorm.DB
}
