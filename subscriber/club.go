package subscriber

import (
	"context"
	log "github.com/micro/go-micro/v2/logger"

	club "club/proto/club"
)

type Club struct{}

func (e *Club) Handle(ctx context.Context, msg *club.Message) error {
	log.Info("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *club.Message) error {
	log.Info("Function Received message: ", msg.Say)
	return nil
}
