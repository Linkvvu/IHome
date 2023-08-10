package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	authPb "github.com/Wuhao-9/IHome/services/authentication/proto"
	"github.com/gin-gonic/gin"
	"go-micro.dev/v4"
	"go-micro.dev/v4/errors"
)

type CaptchaController struct{}

func (ctr *CaptchaController) GetImgCaptcha(ctx *gin.Context) {
	srv := micro.NewService()
	authClient := authPb.NewAuthenticationService("IHome.auth", srv.Client())
	resp, err := authClient.GetImgCaptcha(context.TODO(), &authPb.ImgCaptchaRequ{Uuid: ctx.Param("uuid")})
	if err != nil {
		e := &errors.Error{}
		_ = json.Unmarshal([]byte(err.Error()), e)
		fmt.Println("val: ", e)
		ctx.JSON(http.StatusInternalServerError, e)
		return
	}
	ctx.String(http.StatusOK, string(resp.GetImage()))
	ctx.Writer.Write(resp.GetImage())
}

func (ctr *CaptchaController) GetPhoneCaptcha(ctx *gin.Context) {
	srv := micro.NewService()
	authClient := authPb.NewAuthenticationService("IHome.auth", srv.Client())

	uuid := ctx.Query("id")
	imageCode := ctx.Query("text")
	_, err := authClient.GetSmsCaptcha(context.TODO(), &authPb.SmsCaptchaRequ{Uuid: uuid, ImageCode: imageCode, PhoneNum: ctx.Param("phomeNum")})
	if err != nil {
		e := &errors.Error{}
		_ = json.Unmarshal([]byte(err.Error()), e)
		fmt.Println("val: ", e)
		ctx.JSON(http.StatusInternalServerError, e)
		return
	}
	ctx.Writer.WriteHeader(http.StatusOK)
}
