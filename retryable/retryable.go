package retryable

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/imelon2/orbit-cli/solgen/go/bridgegen"
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

type RetryableData struct {
	From                   common.Address // address
	To                     common.Address // address
	L2CallValue            *big.Int       // uint256
	Deposit                *big.Int       // uint256
	MaxSubmissionCost      *big.Int       // uint256
	ExcessFeeRefundAddress common.Address // address
	CallValueRefundAddress common.Address // address
	GasLimit               *big.Int       // uint256
	MaxFeePerGas           *big.Int       // uint256
	Data                   []byte         // bytes
}

func ParseRetryable(data []byte) (RetryableData, error) {
	parsed, err := bridgegen.InboxMetaData.GetAbi()
	if err != nil {
		return RetryableData{}, fmt.Errorf("failed get InboxMetaData abi : %d", err)
	}
	var sigdata [4]byte
	for i, data := range data[:4] {
		sigdata[i] = data
	}
	errorAbi, err := parsed.ErrorByID(sigdata)
	if err != nil {
		return RetryableData{}, err
	}

	if errorAbi.Name != "RetryableData" {
		return RetryableData{}, fmt.Errorf("no retryable data found in error: : %s", data)
	}

	decodedError, err := errorAbi.Inputs.Unpack(data[4:])
	if err != nil {
		return RetryableData{}, fmt.Errorf("failed to unpack Retryable error data : %v", err)
	}

	return RetryableData{
		From:                   decodedError[0].(common.Address),
		To:                     decodedError[1].(common.Address),
		L2CallValue:            decodedError[2].(*big.Int),
		Deposit:                decodedError[3].(*big.Int),
		MaxSubmissionCost:      decodedError[4].(*big.Int),
		ExcessFeeRefundAddress: decodedError[5].(common.Address),
		CallValueRefundAddress: decodedError[6].(common.Address),
		GasLimit:               decodedError[7].(*big.Int),
		MaxFeePerGas:           decodedError[8].(*big.Int),
		Data:                   decodedError[9].([]byte),
	}, nil
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
