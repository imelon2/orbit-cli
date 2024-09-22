package utils

import (
	"regexp"
	"strings"
)

func IsAddress(addr string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(addr)
}

func IsTransaction(addr string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{64}$")
	return re.MatchString(addr)
}

func IsPrivateKey(addr string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{64}$")
	return re.MatchString(addr)
}

func IsWithSpace(s string) bool {
	// 문자열에 공백이 포함되어 있는지 검사
	return strings.Contains(s, " ")
}
