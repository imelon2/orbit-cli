package arbnetwork

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/imelon2/orbit-cli/common/path"
	"github.com/imelon2/orbit-cli/solgen/go/bridgegen"
	"github.com/imelon2/orbit-cli/solgen/go/rollupgen"
)

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

func genNetworkInfo() (map[string]NetworkInfo, error) {
	networkDir := path.GetContractNetworkDir()
	files, err := os.ReadDir(networkDir)

	if err != nil {
		return nil, fmt.Errorf("failed to open JSON file: %s", err)
	}

	networks := make(map[string]NetworkInfo)

	for _, file := range files {
		jsonFile := filepath.Join(networkDir, file.Name())
		jsonData, err := os.Open(jsonFile)
		if err != nil {
			return nil, fmt.Errorf("failed to open JSON file: %s", err)
		}

		byteValue, err := io.ReadAll(jsonData)
		if err != nil {
			return nil, fmt.Errorf("failed to read JSON file: %s", err)
		}

		var networkInfo NetworkInfo

		err = json.Unmarshal(byteValue, &networkInfo)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal network info JSON: %s", err)
		}

		networks[file.Name()] = networkInfo
	}

	return networks, nil
}

func GetNetworkInfo(childClient *ethclient.Client) (*NetworkInfo, error) {
	networks, err := genNetworkInfo()
	if err != nil {
		return nil, err
	}

	chainID, err := childClient.NetworkID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed get Network ID %v", err)
	}

	for val := range networks {
		if network := networks[val]; network.ChainId == int(chainID.Int64()) {
			return &network, nil
		}
	}

	return nil, fmt.Errorf("there is no network info for chain id %d : %v", chainID, err)
}

func (network NetworkInfo) NewInbox(parentClient *ethclient.Client) (*bridgegen.Inbox, error) {
	return bridgegen.NewInbox(network.EthBridge.Inbox, parentClient)
}

func (network NetworkInfo) NewRollupCore(parentClient *ethclient.Client) (*rollupgen.RollupCore, error) {
	return rollupgen.NewRollupCore(network.EthBridge.Rollup, parentClient)
}

func (network NetworkInfo) NewRollupAdminLogic(parentClient *ethclient.Client) (*rollupgen.RollupAdminLogic, error) {
	return rollupgen.NewRollupAdminLogic(network.EthBridge.Rollup, parentClient)
}

func (network NetworkInfo) NewUpgradeExecutor(parentClient *ethclient.Client) (*UpgradeExecutor, error) {
	return NewUpgradeExecutor(network, parentClient)
}

func (network NetworkInfo) NewBridge(parentClient *ethclient.Client) (*bridgegen.Bridge, error) {
	return bridgegen.NewBridge(network.EthBridge.Bridge, parentClient)
}

func (network NetworkInfo) NewSequencerInbox(parentClient *ethclient.Client) (*bridgegen.SequencerInbox, error) {
	return bridgegen.NewSequencerInbox(network.EthBridge.SequencerInbox, parentClient)
}
