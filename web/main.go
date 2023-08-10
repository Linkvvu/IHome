package main

import (
	"github.com/Wuhao-9/IHome/web/controller/route"
)

func main() {
	router := route.InitRoute()
	router.Run(":8080")
}
