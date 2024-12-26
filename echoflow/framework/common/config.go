package common

import "time"

var (
	//for master
	// WorkerIP []string = []string{"10.21.4.33", "10.21.4.39"}
	WorkerIP      []string = []string{"10.15.9.14", "10.15.9.15", "10.15.9.16", "10.15.9.17"}
	DataFilePath  string   = "./data_0413/nodeinfo" + time.Now().Format("2006010215") + ".json"
	AliveNodePath string   = "./data_0413/aliveNode" + time.Now().Format("2006010215") + ".dat"
	BatchSize     int      = 10

	//for worker
	MasterIP string = "10.15.9.37"
	//for tx
	TransactionType        string = "store"
	TransactionCount       int    = 40000
	TransactionRate        int    = 4000
	Connections            int    = 40
	ChainmakerContractPath string = "./worker/chainmaker/contracts"

	//for all
	ChainType string = "chainmaker"
)
