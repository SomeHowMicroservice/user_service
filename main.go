package main

import (
	"fmt"
	"log"
	"net"

	"github.com/SomeHowMicroservice/shm-be/user/config"
	"github.com/SomeHowMicroservice/shm-be/user/initialization"
	"github.com/SomeHowMicroservice/shm-be/user/server"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Tải cấu hình User Service thất bại: %v", err)
	}

	db, err := initialization.InitDB(cfg)
	if err != nil {
		log.Fatalf("Lỗi kết nối DB ở User Service: %v", err)
	}
	defer db.Close()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.App.GRPCPort))
	if err != nil {
		log.Fatalf("Không thể lắng nghe: %v", err)
	}
	defer lis.Close()

	grpcServer := server.NewGRPCServer(db.Gorm)

	log.Println("Khởi chạy service thành công")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Chạy gRPC server thất bại: %v", err)
	}
}
