package agent

import (
	"github.com/micro/go-micro/v2/registry"
	"github.com/stretchr/testify/mock"
)

type _mock struct {
	mock *mock.Mock
}

func Mock(mock *mock.Mock) _mock {
	return _mock{mock: mock}
}


func (m _mock) GetNextServiceNode() (*registry.Node, error) {
	args := m.mock.Called()
	return args.Get(0).(*registry.Node), args.Error(1)
}
