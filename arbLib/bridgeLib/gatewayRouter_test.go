package bridgelib_test

import (
	"fmt"
	"log"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	bridgelib "github.com/imelon2/orbit-cli/arbLib/bridgeLib"
)

func Test_getDepositInnerData(t *testing.T) {
	encodedData, _, err := bridgelib.GetDepositInnerData(big.NewInt(0), big.NewInt(79138), big.NewInt(600000000), true)
	if err != nil {
		log.Fatal(err)
	}
	// 인코딩된 데이터를 출력 (hex로 변환)
	fmt.Printf("ABI Encode DepositInnerData IsNative true : 0x%x\n", encodedData)

	encodedData, value, err := bridgelib.GetDepositInnerData(big.NewInt(163520), big.NewInt(79138), big.NewInt(600000000), false)
	if err != nil {
		log.Fatal(err)
	}
	// 인코딩된 데이터를 출력 (hex로 변환)
	fmt.Printf("ABI Encode DepositInnerData IsNative false : 0x%x\n", encodedData)
	fmt.Printf("Tatal Fee Amount : %d\n", value)
}

func Test_EncodeOutboundTransferFunc(t *testing.T) {
	client, err := ethclient.Dial("http://localhost:8547")
	if err != nil {
		log.Fatal(err)
	}

	router, err := bridgelib.NewL1GatewayRouter(client, common.HexToAddress("0xE33F71590e7307Cc003C46EC1ae78A6b1D0E2528"))
	if err != nil {
		log.Fatal(err)
	}

	_token := common.HexToAddress("0x520bBaff9939372d64Ba1E6f0483dA23eB22700e")
	_to := common.HexToAddress("0x07c9bf6399012d3dfe6bb878733d4d6426f9dfe0")
	_amount := big.NewInt(1000000000000000000)
	_maxGas := big.NewInt(79138)
	_gasPriceBid := big.NewInt(600000000)

	encodedData, _, err := bridgelib.GetDepositInnerData(big.NewInt(0), _maxGas, _gasPriceBid, true)
	if err != nil {
		log.Fatal(err)
	}

	params := bridgelib.OutboundTransferPrams{
		Token:       _token,
		To:          _to,
		Amount:      _amount,
		MaxGas:      _maxGas,
		GasPriceBid: _gasPriceBid,
		Data:        encodedData,
	}

	calldata, err := router.EncodeOutboundTransferFunc(&params)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n OutboundTransfer Calldata : %s\n", common.Bytes2Hex(calldata))
}
