package consul

import "github.com/micro/go-micro/v2/registry"

type Agent interface {
	GetNextServiceNode() (registry.Node, error)
}
