package model

type Club struct {
	UUID       uuid
	LeaderUUID leaderUUID
}

type ClubInform struct {
	ClubUUID     clubUUID
	Name         name
	ClubConcept  clubConcept
	Introduction introduction
	Field        field
	Location     location
	Floor        floor
	Link         link
	LogoURI      logoURI
}
