package contractgen

import "github.com/ethereum/go-ethereum/common"

type NetworkInfo struct {
	ChainId        int         `json:"chainID"`
	PartnerChainID int         `json:"partnerChainID"`
	EthBridge      EthBridge   `json:"ethBridge"`
	TokenBridge    TokenBridge `json:"tokenBridge"`
	Teleport       Teleport    `json:"teleporter"`
}

type EthBridge struct {
	Bridge         common.Address `json:"bridge"`
	Inbox          common.Address `json:"inbox"`
	Outbox         common.Address `json:"outbox"`
	Rollup         common.Address `json:"rollup"`
	SequencerInbox common.Address `json:"sequencerInbox"`
}

type TokenBridge struct {
	L1CustomGateway common.Address `json:"l1CustomGateway"`
	L1ERC20Gateway  common.Address `json:"l1ERC20Gateway"`
	L1GatewayRouter common.Address `json:"l1GatewayRouter"`
	L1MultiCall     common.Address `json:"l1MultiCall"`
	L1ProxyAdmin    common.Address `json:"l1ProxyAdmin"`
	L1Weth          common.Address `json:"l1Weth"`
	L1WethGateway   common.Address `json:"l1WethGateway"`
	L2CustomGateway common.Address `json:"l2CustomGateway"`
	L2ERC20Gateway  common.Address `json:"l2ERC20Gateway"`
	L2GatewayRouter common.Address `json:"l2GatewayRouter"`
	L2Multicall     common.Address `json:"l2Multicall"`
	L2ProxyAdmin    common.Address `json:"l2ProxyAdmin"`
	L2Weth          common.Address `json:"l2Weth"`
	L2WethGateway   common.Address `json:"l2WethGateway"`
}

type Teleport struct {
	L1Teleporter       common.Address `json:"l1Teleporter"`
	L2ForwarderFactory common.Address `json:"l2ForwarderFactory"`
}

func NewNetworkInfo() {

}
