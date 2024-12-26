package common

import (
	"errors"
)

type Config map[string]int
type ConfigWeight map[string]float64
type Step map[string]int

var (
	ValidConfigName []string = []string{
		CHAINMAKER_TBFT_PROPOSE_TIMEOUT_KEY,
		CHAINMAKER_TBFT_PROPOSE_DELTA_TIMEOUT_KEY,
		CHAINMAKER_TBFT_BLOCK_PER_PROPOSER,
		CHAINMAKER_BATCH_MAX_SIZE_KEY,
		CHAINMAKER_BATCH_CREATE_TIMEOUT_KEY,
		CHAINMAKER_BLOCK_TX_CAPACITY_KEY,
		// CHAINMAKER_MAX_PEER_COUNT_ALLOW_KEY,
		// CHAINMAKER_CONNECTORS_KEY,
		CHAINMAKER_OPPORTUNISTIC_GRAFT_PEERS_KEY,
		CHAINMAKER_DOUT_KEY,
		CHAINMAKER_DLAZY_KEY,
		CHAINMAKER_HEARTBEAT_INTERVAL_KEY,
		CHAINMAKER_OPPORTUNISTIC_GRAFT_TICKS_KEY,
		CHAINMAKER_GOSSIP_RETRANSMISSION_KEY,
	}

	ChainmakerConfigTpl Config = Config{
		CHAINMAKER_TBFT_PROPOSE_TIMEOUT_KEY:       30,
		CHAINMAKER_TBFT_PROPOSE_DELTA_TIMEOUT_KEY: 1,
		CHAINMAKER_TBFT_BLOCK_PER_PROPOSER:        1,
		CHAINMAKER_BATCH_MAX_SIZE_KEY:             50,
		CHAINMAKER_BATCH_CREATE_TIMEOUT_KEY:       50,
		CHAINMAKER_BLOCK_TX_CAPACITY_KEY:          50,
		// CHAINMAKER_MAX_PEER_COUNT_ALLOW_KEY:       6,
		// CHAINMAKER_CONNECTORS_KEY:                 8,
		CHAINMAKER_OPPORTUNISTIC_GRAFT_PEERS_KEY: 2,
		CHAINMAKER_DOUT_KEY:                      1,
		CHAINMAKER_DLAZY_KEY:                     1,
		CHAINMAKER_HEARTBEAT_INTERVAL_KEY:        1,
		CHAINMAKER_OPPORTUNISTIC_GRAFT_TICKS_KEY: 60,
		CHAINMAKER_GOSSIP_RETRANSMISSION_KEY:     3,
	}

	ChainmakerStepTpl Step = Step{
		CHAINMAKER_TBFT_PROPOSE_TIMEOUT_KEY:       1,
		CHAINMAKER_TBFT_PROPOSE_DELTA_TIMEOUT_KEY: 1,
		CHAINMAKER_TBFT_BLOCK_PER_PROPOSER:        1,
		CHAINMAKER_BATCH_MAX_SIZE_KEY:             10,
		CHAINMAKER_BATCH_CREATE_TIMEOUT_KEY:       10,
		CHAINMAKER_BLOCK_TX_CAPACITY_KEY:          50,
		// CHAINMAKER_MAX_PEER_COUNT_ALLOW_KEY:       1,
		// CHAINMAKER_CONNECTORS_KEY:                 1,
		CHAINMAKER_OPPORTUNISTIC_GRAFT_PEERS_KEY: 1,
		CHAINMAKER_DOUT_KEY:                      1,
		CHAINMAKER_DLAZY_KEY:                     1,
		CHAINMAKER_HEARTBEAT_INTERVAL_KEY:        1,
		CHAINMAKER_OPPORTUNISTIC_GRAFT_TICKS_KEY: 1,
		CHAINMAKER_GOSSIP_RETRANSMISSION_KEY:     1,
	}
)

