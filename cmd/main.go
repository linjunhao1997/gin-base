package main

import (
	initlization "gin-base/init"
	"gin-base/internal/web/router"
)

func main() {
	initlization.Initialize()
	router.Router.Run() // listen and serve on 0.0.0.0:8080

}
