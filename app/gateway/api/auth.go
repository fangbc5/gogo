package api

import (
	"encoding/json"
	"gogo/app/gateway/conf"
	"gogo/app/gateway/model"
	"gogo/app/gateway/service"
	"gogo/constant"
	"gogo/core/common"
	"gogo/core/db"
	"gogo/utils"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// Login 登录接口
// @Summary 登录接口
// @Description 登录接口
// @Tags 系统接口
// @Accept mpfd
// @Produce json
// @Param username formData string true "用户名"
// @Param password formData string true "密码"
// @Param code formData string false "验证码"
// @Security ApiKeyAuth
// @Success 200
// @Router /login [POST]
func Login(c *gin.Context) {
	m := make(map[string]string, 3)
	err := c.ShouldBindJSON(&m)
	if err != nil {
		c.JSON(http.StatusOK, common.GetRspAll(http.StatusNoContent, "入参格式错误", nil))
		return
	}
	username := m["username"]
	password := m["password"]
	code := m["code"]
	if utils.IsBlack(username) || utils.IsBlack(password) {
		c.JSON(http.StatusOK, common.GetRspAll(http.StatusNoContent, "用户名或密码为空", nil))
		return
	}
	//如果开启了验证码功能
	if conf.Config.VerifyCode {
		if utils.IsBlack(code) {
			c.JSON(http.StatusOK, common.GetRspAll(http.StatusNoContent, "验证码不能为空", nil))
			return
		}
		verifyCode := db.RedisCache("get", constant.VerifyCode+username)
		if utils.IsNull(verifyCode) {
			c.JSON(http.StatusOK, common.GetRspAll(http.StatusNoContent, "未生成验证码或已过期", nil))
			return
		}
		if verifyCode != code {
			c.JSON(http.StatusOK, common.GetRspAll(http.StatusNoContent, "验证码错误", nil))
			return
		}
	}

	//判断用户名是否为手机号或邮箱
	success := false
	if utils.IsPhone(username) {
		success = getUserByPhone(username, password)
	} else if utils.IsEamil(username) {
		success = getUserByEmail(username, password)
	} else {
		success = getUserByUsername(username, password)
	}
	if success {
		c.JSON(http.StatusOK, common.GetRspAll(http.StatusNoContent, "用户名或密码不存在", nil))
		return
	}

	//生成token
	tokenString, err := getToken(username)
	if err != nil {
		log.Panicln(err)
	}
	if tokenString == "" {
		c.JSON(http.StatusOK, common.GetRspAll(http.StatusNoContent, "token生成失败", nil))
	}
	db.RedisCache("setex", constant.TokenKey+tokenString, conf.Config.TokenLife, username)
	c.JSON(http.StatusOK, common.GetRspAll(http.StatusOK, "登录成功", tokenString))
}

// Logout 登出接口
// @Summary 登出接口
// @Description 登出接口
// @Tags 系统接口
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Security ApiKeyAuth
// @Success 200
// @Router /logout [GET]
func Logout(c *gin.Context) {
	tokenString := strings.ReplaceAll(c.GetHeader(constant.Authorization), constant.Bearer, "")
	db.RedisCache("del", constant.TokenKey+tokenString)
	c.JSON(http.StatusOK, common.GetRspMsg("登出成功"))
}

// Register 注册接口
// @Summary 注册接口
// @Description 注册接口
// @Tags 系统接口
// @Accept json
// @Produce json
// @Param employerType formData string true "1-学员，2-考官"
// @Param employerName formData string true "名称/昵称"
// @Param username formData string true "用户名"
// @Param password formData string true "密码"
// @Security ApiKeyAuth
// @Success 200
// @Router /register [POST]
func Register(c *gin.Context) {
	m := make(map[string]string, 3)
	err := c.ShouldBindJSON(&m)
	if err != nil {
		c.JSON(http.StatusOK, common.GetRspAll(http.StatusNoContent, "参数解析失败，请检查参数是否正确", nil))
	}
	usertype := m["usertype"]
	nickname := m["nickname"]
	username := m["username"]
	password := m["password"]
	if utils.IsBlack(usertype) {
		c.JSON(http.StatusOK, common.GetRspAll(http.StatusNoContent, "用户类型不能为空", nil))
		return
	}
	if utils.IsBlack(username) || utils.IsBlack(password) {
		c.JSON(http.StatusOK, common.GetRspAll(http.StatusNoContent, "用户名或密码不能为空", nil))
		return
	}
	//密码加密
	fromPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusOK, common.GetRspAll(http.StatusNoContent, "密码加密异常", nil))
		return
	}
	employer := model.User{Usertype: usertype, Nickname: nickname, Username: username, Password: string(fromPassword)}
	es := &service.UserService{}
	eid := es.InsertUser(employer)
	if eid > 0 {
		tokenString, err := getToken(username)
		if err != nil {
			log.Panicln(err)
		}
		if tokenString == "" {
			c.JSON(http.StatusOK, common.GetRspAll(http.StatusNoContent, "token生成失败", nil))
		}
		db.RedisCache("setex", constant.TokenKey+tokenString, conf.Config.TokenLife, username)
		c.JSON(http.StatusOK, common.GetRspAll(http.StatusOK, "注册成功", tokenString))
	} else {
		c.JSON(http.StatusOK, common.GetRspAll(http.StatusNoContent, "注册失败", nil))
	}
}
func getToken(username string) (string, error) {
	cla := model.MyClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(conf.Config.JwtExpire * time.Second)}, // 过期时间
			Issuer:    "gogo-jwt",                                                                  // 签发人
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cla)
	return token.SignedString([]byte(conf.Config.JwtSecret)) // 进行签名生成对应的token
}

