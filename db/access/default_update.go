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
	contextForUpdate := make(map[string]interface{}, 8)

	if revisionRecruit.UUID != "" {
		err = errors.RecruitmentUUIDCannotBeChanged
		return
	}

	if revisionRecruit.ClubUUID != "" {
		err = errors.ClubUUIDCannotBeChanged
		return
	}

	if revisionRecruit.RecruitConcept != "" { contextForUpdate[revisionRecruit.RecruitConcept.KeyName()] = revisionRecruit.RecruitConcept }

	if revisionRecruit.StartPeriod != model.StartPeriod(time.Time{}) {
		if revisionRecruit.StartPeriod == model.StartPeriod(revisionRecruit.StartPeriod.NullReplaceValue()) {
			contextForUpdate[revisionRecruit.StartPeriod.KeyName()] = model.StartPeriod(time.Time{})
		} else {
			contextForUpdate[revisionRecruit.StartPeriod.KeyName()] = revisionRecruit.StartPeriod
		}
	}

	if revisionRecruit.EndPeriod != model.EndPeriod(time.Time{}) {
		if revisionRecruit.EndPeriod == model.EndPeriod(revisionRecruit.EndPeriod.NullReplaceValue()) {
			contextForUpdate[revisionRecruit.EndPeriod.KeyName()] = model.EndPeriod(time.Time{})
		} else {
			contextForUpdate[revisionRecruit.EndPeriod.KeyName()] = revisionRecruit.EndPeriod
		}
	}

	updateResult := d.tx.Model(&model.ClubRecruitment{}).Where("uuid = ?", recruitUUID).Updates(contextForUpdate)
	err = updateResult.Error
	rowAffected = updateResult.RowsAffected
	return
}
