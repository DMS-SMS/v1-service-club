package errors

import "errors"

var (
	RecruitmentUUIDCannotBeChanged = errors.New("recruitment uuid cannot be changed")
	ClubUUIDCannotBeChanged = errors.New("club uuid cannot be changed")
)