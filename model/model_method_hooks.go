package model

import "club/model/validate"

const (
	emptyString = ""
	emptyInt = 0

	validClubUUID = "club-111111111111"
	validRecruitmentUUID = "recruitment-111111111111"
	validLeaderUUID = "student-111111111111"
	validName = "DMS"
	validField = "SW 개발"
	validLocation = "2-2반 교실"
	validFloor = 3
	validLogoURI = "logos/club-111111111111"
	validRecruitConcept = "디자인에 좋은 감각이 있는 새로운 1학년 부원을 모집합니다!"
)

func (c *Club) BeforeCreate() error {
	return validate.DBValidator.Struct(c)
}

func (ci *ClubInform) BeforeCreate() error {
	return validate.DBValidator.Struct(ci)
}

func (cm *ClubMember) BeforeCreate() error {
	return validate.DBValidator.Struct(cm)
}

func (cr *ClubRecruitment) BeforeCreate() error {
	return validate.DBValidator.Struct(cr)
}

func (rm *RecruitMember) BeforeCreate() error {
	return validate.DBValidator.Struct(rm)
}

func (c *Club) BeforeUpdate() error {
	clubForValidate := c.DeepCopy()

	if clubForValidate.UUID == emptyString       { clubForValidate.UUID = validClubUUID }
	if clubForValidate.LeaderUUID == emptyString { clubForValidate.LeaderUUID = validLeaderUUID }

	return validate.DBValidator.Struct(clubForValidate)
}

func (ci *ClubInform) BeforeUpdate() error {
	clubInformForValidate := ci.DeepCopy()

	if clubInformForValidate.ClubUUID == emptyString { clubInformForValidate.ClubUUID = validClubUUID }
	if clubInformForValidate.Name == emptyString     { clubInformForValidate.Name = validName }
	if clubInformForValidate.Field == emptyString    { clubInformForValidate.Field = validField }
	if clubInformForValidate.Location == emptyString { clubInformForValidate.Location = validLocation }
	if clubInformForValidate.Floor == emptyInt       { clubInformForValidate.Floor = validFloor }
	if clubInformForValidate.LogoURI == emptyString  { clubInformForValidate.LogoURI = validLogoURI }

	return validate.DBValidator.Struct(clubInformForValidate)
}

func (cr *ClubRecruitment) BeforeUpdate() error {
	recruitmentForValidate := cr.DeepCopy()

	if recruitmentForValidate.UUID == emptyString           { recruitmentForValidate.UUID = validRecruitmentUUID }
	if recruitmentForValidate.ClubUUID == emptyString       { recruitmentForValidate.ClubUUID = validClubUUID }
	if recruitmentForValidate.RecruitConcept == emptyString { recruitmentForValidate.RecruitConcept = validRecruitConcept }

	return validate.DBValidator.Struct(recruitmentForValidate)
}
