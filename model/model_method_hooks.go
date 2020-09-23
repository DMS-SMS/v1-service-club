package model

import "club/model/validate"

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
