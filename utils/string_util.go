package utils

import (
	"log"
	"regexp"
	"strconv"
)

// IsEamil 是否邮箱
func IsEamil(val string) bool {
	match, err := regexp.Match("\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*", []byte(val))
	if err != nil {
		return false
	}
	return match
}

// IsPhone 是否手机号
func IsPhone(val string) bool {
	return regexp.MustCompile("^(0|86|17951)?(13[0-9]|15[012356789]|166|17[3678]|18[0-9]|14[57])[0-9]{8}$").MatchString(val)
}

// IsNumber 是否纯数字
func IsNumber(val string) bool {
	match, err := regexp.MatchString("^[0-9]+$", val)
	if err != nil {
		return false
	}
	return match
}

// IsCardNo 是否身份证号
func IsCardNo(val string) bool {
	match, err := regexp.MatchString("^(^[1-9]\\d{7}((0\\d)|(1[0-2]))(([0|1|2]\\d)|3[0-1])\\d{3}$)|(^[1-9]\\d{5}[1-9]\\d{3}((0\\d)|(1[0-2]))(([0|1|2]\\d)|3[0-1])((\\d{4})|\\d{3}[Xx])$)$", val)
	if err != nil {
		return false
	}
	return match
}

func Str2Uint(val string) uint {
	intVal, err := strconv.Atoi(val)
	if err != nil {
		log.Panicln("string_util.Str2Uint执行错误，入参字符串不是一个有效的数字")
	}
	return uint(intVal)
}
