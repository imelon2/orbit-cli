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

type TransactionHash struct {
	client      *ethclient.Client
	transaction *types.Transaction
	hash        common.Hash
}

func NewTransactionHash(client *ethclient.Client, hash common.Hash) *TransactionHash {
	newTransactionHash := new(TransactionHash)
	newTransactionHash.hash = hash
	newTransactionHash.client = client
	return newTransactionHash
}

func NewTransaction(client *ethclient.Client, tx *types.Transaction) *TransactionHash {
	newTx := new(TransactionHash)
	newTx.transaction = tx
	newTx.client = client
	return newTx
}

func (txHash TransactionHash) Wait() (*types.Receipt, error) {
	return bind.WaitMined(context.Background(), txHash.client, txHash.transaction)
}

func (txHash TransactionHash) GetTransactionByHash() (*types.Transaction, bool, error) {
	tx, isPending, err := txHash.client.TransactionByHash(context.Background(), txHash.hash)
	if err == ethereum.NotFound {
		return nil, isPending, fmt.Errorf("transaction hash %v not found: %d", txHash.hash.Hex(), err)
	} else if err != nil {
		return nil, isPending, fmt.Errorf("failed to get TransactionByHash: %d", err)
	}
	return tx, isPending, nil
}

func (txHash TransactionHash) GetTransactionSender() (*types.Transaction, *types.Receipt, *common.Address, error) {
	tx, _, err := txHash.GetTransactionByHash()
	if err != nil {
		return nil, nil, nil, err
	}

	receipt, err := txHash.GetTransactionReceipt()
	if err != nil {
		return nil, nil, nil, err
	}

	sender, err := txHash.client.TransactionSender(context.Background(), tx, receipt.BlockHash, receipt.TransactionIndex)
	if err == ethereum.NotFound {
		return nil, nil, nil, fmt.Errorf("transaction hash %v not found: %d", txHash.hash.Hex(), err)
	} else if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get TransactionByHash: %d", err)
	}
	return tx, receipt, &sender, nil
}

func (txHash TransactionHash) GetTransactionReceipt() (*types.Receipt, error) {
	tx, err := txHash.client.TransactionReceipt(context.Background(), txHash.hash)
	if err == ethereum.NotFound {
		return nil, fmt.Errorf("transaction hash %v not found: %d", txHash.hash.Hex(), err)
	} else if err != nil {
		return nil, fmt.Errorf("failed to get TransactionReceipt: %d", err)
	}
	return tx, nil
}

func GenerateAuth() (*ethclient.Client, *bind.TransactOpts, string, error) {

	provider, err := prompt.SelectProvider()
	if err != nil {
		return nil, nil, "", fmt.Errorf("GenerateAuth : %v\n", err)
	}

	client, err := ethclient.Dial(provider)
	if err != nil {
		return nil, nil, "", fmt.Errorf("GenerateAuth : %v\n", err)
	}

	_, ks, account, err := prompt.SelectWalletForSign()
	if err != nil {
		return nil, nil, "", fmt.Errorf("GenerateAuth : %v\n", err)
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, nil, "", fmt.Errorf("GenerateAuth - NetworkID failed %v\n", err)
	}

	auth, err := bind.NewKeyStoreTransactorWithChainID(ks, account, chainID)
	if err != nil {
		return nil, nil, "", fmt.Errorf("GenerateAuth - NewKeyStoreTransactorWithChainID failed %v\n", err)
	}

	nonce, err := client.PendingNonceAt(context.Background(), account.Address)
	if err != nil {
		return nil, nil, "", fmt.Errorf("GenerateAuth - PendingNonceAt failed %v\n", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, nil, "", fmt.Errorf("GenerateAuth - SuggestGasPrice failed %v\n", err)
	}

	auth.GasPrice = gasPrice
	return client, auth, provider, nil
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
