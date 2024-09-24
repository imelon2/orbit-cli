package utils

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
)

type Retryable struct {
	To                     common.Address `json:to`
	L2CallValue            *big.Int       `json:l2CallValue`
	Amount                 *big.Int       `json:amount`
	MaxSubmissionCost      *big.Int       `json:maxSubmissionCost`
	ExcessFeeRefundAddress common.Address `json:excessFeeRefundAddress`
	CallValueRefundAddress common.Address `json:callValueRefundAddress`
	GasLimit               *big.Int       `json:gasLimit`
	MaxFeePerGas           *big.Int       `json:maxFeePerGas`
	DataLength             *big.Int       `json:dataLength`
	Data                   string         `json:data`
}

func ConvertBytesToHex(data interface{}) interface{} {
	val := reflect.ValueOf(data)

	switch val.Kind() {
	case reflect.Slice:
		// []byte인지 확인
		if val.Type().Elem().Kind() == reflect.Uint8 {
			// []byte를 hex 문자열로 변환하고 "0x" 접두사 추가
			byteSlice := val.Bytes()
			hexString := "0x" + hex.EncodeToString(byteSlice)
			return hexString
		} else {
			// 다른 슬라이스 타입일 경우, 각 요소를 재귀적으로 변환
			updatedSlice := make([]interface{}, val.Len())
			for i := 0; i < val.Len(); i++ {
				updatedSlice[i] = ConvertBytesToHex(val.Index(i).Interface())
			}
			return updatedSlice
		}
	case reflect.Array:
		// [N]byte인지 확인
		if val.Type().Elem().Kind() == reflect.Uint8 {
			// [N]byte를 []byte로 변환한 후 hex 문자열로 변환하고 "0x" 접두사 추가
			byteArray := make([]byte, val.Len())
			for i := 0; i < val.Len(); i++ {
				byteArray[i] = byte(val.Index(i).Uint())
			}
			hexString := "0x" + hex.EncodeToString(byteArray)
			return hexString
		} else {
			// 다른 배열 타입일 경우, 각 요소를 재귀적으로 변환
			updatedArray := make([]interface{}, val.Len())
			for i := 0; i < val.Len(); i++ {
				updatedArray[i] = ConvertBytesToHex(val.Index(i).Interface())
			}
			return updatedArray
		}
	case reflect.Map:
		// 맵 타입일 경우, 키와 값을 재귀적으로 변환
		updatedMap := make(map[string]interface{})
		for _, key := range val.MapKeys() {
			// 키는 문자열로 변환 (만약 키가 다른 타입이라면 적절히 변환 필요)
			keyStr := fmt.Sprintf("%v", key.Interface())
			updatedMap[keyStr] = ConvertBytesToHex(val.MapIndex(key).Interface())
		}
		return updatedMap
	case reflect.Struct:
		// 구조체 타입일 경우, 필드를 재귀적으로 변환
		updatedStruct := make(map[string]interface{})
		for i := 0; i < val.NumField(); i++ {
			field := val.Type().Field(i)
			fieldName := field.Name
			fieldValue := val.Field(i).Interface()
			updatedStruct[fieldName] = ConvertBytesToHex(fieldValue)
		}
		return updatedStruct
	default:
		// 다른 타입일 경우, 그대로 반환
		return data
	}
}

func FloatToWei(a *big.Float) *big.Int {
	weiFactor := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)

	etherInWei := new(big.Float).Set(a)

	weiAmount := new(big.Float).Mul(etherInWei, new(big.Float).SetInt(weiFactor))

	weiResult := new(big.Int)
	weiAmount.Int(weiResult)

	return weiResult
}

// uint256 값을 디코딩하는 함수
func DecodeUint256(data []byte) (*big.Int, error) {
	if len(data) != 32 {
		return nil, fmt.Errorf("expected 32 bytes, got %d", len(data))
	}
	return new(big.Int).SetBytes(data), nil
}

func ParseRetryableMessage(msg []byte) Retryable {

	buf := bytes.NewReader(msg)

	toBytes := make([]byte, 32)
	l2CallValueBytes := make([]byte, 32)
	amountBytes := make([]byte, 32)
	maxSubmissionCostBytes := make([]byte, 32)
	excessFeeRefundAddressBytes := make([]byte, 32)
	callValueRefundAddressBytes := make([]byte, 32)
	gasLimitBytes := make([]byte, 32)
	maxFeePerGasBytes := make([]byte, 32)
	dataLengthBytes := make([]byte, 32)

	_, err := buf.Read(toBytes)
	if err != nil {
		log.Fatal(err)
	}
	_, err = buf.Read(l2CallValueBytes)
	if err != nil {
		log.Fatal(err)
	}
	_, err = buf.Read(amountBytes)
	if err != nil {
		log.Fatal(err)
	}
	_, err = buf.Read(maxSubmissionCostBytes)
	if err != nil {
		log.Fatal(err)
	}
	_, err = buf.Read(excessFeeRefundAddressBytes)
	if err != nil {
		log.Fatal(err)
	}
	_, err = buf.Read(callValueRefundAddressBytes)
	if err != nil {
		log.Fatal(err)
	}
	_, err = buf.Read(gasLimitBytes)
	if err != nil {
		log.Fatal(err)
	}
	_, err = buf.Read(maxFeePerGasBytes)
	if err != nil {
		log.Fatal(err)
	}
	_, err = buf.Read(dataLengthBytes)
	if err != nil {
		log.Fatal(err)
	}

	l2CallValue, err := DecodeUint256(l2CallValueBytes)
	if err != nil {
		log.Fatal(err)
	}
	amount, err := DecodeUint256(amountBytes)
	if err != nil {
		log.Fatal(err)
	}
	maxSubmissionCost, err := DecodeUint256(maxSubmissionCostBytes)
	if err != nil {
		log.Fatal(err)
	}
	gasLimit, err := DecodeUint256(gasLimitBytes)
	if err != nil {
		log.Fatal(err)
	}
	maxFeePerGas, err := DecodeUint256(maxFeePerGasBytes)
	if err != nil {
		log.Fatal(err)
	}
	dataLength, err := DecodeUint256(dataLengthBytes)
	if err != nil {
		log.Fatal(err)
	}

	remainingData := make([]byte, buf.Len()) // 남은 데이터 길이만큼 버퍼를 할당
	_, err = buf.Read(remainingData)
	if err != nil {
		log.Fatal(err)
	}

	return Retryable{
		common.HexToAddress(hex.EncodeToString(toBytes)),
		l2CallValue,
		amount,
		maxSubmissionCost,
		common.HexToAddress(hex.EncodeToString(excessFeeRefundAddressBytes)),
		common.HexToAddress(hex.EncodeToString(callValueRefundAddressBytes)),
		gasLimit,
		maxFeePerGas,
		dataLength,
		"0x" + hex.EncodeToString(remainingData),
	}
}
