package arbnetwork

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type NodeCreatedEvent struct {
	NodeNum        uint64
	ParentNodeHash [32]byte
	NodeHash       [32]byte
	ExecutionHash  [32]byte
	// Assertion          rollupgen.Assertion
	AfterInboxBatchAcc [32]byte
	WasmModuleRoot     [32]byte
	InboxMaxCount      *big.Int
	TransactionHash    *common.Hash
}

type NodeConfirmedEvent struct {
	NodeNum         uint64
	BlockHash       [32]byte
	SendRoot        [32]byte
	TransactionHash *common.Hash
}

type SequencerBatchDeliveredEvent struct {
	BatchSequenceNumber      *big.Int
	TransactionHash          *common.Hash
	BeforeAcc                [32]byte
	AfterAcc                 [32]byte
	DelayedAcc               [32]byte
	AfterDelayedMessagesRead *big.Int
}

type MaxTimeVariation struct {
	DelayBlocks   *big.Int
	FutureBlocks  *big.Int
	DelaySeconds  *big.Int
	FutureSeconds *big.Int
}
