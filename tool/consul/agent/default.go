package agent

import (
	"github.com/micro/go-micro/v2/client/selector"
	"github.com/micro/go-micro/v2/registry"
)

type _default struct {
	Strategy selector.Strategy
	next     selector.Next
	node     []*registry.Node
}

func Default(fn selector.Strategy) *_default {
	return &_default{
		Strategy: fn,
	}
}
