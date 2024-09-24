package ethlib

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var Callopts = &bind.CallOpts{
	Pending: false, // 트랜잭션이 확정된 상태를 조회
	Context: nil,   // 컨텍스트가 필요한 경우 (예: 시간 초과)
}

func GenerateAuth(client *ethclient.Client, keystore *keystore.KeyStore, account accounts.Account) *bind.TransactOpts {
	chainID, err := client.NetworkID(context.Background())
	nonce, err := client.PendingNonceAt(context.Background(), account.Address)
	if err != nil {
		fmt.Print("nonce Error")
		log.Fatal(err)
	}
	auth, err := bind.NewKeyStoreTransactorWithChainID(keystore, account, chainID)
	auth.Nonce = big.NewInt(int64(nonce))
	// gasPrice, err := client.SuggestGasPrice(context.Background())
	// auth.GasPrice = gasPrice

	return auth
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
