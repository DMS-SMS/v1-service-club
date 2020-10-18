package access

import (
	"club/model"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type _mock struct {
	mock *mock.Mock
}

func Mock(mock *mock.Mock) _mock {
	return _mock{mock: mock}
}

func (m _mock) CreateClub(club *model.Club) (resultClub *model.Club, err error) {
	args := m.mock.Called(club)
	return args.Get(0).(*model.Club), args.Error(1)
}

func (m _mock) CreateClubInform(inform *model.ClubInform) (resultInform *model.ClubInform, err error) {
	args := m.mock.Called(inform)
	return args.Get(0).(*model.ClubInform), args.Error(1)
}

func (m _mock) CreateClubMember(clubMember *model.ClubMember) (resultMember *model.ClubMember, err error) {
	args := m.mock.Called(clubMember)
	return args.Get(0).(*model.ClubMember), args.Error(1)
}

func (m _mock) CreateRecruitment(recruit *model.ClubRecruitment) (resultRecruit *model.ClubRecruitment, err error) {
	args := m.mock.Called(recruit)
	return args.Get(0).(*model.ClubRecruitment), args.Error(1)
}

func (m _mock) CreateRecruitMember(recruitMember *model.RecruitMember) (resultMember *model.RecruitMember, err error) {
	args := m.mock.Called(recruitMember)
	return args.Get(0).(*model.RecruitMember), args.Error(1)
}

func (m _mock) GetClubWithClubUUID(clubUUID string) (*model.Club, error) {
	args := m.mock.Called(clubUUID)
	return args.Get(0).(*model.Club), args.Error(1)
}

func (m _mock) GetClubWithLeaderUUID(leaderUUID string) (*model.Club, error) {
	args := m.mock.Called(leaderUUID)
	return args.Get(0).(*model.Club), args.Error(1)
}

func (m _mock) GetCurrentRecruitmentWithClubUUID(clubUUID string) (*model.ClubRecruitment, error) {
	args := m.mock.Called(clubUUID)
	return args.Get(0).(*model.ClubRecruitment), args.Error(1)
}

func (m _mock) GetCurrentRecruitmentWithRecruitmentUUID(recruitmentUUID string) (*model.ClubRecruitment, error) {
	args := m.mock.Called(recruitmentUUID)
	return args.Get(0).(*model.ClubRecruitment), args.Error(1)
}

func (m _mock) GetClubInformsSortByUpdateTime(offset, limit int, field, name string) ([]*model.ClubInform, error) {
	args := m.mock.Called(offset, limit, field, name)
	return args.Get(0).([]*model.ClubInform), args.Error(1)
}

func (m _mock) GetCurrentRecruitmentsSortByCreateTime(offset, limit int, field, name string) ([]*model.ClubRecruitment, error) {
	args := m.mock.Called(offset, limit, field, name)
	return args.Get(0).([]*model.ClubRecruitment), args.Error(1)
}

func (m _mock) GetClubInformWithClubUUID(clubUUID string) (*model.ClubInform, error) {
	args := m.mock.Called(clubUUID)
	return args.Get(0).(*model.ClubInform), args.Error(1)
}

func (m _mock) GetRecruitmentWithRecruitmentUUID(recruitUUID string) (*model.ClubRecruitment, error) {
	args := m.mock.Called(recruitUUID)
	return args.Get(0).(*model.ClubRecruitment), args.Error(1)
}

func (m _mock) GetClubMembersWithClubUUID(clubUUID string) ([]*model.ClubMember, error) {
	args := m.mock.Called(clubUUID)
	return args.Get(0).([]*model.ClubMember), args.Error(1)
}

func (m _mock) GetRecruitMembersWithRecruitmentUUID(recruitUUID string) ([]*model.RecruitMember, error) {
	args := m.mock.Called(recruitUUID)
	return args.Get(0).([]*model.RecruitMember), args.Error(1)
}

func (m _mock) GetAllClubInforms() ([]*model.ClubInform, error) {
	args := m.mock.Called()
	return args.Get(0).([]*model.ClubInform), args.Error(1)
}

func (m _mock) GetAllCurrentRecruitments() ([]*model.ClubRecruitment, error) {
	args := m.mock.Called()
	return args.Get(0).([]*model.ClubRecruitment), args.Error(1)
}

func (m _mock) ChangeClubLeader(clubUUID, newLeaderUUID string) (error, int64) {
	args := m.mock.Called(clubUUID, newLeaderUUID)
	return args.Error(0), int64(args.Int(1))
}

func (m _mock) ModifyClubInform(clubUUID string, revisionInform *model.ClubInform) (error, int64) {
	args := m.mock.Called(clubUUID, revisionInform)
	return args.Error(0), int64(args.Int(1))
}

func (m _mock) ModifyRecruitment(recruitUUID string, revisionRecruit *model.ClubRecruitment) (error, int64) {
	args := m.mock.Called(recruitUUID, revisionRecruit)
	return args.Error(0), int64(args.Int(1))
}

func (m _mock) DeleteClub(clubUUID string) (error, int64) {
	args := m.mock.Called(clubUUID)
	return args.Error(0), int64(args.Int(1))
}

func (m _mock) DeleteClubInform(clubUUID string) (error, int64) {
	args := m.mock.Called(clubUUID)
	return args.Error(0), int64(args.Int(1))
}

func (m _mock) DeleteClubMember(clubUUID, studentUUID string) (error, int64) {
	args := m.mock.Called(clubUUID, studentUUID)
	return args.Error(0), int64(args.Int(1))
}

func (m _mock) DeleteAllClubMembers(clubUUID string) (error, int64) {
	args := m.mock.Called(clubUUID)
	return args.Error(0), int64(args.Int(1))
}

func (m _mock) DeleteRecruitment(recruitUUID string) (error, int64) {
	args := m.mock.Called(recruitUUID)
	return args.Error(0), int64(args.Int(1))
}

func (m _mock) DeleteAllRecruitMember(recruitUUID string) (error, int64) {
	args := m.mock.Called(recruitUUID)
	return args.Error(0), int64(args.Int(1))
}

func (m _mock) BeginTx() {
	m.mock.Called()
}

func (m _mock) Commit() *gorm.DB {
	return m.mock.Called().Get(0).(*gorm.DB)
}

func (m _mock) Rollback() *gorm.DB {
	return m.mock.Called().Get(0).(*gorm.DB)
}