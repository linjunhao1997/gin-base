package controller

import (
	"gin-base/controller/access"
	"gin-base/router"
)

func Enable() {

	access.EnableController()

	for _, c := range router.Controllers {
		c.PathConfig()
	}
}
