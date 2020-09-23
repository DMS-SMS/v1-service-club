package validate

import "regexp"

const (
	adminUUIDRegexString = "^admin-\\d{12}"
	studentUUIDRegexString = "^student-\\d{12}"
	teacherUUIDRegexString = "^teacher-\\d{12}"
	parentUUIDRegexString = "^parent-\\d{12}"
	clubUUIDRegexString = "^club-\\d{12}"
	recruitmentUUIDRegexString = "^recruitment-\\d{12}"
	timeRegexString = "\\d{4}-\\d{2}-\\d{2}"
)

var (
	adminUUIDRegex = regexp.MustCompile(adminUUIDRegexString)
	studentUUIDRegex = regexp.MustCompile(studentUUIDRegexString)
	teacherUUIDRegex = regexp.MustCompile(teacherUUIDRegexString)
	parentUUIDRegex = regexp.MustCompile(parentUUIDRegexString)
	clubUUIDRegex = regexp.MustCompile(clubUUIDRegexString)
	recruitmentUUIDRegex = regexp.MustCompile(recruitmentUUIDRegexString)
	timeRegex = regexp.MustCompile(timeRegexString)
)
