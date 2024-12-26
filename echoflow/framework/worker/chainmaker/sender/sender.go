package chainmakersender

import (
	"framework/worker/chainmaker/call"
	"framework/worker/chainmaker/config"
	"math/rand"
	"time"

	sdk "chainmaker.org/chainmaker/sdk-go/v2"
)

func SendTxSync(clientList []*sdk.ChainClient, txType string, txRate int, txCount int, connections int) {

	// txNum := txPoolStatus.CommonTxNumInQueue
	// fmt.Println("当前交易池交易数：", txNum)
	sendTime := txCount / txRate
	for i := 0; i < sendTime+2; i++ {
		start := time.Now().UnixMilli()
		if txType == "donothing" {
			for i := 0; i < int(connections); i++ {
				go call.CallContract(clientList[rand.Intn(config.OrgNum*config.ClientNum)], txType, txRate/connections, false)
			}
		} else if txType == "store" {
			for i := 0; i < int(connections); i++ {
				go call.CallStoreContrace(clientList[rand.Intn(config.OrgNum*config.ClientNum)], txRate/connections, false)
			}
		}
		end := time.Now().UnixMilli()
		time.Sleep(1000 - time.Millisecond*time.Duration(end-start))
	}
}

func SendTxAsync(clientList []*sdk.ChainClient, contractName string) error {

	// txNum := txPoolStatus.CommonTxNumInQueue
	// fmt.Println("当前交易池交易数：", txNum)
	var err error
	if contractName == "donothing" {
		err = call.CallContract(clientList[rand.Intn(config.OrgNum*config.ClientNum)], contractName, 1, true)
	} else if contractName == "store" {
		err = call.CallStoreContrace(clientList[rand.Intn(config.OrgNum*config.ClientNum)], 1, true)
	}
	return err
}
