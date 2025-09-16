package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/SomeHowMicroservice/user/config"
	"github.com/SomeHowMicroservice/user/initialization"
	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *grpc.Server
	lis        net.Listener
	db         *initialization.DB
}

func NewServer(cfg *config.Config) (*Server, error) {
	db, err := initialization.InitDB(cfg)
	if err != nil {
		return nil, err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.App.GRPCPort))
	if err != nil {
		return nil, err
	}

	grpcServer := NewGRPCServer(db.Gorm)

	return &Server{
		grpcServer,
		lis,
		db,
	}, nil
}

func (s *Server) Start() error {
	return s.grpcServer.Serve(s.lis)
}

func (s *Server) Shutdown(ctx context.Context) {
	log.Println("Đang shutdown service...")

	if s.grpcServer != nil {
		stopped := make(chan struct{})
		go func() {
			s.grpcServer.GracefulStop()
			close(stopped)
		}()

		select {
		case <-ctx.Done():
			log.Println("Timeout khi dừng gRPC server, force stop...")
			s.grpcServer.Stop()
		case <-stopped:
			log.Println("Đã shutdown gRPC server")
		}
	}
	if s.lis != nil {
		s.lis.Close()
	}
	if s.db != nil {
		s.db.Close()
	}

	log.Println("Shutdown service thành công")
}
