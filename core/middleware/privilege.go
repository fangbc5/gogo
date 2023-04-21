package middleware

import (
	"github.com/fangbc5/gogo/constant"
	"github.com/fangbc5/gogo/core/auth"
	"github.com/fangbc5/gogo/utils"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

func Privilege() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		//method := c.Request.Method
		token := c.GetHeader(constant.Authorization)
		if utils.IsBlack(token) {
			token = c.GetHeader(strings.ToLower(constant.Authorization))
		}
		//白名单
		whiteList := make([]string, 16)
		whiteList = append(whiteList, "/login")
		whiteList = append(whiteList, "/swagger/*")
		whiteList = append(whiteList, "/register")
		//如果是白名单请求直接放行
		inWhiteList := false
		for _, whitePath := range whiteList {
			if ok, _ := regexp.MatchString(whitePath, path); ok && utils.IsNotBlack(whitePath) {
				inWhiteList = true
			}
		}
		if inWhiteList {
			c.Next()
		} else if utils.IsBlack(token) {
			//如果token为空则用户未登录直接返回
			c.JSON(http.StatusOK, model.GetRspAll(http.StatusUnauthorized, "未登录，请登录后访问", nil))
			c.Abort()
		} else {
			//不在白名单中，判断用户是否已登录
			username := db.RedisCache("get", constant.TokenKey+token)
			if utils.IsNull(username) {
				//用户登录状态已过期或token无效
				c.JSON(http.StatusOK, common.GetRspAll(http.StatusUnauthorized, "用户登录状态已过期或token错误，请重新登录", nil))
				c.Abort()
				return
			} else {
				//casbinValid(c,username,path,method)
				ok, msg := tokenValid(token)
				if ok {
					c.Next()
				} else {
					c.Abort()
					c.JSON(http.StatusOK, common.GetRspAll(http.StatusUnauthorized, msg, nil))
					return
				}
			}
		}
	}
}

func casbinValid(sub string, obj string, act string) (bool, string) {
	//用户已登录,使用casbin进行鉴权操作
	e := auth.MakeEnforcer()
	e.EnableLog(true)
	err := e.LoadPolicy()
	if err != nil {
		log.Panicln("Casbin LoadPolicy Error")
	}
	//权限认证
	ok, err := e.Enforce(sub, obj, act)
	if err != nil {
		return false, "Casbin Enforce Error"
	}
	if ok {
		return true, "验证成功"
	} else {
		return false, "验证失败"
	}
}

func tokenValid(token string) (bool, string) {
	claims, err := auth.VerifyToken(token)
	if err != nil {
		return false, err.Error()
	}
	if claims == nil {
		return false, "claims is nil"
	}
	return true, "校验成功"
}
