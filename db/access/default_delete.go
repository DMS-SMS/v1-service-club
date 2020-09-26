package access

import "club/model"

func (d *_default) DeleteClub(clubUUID string) (err error, rowsAffected int64) {
	deleteResult := d.tx.Where("uuid = ?", clubUUID).Delete(&model.Club{})
	err = deleteResult.Error
	rowsAffected = deleteResult.RowsAffected
	return
}

func (d *_default) DeleteClubInform(clubUUID string) (err error, rowsAffected int64) {
	deleteResult := d.tx.Where("club_uuid = ?", clubUUID).Delete(&model.ClubInform{})
	err = deleteResult.Error
	rowsAffected = deleteResult.RowsAffected
	return
}

func (d *_default) DeleteClubMember(clubUUID, studentUUID string) (err error, rowsAffected int64) {
	deleteResult := d.tx.Where("club_uuid = ? AND student_uuid = ?", clubUUID, studentUUID).Delete(&model.ClubMember{})
	err = deleteResult.Error
	rowsAffected = deleteResult.RowsAffected
	return
}

func (d *_default) DeleteRecruitment(recruitUUID string) (err error, rowsAffected int64) {
	deleteResult := d.tx.Where("uuid = ?", recruitUUID).Delete(&model.ClubRecruitment{})
	err = deleteResult.Error
	rowsAffected = deleteResult.RowsAffected
	return
}

func (d *_default) DeleteAllRecruitMember(recruitUUID string) (err error, rowsAffected int64) {
	deleteResult := d.tx.Where("recruitment_uuid = ?", recruitUUID).Delete(&model.RecruitMember{})
	err = deleteResult.Error
	rowsAffected = deleteResult.RowsAffected
	return
}
