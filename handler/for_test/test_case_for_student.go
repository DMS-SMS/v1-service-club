package test

type GetClubsSortByUpdateTimeCase struct {
	UUID string
	Field, Name string
	Start, Count uint32
	XRequestID        string
	SpanContextString string
	ExpectedMethods   map[Method]Returns
	ExpectedStatus    uint32
	ExpectedCode      int32
	ExpectedClubUUID  string
}
