/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/AlecAivazis/survey/v2"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	arblib "github.com/imelon2/orbit-cli/arbLib"
	arboslib "github.com/imelon2/orbit-cli/arbLib/arbosLib"
	bridgelib "github.com/imelon2/orbit-cli/arbLib/bridgeLib"
	standardlib "github.com/imelon2/orbit-cli/arbLib/standardLib"
	"github.com/imelon2/orbit-cli/contractgen"
	ethlib "github.com/imelon2/orbit-cli/ethLib"
	"github.com/imelon2/orbit-cli/prompt"
	"github.com/imelon2/orbit-cli/utils"
	"github.com/spf13/cobra"
)

const (
	l1l2Deposit = iota
	l2l3Deposit
	l1l2Withdraw
	l2l3Withdraw
)

var setERC20BridgeCommand = []string{"l1-l2-deposit", "l2-l3-deposit", "l1-l2-withdraw", "l2-l3-withdraw"}

// erc20Cmd represents the erc20 command
var Erc20Cmd = &cobra.Command{
	Use:   "erc20",
	Short: "Deposit and withdraw erc20 asset by token bridge",
	Run: func(cmd *cobra.Command, args []string) {
		var qs = &survey.Select{
			Message: "Select Command: ",
			Options: setERC20BridgeCommand,
		}

		answerIndex := 0
		err := survey.AskOne(qs, &answerIndex)
		if err != nil {
			log.Fatal(err)
		}

		tokenAddress, err := prompt.EnterAddress("ERC20 Contract")
		if err != nil {
			log.Fatal(err)
		}

		value, err := prompt.EnterValue("deposit value")
		if err != nil {
			log.Fatal(err)
		}

		providers, err := prompt.SelectProviders()
		if err != nil {
			log.Fatal(err)
		}

		_, ks, account, err := prompt.SelectWalletForSign()
		if err != nil {
			log.Fatal(err)
		}

		l1ProviderUrl := providers[0]
		l2ProviderUrl := providers[1]
		l3ProviderUrl := providers[2]

		var (
			parentClient *ethclient.Client
			childClient  *ethclient.Client
		)
		var response *types.Transaction
		switch answerIndex {
		case l1l2Deposit:
			parentClient, err = ethclient.Dial(l1ProviderUrl)
			if err != nil {
				log.Fatal(err)
			}
			childClient, err = ethclient.Dial(l2ProviderUrl)
			if err != nil {
				log.Fatal(err)
			}
			response, err = depositFunc(common.HexToAddress(tokenAddress), value, parentClient, childClient, ks, account)
			if err != nil {
				log.Fatal(err)
			}
		case l2l3Deposit:
			parentClient, err = ethclient.Dial(l2ProviderUrl)
			if err != nil {
				log.Fatal(err)
			}
			childClient, err = ethclient.Dial(l3ProviderUrl)
			if err != nil {
				log.Fatal(err)
			}
			response, err = depositFunc(common.HexToAddress(tokenAddress), value, parentClient, childClient, ks, account)
			if err != nil {
				log.Fatal(err)
			}
		case l1l2Withdraw:
			// response, err = l1l2DepositFunc(l2ProviderUrl, l3ProviderUrl, ks, account)
		case l2l3Withdraw:
			// response, err = l1l2DepositFunc(l2ProviderUrl, l3ProviderUrl, ks, account)
		}

		transaction := ethlib.NewTransaction(parentClient, response)
		fmt.Print("\n\nTransaction Response: \n")
		utils.PrintPrettyJson(response)

		fmt.Print("\n\nWait Mined Transaction ... \n\n")

		txRes, err := transaction.Wait()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Transaction receipt: ")
		utils.PrintPrettyJson(txRes)
	},
}

