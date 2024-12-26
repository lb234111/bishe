package clientList

import (
	"errors"
	"fmt"
	"framework/utils"
	"framework/worker/chainmaker/client"
	"framework/worker/chainmaker/config"
)

//	type Clients struct {
//		ClientList
//		ChainType string
//	}
type ClientList interface {
	SendTxSync(txType string, txRate int, txCount int, connections int) error
	SendTxAsync(txType string) error
	GetBlockHeight() (int, error)
	GetBlockByHeight(int) (interface{}, error)
	DeployContract(contractName string) error
	GetPoolTxNum() (int, error)
	CalTps(start, end int) (float64, error)
	UpdateConfig(configName string, value int) error
}

func NewClients(chainType string) (ClientList, error) {
	switch chainType {
	case "chainmaker":
		fmt.Println("Create Chainmaker Clients")
		clientList, err := client.CreateChainmakerClient(10, len(config.Chain1Admins), "chain1")
		if err != nil {
			utils.PrintError("初始化客户端失败：", err)
			return nil, err
		}
		return &ChainmakerClientList{clientList: clientList}, nil
	default:
		utils.PrintError("未知的链类型")
		err := errors.New("未知的链类型")
		return nil, err
	}
}
