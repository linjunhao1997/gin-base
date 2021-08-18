package access

import (
	"gin-base/pkg/router"
	"github.com/gin-gonic/gin"
)

const (
	SysRolePath = "/SysRoles"
)

type SysRoleController struct {
}

func (c *SysRoleController) HandlerConfig() {
	router.V1.GET(SysRolePath+"/:id", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
