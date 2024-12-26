package main

import (
	"fmt"
	"log"

	"chainmaker.org/chainmaker/pb-go/v3/common"
)

const (
	send_tx_timeout = int64(50)
)

// chooseNetworkType 判断属于哪种参数
func chooseNetworkType(param ParamAdapter) (uint64, error) {
	var err error
	var block_height uint64
	// tbft 算法
	if param.Type == 1 {
		block_height, err = adaptTbftNetwork(param)
		if err != nil {
			fmt.Println("adaptTbftNetwork err: ", err)
			return block_height, err
		}
	}
	// dpos 算法
	if param.Type == 3 {
		block_height, err = adaptDposNetwork(param)
		if err != nil {
			fmt.Println("adaptDposNetwork err: ", err)
			return block_height, err
		}
	}
	// maxbft 算法
	if param.Type == 4 {
		block_height, err = adaptMaxbftNetwork(param)
		if err != nil {
			fmt.Println("adaptMaxbftNetwork err: ", err)
			return block_height, err
		}
	}
	return block_height, nil
}

// adaptTbft 调整 tbft 节点的网络层参数
func adaptTbftNetwork(param ParamAdapter) (uint64, error) {
	client, err := createClientWithConfig(1)
	if err != nil {
		log.Fatal("createClientWithConfig err: ", err)
		return 0, err
	}

	var payload *common.Payload
	mutex.Lock()
	defer mutex.Unlock()

	if param.ParamName == "block_tx_capacity" {
		payload, err = client.CreateChainConfigBlockUpdatePayload(true, uint32(600), uint32(param.ParamValue), uint32(10), uint32(10), uint32(100))
		if err != nil {
			fmt.Println("Create block_tx_capacity Payload err: ", err)
			return 0, err
		}
		fmt.Println("Create block_tx_capacity Payload success")
	} else {
		kvs := []*common.KeyValuePair{
			{
				Key:   param.ParamName,
				Value: []byte(fmt.Sprintf("%d", param.ParamValue)),
			},
		}
		if _, ok := tbftParam[param.ParamName]; ok {
			payload, err = client.CreateChainConfigConsensusExtUpdatePayload(kvs)
			if err != nil {
				return 0, err
			}
		} else {
			payload, err = client.CreateChainConfigConsensusExtAddPayload(kvs)
			if err != nil {
				return 0, err
			}
		}
	}

	tbftParam[param.ParamName] = ""
	// 第一种方式获取多签
	endorsers := signPayload(client, payload, Chain1Admins...)
	fmt.Println("Create endorsers success")
	// 发送配置更新请求
	resp, err := client.SendChainConfigUpdateRequest(payload, endorsers, send_tx_timeout, true)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("block_height=", resp.TxBlockHeight)
	return resp.TxBlockHeight, nil
}

// adaptMaxbft 调整 maxbft 算法
func adaptMaxbftNetwork(param ParamAdapter) (uint64, error) {
	client, err := createClientWithConfig(1)
	if err != nil {
		fmt.Println("createClientWithConfig err: ", err)
		return 0, err
	}

	var payload *common.Payload
	mutex.Lock()
	defer mutex.Unlock()

	if param.ParamName == "block_tx_capacity" {
		payload, err = client.CreateChainConfigBlockUpdatePayload(true, uint32(600), uint32(param.ParamValue), uint32(10), uint32(10), uint32(100))
		if err != nil {
			fmt.Println("Create block_tx_capacity Payload err: ", err)
			return 0, err
		}
		fmt.Println("Create block_tx_capacity Payload success")
	} else {
		kvs := []*common.KeyValuePair{
			{
				Key:   param.ParamName,
				Value: []byte(fmt.Sprintf("%d", param.ParamValue)),
			},
		}
		if _, ok := tbftParam[param.ParamName]; ok {
			payload, err = client.CreateChainConfigConsensusExtUpdatePayload(kvs)
			if err != nil {
				return 0, err
			}
		} else {
			payload, err = client.CreateChainConfigConsensusExtAddPayload(kvs)
			if err != nil {
				return 0, err
			}
		}
	}
	tbftParam[param.ParamName] = ""
	// 第一种方式获取多签
	endorsers := signPayload(client, payload, UserNameOrg1Admin1, UserNameOrg2Admin1,
		UserNameOrg3Admin1, UserNameOrg4Admin1, UserNameOrg5Admin1, UserNameOrg6Admin1, UserNameOrg7Admin1, UserNameOrg8Admin1, UserNameOrg9Admin1, UserNameOrg10Admin1)
	// 发送配置更新请求
	resp, err := client.SendChainConfigUpdateRequest(payload, endorsers, send_tx_timeout, true)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("block_height=", resp.TxBlockHeight)
	return resp.TxBlockHeight, nil
}

// adaptDpos 调整 dpos 算法
func adaptDposNetwork(param ParamAdapter) (uint64, error) {
	client, err := createClientWithConfig(3)
	if err != nil {
		fmt.Println("createClientWithConfig err: ", err)
		return 0, err
	}
	payload, err := client.SetEpochValidatorNumber(int(param.ParamValue))
	if err != nil {
		return 0, err
	}
	// 第一种方式获取多签
	endorsers := signPayload(client, payload, UserNameOrg1Admin1, UserNameOrg2Admin1,
		UserNameOrg3Admin1, UserNameOrg4Admin1, UserNameOrg5Admin1, UserNameOrg6Admin1, UserNameOrg7Admin1, UserNameOrg8Admin1, UserNameOrg9Admin1, UserNameOrg10Admin1)
	// 发送配置更新请求
	resp, err := client.SendPayloadRequest(payload, endorsers, send_tx_timeout, true)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("block_height=", resp.TxBlockHeight)
	return resp.TxBlockHeight, nil
}
