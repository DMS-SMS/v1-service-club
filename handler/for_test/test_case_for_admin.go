package for_test

type Method string
type Returns []interface{}

type CreateNewClubCase struct {
	UUID              string
	Name, LeaderUUID  string
	MemberUUIDs       []string
	Field, Location   string
	Floor             uint32
	Logo              []byte
	ClubUUID          string
	XRequestID        string
	SpanContextString string
	ExpectedMethods   map[Method]Returns
	ExpectedStatus    uint32
	ExpectedCode      int32
	ExpectedClubUUID  string
}
