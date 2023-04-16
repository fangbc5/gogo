//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"gogo/core"
	"gogo/core/config"
	_ "gogo/docs"
	"gogo/routers"
	"log"

	"github.com/gin-gonic/gin"
)

// @title 考试系统接口文档
// @version v1.0.0
// @description 应用程序API管理
// @termsOfService http://swagger.io/terms/
// @contact.name fangbc5
// @contact.url http://chinaexam.top
// @contact.email fangbc5@163.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host 127.0.0.1:8080
func main() {
	r := gin.Default()
	//初始化操作
	core.Load(r)
	//加载路由
	routers.Load(r)
	//服务端口
	server := "127.0.0.1" + ":" + config.All.Server.Port
	log.Println("接口文档：http://" + server + "/swagger/index.html")
	//启动服务
	if err := r.Run(server); err != nil {
		fmt.Printf("startup service failed, err:%v\n", err)
	}
}
