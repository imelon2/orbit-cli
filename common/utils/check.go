package utils

import (
	"regexp"
	"strings"
)

func IsAddress(addr string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(addr)
}

func IsTransactionHash(addr string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{64}$")
	return re.MatchString(addr)
}

func IsBytes(addr string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]")
	return re.MatchString(addr)
}

func IsPrivateKey(addr string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{64}$")
	return re.MatchString(addr)
}

// 문자열에 공백이 포함되어 있는지 검사
func IsWithSpace(s string) bool {
	return strings.Contains(s, " ")
}
