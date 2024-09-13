package utils

import (
	"regexp"
)

func IsAddress(addr string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(addr)
}

func IsTransaction(addr string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{64}$")
	return re.MatchString(addr)
}
