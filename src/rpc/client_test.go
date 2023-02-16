package rpc

import (
	"context"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
	pb "pokemonscan-pokeball/src/proto/pokeball"
	"testing"
	"time"
)

type taskService struct {
	pb.UnimplementedTaskServiceServer
}

func (s *taskService) ReportCompletionStatus(ctx context.Context, args *pb.CompletionStatusArgs) (*pb.CompletionStatusReply, error) {
	timeStr := time.Now().Format("2006-01-02 15:04:05") //当前时间的字符串，2006-01-02 15:04:05据说是golang的诞生时间，固定写法
	log.Info(timeStr)
	return nil, status.Error(codes.Unavailable, "Unavailable")
}

func TestRetry(t *testing.T) {

	listen, err := net.Listen("tcp", ":6414")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTaskServiceServer(s, &taskService{})
	go func() {
		if err := s.Serve(listen); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}

	}()
	log.Info("run server 6414 ")

	scsource := `{
		"methodConfig": [{
		  "name": [{"service": "pokemon.proto.pokeball.TaskService","method":"ReportCompletionStatus"}],
		  "retryPolicy": {
			  "MaxAttempts": 10,
			  "InitialBackoff": "30s",
			  "MaxBackoff": "60s",
			  "BackoffMultiplier": 2.0,
			  "RetryableStatusCodes": [ "UNAVAILABLE" ,"UNKNOWN" ]
		  }
		}]}`

	conn, err := grpc.Dial("127.0.0.1:6414", grpc.WithInsecure(), grpc.WithDefaultServiceConfig(scsource))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := pb.NewTaskServiceClient(conn)
	_, err = client.ReportCompletionStatus(context.Background(), &pb.CompletionStatusArgs{TaskId: 1})
	if err != nil {
		log.Errorf("ReportCompletionStatus err: %v", err)
	}

	defer conn.Close()
}
