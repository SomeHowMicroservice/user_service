package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os/signal"
	"syscall"
	"time"

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

func (s *Server) GracefulShutdown(ch <-chan error) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-ch:
		log.Printf("Chạy service thất bại: %v", err)
	case <-ctx.Done():
		log.Println("Có tín hiệu dừng server")
	}

	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	s.Shutdown(shutdownCtx)
}
