package handler

import (
	proto "club/proto/golang/club"
	"context"
	log "github.com/micro/go-micro/v2/logger"
)

func (d _default) ChangeAllServiceNodes(ctx context.Context, req *proto.Empty, resp *proto.Empty) (_ error) {
	err := d.consulAgent.ChangeAllServiceNodes()
	log.Infof("change all service nodes!, err: %v", err)
	return
}
