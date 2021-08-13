package main

import (
	"gin-base/pkg/controller/access"
	"gin-base/pkg/router"
)

func main() {
	root := router.Root
	router.AppendController(&access.SysUserController{}, &access.SysRoleController{})
	router.ConfigHandler()
	root.Run() // listen and serve on 0.0.0.0:8080
}
