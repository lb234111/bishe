// coding:utf-8
// 部署合约
package deploy

import (
	"errors"
	"fmt"
	"framework/worker/chainmaker/config"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	sdkutils "chainmaker.org/chainmaker/sdk-go/v2/utils"
)

func DeployContract(
	client *sdk.ChainClient,
	contractName string,
	version string,
	byteCodePath string,
	runtime common.RuntimeType,
	kvs []*common.KeyValuePair,
	withSyncResult bool,
	createContractTimeout int64,
	usernames ...string,
) (*common.TxResponse, error) {
	// todo:这是干什么的
	payload, err := client.CreateContractCreatePayload(contractName, version, byteCodePath, runtime, kvs)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	// 获取背书节点
	endorsers, err := getEndorsers(payload, usernames...)
	if err != nil {
		return nil, fmt.Errorf("getEndorsers error: " + err.Error())
	}
	// todo:干什么的
	resp, err := client.SendContractManageRequest(payload, endorsers, createContractTimeout, withSyncResult)
	if err != nil {
		return resp, fmt.Errorf("SendContractManageRequest error:" + err.Error())
	}
	// 检查回执
	err = checkProposalRequestResp(resp, true)
	if err != nil {
		return resp, fmt.Errorf("checkProposalRequestResp error:" + err.Error())
	}

	// log
	fmt.Println("deploy contract success\n[code]:" + string(resp.Code) + "\n[message]:" + resp.Message + "\n[contractResult]:" + resp.ContractResult.String() + "\n[txid]:" + resp.TxId)
	return resp, nil
}

// 检查回执
func checkProposalRequestResp(resp *common.TxResponse, needContractResult bool) error {
	if resp.Code != common.TxStatusCode_SUCCESS {
		if resp.Message == "" {
			resp.Message = resp.Code.String()
		}
		return fmt.Errorf(resp.Message)
	}
	// todo:这是检查什么
	if needContractResult && resp.ContractResult == nil {
		return fmt.Errorf("contract result is nil")
	}
	// 0是成功代码
	if resp.ContractResult != nil && resp.ContractResult.Code != 0 {
		return fmt.Errorf(resp.ContractResult.Message)
	}
	return nil
}

// 获取背书节点
func getEndorsers(payload *common.Payload, usernames ...string) ([]*common.EndorsementEntry, error) {
	var endorsers []*common.EndorsementEntry

	for _, name := range usernames {
		// 获取管理员证书路径
		users := config.Admins
		u, ok := users[name]
		if !ok {
			return nil, errors.New("user not found")
		}

		var err error
		var entry *common.EndorsementEntry
		p11Handle := sdk.GetP11Handle()
		if p11Handle != nil {
			entry, err = sdkutils.MakeEndorserWithPathAndP11Handle(u.SignKeyPath, u.SignCrtPath, p11Handle, payload)
			if err != nil {
				return nil, err
			}
		} else {
			entry, err = sdkutils.MakeEndorserWithPath(u.SignKeyPath, u.SignCrtPath, payload)
			if err != nil {
				return nil, err
			}
		}

		endorsers = append(endorsers, entry)
	}

	return endorsers, nil
}
