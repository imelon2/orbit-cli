package arbnetwork

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/imelon2/orbit-cli/solgen/go/bridgegen"
	"github.com/imelon2/orbit-cli/solgen/go/rollupgen"
	"github.com/imelon2/orbit-cli/solgen/go/srcgen"
)

type UpgradeExecutor struct {
	UpgradeExecutor  *srcgen.UpgradeExecutor
	RollupAdminLogic *rollupgen.RollupAdminLogic
	RollupUserLogic  *rollupgen.RollupUserLogic
	SequencerInbox   *bridgegen.SequencerInbox
	Bridge           *bridgegen.Bridge
}

func NewUpgradeExecutor(network NetworkInfo, client *ethclient.Client) (*UpgradeExecutor, error) {
	rollupUserLogic, err := rollupgen.NewRollupUserLogic(network.EthBridge.Rollup, client)
	if err != nil {
		return nil, err
	}
	rollupAdminLogic, err := rollupgen.NewRollupAdminLogic(network.EthBridge.Rollup, client)
	if err != nil {
		return nil, err
	}
	sequencerInbox, err := bridgegen.NewSequencerInbox(network.EthBridge.SequencerInbox, client)
	if err != nil {
		return nil, err
	}
	bridge, err := bridgegen.NewBridge(network.EthBridge.Bridge, client)
	if err != nil {
		return nil, err
	}

	Callopts := &bind.CallOpts{
		Pending: false,
		Context: nil,
	}
	executor, err := rollupUserLogic.RollupUserLogicCaller.Owner(Callopts)
	if err != nil {
		return nil, err
	}
	upgradeExecutor, err := srcgen.NewUpgradeExecutor(executor, client)
	if err != nil {
		return nil, err
	}
	return &UpgradeExecutor{
		UpgradeExecutor:  upgradeExecutor,
		RollupAdminLogic: rollupAdminLogic,
		RollupUserLogic:  rollupUserLogic,
		SequencerInbox:   sequencerInbox,
		Bridge:           bridge,
	}, nil
}

func (wallet UpgradeExecutor) SetMinimumAssertionPeriod(opts *bind.TransactOpts, newPeriod *big.Int) (*types.Transaction, error) {
	cacheGas := opts.GasLimit
	simulation := func(gasLimit uint64) {
		opts.NoSend = !opts.NoSend
		opts.GasLimit = gasLimit
	}
	simulation(1) // 1 for skip estimate gas
	txRes, _ := wallet.RollupAdminLogic.SetMinimumAssertionPeriod(opts, newPeriod)
	simulation(cacheGas)

	return wallet.UpgradeExecutor.ExecuteCall(opts, *txRes.To(), txRes.Data())
}

func (wallet UpgradeExecutor) SetConfirmPeriodBlocks(opts *bind.TransactOpts, newConfirmPeriod uint64) (*types.Transaction, error) {
	cacheGas := opts.GasLimit
	simulation := func(gasLimit uint64) {
		opts.NoSend = !opts.NoSend
		opts.GasLimit = gasLimit
	}
	simulation(1) // 1 for skip estimate gas
	txRes, _ := wallet.RollupAdminLogic.SetConfirmPeriodBlocks(opts, newConfirmPeriod)
	simulation(cacheGas)

	return wallet.UpgradeExecutor.ExecuteCall(opts, *txRes.To(), txRes.Data())
}

func (wallet UpgradeExecutor) SetMaxTimeVariation(opts *bind.TransactOpts, maxTimeVariation bridgegen.ISequencerInboxMaxTimeVariation) (*types.Transaction, error) {
	cacheGas := opts.GasLimit
	simulation := func(gasLimit uint64) {
		opts.NoSend = !opts.NoSend
		opts.GasLimit = gasLimit
	}
	simulation(1) // 1 for skip estimate gas
	txRes, _ := wallet.SequencerInbox.SetMaxTimeVariation(opts, maxTimeVariation)
	simulation(cacheGas)

	return wallet.UpgradeExecutor.ExecuteCall(opts, *txRes.To(), txRes.Data())
}

func (wallet UpgradeExecutor) SetIsBatchPoster(opts *bind.TransactOpts, addr common.Address, isBatchPoster_ bool) (*types.Transaction, error) {
	cacheGas := opts.GasLimit
	simulation := func(gasLimit uint64) {
		opts.NoSend = !opts.NoSend
		opts.GasLimit = gasLimit
	}
	simulation(1) // 1 for skip estimate gas
	txRes, _ := wallet.SequencerInbox.SetIsBatchPoster(opts, addr, isBatchPoster_)
	simulation(cacheGas)

	return wallet.UpgradeExecutor.ExecuteCall(opts, *txRes.To(), txRes.Data())
}

func (wallet UpgradeExecutor) SetSequencerReportedSubMessageCount(opts *bind.TransactOpts, newMsgCount *big.Int) (*types.Transaction, error) {
	cacheGas := opts.GasLimit
	simulation := func(gasLimit uint64) {
		opts.NoSend = !opts.NoSend
		opts.GasLimit = gasLimit
	}
	simulation(1) // 1 for skip estimate gas
	txRes, _ := wallet.Bridge.SetSequencerReportedSubMessageCount(opts, newMsgCount)
	simulation(cacheGas)

	return wallet.UpgradeExecutor.ExecuteCall(opts, *txRes.To(), txRes.Data())
}
