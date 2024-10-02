package arbmath

import (
	"fmt"
	"math/big"
)

func PercentIncrease(num *big.Int, increase *big.Int) *big.Int {
	value := new(big.Int).Set(num)
	value.Mul(value, increase)
	value.Div(value, big.NewInt(100))
	value.Add(value, num)
	return value
}

func ParseEther(etherAmount string) (*big.Int, error) {
	// etherAmount를 big.Float로 변환
	ether, ok := new(big.Float).SetString(etherAmount)
	if !ok {
		return nil, fmt.Errorf("invalid ether amount")
	}

	// 1 Ether = 10^18 Wei 이므로 10^18을 정의
	weiMultiplier := new(big.Float).SetFloat64(1e18)

	// Ether * 10^18 (Wei 단위로 변환)
	wei := new(big.Float).Mul(ether, weiMultiplier)

	// big.Float 값을 big.Int로 변환
	weiInt := new(big.Int)
	wei.Int(weiInt) // 소수점 이하 절삭

	return weiInt, nil
}
