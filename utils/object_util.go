package utils

import (
	"reflect"
	"strings"
)

func IsNull(obj interface{}) bool {
	kind := reflect.TypeOf(obj).Kind()
	if obj == nil {
		return true
	} else if kind == 22 {
		return reflect.ValueOf(obj).IsNil()
	} else if kind == reflect.String {
		str := reflect.ValueOf(obj).String()
		return IsBlack(str)
	}
	return false
}

func IsNotNull(obj interface{}) bool {
	return !IsNull(obj)
}

func IsBlack(str string) bool {
	str = strings.TrimSpace(str)
	return str == ""
}

func IsNotBlack(str string) bool {
	return !IsBlack(str)
}

func IsEmpty(obj interface{}) bool {
	if obj == nil {
		return true
	}
	if reflect.TypeOf(obj).Kind() == reflect.Slice || reflect.TypeOf(obj).Kind() == reflect.Map {
		dataSlice := reflect.ValueOf(obj)
		return dataSlice.Len() == 0
	}
	return false
}

func IsNotEmpty(obj interface{}) bool {
	return !IsEmpty(obj)
}
