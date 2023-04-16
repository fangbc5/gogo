package routers

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"gogo/core/render"
	"net/http"
)

type option func(*gin.Engine)

var options []option

func include(opts ...option) {
	options = append(options, opts...)
}

func Load(r *gin.Engine) *gin.Engine {
	// 绑定路由
	include(index)
	//include(employer.Routers) //用户模块
	// include(dir.Routers)      //目录模块
	// include(resource.Routers) //资源管理
	// include(paper.Routers)    //试卷模块
	// include(bank.Routers)     //试卷模块
	// include(test.Routers)     //试卷模块
	// include(notify.Routers)   //用户通知
	for _, opt := range options {
		opt(r)
	}
	return r
}

func index(e *gin.Engine) {
	e.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", render.Context{"aaa": 123})
	})
	e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
