package main

import (
	"errors"
	"fmt"
	"log"

	"chainmaker.org/chainmaker/common/v3/crypto"
	"chainmaker.org/chainmaker/pb-go/v3/common"
	sdk "chainmaker.org/chainmaker/sdk-go/v3"
	sdkutils "chainmaker.org/chainmaker/sdk-go/v3/utils"
)

// chooseConsensusType 判断属于哪种算法
func chooseConsensusType(param ParamAdapter) (uint64, error) {
	var err error
	var block_height uint64
	// tbft 算法
	if param.Type == 1 {
		block_height, err = adaptTbft(param)
		if err != nil {
			fmt.Println("adaptTbft err: ", err)
			return 0, err
		}
	}
	// dpos 算法
	if param.Type == 3 {
		block_height, err = adaptDpos(param)
		if err != nil {
			fmt.Println("adaptDpos err: ", err)
			return 0, err
		}
	}
	// maxbft 算法
	if param.Type == 4 {
		block_height, err = adaptMaxbft(param)
		if err != nil {
			fmt.Println("adaptMaxbft err: ", err)
			return 0, err
		}
	}
	return block_height, nil
}

// adaptMaxbft 调整 maxbft 算法
func adaptMaxbft(param ParamAdapter) (uint64, error) {
	client, err := createClientWithConfig(1)
	if err != nil {
		fmt.Println("createClientWithConfig err: ", err)
		return 0, err
	}
	kvs := []*common.KeyValuePair{
		{
			Key:   param.ParamName,
			Value: []byte(fmt.Sprintf("%d", param.ParamValue)),
		},
	}
	var payload *common.Payload
	mutex.Lock()
	defer mutex.Unlock()
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
	tbftParam[param.ParamName] = ""
	// 第一种方式获取多签
	endorsers := signPayload(client, payload, UserNameOrg1Admin1, UserNameOrg2Admin1,
		UserNameOrg3Admin1, UserNameOrg4Admin1)
	// 发送配置更新请求
	resp, err := client.SendChainConfigUpdateRequest(payload, endorsers, send_tx_timeout, true)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("block_height=", resp.TxBlockHeight)
	return resp.TxBlockHeight, err
}

// adaptDpos 调整 dpos 算法
func adaptDpos(param ParamAdapter) (uint64, error) {
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
		UserNameOrg3Admin1, UserNameOrg4Admin1)
	// 发送配置更新请求
	resp, err := client.SendPayloadRequest(payload, endorsers, send_tx_timeout, true)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("block_height=", resp.TxBlockHeight)
	return resp.TxBlockHeight, nil
}

// adaptTbft 调整 tbft 算法
func adaptTbft(param ParamAdapter) (uint64, error) {
	client, err := createClientWithConfig(1)
	if err != nil {
		fmt.Println("createClientWithConfig err: ", err)
		return 0, err
	}
	var kvs []*common.KeyValuePair
	if param.ParamName == "TBFT_blocks_per_proposer" {
		kvs = []*common.KeyValuePair{
			{
				Key:   param.ParamName,
				Value: []byte(fmt.Sprintf("%d", param.ParamValue)),
			},
		}
	} else {
		kvs = []*common.KeyValuePair{
			{
				Key:   param.ParamName,
				Value: []byte(fmt.Sprintf("%d", param.ParamValue) + "s"),
			},
		}
	}
	var payload *common.Payload
	mutex.Lock()
	defer mutex.Unlock()
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
	tbftParam[param.ParamName] = ""
	// 第一种方式获取多签
	endorsers := signPayload(client, payload, UserNameOrg1Admin1, UserNameOrg2Admin1,
		UserNameOrg3Admin1, UserNameOrg4Admin1)
	// 发送配置更新请求
	resp, err := client.SendChainConfigUpdateRequest(payload, endorsers, send_tx_timeout, true)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("block_height=", resp.TxBlockHeight)
	// 第二种方式（暂时没走通）
	// endorser, err := client.SignChainConfigPayload(payload)
	// if err != nil {
	// 	return err
	// }
	// var endorsers []*common.EndorsementEntry
	// endorsers = append(endorsers, endorser)
	// rsp, err := client.SendChainConfigUpdateRequest(payload, endorsers, -1, true)
	// if err != nil {
	// 	return err
	// }
	// fmt.Println(rsp)
	return resp.TxBlockHeight, err
}

// signPayload 签名 Payload
func signPayload(client *sdk.ChainClient, payload *common.Payload, usernames ...string) []*common.EndorsementEntry {
	// 各组织Admin权限用户签名
	endorsers, err := getEndorsersWithAuthType(client.GetHashType(), client.GetAuthType(), payload, usernames...)
	if err != nil {
		fmt.Println("#####1######")
		log.Fatalln(err)
	}
	return endorsers
}

// getEndorsersWithAuthType 多签
func getEndorsersWithAuthType(hashType crypto.HashType, authType sdk.AuthType,
	payload *common.Payload, usernames ...string) ([]*common.EndorsementEntry, error) {
	var endorsers []*common.EndorsementEntry
	for _, name := range usernames {
		var entry *common.EndorsementEntry
		var err error
		switch authType {
		case sdk.PermissionedWithCert:
			u, ok := users[name]
			if !ok {
				return nil, errors.New("user not found")
			}
			entry, err = sdkutils.MakeEndorserWithPath(u.SignKeyPath, u.SignCrtPath, payload)
			if err != nil {
				fmt.Println("#####2######")
				return nil, err
			}
		case sdk.PermissionedWithKey:
			u, ok := permissionedPkUsers[name]
			if !ok {
				return nil, errors.New("user not found")
			}
			entry, err = sdkutils.MakePkEndorserWithPath(u.SignKeyPath, hashType, u.OrgId, payload)
			if err != nil {
				return nil, err
			}
		case sdk.Public:
			u, ok := pkUsers[name]
			if !ok {
				return nil, errors.New("user not found")
			}
			entry, err = sdkutils.MakePkEndorserWithPath(u.SignKeyPath, hashType, "", payload)
			if err != nil {
				return nil, err
			}
		default:
			return nil, errors.New("invalid authType")
		}
		endorsers = append(endorsers, entry)
	}
	return endorsers, nil
}
