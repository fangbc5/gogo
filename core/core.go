package core

import (
	"fmt"
	"gogo/constant"
	"gogo/core/config"
	"gogo/core/db"
	"gogo/core/middleware"
	"gogo/core/render"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
)

var ConsulClient *api.Client

func Load(r *gin.Engine,mysql config.Mysql,redis config.Redis) {
	//连接数据库
	db.MysqlConn(mysql)
	db.RedisConn(redis)
	//db.BigCacheConn()
	//设置模式
	SetMode()
	//加载中间件
	middleware.Load(r)
	// 加载模版
	tmplLoad(r)
}

func SetMode() {
	switch os.Getenv("MODE") {
	case constant.ReleaseMode:
		gin.SetMode(gin.ReleaseMode)

	case constant.DebugMode:
		gin.SetMode(gin.DebugMode)

	case constant.TestMode:
		gin.SetMode(gin.TestMode)

	default:
		gin.SetMode(gin.DebugMode)
	}
}

func tmplLoad(r *gin.Engine) {
	if gin.IsDebugging() {
		r.HTMLRender = render.NewDebug("web/template")
	} else {
		r.HTMLRender = render.NewProduction("web/template")
	}
	r.Static("/static", "web/static")
}

func ConsulConn(serviceName string, healthUrl string) {

	//获取consul客户端
	consulClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		fmt.Println("consul.api.NewClient:", err)
	}

	//注册到consul
	consulClient.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      constant.Namespace + "-" + serviceName,
		Name:    constant.Namespace + "-" + serviceName,
		Tags:    []string{"gogo-micro"},
		Address: config.All.Consul.Host,
		Port:    config.All.Consul.Port,
		Check: &api.AgentServiceCheck{
			CheckID:  "health check",
			HTTP:     healthUrl,
			Timeout:  "30s",
			Interval: "10s",
		},
	})
	ConsulClient = consulClient
}

func ConsulConnClose(serviceID string) {
	ConsulClient.Agent().ServiceDeregister(serviceID)
}

func GrpcLoad(r *gin.Engine, g *grpc.Server) {
	r.Use(func(ctx *gin.Context) {
		// 判断协议是否为http/2
		// 判断是否是grpc
		if ctx.Request.ProtoMajor == 2 &&
			strings.HasPrefix(ctx.GetHeader("Content-Type"), "application/grpc") {
			// 按grpc方式来请求
			g.ServeHTTP(ctx.Writer, ctx.Request)
			// 不要再往下请求了,防止继续链式调用拦截器
			ctx.Abort()
			return
		}
		// 当作普通api
		ctx.Next()
	})
}
