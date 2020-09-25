package model

import (
	"club/tool/random"
	"database/sql/driver"
	"fmt"
	"time"
)

var (
	nullReplaceValueForStartPeriod time.Time
	nullReplaceValueForEndPeriod time.Time
)

func init() {
	nullReplaceValueForStartPeriod = time.Date(0, 0, 0, 0, random.IntWithLength(11), 0, 0, time.UTC)
	nullReplaceValueForEndPeriod = time.Date(0, 0, 0, 0, random.IntWithLength(11), 0, 0, time.UTC)
}

// UUID 필드에서 사용할 사용자 정의 타입
type uuid string
func UUID(s string) uuid { return uuid(s) }
func (u uuid) Value() (driver.Value, error) { return string(u), nil }
func (u *uuid) Scan(src interface{}) (err error) { *u = uuid(src.([]uint8)); return }
func (u uuid) KeyName() string { return "uuid" }

// LeaderUUID 필드에서 사용할 사용자 정의 타입
type leaderUUID string
func LeaderUUID(s string) leaderUUID { return leaderUUID(s) }
func (lu leaderUUID) Value() (driver.Value, error) { return string(lu), nil }
func (lu *leaderUUID) Scan(src interface{}) (err error) { *lu = leaderUUID(src.([]uint8)); return }
func (lu leaderUUID) KeyName() string { return "leader_uuid" }

// ClubUUID 필드에서 사용할 사용자 정의 타입
type clubUUID string
func ClubUUID(s string) clubUUID { return clubUUID(s) }
func (cu clubUUID) Value() (driver.Value, error) { return string(cu), nil }
func (cu *clubUUID) Scan(src interface{}) (err error) { *cu = clubUUID(src.([]uint8)); return }
func (cu clubUUID) KeyName() string { return "club_uuid" }

// Name 필드에서 사용할 사용자 정의 타입
type name string
func Name(s string) name { return name(s) }
func (n name) Value() (driver.Value, error) { return string(n), nil }
func (n *name) Scan(src interface{}) (err error) { *n = name(src.([]uint8)); return }
func (n name) KeyName() string { return "name" }

// ClubConcept 필드에서 사용할 사용자 정의 타입
type clubConcept string
func ClubConcept(s string) clubConcept { return clubConcept(s) }
func (cc clubConcept) Value() (driver.Value, error) { return string(cc), nil }
func (cc *clubConcept) Scan(src interface{}) (err error) { *cc = clubConcept(src.([]uint8)); return }
func (cc clubConcept) KeyName() string { return "club_concept" }

// Introduction 필드에서 사용할 사용자 정의 타입
type introduction string
func Introduction(s string) introduction { return introduction(s) }
func (i introduction) Value() (driver.Value, error) { return string(i), nil }
func (i *introduction) Scan(src interface{}) (err error) { *i = introduction(src.([]uint8)); return }
func (i introduction) KeyName() string { return "introduction" }

// Introduction 필드에서 사용할 사용자 정의 타입
type field string
func Field(s string) field { return field(s) }
func (f field) Value() (driver.Value, error) { return string(f), nil }
func (f *field) Scan(src interface{}) (err error) { *f = field(src.([]uint8)); return }
func (f field) KeyName() string { return "field" }

// Introduction 필드에서 사용할 사용자 정의 타입
type location string
func Location(s string) location { return location(s) }
func (l location) Value() (driver.Value, error) { return string(l), nil }
func (l *location) Scan(src interface{}) (err error) { *l = location(src.([]uint8)); return }
func (l location) KeyName() string { return "location" }

// Floor 필드에서 사용할 사용자 정의 타입
type floor int64
func Floor(i int64) floor { return floor(i) }
func (f floor) Value() (driver.Value, error) { return int64(f), nil }
func (f *floor) Scan(src interface{}) (err error) { *f = floor(src.(int64)); return }
func (f floor) KeyName() string { return "floor" }

// Link 필드에서 사용할 사용자 정의 타입
type link string
func Link(s string) link { return link(s) }
func (l link) Value() (driver.Value, error) { return string(l), nil }
func (l *link) Scan(src interface{}) (err error) { *l = link(src.([]uint8)); return }
func (l link) KeyName() string { return "link" }

