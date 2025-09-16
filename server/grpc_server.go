package server

import (
	"time"

	"github.com/SomeHowMicroservice/user/container"
	userpb "github.com/SomeHowMicroservice/user/protobuf/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"gorm.io/gorm"
)

func NewGRPCServer(db *gorm.DB) *grpc.Server {
	kaParams := keepalive.ServerParameters{
		Time:                  5 * time.Minute,
		Timeout:               20 * time.Second,
		MaxConnectionIdle:     0,
		MaxConnectionAge:      0,
		MaxConnectionAgeGrace: 0,
	}

	kaPolicy := keepalive.EnforcementPolicy{
		MinTime:             1 * time.Minute,
		PermitWithoutStream: true,
	}

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(kaParams),
		grpc.KeepaliveEnforcementPolicy(kaPolicy),
	)

	userContainer := container.NewContainer(db, grpcServer)

	userpb.RegisterUserServiceServer(grpcServer, userContainer.GRPCHandler)

	return grpcServer
}
