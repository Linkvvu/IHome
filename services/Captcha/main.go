package main

import (
	"github.com/Wuhao-9/IHome/services/Captcha/handler"
	pb "github.com/Wuhao-9/IHome/services/Captcha/proto"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(service.Name("Captcha"), service.Address(":12345"))

	// Register handler
	pb.RegisterCaptchaHandler(srv.Server(), new(handler.CaptchaSrvImpl))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
