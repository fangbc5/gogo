package middleware

import "github.com/gin-gonic/gin"

var middleWares []gin.HandlerFunc

func Load(r *gin.Engine) *gin.Engine {
	//增加全局中间件
	include(Cors())
	//认证中间件，认证功能交给istio
	//include(Privilege())
	for _, ware := range middleWares {
		r.Use(ware)
	}
	return r
}
func include(opts ...gin.HandlerFunc) {
	middleWares = append(middleWares, opts...)
}
