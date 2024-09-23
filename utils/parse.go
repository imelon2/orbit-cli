package utils

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"reflect"
)

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
