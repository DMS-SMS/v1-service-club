package model

import (
	"club/model/validate"
	"club/tool/mysqlerr"
	"fmt"
	"gorm.io/gorm"
)

const (
	emptyString = ""
	emptyInt = 0

	validClubUUID = "club-111111111111"
	validRecruitmentUUID = "recruitment-111111111111"
	validLeaderUUID = "student-111111111111"
	validName = "DMS"
	validField = "SW 개발"
	validLocation = "2-2반 교실"
	validFloor = "3"
	validLogoURI = "logos/club-111111111111"
	validRecruitConcept = "디자인에 좋은 감각이 있는 새로운 1학년 부원을 모집합니다!"
)

func (c *Club) BeforeCreate(tx *gorm.DB) (err error) {
	if err = validate.DBValidator.Struct(c); err != nil {
		return
	}

	selectResult := tx.Where("leader_uuid = ?", c.LeaderUUID).Find(&Club{})
	if selectResult.RowsAffected != 0 {
		err = mysqlerr.DuplicateEntry(ClubInstance.LeaderUUID.KeyName(), string(c.LeaderUUID))
	}
	return
}

func (ci *ClubInform) BeforeCreate(tx *gorm.DB) (err error) {
	if err = validate.DBValidator.Struct(ci); err != nil {
		return
	}

	if tx.Where("name = ?", ci.Name).Find(&ClubInform{}).RowsAffected != 0 {
		err = mysqlerr.DuplicateEntry(ClubInformInstance.Name.KeyName(), string(ci.Name))
		return
	}

	if tx.Where("location = ?", ci.Location).Find(&ClubInform{}).RowsAffected != 0 {
		err = mysqlerr.DuplicateEntry(ClubInformInstance.Location.KeyName(), string(ci.Location))
		return
	}

	return
}

func (cm *ClubMember) BeforeCreate(tx *gorm.DB) (err error) {
	if err = validate.DBValidator.Struct(cm); err != nil {
		return
	}

	selectResult := tx.Where("club_uuid = ? AND student_uuid = ?", cm.ClubUUID, cm.StudentUUID).Find(&ClubMember{})
	if selectResult.RowsAffected != 0 {
		err = mysqlerr.DuplicateEntry(ClubMemberInstance.StudentUUID.KeyName(), fmt.Sprintf("%s.%s", cm.ClubUUID, cm.StudentUUID))
	}
	return
}

func (cr *ClubRecruitment) BeforeCreate(tx *gorm.DB) error {
	return validate.DBValidator.Struct(cr)
}

func (rm *RecruitMember) BeforeCreate(tx *gorm.DB) error {
	return validate.DBValidator.Struct(rm)
}

func (c *Club) BeforeUpdate(tx *gorm.DB) error {
	clubForValidate := c.DeepCopy()

	if clubForValidate.UUID == emptyString       { clubForValidate.UUID = validClubUUID }
	if clubForValidate.LeaderUUID == emptyString { clubForValidate.LeaderUUID = validLeaderUUID }

	return validate.DBValidator.Struct(clubForValidate)
}

func (ci *ClubInform) BeforeUpdate(tx *gorm.DB) error {
	clubInformForValidate := ci.DeepCopy()

	if clubInformForValidate.ClubUUID == emptyString { clubInformForValidate.ClubUUID = validClubUUID }
	if clubInformForValidate.Name == emptyString     { clubInformForValidate.Name = validName }
	if clubInformForValidate.Field == emptyString    { clubInformForValidate.Field = validField }
	if clubInformForValidate.Location == emptyString { clubInformForValidate.Location = validLocation }
	if clubInformForValidate.Floor == emptyString    { clubInformForValidate.Floor = validFloor }
	if clubInformForValidate.LogoURI == emptyString  { clubInformForValidate.LogoURI = validLogoURI }

	return validate.DBValidator.Struct(clubInformForValidate)
}

func (cr *ClubRecruitment) BeforeUpdate(tx *gorm.DB) error {
	recruitmentForValidate := cr.DeepCopy()

	if recruitmentForValidate.UUID == emptyString           { recruitmentForValidate.UUID = validRecruitmentUUID }
	if recruitmentForValidate.ClubUUID == emptyString       { recruitmentForValidate.ClubUUID = validClubUUID }
	if recruitmentForValidate.RecruitConcept == emptyString { recruitmentForValidate.RecruitConcept = validRecruitConcept }

	return validate.DBValidator.Struct(recruitmentForValidate)
}
