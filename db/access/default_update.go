package access

import (
	"club/db/access/errors"
	"club/model"
	"time"
)

func (d *_default) ChangeClubLeader(clubUUID, newLeaderUUID string) (err error, rowAffected int64) {
	updateResult := d.tx.Model(&model.Club{}).Where("uuid = ?", clubUUID).Updates(&model.Club{
		LeaderUUID: model.LeaderUUID(newLeaderUUID),
	})
	err = updateResult.Error
	rowAffected = updateResult.RowsAffected
	return
}

func (d *_default) ModifyClubInform(clubUUID string, revisionInform *model.ClubInform) (err error, rowAffected int64) {
	if revisionInform.ClubUUID != "" {
		err = errors.ClubUUIDCannotBeChanged
		return
	}

	updateResult := d.tx.Model(&model.ClubInform{}).Where("club_uuid = ?", clubUUID).Updates(revisionInform)
	err = updateResult.Error
	rowAffected = updateResult.RowsAffected
	return
}

func (d *_default) ModifyRecruitment(recruitUUID string, revisionRecruit *model.ClubRecruitment) (err error, rowAffected int64) {
	if revisionRecruit.UUID != "" {
		err = errors.RecruitmentUUIDCannotBeChanged
		return
	}

	if revisionRecruit.ClubUUID != "" {
		err = errors.ClubUUIDCannotBeChanged
		return
	}

	var updateAttrs []interface{}

	if revisionRecruit.RecruitConcept != "" {
		updateAttrs = append(updateAttrs, model.ClubRecruitmentInstance.RecruitConcept.KeyName())
	}

	if revisionRecruit.StartPeriod != model.StartPeriod(time.Time{}) {
		updateAttrs = append(updateAttrs, model.ClubRecruitmentInstance.StartPeriod.KeyName())
		if revisionRecruit.StartPeriod == model.StartPeriod(revisionRecruit.StartPeriod.NullReplaceValue()) {
			revisionRecruit.StartPeriod = model.StartPeriod(time.Time{})
		}
	}

	if revisionRecruit.EndPeriod != model.EndPeriod(time.Time{}) {
		updateAttrs = append(updateAttrs, model.ClubRecruitmentInstance.EndPeriod.KeyName())
		if revisionRecruit.EndPeriod == model.EndPeriod(revisionRecruit.EndPeriod.NullReplaceValue()) {
			revisionRecruit.EndPeriod = model.EndPeriod(time.Time{})
		}
	}

	selectedTx := d.tx.Model(&model.ClubRecruitment{})
	if len(updateAttrs) != 0 {
		argsAttr := updateAttrs[1:]
		selectedTx = selectedTx.Select(updateAttrs[0], argsAttr...)
	} else {
		selectedTx = selectedTx.Select("")
	}

	updateResult := selectedTx.Where("uuid = ?", recruitUUID).Updates(revisionRecruit)
	err = updateResult.Error
	rowAffected = updateResult.RowsAffected
	return
}
