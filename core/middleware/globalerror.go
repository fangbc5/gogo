package middleware

import (
	"fmt"
	"net/http"

	"github.com/fangbc5/gogo/core/common"
	"github.com/fangbc5/gogo/core/exception"
	"github.com/gin-gonic/gin"
)


func GlobalError() gin.HandlerFunc {
	return func(c *gin.Context) {
		exception.Try(func() {
			c.Next() //  处理请求
		},func(err interface{}) {
			c.Abort()
			c.JSON(http.StatusOK,common.GetFailMsg(fmt.Sprintf("error: %v",err)))
		})
	}
}