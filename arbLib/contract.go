package arblib

import (
	"github.com/ethereum/go-ethereum/ethclient"
	bridgelib "github.com/imelon2/orbit-cli/arbLib/bridgeLib"
	rolluplib "github.com/imelon2/orbit-cli/arbLib/rollupLib"
	"github.com/imelon2/orbit-cli/contractgen"
)

type ArbContract struct {
	Contracts *contractgen.NetworkInfo
	client    *ethclient.Client
}

func NewContractLib(contracts *contractgen.NetworkInfo, client *ethclient.Client) ArbContract {
	return ArbContract{
		Contracts: contracts,
		client:    client,
	}
}

func (arb ArbContract) NewInbox() (rolluplib.Inbox, error) {
	return rolluplib.NewInbox(arb.client, arb.Contracts.EthBridge.Inbox)
}

func (arb ArbContract) NewBridge() (rolluplib.Bridge, error) {
	return rolluplib.NewBridge(arb.client, arb.Contracts.EthBridge.Bridge)
}

func (arb ArbContract) NewERC20Bridge() (rolluplib.ERC20Bridge, error) {
	return rolluplib.NewERC20Bridge(arb.client, arb.Contracts.EthBridge.Bridge)
}

func (arb ArbContract) NewL1GatewayRouter() (bridgelib.Router, error) {
	return bridgelib.NewL1GatewayRouter(arb.client, arb.Contracts.TokenBridge.L1GatewayRouter)
}
