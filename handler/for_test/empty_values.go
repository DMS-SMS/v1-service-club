package test

import (
	clubproto "club/proto/golang/club"
	"club/tool/random"
)

const (
	EmptyString = ""
	EmptyInt = 0
)

var (
	EmptyReplaceValueForString string
	EmptyReplaceValueForMemberUUIDs []string
	emptyReplaceValueForMemberUUIDsLen int
	EmptyReplaceValueForRecruitMembers []*clubproto.RecruitMember
	emptyReplaceValueForRecruitMembersLen int
)

func init() {
	EmptyReplaceValueForString = random.StringConsistOfIntWithLength(10)
	emptyReplaceValueForMemberUUIDsLen = random.IntWithLengthWithoutZero(3)
	EmptyReplaceValueForMemberUUIDs = make([]string, emptyReplaceValueForMemberUUIDsLen)
	emptyReplaceValueForRecruitMembersLen = random.IntWithLengthWithoutZero(3)
	EmptyReplaceValueForRecruitMembers = make([]*clubproto.RecruitMember, emptyReplaceValueForRecruitMembersLen)
}
