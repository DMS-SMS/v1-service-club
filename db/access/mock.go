package access

import "github.com/stretchr/testify/mock"

type _mock struct {
	mock *mock.Mock
}

func Mock(mock *mock.Mock) _mock {
	return _mock{mock: mock}
}
