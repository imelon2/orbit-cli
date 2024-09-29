package arblib

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/imelon2/orbit-cli/solgen/go/gatewaygen"
)

type Router struct {
	Router *gatewaygen.L1GatewayRouter
	Client *ethclient.Client
	Auth   *bind.TransactOpts
}

type OutboundTransferPrams struct {
	Token       common.Address
	To          common.Address
	Amount      *big.Int
	MaxGas      *big.Int
	GasPriceBid *big.Int
	Data        []byte
}

func NewRouter(client *ethclient.Client, addr common.Address) (Router, error) {
	router, err := gatewaygen.NewL1GatewayRouter(addr, client)
	if err != nil {
		return Router{}, fmt.Errorf("failed new router : %d", err)
	}

	return Router{
		Router: router,
		Client: client,
	}, nil
}

func (r Router) DepositFunc(params *OutboundTransferPrams) (*types.Transaction, error) {
	r.Auth.NoSend = true // for estimate
	return r.Router.OutboundTransfer(r.Auth, params.Token, params.To, params.Amount, params.MaxGas, params.GasPriceBid, params.Data)
}
