package bridgelib

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	arboslib "github.com/imelon2/orbit-cli/arbLib/arbosLib"
	rolluplib "github.com/imelon2/orbit-cli/arbLib/rollupLib"
	ethlib "github.com/imelon2/orbit-cli/ethLib"
	"github.com/imelon2/orbit-cli/retryable"
	"github.com/imelon2/orbit-cli/solgen/go/gatewaygen"
	"github.com/imelon2/orbit-cli/utils"
	arbmath "github.com/imelon2/orbit-cli/utils/arbMath"
)

type Router struct {
	Router    *gatewaygen.L1GatewayRouter
	RouterRaw *bind.BoundContract
	Client    *ethclient.Client
}

type OutboundTransferPrams struct {
	Token       common.Address
	To          common.Address
	Amount      *big.Int
	MaxGas      *big.Int
	GasPriceBid *big.Int
	Data        []byte
}

func NewL1GatewayRouter(client *ethclient.Client, addr common.Address) (Router, error) {
	router, err := gatewaygen.NewL1GatewayRouter(addr, client)

	if err != nil {
		return Router{}, fmt.Errorf("failed new router : %d", err)
	}

	parsed, err := gatewaygen.L1GatewayRouterMetaData.GetAbi()
	if err != nil {
		return Router{}, fmt.Errorf("failed get GatewayRouter abi : %d", err)
	}

	bound := bind.NewBoundContract(addr, *parsed, client, client, client)
	if err != nil {
		return Router{}, fmt.Errorf("failed new routerHandler : %d", err)
	}

	return Router{
		Router:    router,
		Client:    client,
		RouterRaw: bound,
	}, nil
}

func (r Router) GetGateway(_token common.Address) (common.Address, error) {
	return r.Router.GetGateway(ethlib.Callopts, _token)
}

func (r Router) OutboundTransfer(params *OutboundTransferPrams, auth *bind.TransactOpts) (*types.Transaction, error) {
	return r.Router.OutboundTransfer(auth, params.Token, params.To, params.Amount, params.MaxGas, params.GasPriceBid, params.Data)
}

func (r Router) EncodeOutboundTransferFunc(params *OutboundTransferPrams) ([]byte, error) {
	parsed, err := gatewaygen.L1GatewayRouterMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed get GatewayRouter abi : %d", err)
	}
	calldataLib := ethlib.NewCalldata(parsed, nil)
	return calldataLib.Abi.Pack("outboundTransfer", params.Token, params.To, params.Amount, params.MaxGas, params.GasPriceBid, params.Data)
}

func (r Router) DepositFuncCall(result *[]interface{}, opts *bind.CallOpts, method string, params *OutboundTransferPrams) error {
	return r.RouterRaw.Call(opts, result, method, params.Token, params.To, params.Amount, params.MaxGas, params.GasPriceBid, params.Data)
}

// ----------------------------------------------------------------------------------------------- //
// ------------------------------------------- Library ------------------------------------------- //
// ----------------------------------------------------------------------------------------------- //

func GetDepositInnerData(maxSubmissionCost *big.Int, gasLimit *big.Int, maxFeePerGas *big.Int, isNativeToken bool) ([]byte, *big.Int, error) {
	uint256Type, err := abi.NewType("uint256", "", nil)
	if err != nil {
		return nil, nil, err
	}

	bytesType, err := abi.NewType("bytes", "", nil)
	if err != nil {
		return nil, nil, err
	}

	value := new(big.Int).Set(gasLimit)
	value.Mul(value, maxFeePerGas)
	value.Add(value, maxSubmissionCost)

	if isNativeToken {
		arguments := abi.Arguments{
			{Type: uint256Type},
			{Type: bytesType},
			{Type: uint256Type},
		}

		encodedData, err := arguments.Pack(maxSubmissionCost, []byte(""), value)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to ABI encode DepositInnerData: %v", err)
		}

		return encodedData, value, err
	} else {
		arguments := abi.Arguments{
			{Type: uint256Type},
			{Type: bytesType},
		}
		encodedData, err := arguments.Pack(maxSubmissionCost, []byte(""))
		if err != nil {
			return nil, nil, fmt.Errorf("failed to ABI encode DepositInnerData: %v", err)
		}

		return encodedData, value, err

	}
}

func GetOutboundTransferPrams(token common.Address, to common.Address, amount *big.Int, childGasPrice *big.Int, l1Router *Router, inbox *rolluplib.Inbox, bridge *rolluplib.ERC20Bridge, nodeInterface *arboslib.NodeInterface, auth *bind.TransactOpts) (*OutboundTransferPrams, error) {
	orgin := auth.NoSend
	defer func() {
		auth.NoSend = orgin
	}()

	isFeeToken := false
	if feeToken, _ := bridge.GetFeeToken(); feeToken.Hex() != common.HexToAddress("0x00").Hex() {
		isFeeToken = true
	}

	zombieEncodedData, _, err := GetDepositInnerData(big.NewInt(0), utils.ErrorTrigger_MaxFeePerGas, utils.ErrorTrigger_GasLimit, isFeeToken)
	if err != nil {
		return nil, err
	}

	zombieParams := OutboundTransferPrams{
		Token:       token,
		To:          to,
		Amount:      amount,
		MaxGas:      utils.ErrorTrigger_MaxFeePerGas,
		GasPriceBid: utils.ErrorTrigger_GasLimit,
		Data:        zombieEncodedData,
	}

	var retryableData retryable.RetryableData
	auth.NoSend = true // eth_call
	if !isFeeToken {
		auth.Value = big.NewInt(2)
	}

	_, errData := l1Router.OutboundTransfer(&zombieParams, auth)
	rpcErr, ok := errData.(rpc.DataError)
	if ok {
		errorData := rpcErr.ErrorData()
		stringData := errorData.(string)
		bytesData, _ := hex.DecodeString(stringData[2:])

		_retryable, err := retryable.ParseRetryable(bytesData)
		if err != nil {
			return nil, err
		}

		retryableData = _retryable
	} else {
		return nil, fmt.Errorf("Failed to call contract: %v\n", err)
	}

	block, err := inbox.Client.BlockByNumber(context.Background(), nil /* Latest */)
	if err != nil {
		return nil, err
	}

	submissionFee, err := inbox.EstimateSubmissionFee(new(big.Int).SetInt64(int64(len(retryableData.Data))), block.BaseFee())
	if err != nil {
		return nil, err
	}

	maxSubmissionFee := arbmath.PercentIncrease(submissionFee, utils.DEFAULT_SUBMISSION_FEE_PERCENT_INCREASE)
	maxFeePerGas := arbmath.PercentIncrease(childGasPrice, utils.DEFAULT_GAS_PRICE_PERCENT_INCREASE)

	gasLimit, err := nodeInterface.EstimateRetryableTicket(retryableData, nil)
	if err != nil {
		return nil, err
	}
	maxGsg := new(big.Int).SetUint64(gasLimit)
	realEncodedData, value, err := GetDepositInnerData(maxSubmissionFee, maxGsg, maxFeePerGas, isFeeToken)

	if err != nil {
		return nil, err
	}

	if !isFeeToken {
		auth.Value = value
	}

	params := OutboundTransferPrams{
		Token:       token,
		To:          to,
		Amount:      amount,
		MaxGas:      maxGsg,
		GasPriceBid: maxFeePerGas,
		Data:        realEncodedData,
	}

	return &params, nil
}
