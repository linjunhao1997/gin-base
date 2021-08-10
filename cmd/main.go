package main

import (
	"gin-base/controller"
	"gin-base/router"
)

func main() {

	controller.Enable()

	root := router.Root
	root.Run() // listen and serve on 0.0.0.0:8080
}
