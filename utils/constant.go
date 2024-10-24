package utils

import "math/big"

var (
	ErrorTrigger_GasLimit          = big.NewInt(1)
	ErrorTrigger_MaxFeePerGas      = big.NewInt(1)
	ErrorTrigger_MaxSubmissionCost = big.NewInt(1)
)

var (
	DEFAULT_SUBMISSION_FEE_PERCENT_INCREASE = big.NewInt(300)
	DEFAULT_GAS_PRICE_PERCENT_INCREASE      = big.NewInt(500)
)

var (
	MAX_EVENT_BLOCK = uint64(5000)
)
