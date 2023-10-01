package api

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	pb "therealbroker/api/proto"
	"therealbroker/pkg/broker"
)

type server struct {
	pb.UnimplementedBrokerServer
	broker broker.Broker
}

func (s *server) Publish(ctx context.Context, request *pb.PublishRequest) (*pb.PublishResponse, error) {
	newMsg := broker.NewCreateMessageDTO(request.Subject, string(request.Body), request.ExpirationSeconds)
	id, err := s.broker.Publish(ctx, newMsg)
	if err != nil {
		log.Printf("Publish failed %v", err)
	}
	return &pb.PublishResponse{
		Id: int32(id),
	}, nil

}

func (s *server) Subscribe(request *pb.SubscribeRequest, srv pb.Broker_SubscribeServer) error {
	return nil
}

func (s *server) Fetch(ctx context.Context, request *pb.FetchRequest) (*pb.MessageResponse, error) {
	msg, err := s.broker.Fetch(ctx, request.Subject, int(request.Id))
	if err != nil {
		return nil, err
	}
	return &pb.MessageResponse{
		Body: []byte(msg.Body),
	}, nil
}

func StartGrpcServer() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Starting server failed %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterBrokerServer(s, &server{})
	err = s.Serve(listener)
	if err != nil {
		log.Fatalf("Starting server failed %v", err)
	}
}
