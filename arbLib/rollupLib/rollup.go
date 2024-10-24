package rolluplib

import (
	"context"
	"fmt"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	ethlib "github.com/imelon2/orbit-cli/ethLib"
	"github.com/imelon2/orbit-cli/solgen/go/rollupgen"
	"github.com/imelon2/orbit-cli/utils"
)

type Rollup struct {
	Rollup    *rollupgen.RollupCore
	RollupRaw *bind.BoundContract
	Client    *ethclient.Client
}

type NodeConfirmedEvent struct {
	NodeNum         uint64
	BlockHash       [32]byte
	SendRoot        [32]byte
	TransactionHash *common.Hash
}

type NodeCreatedEvent struct {
	NodeNum        uint64
	ParentNodeHash [32]byte
	NodeHash       [32]byte
	ExecutionHash  [32]byte
	// Assertion          rollupgen.Assertion
	AfterInboxBatchAcc [32]byte
	WasmModuleRoot     [32]byte
	InboxMaxCount      *big.Int
	TransactionHash    *common.Hash
}

func NewRollup(client *ethclient.Client, addr common.Address) (Rollup, error) {
	rollup, err := rollupgen.NewRollupCore(addr, client)

	if err != nil {
		return Rollup{}, fmt.Errorf("failed new Rollup : %d", err)
	}
	parsed, err := rollupgen.RollupCoreMetaData.GetAbi()
	if err != nil {
		return Rollup{}, fmt.Errorf("failed get Rollup abi : %d", err)
	}

	bound := bind.NewBoundContract(addr, *parsed, client, client, client)
	if err != nil {
		return Rollup{}, fmt.Errorf("failed new Rollup Handler : %d", err)
	}

	return Rollup{
		Rollup:    rollup,
		RollupRaw: bound,
		Client:    client,
	}, nil
}

func (r Rollup) LatestConfirmed() (uint64, error) {
	return r.Rollup.LatestConfirmed(ethlib.Callopts)
}

func (r Rollup) LatestNodeCreated() (uint64, error) {
	return r.Rollup.LatestNodeCreated(ethlib.Callopts)
}

func (r Rollup) GetNodeCreated(count *big.Int) ([]NodeCreatedEvent, error) {

	events := make([]NodeCreatedEvent, 0)
	num, err := r.Client.BlockNumber(context.Background())
	if err != nil {
		return nil, fmt.Errorf("fail get block number : %s", err)
	}

	from, to := uint64(0), num

	for len(events) < int(count.Int64()) {
		if to >= utils.MAX_EVENT_BLOCK {
			from = to - utils.MAX_EVENT_BLOCK
		} else {
			from = 0
		}

		fmt.Printf("Search Batch Tx from L2 Block %d ~ %d\n", from, to)
		opt := bind.FilterOpts{
			Start:   from,
			End:     &to,
			Context: nil,
		}

		iterator, err := r.Rollup.RollupCoreFilterer.FilterNodeCreated(&opt, nil, nil, nil)
		if err != nil {
			return nil, fmt.Errorf("fail FilterSequencerBatchDelivered : %s", err)
		}

		_events := make([]NodeCreatedEvent, 0)
		for iterator.Next() {
			event := iterator.Event
			e := NodeCreatedEvent{
				NodeNum:        event.NodeNum,
				ParentNodeHash: event.ParentNodeHash,
				NodeHash:       event.NodeHash,
				ExecutionHash:  event.ExecutionHash,
				// Assertion:          event.Assertion,
				AfterInboxBatchAcc: event.AfterInboxBatchAcc,
				WasmModuleRoot:     event.WasmModuleRoot,
				InboxMaxCount:      event.InboxMaxCount,
				TransactionHash:    &event.Raw.TxHash,
			}

			_events = append(_events, e)
		}

		sort.Slice(_events, func(i, j int) bool {
			return _events[i].NodeNum > _events[j].NodeNum
		})

		for _, data := range _events {
			events = append(events, data)
			if len(events) == int(count.Int64()) {
				break
			}
		}
		to = to - utils.MAX_EVENT_BLOCK
	}

	return events, nil
}

func (r Rollup) GetNodeConfirmed(count *big.Int) ([]NodeConfirmedEvent, error) {

	events := make([]NodeConfirmedEvent, 0)
	num, err := r.Client.BlockNumber(context.Background())
	if err != nil {
		return nil, fmt.Errorf("fail get block number : %s", err)
	}

	from, to := uint64(0), num

	for len(events) < int(count.Int64()) {
		if to >= utils.MAX_EVENT_BLOCK {
			from = to - utils.MAX_EVENT_BLOCK
		} else {
			from = 0
		}

		fmt.Printf("Search Batch Tx from L2 Block %d ~ %d\n", from, to)
		opt := bind.FilterOpts{
			Start:   from,
			End:     &to,
			Context: nil,
		}

		iterator, err := r.Rollup.RollupCoreFilterer.FilterNodeConfirmed(&opt, nil)
		if err != nil {
			return nil, fmt.Errorf("fail FilterSequencerBatchDelivered : %s", err)
		}

		_events := make([]NodeConfirmedEvent, 0)
		for iterator.Next() {
			event := iterator.Event
			e := NodeConfirmedEvent{
				NodeNum:         event.NodeNum,
				BlockHash:       event.BlockHash,
				SendRoot:        event.SendRoot,
				TransactionHash: &event.Raw.TxHash,
			}

			_events = append(_events, e)
		}

		sort.Slice(_events, func(i, j int) bool {
			return _events[i].NodeNum > _events[j].NodeNum
		})

		for _, data := range _events {
			events = append(events, data)
			if len(events) == int(count.Int64()) {
				break
			}
		}
		to = to - utils.MAX_EVENT_BLOCK
	}

	return events, nil
}
