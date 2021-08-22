package main

import (
	initialize "gin-base/init"
	"gin-base/pkg/router"
)

type A interface {
	GetID()
}

type B struct {
	name string
}

func (B) GetID() {

}

func main() {
	initialize.Load()
	root := router.Root
	root.Run() // listen and serve on 0.0.0.0:8080
}