const (
	CHAINMAKER_TBFT_PROPOSE_TIMEOUT_KEY       string = "TBFT_propose_timeout"
	CHAINMAKER_TBFT_PROPOSE_DELTA_TIMEOUT_KEY string = "TBFT_propose_delta_timeout"
	CHAINMAKER_TBFT_BLOCK_PER_PROPOSER        string = "TBFT_blocks_per_proposer"
	CHAINMAKER_BATCH_MAX_SIZE_KEY             string = "BatchMaxSize"
	CHAINMAKER_BATCH_CREATE_TIMEOUT_KEY       string = "BatchCreateTimeout"
	CHAINMAKER_BLOCK_TX_CAPACITY_KEY          string = "block_tx_capacity"
	CHAINMAKER_MAX_PEER_COUNT_ALLOW_KEY       string = "MaxPeerCountAllow"
	CHAINMAKER_CONNECTORS_KEY                 string = "Connectors"
	CHAINMAKER_OPPORTUNISTIC_GRAFT_PEERS_KEY  string = "OpportunisticGraftPeers"
	CHAINMAKER_DOUT_KEY                       string = "Dout"
	CHAINMAKER_DLAZY_KEY                      string = "Dlazy"
	CHAINMAKER_HEARTBEAT_INTERVAL_KEY         string = "HeartbeatInterval"
	CHAINMAKER_OPPORTUNISTIC_GRAFT_TICKS_KEY  string = "OpportunisticGraftTicks"
	CHAINMAKER_GOSSIP_RETRANSMISSION_KEY      string = "GossipRetransmission"
)

const (
	THRESHOLD float64 = 0.05 //TPS优化阈值

)

func CreateInitConfig(chainType string) (Config, error) {
	switch chainType {
	case "chainmaker":
		// var chainmakerInitConfig Config = Config{
		// 	CHAINMAKER_TBFT_PROPOSE_TIMEOUT_KEY:       30,
		// 	CHAINMAKER_TBFT_PROPOSE_DELTA_TIMEOUT_KEY: 1,
		// 	CHAINMAKER_BATCH_MAX_SIZE_KEY:             20,
		// 	CHAINMAKER_BATCH_CREATE_TIMEOUT_KEY:       50,
		// 	CHAINMAKER_BLOCK_TX_CAPACITY_KEY:          100,
		// 	CHAINMAKER_MAX_PEER_COUNT_ALLOW_KEY:       6,
		// 	CHAINMAKER_CONNECTORS_KEY:                 8,
		// 	CHAINMAKER_OPPORTUNISTIC_GRAFT_PEERS_KEY:  2,
		// 	CHAINMAKER_DOUT_KEY:                       1,
		// 	CHAINMAKER_DLAZY_KEY:                      1,
		// 	CHAINMAKER_HEARTBEAT_INTERVAL_KEY:         1,
		// 	CHAINMAKER_OPPORTUNISTIC_GRAFT_TICKS_KEY:  60,
		// 	CHAINMAKER_GOSSIP_RETRANSMISSION_KEY:      3,
		// }
		chainmakerInitConfig := make(Config)
		for _, configName := range ValidConfigName {
			chainmakerInitConfig[configName] = ChainmakerConfigTpl[configName]
		}
		return chainmakerInitConfig, nil
	default:
		return nil, errors.New("未知的链类型")
	}
}

func GetStep(chainType string) (Step, error) {
	switch chainType {
	case "chainmaker":
		// var chainmakerStep Step = Step{
		// 	CHAINMAKER_TBFT_PROPOSE_TIMEOUT_KEY:       1,
		// 	CHAINMAKER_TBFT_PROPOSE_DELTA_TIMEOUT_KEY: 1,
		// 	CHAINMAKER_BATCH_MAX_SIZE_KEY:             10,
		// 	CHAINMAKER_BATCH_CREATE_TIMEOUT_KEY:       10,
		// 	CHAINMAKER_BLOCK_TX_CAPACITY_KEY:          20,
		// 	CHAINMAKER_MAX_PEER_COUNT_ALLOW_KEY:       1,
		// 	CHAINMAKER_CONNECTORS_KEY:                 1,
		// 	CHAINMAKER_OPPORTUNISTIC_GRAFT_PEERS_KEY:  1,
		// 	CHAINMAKER_DOUT_KEY:                       1,
		// 	CHAINMAKER_DLAZY_KEY:                      1,
		// 	CHAINMAKER_HEARTBEAT_INTERVAL_KEY:         1,
		// 	CHAINMAKER_OPPORTUNISTIC_GRAFT_TICKS_KEY:  1,
		// 	CHAINMAKER_GOSSIP_RETRANSMISSION_KEY:      1,
		// }
		chainmakerInitStep := make(Step)
		for _, configName := range ValidConfigName {
			chainmakerInitStep[configName] = ChainmakerStepTpl[configName]
		}
		return chainmakerInitStep, nil
	default:
		return nil, errors.New("未知的链类型")
	}
}

func CopyConfig(dst, src Config) {
	for k, v := range src {
		dst[k] = v
	}
}
func CopyStep(dst, src Step) {
	for k, v := range src {
		dst[k] = v
	}
}
func CopyConfigWeight(dst, src ConfigWeight) {
	for k, v := range src {
		dst[k] = v
	}
}
