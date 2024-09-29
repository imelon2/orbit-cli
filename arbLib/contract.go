package arblib

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/imelon2/orbit-cli/contractgen"
	ethlib "github.com/imelon2/orbit-cli/ethLib"
	"github.com/imelon2/orbit-cli/solgen/go/bridgegen"
)

type ArbContract struct {
	Contracts contractgen.NetworkInfo
	client    *ethclient.Client
}

func NewContractLib() {

}

func (arb ArbContract) newInbox() (*bridgegen.Inbox, error) {
	return bridgegen.NewInbox(arb.Contracts.EthBridge.Inbox, arb.client)
}

func (arb ArbContract) EstimateSubmissionFee(dataLength *big.Int, baseFee *big.Int) (*big.Int, error) {
	inbox, err := arb.newInbox()
	if err != nil {
		return nil, fmt.Errorf("failed bind inbox contract: %d", err)
	}

	return inbox.CalculateRetryableSubmissionFee(ethlib.Callopts, dataLength, baseFee)
}
