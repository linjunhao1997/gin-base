package initialize

import (
	"fmt"
	"gin-base/global/db"
	"gin-base/global/mid"
	model "gin-base/model/access"
	"gin-base/pkg/base"
	"gin-base/pkg/router"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var identityKey = "id"

func JwtMiddleware() {
	newJwtMiddleware()
	router.V1.Use(mid.JwtMiddleware.MiddlewareFunc())
}

func newJwtMiddleware() {
	middleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "auth",
		Key:         []byte("2021"),
		Timeout:     time.Hour * 24 * 7,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*model.SysUser); ok {
				return jwt.MapClaims{
					identityKey: v.UserName,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &model.SysUser{
				UserName: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var userInfo UserInfo
			if err := c.ShouldBind(&userInfo); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			username := userInfo.Username
			password := userInfo.Password

			user := &model.SysUser{
				UserName: username,
				Password: password,
			}

			err := db.DB.Model(user).Where("username = ? and password = ?", user.UserName, user.Password).Take(&user).Error
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			} else {
				return user, nil
			}
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*model.SysUser); ok && v.UserName == "admin" {
				return true
			}

			return false

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			g := base.Gin{C: c}
			g.RespUnauthorized("")
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		panic(fmt.Sprintf("初始化Jwt中间件失败: %v", err))
	}

	// When you use jwt.New(), the function is already automatically called for checking,
	// which means you don't need to call it again.
	errInit := middleware.MiddlewareInit()

	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	mid.JwtMiddleware = middleware
}
