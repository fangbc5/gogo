package utils

import (
	"fmt"
	"gogo/constant"
	"testing"
)

func TestIsPhone(t *testing.T) {
	fmt.Println(IsEamil("fangbc5@163.com"))
}

func TestObject(t *testing.T) {
	fmt.Println(IsNull(""))
}

func TestRpcGet(t *testing.T) {
	get := RpcGet(constant.ResourceApp, "/health", "", map[string][]string{})
	fmt.Println(get)
}

func TestRpcPost(t *testing.T) {
	post := RpcPost(constant.DirServerAddr, "/dir/dirRoot/getRoot", "{\"rootType\": 102,\"resType\": 1}", map[string][]string{})
	fmt.Println(post)
}
