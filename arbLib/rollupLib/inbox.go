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

type Inbox struct {
	Inbox    *bridgegen.Inbox
	InboxRaw *bind.BoundContract
	Client   *ethclient.Client
}

func NewInbox(client *ethclient.Client, addr common.Address) (Inbox, error) {
	inbox, err := bridgegen.NewInbox(addr, client)

	if err != nil {
		return Inbox{}, fmt.Errorf("failed new inbox : %d", err)
	}
	parsed, err := bridgegen.InboxMetaData.GetAbi()
	if err != nil {
		return Inbox{}, fmt.Errorf("failed get inbox abi : %d", err)
	}

	bound := bind.NewBoundContract(addr, *parsed, client, client, client)
	if err != nil {
		return Inbox{}, fmt.Errorf("failed new Inbox Handler : %d", err)
	}

	return Inbox{
		Inbox:    inbox,
		InboxRaw: bound,
		Client:   client,
	}, nil
}

func (inbox Inbox) EstimateSubmissionFee(dataLength *big.Int, baseFee *big.Int) (*big.Int, error) {
	return inbox.Inbox.CalculateRetryableSubmissionFee(ethlib.Callopts, dataLength, baseFee)
}
