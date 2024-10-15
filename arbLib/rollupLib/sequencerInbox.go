package rolluplib

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/imelon2/orbit-cli/solgen/go/bridgegen"
)

type SequencerBatchDeliveredEvent struct {
	BatchSequenceNumber *big.Int
	TransactionHash     *common.Hash
	BeforeAcc           [32]byte
	AfterAcc            [32]byte
	DelayedAcc          [32]byte
}

type SequencerInbox struct {
	SequencerInbox    *bridgegen.SequencerInbox
	SequencerInboxRaw *bind.BoundContract
	Address           common.Address
	Client            *ethclient.Client
}

func NewSequencerInbox(client *ethclient.Client, addr common.Address) (SequencerInbox, error) {
	sequencerInbox, err := bridgegen.NewSequencerInbox(addr, client)
	if err != nil {
		return SequencerInbox{}, fmt.Errorf("failed new SequencerInbox : %d", err)
	}
	parsed, err := bridgegen.SequencerInboxMetaData.GetAbi()
	if err != nil {
		return SequencerInbox{}, fmt.Errorf("failed get SequencerInbox abi : %d", err)
	}

	bound := bind.NewBoundContract(addr, *parsed, client, client, client)
	if err != nil {
		return SequencerInbox{}, fmt.Errorf("failed new SequencerInbox Handler : %d", err)
	}

	return SequencerInbox{
		SequencerInbox:    sequencerInbox,
		SequencerInboxRaw: bound,
		Address:           addr,
		Client:            client,
	}, nil
}

func (si SequencerInbox) GetBatchData(count *big.Int) ([]SequencerBatchDeliveredEvent, error) {

	MAX_EVENT_BLOCK := uint64(5000)

	batchTransactions := make([]SequencerBatchDeliveredEvent, 0)

	num, err := si.Client.BlockNumber(context.Background())
	if err != nil {
		return nil, fmt.Errorf("fail get block number : %s", err)
	}

	from, to := uint64(0), num

	for len(batchTransactions) < int(count.Int64()) {
		if to >= MAX_EVENT_BLOCK {
			from = to - MAX_EVENT_BLOCK
		} else {
			from = 0
		}

		fmt.Printf("Search Batch Tx from L2 Block %d ~ %d\n", from, to)
		opt := bind.FilterOpts{
			Start:   from,
			End:     &to,
			Context: nil,
		}

		iterator, err := si.SequencerInbox.SequencerInboxFilterer.FilterSequencerBatchDelivered(&opt, nil, nil, nil)
		if err != nil {
			return nil, fmt.Errorf("fail FilterSequencerBatchDelivered : %s", err)
		}

		for iterator.Next() {
			event := iterator.Event
			e := SequencerBatchDeliveredEvent{
				BatchSequenceNumber: event.BatchSequenceNumber,
				TransactionHash:     &event.Raw.TxHash,
				BeforeAcc:           event.BeforeAcc,
				AfterAcc:            event.AfterAcc,
				DelayedAcc:          event.DelayedAcc,
			}

			batchTransactions = append(batchTransactions, e)

			if len(batchTransactions) == int(count.Int64()) {
				break
			}
		}

		to = to - MAX_EVENT_BLOCK
	}

	return batchTransactions, nil
}
