package for_test

import (
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
)

func init() {
	EmptyReplaceValueForString = random.StringConsistOfIntWithLength(10)
	emptyReplaceValueForMemberUUIDsLen = random.IntWithLengthWithoutZero(3)
	EmptyReplaceValueForMemberUUIDs = make([]string, emptyReplaceValueForMemberUUIDsLen)
}
