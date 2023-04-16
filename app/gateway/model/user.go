package model

/*
*
此处定义user结构体
针对不同系统的user表应自行转换成该user格式
*/
type User struct {
	Id       int
	Username string
	Password string
	Nickname string
	Avator   string
	Email    string
	Phone    string
	Gender   int
	Country  string
	Province string
	City     string
	//用户类型（system、weixin、qq、sina、alipay）
	Usertype string
	OpenId   string
	UnionId  string
	//实名信息
	Realname string
	Idcard   string
}
