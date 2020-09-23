package model

import "database/sql/driver"

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
