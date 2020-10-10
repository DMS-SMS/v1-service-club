package consul

import "github.com/micro/go-micro/v2/registry"

type Agent interface {
	GetNextServiceNode(service string) (*registry.Node, error)
}
