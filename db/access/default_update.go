package access

import (
	"club/db/access/errors"
	"club/model"
	"time"
)

func (d *_default) ChangeClubLeader(clubUUID, newLeaderUUID string) (err error) {
	err = d.tx.Model(model.Club{}).Where("uuid = ?", clubUUID).Update("leader_uuid", newLeaderUUID).Error
	return
}

func (d *_default) ModifyClubInform(clubUUID string, revisionInform *model.ClubInform) (err error) {
	contextForUpdate := make(map[string]interface{}, 8)

	if revisionInform.ClubUUID != "" {
		err = errors.ClubUUIDCannotBeChanged
		return
	}

	if revisionInform.Name != ""         { contextForUpdate[revisionInform.Name.KeyName()] = revisionInform.Name }
	if revisionInform.ClubConcept != ""  { contextForUpdate[revisionInform.ClubConcept.KeyName()] = revisionInform.ClubConcept }
	if revisionInform.Introduction != "" { contextForUpdate[revisionInform.Introduction.KeyName()] = revisionInform.Introduction }
	if revisionInform.Field != ""        { contextForUpdate[revisionInform.Field.KeyName()] = revisionInform.Field }
	if revisionInform.Location != ""     { contextForUpdate[revisionInform.Location.KeyName()] = revisionInform.Location }
	if revisionInform.Floor != 0         { contextForUpdate[revisionInform.Floor.KeyName()] = revisionInform.Floor }
	if revisionInform.Link != ""         { contextForUpdate[revisionInform.Link.KeyName()] = revisionInform.Link }
	if revisionInform.LogoURI != ""      { contextForUpdate[revisionInform.LogoURI.KeyName()] = revisionInform.LogoURI }

	err = d.tx.Model(&model.ClubInform{}).Where("club_uuid = ?", clubUUID).Updates(contextForUpdate).Error
	return
}

func (d *_default) ModifyRecruitment(recruitUUID string, revisionRecruit *model.ClubRecruitment) (err error) {
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

	err = d.tx.Model(&model.ClubRecruitment{}).Where("recruitment_uuid = ?", recruitUUID).Updates(contextForUpdate).Error
	return
}