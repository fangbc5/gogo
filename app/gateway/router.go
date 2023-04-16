package gateway

import (
	"gogo/app/gateway/api"
	"gogo/core/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Routers(e *gin.Engine) {
	gateway := e.Group("/gateway")
	gateway.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, common.GetRsp())
	})
	auth := e.Group("/auth")
	auth.POST("/login", api.Login)
	auth.POST("/register", api.Register)
	auth.GET("/logout", api.Logout)
	auth.GET("/getUserInfo", api.GetUserInfo)
	auth.GET("/getLoginCode", api.LoginCode)
}
