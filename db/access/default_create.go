package access

import "club/model"

func (d *_default) CreateClub(club *model.Club) (*model.Club, error) {
	err := d.tx.Create(club).Error
	return club, err
}

func (d *_default) CreateClubInform(inform *model.ClubInform) (*model.ClubInform, error) {
	err := d.tx.Create(inform).Error
	return inform, err
}

func (d *_default) CreateClubMember(member *model.ClubMember) (*model.ClubMember, error) {
	err := d.tx.Create(member).Error
	return member, err
}

func (d *_default) CreateRecruitment(recruitment *model.ClubRecruitment) (*model.ClubRecruitment, error) {
	err := d.tx.Create(recruitment).Error
	return recruitment, err
}

func (d *_default) CreateRecruitMember(member *model.RecruitMember) (*model.RecruitMember, error) {
	err := d.tx.Create(member).Error
	return member, err
}