// LogoURI 필드에서 사용할 사용자 정의 타입
type logoURI string
func LogoURI(s string) logoURI { return logoURI(s) }
func (lu logoURI) Value() (driver.Value, error) { return string(lu), nil }
func (lu *logoURI) Scan(src interface{}) (err error) { *lu = logoURI(src.([]uint8)); return }
func (lu logoURI) KeyName() string { return "logo_uri" }

// StudentUUID 필드에서 사용할 사용자 정의 타입
type studentUUID string
func StudentUUID(s string) studentUUID { return studentUUID(s) }
func (su studentUUID) Value() (driver.Value, error) { return string(su), nil }
func (su *studentUUID) Scan(src interface{}) (err error) { *su = studentUUID(src.([]uint8)); return }
func (su studentUUID) KeyName() string { return "student_uuid" }

// RecruitConcept 필드에서 사용할 사용자 정의 타입
type recruitConcept string
func RecruitConcept(s string) recruitConcept { return recruitConcept(s) }
func (rc recruitConcept) Value() (driver.Value, error) { return string(rc), nil }
func (rc *recruitConcept) Scan(src interface{}) (err error) { *rc = recruitConcept(src.([]uint8)); return }
func (rc recruitConcept) KeyName() string { return "recruit_concept" }

// StartPeriod 필드에서 사용할 사용자 정의 타입
type startPeriod time.Time
func StartPeriod(t time.Time) startPeriod { return startPeriod(t) }
func (sp startPeriod) Value() (value driver.Value, err error) {
	start := time.Time(sp)
	if (start == time.Time{}) {
		value = nil
	} else {
		value = fmt.Sprintf("%04d-%02d-%02d", start.Year(), start.Month(), start.Day())
	}
	return
}
func (sp *startPeriod) Scan(src interface{}) (err error) { *sp = startPeriod(src.(time.Time)); return }
func (sp startPeriod) KeyName() string { return "start_period" }
func (sp startPeriod) NullReplaceValue() time.Time { return nullReplaceValueForStartPeriod  }

// EndPeriod 필드에서 사용할 사용자 정의 타입
type endPeriod time.Time
func EndPeriod(t time.Time) endPeriod { return endPeriod(t) }
func (ep endPeriod) Value() (value driver.Value, err error) {
	end := time.Time(ep)
	if (end == time.Time{}) {
		value = nil
	} else {
		value = fmt.Sprintf("%04d-%02d-%02d", end.Year(), end.Month(), end.Day())
	}
	return
}
func (ep *endPeriod) Scan(src interface{}) (err error) { *ep = endPeriod(src.(time.Time)); return }
func (ep endPeriod) KeyName() string { return "end_period" }
func (ep endPeriod) NullReplaceValue() time.Time { return nullReplaceValueForEndPeriod  }

// RecruitmentUUID 필드에서 사용할 사용자 정의 타입
type recruitmentUUID string
func RecruitmentUUID(s string) recruitmentUUID { return recruitmentUUID(s) }
func (ru recruitmentUUID) Value() (driver.Value, error) { return string(ru), nil }
func (ru *recruitmentUUID) Scan(src interface{}) (err error) { *ru = recruitmentUUID(src.([]uint8)); return }
func (ru recruitmentUUID) KeyName() string { return "recruitment_uuid" }

// Grade 필드에서 사용할 사용자 정의 타입
type grade int64
func Grade(i int64) grade { return grade(i) }
func (g grade) Value() (value driver.Value, err error) { return int64(g), nil }
func (g *grade) Scan(src interface{}) (_ error) { *g = grade(src.(int64)); return }
func (g grade) KeyName() string { return "grade" }

// StudentNumber 필드에서 사용할 사용자 정의 타입
type number int64
func Number(i int64) number { return number(i) }
func (n number) Value() (driver.Value, error) { return int64(n), nil }
func (n *number) Scan(src interface{}) (err error) { *n = number(src.(int64)); return }
func (n number) KeyName() string { return "number" }
