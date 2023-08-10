package main

import (
	"github.com/Wuhao-9/IHome/services/authentication/handler"
	pb "github.com/Wuhao-9/IHome/services/authentication/proto"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
)

func main() {
	// Create service
	srv := micro.NewService(
		micro.Name("IHome.auth"),
		micro.Address(":12345"),
	)
	// Register handler
	pb.RegisterAuthenticationHandler(srv.Server(), new(handler.AuthSrvImpl))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}