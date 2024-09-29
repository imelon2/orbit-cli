package arblib

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Retryable struct {
	From                   common.Address `json:from,omitempty`
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

const (
	ErrorTrigger_GasLimit          = 1
	ErrorTrigger_MaxFeePerGas      = 1
	ErrorTrigger_MaxSubmissionCost = 1
)

// func PpopulateRetryable(client *ethclient.Client, from common.Address, to *common.Address, data []byte) {
// 	zombieMsg := ethereum.CallMsg{
// 		From:      from,
// 		To:        to,
// 		Gas:       tx.Gas(),
// 		GasPrice:  tx.GasPrice(),
// 		GasFeeCap: tx.GasFeeCap(),
// 		GasTipCap: tx.GasTipCap(),
// 		Value:     new(big.Int),
// 		Data:      data,
// 	}
// 	client.CallContract(context)
// }

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
		common.Address{},
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

// uint256 값을 디코딩하는 함수
func DecodeUint256(data []byte) (*big.Int, error) {
	if len(data) != 32 {
		return nil, fmt.Errorf("expected 32 bytes, got %d", len(data))
	}
	return new(big.Int).SetBytes(data), nil
}
