package main

import (
	"log"

	"github.com/SomeHowMicroservice/user/config"
	"github.com/SomeHowMicroservice/user/server"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Tải cấu hình User Service thất bại: %v", err)
	}

	s, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("Khởi tạo service thất bại: %v", err)
	}

	ch := make(chan error, 1)

	go func() {
		if err := s.Start(); err != nil {
			ch <- err
		}
	}()

	log.Println("Chạy service thành công")

	s.GracefulShutdown(ch)
}
