package utils

import (
	"fmt"
	"math/big"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

const ADDRESS_ALIAS_OFFSET = "1111000000000000000000000000000000001111" // 0x는 SetString에서 처리하지 않음
const ADDRESS_BIT_LENGTH = 160
const ADDRESS_NIBBLE_LENGTH = ADDRESS_BIT_LENGTH / 4

func IsAddress(addr string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(addr)
}

func IsTransaction(addr string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{64}$")
	return re.MatchString(addr)
}

func IsBytes(addr string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]$")
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

func Alias(address common.Address, forward bool) common.Address {
	AddressAliasOffset, success := new(big.Int).SetString(ADDRESS_ALIAS_OFFSET, 16)

	if !success {
		panic("Error initializing AddressAliasOffset")
	}

	// "0x" 접두사 제거
	addressStr := strings.TrimPrefix(address.Hex(), "0x")
	// address를 big.Int로 변환
	originalAddress := new(big.Int)
	originalAddress.SetString(addressStr, 16)

	var result *big.Int
	if forward {
		result = new(big.Int).Add(originalAddress, AddressAliasOffset)
	} else {
		result = new(big.Int).Sub(originalAddress, AddressAliasOffset)
	}

	resultUint := new(big.Int).And(result, big.NewInt(0).Sub(new(big.Int).Lsh(big.NewInt(1), ADDRESS_BIT_LENGTH), big.NewInt(1)))

	resultStr := resultUint.Text(16)
	resultStr = strings.ToLower(resultStr)                            // 소문자로 변환
	resultStr = fmt.Sprintf("%0*s", ADDRESS_NIBBLE_LENGTH, resultStr) // Nibble 길이(40자리)에 맞게 0 패딩

	return common.HexToAddress(resultStr)
}

func SafeGetAddressString(slice []common.Address, index int) string {
	if index >= 0 && index < len(slice) {
		return slice[index].Hex()
	}
	return ""
}

func SafeGetLongestArray[T any](arrays ...[]T) int {
	var longest []T

	for _, array := range arrays {
		if len(array) > len(longest) {
			longest = array
		}
	}

	return len(longest)
}
