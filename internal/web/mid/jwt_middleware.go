package mid

import (
	"fmt"
	accessmodel "gin-base/internal/model/access"
	"gin-base/internal/pkg/db"
	"gin-base/internal/web/base"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"time"
)

type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewJwtMiddleware() *jwt.GinJWTMiddleware {
	middleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "auth",
		Key:         []byte("2021"),
		Timeout:     time.Hour * 24 * 7,
		IdentityKey: "userInfo",
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*accessmodel.SysUser); ok {
				return jwt.MapClaims{
					"id":       v.ID,
					"username": v.Username,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			id := int(claims["id"].(float64))
			return &accessmodel.SysUser{
				ID:       id,
				Username: claims["username"].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var userInfo UserInfo
			if err := c.ShouldBind(&userInfo); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			username := userInfo.Username
			password := userInfo.Password

			user := &accessmodel.SysUser{
				Username: username,
				Password: password,
			}

			err := db.G().Model(user).Where("username = ? and password = ?", user.Username, user.Password).Take(&user).Error
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			} else {
				return user, nil
			}
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(*accessmodel.SysUser); ok {
				return true
			}

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

	return middleware
}

func NewAccessGinHandlerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		g := base.Gin{C: c}
		sysUser := g.EnsureSysUser()
		if sysUser.Username == "admin" {
			c.Next()
			return
		}
		ok, err := db.G().Enforce(strconv.Itoa(sysUser.ID), c.Request.RequestURI, c.Request.Method)
		if err != nil {
			g.Abort(err)
			return
		} else if !ok {
			g.RespForbidden("")
			return

		} else {
			c.Next()
		}
	}
}
