package arboslib

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/imelon2/orbit-cli/retryable"
	"github.com/imelon2/orbit-cli/solgen/go/node_interfacegen"
	arbmath "github.com/imelon2/orbit-cli/utils/arbMath"
)

type NodeInterface struct {
	NodeInterface    *node_interfacegen.NodeInterface
	NodeInterfaceRaw *bind.BoundContract
	Client           *ethclient.Client
	Address          *common.Address
}

func NewNodeInterface(client *ethclient.Client) (NodeInterface, error) {
	nodeInterface, err := node_interfacegen.NewNodeInterface(types.NodeInterfaceAddress, client)
	if err != nil {
		return NodeInterface{}, fmt.Errorf("failed new NodeInterface : %d", err)
	}

	parsed, err := node_interfacegen.NodeInterfaceMetaData.GetAbi()
	if err != nil {
		return NodeInterface{}, fmt.Errorf("failed get NodeInterface abi : %d", err)
	}

	bound := bind.NewBoundContract(types.NodeInterfaceAddress, *parsed, client, client, client)
	if err != nil {
		return NodeInterface{}, fmt.Errorf("failed new NodeInterface Handler : %d", err)
	}

	return NodeInterface{
		NodeInterface:    nodeInterface,
		NodeInterfaceRaw: bound,
		Client:           client,
		Address:          &types.NodeInterfaceAddress,
	}, nil
}

func (n NodeInterface) EstimateRetryableTicket(retryableData retryable.RetryableData, senderDeposit *big.Int) (uint64, error) {
	deposit := big.NewInt(0)
	if senderDeposit == nil {
		wei, _ := arbmath.ParseEther("100")
		deposit.Add(deposit, wei)
		deposit.Add(deposit, retryableData.L2CallValue)
	} else {
		deposit.Add(deposit, senderDeposit)
	}

	parsed, _ := node_interfacegen.NodeInterfaceMetaData.GetAbi()
	// if err != nil {
	// 	return nil, fmt.Errorf("failed get GatewayRouter abi : %d", err)
	// }

	calldata, err := parsed.Pack("estimateRetryableTicket", retryableData.From, deposit, retryableData.To, retryableData.L2CallValue, retryableData.ExcessFeeRefundAddress, retryableData.CallValueRefundAddress, retryableData.Data)
	if err != nil {
		return 0, fmt.Errorf("failed get encoded EstimateRetryableTicket Calldata : %d", err)
	}

	msg := ethereum.CallMsg{
		To:   n.Address,
		Data: calldata,
	}

	return n.Client.EstimateGas(context.Background(), msg)
}
