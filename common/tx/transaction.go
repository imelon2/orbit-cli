package tx

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/imelon2/orbit-cli/common/logs"
)

var (
	MAX_EVENT_BLOCK = uint64(5000)
)

type TxLib struct {
	Client *ethclient.Client
	Hash   *common.Hash
}

func NewTxLib(client *ethclient.Client, hash *common.Hash) *TxLib {
	return &TxLib{
		Client: client,
		Hash:   hash,
	}
}

func (t *TxLib) GetTransactionByHash() (*types.Transaction, bool, error) {
	tx, isPending, err := t.Client.TransactionByHash(context.Background(), *t.Hash)
	if err == ethereum.NotFound {
		return nil, isPending, fmt.Errorf("get transaction hash %s not found: %d", t.Hash.Hex(), err)
	} else if err != nil {
		return nil, isPending, fmt.Errorf("failed to get TransactionByHash: %d", err)
	}
	return tx, isPending, nil
}

func (t *TxLib) GetTransactionReceipt() (*types.Receipt, error) {
	tx, err := t.Client.TransactionReceipt(context.Background(), *t.Hash)
	if err == ethereum.NotFound {
		return nil, fmt.Errorf("get transaction receipt hash %s not found: %d", t.Hash.Hex(), err)
	} else if err != nil {
		return nil, fmt.Errorf("failed to get TransactionReceipt: %d", err)
	}
	return tx, nil
}

func (t *TxLib) GetTransactionAll() (*types.Transaction, *types.Receipt, *common.Address, error) {
	tx, _, err := t.GetTransactionByHash()
	if err != nil {
		return nil, nil, nil, err
	}

	receipt, err := t.GetTransactionReceipt()
	if err != nil {
		return nil, nil, nil, err
	}

	sender, err := t.Client.TransactionSender(context.Background(), tx, receipt.BlockHash, receipt.TransactionIndex)
	if err == ethereum.NotFound {
		return nil, nil, nil, fmt.Errorf("get transaction sender hash %s not found: %d", t.Hash.Hex(), err)
	} else if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get TransactionByHash: %d", err)
	}
	return tx, receipt, &sender, nil
}

func (t *TxLib) GetTransactionReturn() ([]byte, *rpc.DataError, *bool, error) {
	tx, receipt, sender, err := t.GetTransactionAll()
	if err != nil {
		return nil, nil, nil, err
	}

	callMsg := ethereum.CallMsg{
		From:     *sender,
		To:       tx.To(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
		Data:     tx.Data(),
		Value:    tx.Value(),
	}

	data, errData := t.Client.CallContract(context.Background(), callMsg, receipt.BlockNumber)
	rpcErr, ok := errData.(rpc.DataError)
	if !ok {
		return nil, nil, nil, fmt.Errorf("failed get rpc DataError")
	}
	status := receipt.Status == 1
	return data, &rpcErr, &status, nil
}

func SearchEvent(count int, client *ethclient.Client, eventFunc func(opt bind.FilterOpts) ([]interface{}, error)) ([]interface{}, error) {
	fmt.Print("\n")
	num, err := client.BlockNumber(context.Background())
	if err != nil {
		return nil, fmt.Errorf("fail get block number : %s", err)
	}

	result := make([]interface{}, 0)
	limit := 0
	start, end := uint64(0), num
	for len(result) < count {
		if end >= MAX_EVENT_BLOCK {
			start = end - MAX_EVENT_BLOCK
		} else {
			start = 0
		}
		limit = limit + int(end-start)
		logs.PrintBlockScope(int(num), limit)

		opt := bind.FilterOpts{
			Start:   start,
			End:     &end,
			Context: nil,
		}

		events, err := eventFunc(opt)
		if err != nil {
			return nil, fmt.Errorf("fail get event func : %s", err)
		}

		// return on count
		for _, data := range events {
			result = append(result, data)
			if len(result) == count {
				break
			}
		}

		end = end - MAX_EVENT_BLOCK
	}

	return result, nil
}