func getUserByUsername(username string, password string) (success bool) {
	m := make(map[string]string, 2)
	db.MysqlClient().Raw("select username,password from "+conf.Config.UserTableName+" where username = ? limit 1", username).Scan(&m)
	if utils.IsNotBlack(m["password"]) {
		err := bcrypt.CompareHashAndPassword([]byte(m["password"]), []byte(password))
		if err == nil {
			return true
		}
	}
	return false
}

func getUserByPhone(phone string, password string) (success bool) {
	m := make(map[string]string, 2)
	db.MysqlClient().Raw("select username,password from "+conf.Config.UserTableName+" where phone = ? limit 1", phone, password).Scan(&m)
	if utils.IsNotBlack(m["password"]) {
		err := bcrypt.CompareHashAndPassword([]byte(m["password"]), []byte(password))
		if err == nil {
			return true
		}
	}
	return false
}

func getUserByEmail(email string, password string) (success bool) {
	m := make(map[string]string, 2)
	db.MysqlClient().Raw("select username,password from "+conf.Config.UserTableName+" where email = ? limit 1", email, password).Scan(&m)
	if utils.IsNotBlack(m["password"]) {
		err := bcrypt.CompareHashAndPassword([]byte(m["password"]), []byte(password))
		if err == nil {
			return true
		}
	}
	return false
}

// LoginCode 获取验证码
// @Summary 获取验证码
// @Description 获取验证码
// @Tags 系统接口
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Security ApiKeyAuth
// @Success 200
// @Router /getLoginCode [GET]
func LoginCode(c *gin.Context) {
	tokenString := strings.ReplaceAll(c.GetHeader(constant.Authorization), constant.Bearer, "")
	username := db.RedisCache("get", constant.TokenKey+tokenString).(string)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := ""
	for i := 0; i < 6; i++ {
		n := r.Intn(10)
		code += strconv.Itoa(n)
	}
	//存入缓存
	db.RedisCache("setex", constant.VerifyCode+username, 300, code)
	c.JSON(http.StatusOK, common.GetRspData(code))
}

// GetUserInfo 获取用户信息
// @Summary 获取用户信息
// @Description 获取用户信息
// @Tags 系统接口
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Security ApiKeyAuth
// @Success 200
// @Router /getUserInfo [GET]
func GetUserInfo(c *gin.Context) {
	employer := GetUser(c)
	c.JSON(http.StatusOK, common.GetRspData(employer))
}

func GetUser(c *gin.Context) *model.User {
	tokenString := c.GetHeader(constant.Authorization)
	tokenString = strings.ReplaceAll(tokenString, constant.Bearer, "")
	employer := &model.User{}
	//校验token并获取username
	username := db.RedisCache("get", constant.TokenKey+tokenString)
	//使用用户名
	user := db.RedisCache("get", username)
	if utils.IsNull(user) {
		//从数据库中加载，放如缓存
		db.MysqlClient().Raw("select * from "+conf.Config.UserTableName+" where username = ?", username).Scan(employer)
		marshal, _ := json.Marshal(employer)
		db.RedisCache("set", username, string(marshal))
	} else {
		json.Unmarshal([]byte(user.(string)), employer)
	}
	return employer
}
