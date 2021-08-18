package main

import (
	"gin-base/init"
	"gin-base/pkg/router"
)

func main() {
	initialize.Load()
	root := router.Root
	root.Run() // listen and serve on 0.0.0.0:8080
}
