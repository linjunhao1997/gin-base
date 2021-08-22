package global

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/casbin/casbin/v2"
	"gorm.io/gorm"
)

var DB *gorm.DB
var Enforcer *casbin.Enforcer
var JwtMiddleware *jwt.GinJWTMiddleware
