package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type UsersController struct {
}

func (ctr *UsersController) Register(ctx *gin.Context) {
	var registerParams struct {
		Mobild  string `json:"mobile"`
		Pwd     string `json:"password"`
		SmsCode string `json:"sms_code"`
	}

	ctx.Bind(&registerParams)
	fmt.Println(registerParams)
}
