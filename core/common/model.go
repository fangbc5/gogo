package common

import "net/http"

const (
	MaxPageSize = 1000
)

type PageReq struct {
	Page int `json:"-" gorm:"-"`
	Size int `json:"-" gorm:"-"`
}

type PageRsp struct {
	Rsp
	Total int64 `json:"total"`
	Pages int   `json:"pages"`
}

type Rsp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func GetRsp() *Rsp {
	return &Rsp{http.StatusOK, "ok", nil}
}
func GetRspMsg(msg string) *Rsp {
	return &Rsp{http.StatusOK, msg, nil}
}
func GetRspData(r interface{}) *Rsp {
	return &Rsp{http.StatusOK, "ok", r}
}

func GetFailMsg(msg string) *Rsp {
	return &Rsp{http.StatusNoContent, msg, nil}
}

func GetRspAll(code int, msg string, r interface{}) *Rsp {
	return &Rsp{code, msg, r}
}

func GetPageRsp() *PageRsp {
	return &PageRsp{*GetRsp(), 0, 0}
}

func MakePageArg() *PageReq {
	return &PageReq{1, 10}
}

func PageArg(page int, size int) *PageReq {
	return &PageReq{page, size}
}
