package access

import (
	"club/model"
	"gorm.io/gorm"
	"time"
)

func (d *_default) GetClubWithClubUUID(clubUUID string) (club *model.Club, err error) {
	club = new(model.Club)
	selectResult := d.tx.Where("uuid = ?", clubUUID).Find(club)
	err = selectResult.Error
	if selectResult.RowsAffected == 0 && err == nil {
		err = gorm.ErrRecordNotFound
	}
	return
}

func (d *_default) GetClubWithLeaderUUID(leaderUUID string) (club *model.Club, err error) {
	club = new(model.Club)
	selectResult := d.tx.Where("leader_uuid = ?", leaderUUID).Find(club)
	err = selectResult.Error
	if selectResult.RowsAffected == 0 && err == nil {
		err = gorm.ErrRecordNotFound
	}
	return
}

func (d *_default) GetCurrentRecruitmentWithClubUUID(clubUUID string) (recruit *model.ClubRecruitment, err error) {
	recruit = new(model.ClubRecruitment)
	now := time.Now()

	fromSubQuery := d.tx.Table(model.ClubRecruitmentInstance.TableName()).Where("club_uuid = ?", clubUUID)
	selectedTx := d.tx.Table("(?) as club_recruitments", fromSubQuery)
	selectResult := selectedTx.Where("club_recruitments.end_period >= ?", now).Or("club_recruitments.end_period IS NULL").Find(recruit)

	err = selectResult.Error
	if selectResult.RowsAffected == 0 && err == nil {
		err = gorm.ErrRecordNotFound
	}

	return
}

func (d *_default) GetClubInformsSortByUpdateTime(offset, limit int, field, name string) (clubInforms []*model.ClubInform, err error) {
	selectedTx := d.tx.Table(model.ClubInformInstance.TableName())
	if field != "" {
		selectedTx = selectedTx.Where("field LIKE ?", "%"+field+"%")
	}
	if name != "" {
		selectedTx = selectedTx.Where("name LIKE ?", "%"+name+"%")
	}

	clubInforms = make([]*model.ClubInform, limit)
	err = selectedTx.Order("updated_at desc").Limit(limit).Offset(offset).Find(&clubInforms).Error

	if len(clubInforms) == 0 && err == nil {
		err = gorm.ErrRecordNotFound
	}

	return
}

func (d *_default) GetCurrentRecruitmentsSortByCreateTime(offset, limit int, field, name string) (recruits []*model.ClubRecruitment, err error) {
	selectedTX := d.tx.New()
	if field != "" {
		selectedTX = selectedTX.Where("field = ?", field)
	}
	if name != "" {
		selectedTX = selectedTX.Where("name = ?", name)
	}

	recruits = make([]*model.ClubRecruitment, limit)
	selectedTX = selectedTX.Where("end_period >= ?", time.Now()).Or("end_period IS NULL")
	err = selectedTX.Order("created_at desc").Limit(limit).Offset(offset).Find(&recruits).Error

	if len(recruits) == 0 && err == nil {
		err = gorm.ErrRecordNotFound
	}

	return
}

func (d *_default) GetClubInformWithClubUUID(clubUUID string) (inform *model.ClubInform, err error) {
	inform = new(model.ClubInform)
	selectResult := d.tx.Where("club_uuid = ?", clubUUID).Find(inform)
	err = selectResult.Error
	if selectResult.RowsAffected == 0 && err == nil {
		err = gorm.ErrRecordNotFound
	}
	return
}

func (d *_default) GetRecruitmentWithRecruitmentUUID(recruitUUID string) (recruit *model.ClubRecruitment, err error) {
	recruit = new(model.ClubRecruitment)
	selectResult := d.tx.Where("uuid = ?", recruitUUID).Find(recruit)
	err = selectResult.Error
	if selectResult.RowsAffected == 0 && err == nil {
		err = gorm.ErrRecordNotFound
	}

	return
}

func (d *_default) GetClubMembersWithClubUUID(clubUUID string) (members []*model.ClubMember, err error) {
	members = make([]*model.ClubMember, 5, 5)
	err = d.tx.Where("club_uuid = ?", clubUUID).Find(&members).Error

	if len(members) == 0 && err == nil {
		err = gorm.ErrRecordNotFound
	}

	return
}

func (d *_default) GetRecruitMembersWithRecruitmentUUID(recruitUUID string) (members []*model.RecruitMember, err error) {
	members = make([]*model.RecruitMember, 5, 5)
	err = d.tx.Where("recruitment_uuid = ?", recruitUUID).Find(&members).Error

	if len(members) == 0 && err == nil {
		err = gorm.ErrRecordNotFound
	}

	return
}

func (d *_default) GetAllClubInforms() (informs []*model.ClubInform, err error) {
	informs = make([]*model.ClubInform, 10, 10)
	err = d.tx.Find(&informs).Error

	if len(informs) == 0 && err == nil {
		err = gorm.ErrRecordNotFound
	}

	return
}

func (d *_default) GetAllCurrentRecruitments() (recruitments []*model.ClubRecruitment, err error) {
	recruitments = make([]*model.ClubRecruitment, 10, 10)
	err = d.tx.Where("end_period >= ?", time.Now()).Or("end_period IS NULL").Find(&recruitments).Error

	if len(recruitments) == 0 && err == nil {
		err = gorm.ErrRecordNotFound
	}

	return
}
