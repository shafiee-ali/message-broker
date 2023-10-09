package api

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
	pb "therealbroker/api/proto"
	"therealbroker/pkg/broker"
	"therealbroker/pkg/metrics"
	"time"
)

type server struct {
	pb.UnimplementedBrokerServer
	broker  broker.Broker
	metrics metrics.BrokerMetrics
}

func NewServer(b broker.Broker, metrics metrics.BrokerMetrics) *server {
	return &server{
		broker:  b,
		metrics: metrics,
	}
}

func (s *server) Publish(ctx context.Context, request *pb.PublishRequest) (*pb.PublishResponse, error) {
	currentTime := time.Now()
	defer s.metrics.RpcMethodLatency.WithLabelValues(metrics.PUBLISH_METHOD).Observe(float64(time.Since(currentTime).Nanoseconds()))
	newMsg := broker.NewCreateMessageDTO(request.Subject, string(request.Body), request.ExpirationSeconds)
	id, err := s.broker.Publish(ctx, newMsg)
	if err != nil {
		s.metrics.RpcMethodCount.WithLabelValues(metrics.PUBLISH_METHOD, metrics.ERROR_STATUS).Inc()
		log.Printf("Publish failed %v", err)
		return nil, err
	}
	s.metrics.RpcMethodCount.WithLabelValues(metrics.PUBLISH_METHOD, metrics.SUCCESS_STATUS).Inc()
	return &pb.PublishResponse{
		Id: int32(id),
	}, nil

}

func (s *server) Subscribe(request *pb.SubscribeRequest, srv pb.Broker_SubscribeServer) error {
	currentTime := time.Now()
	defer s.metrics.RpcMethodLatency.WithLabelValues(metrics.SUBSCRIBE_METHOD).Observe(float64(time.Since(currentTime).Nanoseconds()))
	ch, err := s.broker.Subscribe(srv.Context(), request.Subject)
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for {
			select {
			case msg := <-ch:
				srv.Send(
					&pb.MessageResponse{
						Body: []byte(msg.Body),
					},
				)
			case <-srv.Context().Done():
				wg.Done()
			}

		}
	}()
	wg.Wait()
	return err
}

func (s *server) Fetch(ctx context.Context, request *pb.FetchRequest) (*pb.MessageResponse, error) {
	currentTime := time.Now()
	defer s.metrics.RpcMethodLatency.WithLabelValues(metrics.FETCH_METHOD).Observe(float64(time.Since(currentTime).Nanoseconds()))
	msg, err := s.broker.Fetch(ctx, request.Subject, int(request.Id))
	if err != nil {
		s.metrics.RpcMethodCount.WithLabelValues(metrics.FETCH_METHOD, metrics.ERROR_STATUS).Inc()
		return nil, err
	}
	s.metrics.RpcMethodCount.WithLabelValues(metrics.FETCH_METHOD, metrics.SUCCESS_STATUS).Inc()
	return &pb.MessageResponse{
		Body: []byte(msg.Body),
	}, nil
}

func StartGrpcServer(b broker.Broker, metrics metrics.BrokerMetrics) {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Starting server failed %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterBrokerServer(s, NewServer(b, metrics))
	err = s.Serve(listener)
	if err != nil {
		log.Fatalf("Starting server failed %v", err)
	}
}
