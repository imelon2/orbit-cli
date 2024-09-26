package ethlib

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/imelon2/orbit-cli/prompt"
)

var L1PricerFundsPool = common.HexToAddress("0xa4B00000000000000000000000000000000000F6")

var Callopts = &bind.CallOpts{
	Pending: false, // 트랜잭션이 확정된 상태를 조회
	Context: nil,   // 컨텍스트가 필요한 경우 (예: 시간 초과)
}

func GenerateAuth() (*ethclient.Client, *bind.TransactOpts, error) {

	provider, err := prompt.SelectProvider()
	if err != nil {
		return nil, nil, fmt.Errorf("GenerateAuth - SelectProvider failed %v\n", err)
	}

	client, err := ethclient.Dial(provider)
	if err != nil {
		return nil, nil, fmt.Errorf("GenerateAuth - ethclient failed %v\n", err)
	}

	_, ks, account, err := prompt.SelectWalletForSign()
	if err != nil {
		return nil, nil, fmt.Errorf("GenerateAuth - SelectWalletForSign failed %v\n", err)
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("GenerateAuth - NetworkID failed %v\n", err)
	}

	auth, err := bind.NewKeyStoreTransactorWithChainID(ks, account, chainID)
	if err != nil {
		return nil, nil, fmt.Errorf("GenerateAuth - NewKeyStoreTransactorWithChainID failed %v\n", err)
	}

	nonce, err := client.PendingNonceAt(context.Background(), account.Address)
	if err != nil {
		return nil, nil, fmt.Errorf("GenerateAuth - PendingNonceAt failed %v\n", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("GenerateAuth - SuggestGasPrice failed %v\n", err)
	}

	auth.GasPrice = gasPrice

	return client, auth, nil
}

func AsMessage(tx *types.Transaction, from common.Address) ethereum.CallMsg {
	return ethereum.CallMsg{
		From:      from,
		To:        tx.To(),
		Gas:       tx.Gas(),
		GasPrice:  tx.GasPrice(),
		GasFeeCap: tx.GasFeeCap(),
		GasTipCap: tx.GasTipCap(),
		Value:     tx.Value(),
		Data:      tx.Data(),
	}
}
