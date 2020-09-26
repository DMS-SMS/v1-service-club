package access

import (
	"club/model"
	"gorm.io/gorm"
)

type None struct {}

func (n None) CreateClub(club *model.Club) (_ *model.Club, _ error) { return }
func (n None) CreateClubInform(inform *model.ClubInform) (_ *model.ClubInform, _ error) { return }
func (n None) CreateClubMember(clubMember *model.ClubMember) (_ *model.ClubMember, _ error) { return }
func (n None) CreateRecruitment(recruit *model.ClubRecruitment) (_ *model.ClubRecruitment, _ error) { return }
func (n None) CreateRecruitMember(recruitMember *model.RecruitMember) (_ *model.RecruitMember, _ error) { return }

func (n None) GetClubWithClubUUID(clubUUID string) (_ *model.Club, _ error) { return }
func (n None) GetClubWithLeaderUUID(leaderUUID string) (_ *model.Club, _ error) { return }
func (n None) GetCurrentRecruitmentWithClubUUID(clubUUID string) (_ *model.ClubRecruitment, _ error) { return }
func (n None) GetClubInformsSortByUpdateTime(offset, limit int, field, name string) (_ []*model.ClubInform, _ error) { return }
func (n None) GetCurrentRecruitmentsSortByCreateTime(offset, limit int, field, name string) (_ []*model.ClubRecruitment, _ error) { return }
func (n None) GetClubInformWithClubUUID(clubUUID string) (_ *model.ClubInform, _ error) { return }
func (n None) GetRecruitmentWithRecruitmentUUID(recruitUUID string) (_ *model.ClubRecruitment, _ error) { return }
func (n None) GetClubMembersWithClubUUID(clubUUID string) (_ []*model.ClubMember, _ error) { return }
func (n None) GetRecruitMembersWithRecruitmentUUID(recruitUUID string) (_ []*model.RecruitMember, _ error) { return }
func (n None) GetAllClubInforms() (_ []*model.ClubInform, _ error) { return }
func (n None) GetAllCurrentRecruitments() (_ []*model.ClubRecruitment, _ error) { return }

func (n None) ChangeClubLeader(clubUUID, newLeaderUUID string) (_ error, _ int64) { return }
func (n None) ModifyClubInform(clubUUID string, revisionInform *model.ClubInform) (_ error, _ int64) { return }
func (n None) ModifyRecruitment(recruitUUID string, revisionRecruit *model.ClubRecruitment) (_ error, _ int64) { return }

func (n None) DeleteClub(clubUUID string) (_ error, _ int64) { return }
func (n None) DeleteClubInform(clubUUID string) (_ error, _ int64) { return }
func (n None) DeleteClubMember(clubUUID, studentUUID string) (_ error, _ int64) { return }
func (n None) DeleteRecruitment(recruitUUID string) (_ error, _ int64) { return }
func (n None) DeleteAllRecruitMember(recruitUUID string) (_ error, _ int64) { return }

func (n None) BeginTx() { return }
func (n None) Commit() (_ *gorm.DB) { return }
func (n None) Rollback() (_ *gorm.DB) { return }
