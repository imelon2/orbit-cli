package main

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	types "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	precompilesgen "github.com/imelon2/orbit-cli/solgen/go/precompilesgen"
)

// var ArbosAddress = common.HexToAddress("0xa4b05")
// var ArbosStateAddress = common.HexToAddress("0xA4B05FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
// var ArbSysAddress = common.HexToAddress("0x64")
// var ArbInfoAddress = common.HexToAddress("0x65")
// var ArbAddressTableAddress = common.HexToAddress("0x66")
// var ArbBLSAddress = common.HexToAddress("0x67")
// var ArbFunctionTableAddress = common.HexToAddress("0x68")
// var ArbosTestAddress = common.HexToAddress("0x69")
// var ArbGasInfoAddress = common.HexToAddress("0x6c")
// var ArbOwnerPublicAddress = common.HexToAddress("0x6b")
// var ArbAggregatorAddress = common.HexToAddress("0x6d")
// var ArbRetryableTxAddress = common.HexToAddress("0x6e")
// var ArbStatisticsAddress = common.HexToAddress("0x6f")
// var ArbOwnerAddress = common.HexToAddress("0x70")
// var ArbWasmAddress = common.HexToAddress("0x71")
// var ArbWasmCacheAddress = common.HexToAddress("0x72")
// var NodeInterfaceAddress = common.HexToAddress("0xc8")
// var NodeInterfaceDebugAddress = common.HexToAddress("0xc9")
// var ArbDebugAddress = common.HexToAddress("0xff")

func main() {
	client, err := ethclient.Dial("http://localhost:8547")
	if err != nil {
		log.Fatal(err)
	}

	ArbGasInfo, err := precompilesgen.NewArbGasInfo(types.ArbGasInfoAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	opts := &bind.CallOpts{
		Pending: false, // 트랜잭션이 확정된 상태를 조회
		// BlockNumber: nil,   // 현재 상태 조회
		// From:        ,   // 상태 조회에는 보통 사용되지 않음
		// BlockHash:        ,   // 상태 조회에는 보통 사용되지 않음
		Context: nil, // 컨텍스트가 필요한 경우 (예: 시간 초과)
	}

	result, err := ArbGasInfo.GetCurrentTxL1GasFees(opts)
	// result, err := ArbGasInfo.GetL1BaseFeeEstimateInertia(opts)

	fmt.Print(result)

}
