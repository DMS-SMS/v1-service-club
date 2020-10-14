package test

type AddClubMemberCase struct {
	UUID                  string
	ClubUUID, StudentUUID string
	XRequestID            string
	SpanContextString     string
	ExpectedMethods       map[Method]Returns
	ExpectedStatus        uint32
	ExpectedCode          int32
}
