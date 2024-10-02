package standardlib

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	ethlib "github.com/imelon2/orbit-cli/ethLib"
	"github.com/imelon2/orbit-cli/solgen/go/ERC20gen"
)

type ERC20 struct {
	ERC20    *ERC20gen.ERC20
	ERC20Raw *bind.BoundContract
	Client   *ethclient.Client
}

func NewBridge(client *ethclient.Client, addr common.Address) (ERC20, error) {
	erc20, err := ERC20gen.NewERC20(addr, client)

	if err != nil {
		return ERC20{}, fmt.Errorf("failed new erc20 : %d", err)
	}
	parsed, err := ERC20gen.ERC20MetaData.GetAbi()
	if err != nil {
		return ERC20{}, fmt.Errorf("failed get erc20 abi : %d", err)
	}

	bound := bind.NewBoundContract(addr, *parsed, client, client, client)
	if err != nil {
		return ERC20{}, fmt.Errorf("failed new erc20 Handler : %d", err)
	}

	return ERC20{
		ERC20:    erc20,
		ERC20Raw: bound,
		Client:   client,
	}, nil
}

func (e ERC20) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return e.ERC20.Allowance(ethlib.Callopts, owner, spender)
}

func (e ERC20) Approve(spender common.Address, amount *big.Int, auth *bind.TransactOpts) (*types.Transaction, error) {
	return e.ERC20.Approve(auth, spender, amount)
}
