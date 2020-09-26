package access

import "club/model"

func (d *_default) CreateClub(club *model.Club) (*model.Club, error) {
	result := d.tx.Create(club)
	return result.Value.(*model.Club), result.Error
}

func (d *_default) CreateClubInform(inform *model.ClubInform) (*model.ClubInform, error) {
	result := d.tx.Create(inform)
	return result.Value.(*model.ClubInform), result.Error
}

func (d *_default) CreateClubMember(member *model.ClubMember) (*model.ClubMember, error) {
	result := d.tx.Create(member)
	return result.Value.(*model.ClubMember), result.Error
}

func (d *_default) CreateRecruitment(recruitment *model.ClubRecruitment) (*model.ClubRecruitment, error) {
	result := d.tx.Create(recruitment)
	return result.Value.(*model.ClubRecruitment), result.Error
}

func (d *_default) CreateRecruitMember(member *model.RecruitMember) (*model.RecruitMember, error) {
	result := d.tx.Create(member)
	return result.Value.(*model.RecruitMember), result.Error
}
