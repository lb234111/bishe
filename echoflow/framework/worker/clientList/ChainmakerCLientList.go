package clientList

import (
	"fmt"
	"framework/utils"
	"framework/worker/chainmaker/config"
	"framework/worker/chainmaker/deploy"
	chainmakersender "framework/worker/chainmaker/sender"

	frameworkcommon "framework/common"

	pbcommon "chainmaker.org/chainmaker/pb-go/v2/common"
	chainmaker_sdk "chainmaker.org/chainmaker/sdk-go/v2"
)

type ChainmakerClientList struct {
	clientList []*chainmaker_sdk.ChainClient
}

func (chainmakerclientList *ChainmakerClientList) SendTxSync(txType string, txRate int, txCount int, connections int) error {
	clientList := chainmakerclientList.clientList
	chainmakersender.SendTxSync(clientList, txType, txRate, txCount, connections)
	return nil

}
func (chainmakerclientList *ChainmakerClientList) SendTxAsync(contractName string) error {
	clientList := chainmakerclientList.clientList
	err := chainmakersender.SendTxAsync(clientList, contractName)
	return err

}

func (chainmakerclientList *ChainmakerClientList) GetBlockHeight() (int, error) {

	clientList := chainmakerclientList.clientList
	height, err := clientList[0].GetCurrentBlockHeight()
	if err != nil {
		utils.PrintError("获取区块高度失败：", err)
		return 0, err
	}
	return int(height), nil
}

func (chainmakerClientList *ChainmakerClientList) GetBlockByHeight(height int) (interface{}, error) {
	clientList := chainmakerClientList.clientList
	blockInfo, err := clientList[0].GetBlockByHeight(uint64(height), false)
	if err != nil {
		return nil, err
	}
	return blockInfo.Block, nil
}

func (chainmakerclientList *ChainmakerClientList) GetPoolTxNum() (int, error) {
	clientList := chainmakerclientList.clientList
	txPoolStatus, err := clientList[0].GetPoolStatus()
	if err != nil {
		utils.PrintError("获取交易池状态失败", err)
	}
	txNum := txPoolStatus.CommonTxNumInQueue
	return int(txNum), nil

}

func (chainmakerclientList *ChainmakerClientList) CalTps(start, end int) (float64, error) {
	txNum := 0
	clientList := chainmakerclientList.clientList
	for i := start; i <= end; i++ {
		block, err := clientList[i%(config.OrgNum*config.ClientNum)].GetBlockByHeight(uint64(i), false)
		if err != nil {
			utils.PrintError("获取区块失败：", err)
			i -= 1
			continue
		}
		if block.Block == nil {
			continue
		}
		txNum += int(block.Block.Header.TxCount)
	}
	pre_height_block, err := clientList[0].GetBlockByHeight(uint64(start), false)
	if err != nil {
		utils.PrintError("获取区块失败", err)
	}
	after_height_block, err := clientList[0].GetBlockByHeight(uint64(end), false)
	if err != nil {
		utils.PrintError("获取区块失败", err)
	}
	tps := (float64(txNum) / float64(after_height_block.Block.Header.GetBlockTimestamp()-pre_height_block.Block.Header.GetBlockTimestamp())) * 1000.0
	return tps, nil
}

func (chainmakerclientList *ChainmakerClientList) DeployContract(contractName string) error {
	clientList := chainmakerclientList.clientList
	_, err := clientList[0].GetContractInfo(contractName)
	if err == nil {
		utils.PrintWarn(contractName, "合约存在:")
		return nil
	}
	aProposalArgs := []*pbcommon.KeyValuePair{}
	contract7zPath := fmt.Sprintf(frameworkcommon.ChainmakerContractPath + "/" + contractName + "/" + contractName + ".7z")
	_, err = deploy.DeployContract(clientList[0],
		contractName, "v0",
		contract7zPath, pbcommon.RuntimeType_DOCKER_GO,
		aProposalArgs,
		true,
		10,
		config.Chain1Admins...)
	if err != nil {
		utils.PrintError("部署合约出错", err)
		return fmt.Errorf(err.Error())
	}
	utils.PrintInfo("deploy contract success")
	return nil
}

func (chainmakerClientList *ChainmakerClientList) UpdateConfig(configName string, value int) error {
	// fmt.Print(configName, value)
	return nil
}
