package rolluplib

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	ethlib "github.com/imelon2/orbit-cli/ethLib"
	"github.com/imelon2/orbit-cli/solgen/go/bridgegen"
)

type Bridge struct {
	Bridge    *bridgegen.Bridge
	BridgeRaw *bind.BoundContract
	Client    *ethclient.Client
}

type ERC20Bridge struct {
	ERC20Bridge    *bridgegen.ERC20Bridge
	ERC20BridgeRaw *bind.BoundContract
	Client         *ethclient.Client
}

func NewBridge(client *ethclient.Client, addr common.Address) (Bridge, error) {
	bridge, err := bridgegen.NewBridge(addr, client)

	if err != nil {
		return Bridge{}, fmt.Errorf("failed new Bridge : %d", err)
	}
	parsed, err := bridgegen.BridgeMetaData.GetAbi()
	if err != nil {
		return Bridge{}, fmt.Errorf("failed get Bridge abi : %d", err)
	}

	bound := bind.NewBoundContract(addr, *parsed, client, client, client)
	if err != nil {
		return Bridge{}, fmt.Errorf("failed new Bridge Handler : %d", err)
	}

	return Bridge{
		Bridge:    bridge,
		BridgeRaw: bound,
		Client:    client,
	}, nil
}

func NewERC20Bridge(client *ethclient.Client, addr common.Address) (ERC20Bridge, error) {
	bridge, err := bridgegen.NewERC20Bridge(addr, client)

	if err != nil {
		return ERC20Bridge{}, fmt.Errorf("failed new ERC20Bridge : %d", err)
	}
	parsed, err := bridgegen.ERC20BridgeMetaData.GetAbi()
	if err != nil {
		return ERC20Bridge{}, fmt.Errorf("failed get ERC20Bridge abi : %d", err)
	}

	bound := bind.NewBoundContract(addr, *parsed, client, client, client)
	if err != nil {
		return ERC20Bridge{}, fmt.Errorf("failed new ERC20Bridge Handler : %d", err)
	}

	return ERC20Bridge{
		ERC20Bridge:    bridge,
		ERC20BridgeRaw: bound,
		Client:         client,
	}, nil
}

func (b ERC20Bridge) GetFeeToken() (common.Address, error) {
	return b.ERC20Bridge.ERC20BridgeCaller.NativeToken(ethlib.Callopts)
}

func (b ERC20Bridge) SequencerMessageCount() (*big.Int, error) {
	return b.ERC20Bridge.ERC20BridgeCaller.SequencerMessageCount(ethlib.Callopts)
}