func depositFunc(erc20Address common.Address, depositAmount *big.Int, parentClient *ethclient.Client, childClient *ethclient.Client, ks *keystore.KeyStore, account accounts.Account) (*types.Transaction, error) {
	// --------------------- generate auth --------------------- //
	chainID, err := parentClient.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth, err := bind.NewKeyStoreTransactorWithChainID(ks, account, chainID)
	if err != nil {
		log.Fatal(err)
	}

	network, err := contractgen.GetNetworkInfo(childClient)
	if err != nil {
		log.Fatal(err)
	}

	nonce, err := parentClient.PendingNonceAt(context.Background(), account.Address)
	if err != nil {
		log.Fatal(err)
	}
	gasPrice, err := parentClient.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasPrice = gasPrice

	childGasPrice, err := childClient.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// --------------------- generate contract libs --------------------- //
	arb := arblib.NewContractLib(&network, parentClient)
	l1Router, err := arb.NewL1GatewayRouter()
	if err != nil {
		log.Fatal(err)
	}
	inbox, err := arb.NewInbox()
	if err != nil {
		log.Fatal(err)
	}
	bridge, err := arb.NewERC20Bridge()
	if err != nil {
		log.Fatal(err)
	}
	nodeInterface, _ := arboslib.NewNodeInterface(childClient)
	erc20, err := standardlib.NewBridge(parentClient, erc20Address)
	if err != nil {
		log.Fatal(err)
	}
	// --------------------------------------------------------------------- //

	// --------------------- check Approve ERC20 Token --------------------- //
	l1Gateway, err := l1Router.GetGateway(erc20Address)
	if err != nil {
		log.Fatal(err)
	}

	allowance, err := erc20.Allowance(auth.From, l1Gateway)
	if err != nil {
		log.Fatal(err)
	}

	if allowance.Cmp(depositAmount) < 0 {
		approveRes, err := erc20.Approve(l1Gateway, abi.MaxUint256, auth)
		if err != nil {
			log.Fatal(err)
		}
		transaction := ethlib.NewTransaction(parentClient, approveRes)
		fmt.Print("\n\nWait Mined Transaction ... \n\n")

		txRes, err := transaction.Wait()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Transaction receipt: ")
		utils.PrintPrettyJson(txRes)

		nonce, err := parentClient.PendingNonceAt(context.Background(), account.Address)
		if err != nil {
			log.Fatal(err)
		}

		auth.Nonce = big.NewInt(int64(nonce))
	}
	// --------------------------------------------------------------------- //

	// ----------------------- check Approve FeeToken ---------------------- //
	if feeToken, _ := bridge.GetFeeToken(); feeToken.Hex() != common.HexToAddress("0x00").Hex() {
		feeTokenGateway, err := l1Router.GetGateway(feeToken)
		if err != nil {
			log.Fatal(err)
		}

		feeTokenErc20, err := standardlib.NewBridge(parentClient, feeToken)
		if err != nil {
			log.Fatal(err)
		}

		allowance, err := feeTokenErc20.Allowance(auth.From, feeTokenGateway)
		if err != nil {
			log.Fatal(err)
		}

		if allowance.Cmp(big.NewInt(0)) == 0 {
			approveRes, err := feeTokenErc20.Approve(feeTokenGateway, abi.MaxUint256, auth)
			if err != nil {
				log.Fatal(err)
			}
			transaction := ethlib.NewTransaction(parentClient, approveRes)
			fmt.Print("\n\nWait Mined Transaction ... \n\n")

			txRes, err := transaction.Wait()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Transaction receipt: ")
			utils.PrintPrettyJson(txRes)

			nonce, err := parentClient.PendingNonceAt(context.Background(), account.Address)
			if err != nil {
				log.Fatal(err)
			}

			auth.Nonce = big.NewInt(int64(nonce))
		}
	}
	// --------------------------------------------------------------------- //

	// --------------------- send OutboundTransfer() --------------------- //
	params, err := bridgelib.GetOutboundTransferPrams(erc20Address, auth.From, depositAmount, childGasPrice, &l1Router, &inbox, &bridge, &nodeInterface, auth)
	if err != nil {
		log.Fatal(err)
	}

	auth.GasLimit = 500000 //@TODO

	return l1Router.OutboundTransfer(params, auth)
}
