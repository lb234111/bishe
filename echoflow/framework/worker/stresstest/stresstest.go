package stresstest

import (
	"framework/common"
	"framework/utils"
	"framework/worker/clientList"
	"log"
	"time"
)

type Stresstest struct {
	clients          clientList.ClientList
	transactionType  string
	transactionCount int
	transactionRate  int
	connections      int
}

var logger *log.Logger

func init() {
	logger = utils.InitLogger("txloadtest")
}
func NewStresstest(clients clientList.ClientList) *Stresstest {
	logger.Println("部署合约")
	clients.DeployContract("donothing")
	clients.DeployContract("store")
	return &Stresstest{
		clients:          clients,
		transactionType:  common.TransactionType,
		transactionCount: common.TransactionCount,
		transactionRate:  common.TransactionRate,
		connections:      common.Connections,
	}
}

func (stresstest *Stresstest) SendTx() (start, end int, err error) {
	clients := stresstest.clients
	txType := stresstest.transactionType
	txRate := stresstest.transactionRate
	txCount := stresstest.transactionCount
	connections := stresstest.connections
	logger.Println("开始发送交易，调用合约", txType)

	// go clients.SendTxSync(txType, txRate, txCount, connections)
	go clients.SendTxSync(txType, txRate, txCount, connections)
	time.Sleep(time.Second)
	start, err = clients.GetBlockHeight()
	logger.Println("初始高度", start)
	if err != nil {
		utils.PrintError("获取交易前高度失败")
		return 0, 0, err
	}
	time.Sleep(time.Second * time.Duration(txCount/txRate))
	logger.Println("发送交易完成")
	for {
		end, err = clients.GetBlockHeight()
		if err != nil {
			utils.PrintError("获取交易后高度失败")
			return 0, 0, err
		}
		if end > start {
			break
		}
	}
	logger.Println("结束高度", end)
	logger.Println("等待交易处理完成")
	for {
		txNum, err := clients.GetPoolTxNum()
		if err != nil {
			utils.PrintError(err)
			return 0, 0, err
		}
		utils.PrintInfo("交易池数量", txNum)
		if txNum == 0 {
			break
		}
		time.Sleep(2 * time.Second)
	}
	logger.Println("交易处理完成")

	return
}
