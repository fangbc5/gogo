package middleware

import (
	"errors"
	"github.com/fangbc5/gogo/constant"
	"github.com/fangbc5/gogo/core/auth"
	"github.com/fangbc5/gogo/core/common"
	"github.com/fangbc5/gogo/plugins/auth/jwt"
	"github.com/fangbc5/gogo/utils"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

func Privilege(publicKey string) gin.HandlerFunc {
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
			c.JSON(http.StatusUnauthorized, common.GetRsp(http.StatusUnauthorized, "未登录，请登录后访问", nil))
			c.Abort()
		} else {
			//校验token
			jwtImpl := jwt.NewAuth(auth.WithPublicKey(publicKey))
			if _, err := jwtImpl.Inspect(token); err != nil {
				c.JSON(http.StatusUnauthorized, common.GetRsp(http.StatusUnauthorized, "token无效", nil))
				c.Abort()
			}
			c.Next()
		}
	}
}

func casbinValid(sub string, obj string, act string) error {
	//用户已登录,使用casbin进行鉴权操作
	e := auth.MakeEnforcer()
	e.EnableLog(true)
	if err := e.LoadPolicy(); err != nil {
		return errors.New("casbin LoadPolicy Error")
	}
	//权限认证
	if ok, err := e.Enforce(sub, obj, act); !ok || err != nil {
		return errors.New("casbin Enforce Error")
	}
	return nil
}
