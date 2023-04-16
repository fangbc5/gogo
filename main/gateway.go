//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"gogo/app/gateway"
	"gogo/app/gateway/conf"
	"gogo/app/gateway/rpc/provider"
	"gogo/constant"
	"gogo/core"
	"gogo/core/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {

	//加载配置
	config.LoadLocalConfig("app/gateway/conf", conf.Config)
	config.LoadRemoteConfig(conf.Config.Consul.Config, conf.Config)
	//consul注册
	core.ConsulConn(constant.DirApp, "http://127.0.0.1:8080/gateway/health")
	defer core.ConsulConnClose(constant.Namespace + "-" + constant.DirApp)
	//加载核心
	r := gin.Default()
	core.Load(r,conf.Config.Mysql,conf.Config.Redis)
	//加载grpc服务
	grpcServer := provider.GetGrpcServer()
	core.GrpcLoad(r, grpcServer)
	//加载路由
	gateway.Routers(r)
	//启动服务
	handler := h2c.NewHandler(r, &http2.Server{})
	server := &http.Server{
		Addr:    ":" + conf.Config.Server.Port,
		Handler: handler,
	}

	//http启动
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("startup service failed, err:%v\n", err)
	}

	//ssl方式启动
	//if err := r.RunTLS(":8080", "你的SSL证书路径", "你的SSL私钥路径"); err != nil {
	//	fmt.Printf("startup service failed, err:%v\n", err)
	//}

	//直接gin启动
	//if err := r.Run(server); err != nil {
	//	fmt.Printf("startup service failed, err:%v\n", err)
	//}
}
