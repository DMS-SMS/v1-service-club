// add file in v.1.0.5
// this file declare method that handling event about consul in _default struct

package handler

import (
	"github.com/aws/aws-sdk-go/service/sqs"
	log "github.com/micro/go-micro/v2/logger"
)

func (d *_default) ChangeConsulNodes(message *sqs.Message) (err error) {
	err = d.consulAgent.ChangeAllServiceNodes()
	log.Infof("change all service nodes!, err: %v", err)
	return
}
