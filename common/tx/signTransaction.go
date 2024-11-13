package tx

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GetAuthByKeystore(ks *keystore.KeyStore, account accounts.Account, client *ethclient.Client) (*bind.TransactOpts, error) {
	newContext := context.Background()
	chainID, err := client.NetworkID(newContext)
	if err != nil {
		return nil, fmt.Errorf("failed get Network ID %v", err)
	}

	auth, err := bind.NewKeyStoreTransactorWithChainID(ks, account, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed NewKeyStoreTransactorWithChainID %v", err)
	}

	nonce, err := client.PendingNonceAt(newContext, account.Address)
	if err != nil {
		return nil, fmt.Errorf("failed PendingNonceAt %v", err)
	}

	gasPrice, err := client.SuggestGasPrice(newContext)
	if err != nil {
		return nil, fmt.Errorf("failed SuggestGasPrice %v", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasPrice = gasPrice.Mul(gasPrice, big.NewInt(2))
	return auth, nil
}
