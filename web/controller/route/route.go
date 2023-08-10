package route

import (
	"github.com/Wuhao-9/IHome/web/controller"
	"github.com/gin-gonic/gin"
)

func InitRoute() *gin.Engine {
	router := gin.Default()
	router.Static("/static", "view")
	router.StaticFile("/", "view/index.html")

	v1 := router.Group("api/v1.0")
	{
		captchaCtr := (&controller.CaptchaController{})
		v1.GET("/imagecode/:uuid", captchaCtr.GetImgCaptcha)
		v1.GET("/smscode/:phomeNum", captchaCtr.GetPhoneCaptcha)
		v1.POST("/users", (&controller.UsersController{}).Register)
	}

	return router
}
