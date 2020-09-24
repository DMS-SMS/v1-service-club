package access

import (
	"club/model"
)

func (d *_default) GetClubWithClubUUID(clubUUID string) (club *model.Club, err error) {
	club = new(model.Club)
	err = d.tx.Where("uuid = ?", clubUUID).Find(club).Error
	return
}

func (d *_default) GetClubWithLeaderUUID(leaderUUID string) (club *model.Club, err error) {
	club = new(model.Club)
	err = d.tx.Where("leader_uuid = ?", leaderUUID).Find(club).Error
	return
}

func (d *_default) GetRecruitmentWithClubUUID(clubUUID string) (recruit *model.ClubRecruitment, err error) {
	recruit = new(model.ClubRecruitment)
	err = d.tx.Where("club_uuid = ?", clubUUID).Find(recruit).Error
	return
}

func (d *_default) GetClubInformWithClubUUID(clubUUID string) (inform *model.ClubInform, err error) {
	inform = new(model.ClubInform)
	err = d.tx.Where("club_uuid = ?", clubUUID).Find(inform).Error
	return
}

func (d *_default) GetRecruitmentWithRecruitmentUUID(recruitUUID string) (recruit *model.ClubRecruitment, err error) {
	recruit = new(model.ClubRecruitment)
	err = d.tx.Where("uuid = ?", recruitUUID).Find(recruit).Error
	return
}

