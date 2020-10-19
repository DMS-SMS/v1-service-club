package handler

import (
	proto "club/proto/golang/club"
	"context"
)

type None struct {}

func (n None) CreateNewClub(context.Context, *proto.CreateNewClubRequest, *proto.CreateNewClubResponse) (err error) { return }

func (n None) AddClubMember(context.Context, *proto.AddClubMemberRequest, *proto.AddClubMemberResponse) (err error) { return }
func (n None) DeleteClubMember(context.Context, *proto.DeleteClubMemberRequest, *proto.DeleteClubMemberResponse) (err error) { return }
func (n None) ChangeClubLeader(context.Context, *proto.ChangeClubLeaderRequest, *proto.ChangeClubLeaderResponse) (err error) { return }
func (n None) ModifyClubInform(context.Context, *proto.ModifyClubInformRequest, *proto.ModifyClubInformResponse) (err error) { return }
func (n None) DeleteClubWithUUID(context.Context, *proto.DeleteClubWithUUIDRequest, *proto.DeleteClubWithUUIDResponse) (err error) { return }
func (n None) RegisterRecruitment(context.Context, *proto.RegisterRecruitmentRequest, *proto.RegisterRecruitmentResponse) (err error) { return }
func (n None) ModifyRecruitment(context.Context, *proto.ModifyRecruitmentRequest, *proto.ModifyRecruitmentResponse) (err error) { return }
func (n None) DeleteRecruitmentWithUUID(context.Context, *proto.DeleteRecruitmentWithUUIDRequest, *proto.DeleteRecruitmentWithUUIDResponse) (err error) { return }

func (n None) GetClubsSortByUpdateTime(context.Context, *proto.GetClubsSortByUpdateTimeRequest, *proto.GetClubsSortByUpdateTimeResponse) (err error) { return }
func (n None) GetRecruitmentsSortByCreateTime(context.Context, *proto.GetRecruitmentsSortByCreateTimeRequest, *proto.GetRecruitmentsSortByCreateTimeResponse) (err error) { return }
func (n None) GetClubInformWithUUID(context.Context, *proto.GetClubInformWithUUIDRequest, *proto.GetClubInformWithUUIDResponse) (err error) { return }
func (n None) GetClubInformsWithUUIDs(context.Context, *proto.GetClubInformsWithUUIDsRequest, *proto.GetClubInformsWithUUIDsResponse) (err error) { return }
func (n None) GetRecruitmentInformWithUUID(context.Context, *proto.GetRecruitmentInformWithUUIDRequest, *proto.GetRecruitmentInformWithUUIDResponse) (err error) { return }
func (n None) GetRecruitmentUUIDWithClubUUID(context.Context, *proto.GetRecruitmentUUIDWithClubUUIDRequest, *proto.GetRecruitmentUUIDWithClubUUIDResponse) (err error) { return }
func (n None) GetRecruitmentUUIDsWithClubUUIDs(context.Context, *proto.GetRecruitmentUUIDsWithClubUUIDsRequest, *proto.GetRecruitmentUUIDsWithClubUUIDsResponse) (err error) { return }
func (n None) GetAllClubFields(context.Context, *proto.GetAllClubFieldsRequest, *proto.GetAllClubFieldsResponse) (err error) { return }
func (n None) GetTotalCountOfClubs(context.Context, *proto.GetTotalCountOfClubsRequest, *proto.GetTotalCountOfClubsResponse) (err error) { return }
func (n None) GetTotalCountOfCurrentRecruitments(context.Context, *proto.GetTotalCountOfCurrentRecruitmentsRequest, *proto.GetTotalCountOfCurrentRecruitmentsResponse) (err error) { return }
func (n None) GetClubUUIDWithLeaderUUID(context.Context, *proto.GetClubUUIDWithLeaderUUIDRequest, *proto.GetClubUUIDWithLeaderUUIDResponse) (err error) { return }
