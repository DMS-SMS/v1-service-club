package handler

import (
	"context"

	log "github.com/micro/go-micro/v2/logger"

	club "club/proto/club"
)

type Club struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Club) Call(ctx context.Context, req *club.Request, rsp *club.Response) error {
	log.Info("Received Club.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Club) Stream(ctx context.Context, req *club.StreamingRequest, stream club.Club_StreamStream) error {
	log.Infof("Received Club.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Infof("Responding: %d", i)
		if err := stream.Send(&club.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *Club) PingPong(ctx context.Context, stream club.Club_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Infof("Got ping %v", req.Stroke)
		if err := stream.Send(&club.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
