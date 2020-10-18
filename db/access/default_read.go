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

func (d *_default) GetCurrentRecruitmentWithRecruitmentUUID(recruitmentUUID string) (recruit *model.ClubRecruitment, err error) {
	recruit = new(model.ClubRecruitment)
	now := time.Now()

	fromSubQuery := d.tx.Table(model.ClubRecruitmentInstance.TableName()).Where("uuid = ?", recruitmentUUID)
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
	fromSubQuery := d.tx.Table(model.ClubRecruitmentInstance.TableName()).Select("club_recruitments.*")
	fromSubQuery = fromSubQuery.Joins("JOIN club_informs ON club_informs.club_uuid = club_recruitments.club_uuid")
	fromSubQuery = fromSubQuery.Where("club_informs.deleted_at IS NULL")

	if field != "" {
		fromSubQuery = fromSubQuery.Where("club_informs.field LIKE ?", "%"+field+"%")
	}
	if name != "" {
		fromSubQuery = fromSubQuery.Where("club_informs.name LIKE ?", "%"+name+"%")
	}

	recruits = make([]*model.ClubRecruitment, limit)
	selectedTX := d.tx.Table("(?) as club_recruitments", fromSubQuery)
	selectedTX = selectedTX.Where("club_recruitments.end_period >= ?", time.Now()).Or("club_recruitments.end_period IS NULL")
	err = selectedTX.Order("club_recruitments.created_at desc").Limit(limit).Offset(offset).Find(&recruits).Error

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

func (d *_default) GetClubMembersWithClubUUID(clubUUID string) ([]*model.ClubMember, error) {
	var members []*model.ClubMember
	err := d.tx.Where("club_uuid = ?", clubUUID).Find(&members).Error

	if len(members) == 0 && err == nil {
		err = gorm.ErrRecordNotFound
	}

	return members, err
}

func (d *_default) GetRecruitMembersWithRecruitmentUUID(recruitUUID string) ([]*model.RecruitMember, error) {
	var members []*model.RecruitMember
	err := d.tx.Where("recruitment_uuid = ?", recruitUUID).Find(&members).Error

	if len(members) == 0 && err == nil {
		err = gorm.ErrRecordNotFound
	}

	return members, err
}

func (d *_default) GetAllClubInforms() ([]*model.ClubInform, error) {
	joinedTx := d.tx.Table(model.ClubInformInstance.TableName()).Joins("JOIN clubs ON clubs.uuid = club_informs.club_uuid")
	joinedTx = joinedTx.Where("clubs.deleted_at IS NULL")

	var informs []*model.ClubInform
	err := joinedTx.Find(&informs).Error

	if len(informs) == 0 && err == nil {
		err = gorm.ErrRecordNotFound
	}

	return informs, err
}

func (d *_default) GetAllCurrentRecruitments() ([]*model.ClubRecruitment, error) {
	fromSubQuery := d.tx.Table(model.ClubRecruitmentInstance.TableName()).Select("club_recruitments.*")
	fromSubQuery = fromSubQuery.Joins("JOIN clubs ON clubs.uuid = club_recruitments.club_uuid").Where("clubs.deleted_at IS NULL")

	var recruitments []*model.ClubRecruitment
	selectedTx := d.tx.Table("(?) AS club_recruitments", fromSubQuery)
	err := selectedTx.Where("club_recruitments.end_period >= ?", time.Now()).Or("club_recruitments.end_period IS NULL").Find(&recruitments).Error

	if len(recruitments) == 0 && err == nil {
		err = gorm.ErrRecordNotFound
	}

	return recruitments, err
}
