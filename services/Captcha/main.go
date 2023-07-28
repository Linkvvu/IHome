package main

import (
	"github.com/Wuhao-9/IHome/services/Captcha/handler"
	pb "github.com/Wuhao-9/IHome/services/Captcha/proto"
	"github.com/micro/micro/v3/service/logger"
	"go-micro.dev/v4"
)

func main() {
	// Create service
	srv := micro.NewService(
		micro.Name("Captcha"),
		micro.Address(":12345"),
	)
	// Register handler
	pb.RegisterCaptchaHandler(srv.Server(), new(handler.CaptchaSrvImpl))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
