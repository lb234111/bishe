/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package client

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"chainmaker.org/chainmaker-go/tools/cmc/types"
	"chainmaker.org/chainmaker-go/tools/cmc/util"
	"chainmaker.org/chainmaker/common/v2/evmutils/abi"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	"github.com/spf13/cobra"
)

// CHECK_PROPOSAL_RESPONSE_FAILED_FORMAT define CHECK_PROPOSAL_RESPONSE_FAILED_FORMAT error fmt
const CHECK_PROPOSAL_RESPONSE_FAILED_FORMAT = "checkProposalRequestResp failed, %s"

// SEND_CONTRACT_MANAGE_REQUEST_FAILED_FORMAT define SEND_CONTRACT_MANAGE_REQUEST_FAILED_FORMAT error fmt
const SEND_CONTRACT_MANAGE_REQUEST_FAILED_FORMAT = "SendContractManageRequest failed, %s"

// ADMIN_ORGID_KEY_CERT_LENGTH_NOT_EQUAL_FORMAT define ADMIN_ORGID_KEY_CERT_LENGTH_NOT_EQUAL_FORMAT error fmt
const ADMIN_ORGID_KEY_CERT_LENGTH_NOT_EQUAL_FORMAT = "admin orgId & key & cert list length not equal, " +
	"[keys len: %d]/[certs len:%d]"

// ADMIN_ORGID_KEY_LENGTH_NOT_EQUAL_FORMAT define ADMIN_ORGID_KEY_LENGTH_NOT_EQUAL_FORMAT error fmt
const ADMIN_ORGID_KEY_LENGTH_NOT_EQUAL_FORMAT = "admin orgId & key list length not equal, " +
	"[keys len: %d]/[org-ids len:%d]"

var (
	//errAdminOrgIdKeyCertIsEmpty = errors.New("admin orgId or key or cert list is empty")
	noConstructorErrMsg = "contract does not have a constructor"
)

// UserContract define user contract
type UserContract struct {
	ContractName string
	Method       string
	Params       map[string]string
}

// userContractCMD user contract command
// @return *cobra.Command
func userContractCMD() *cobra.Command {
	userContractCmd := &cobra.Command{
		Use:   "user",
		Short: "user contract command",
		Long:  "user contract command",
	}

	userContractCmd.AddCommand(createUserContractCMD())
	userContractCmd.AddCommand(invokeContractTimesCMD())
	userContractCmd.AddCommand(invokeUserContractCMD())
	userContractCmd.AddCommand(upgradeUserContractCMD())
	userContractCmd.AddCommand(freezeUserContractCMD())
	userContractCmd.AddCommand(unfreezeUserContractCMD())
	userContractCmd.AddCommand(revokeUserContractCMD())
	userContractCmd.AddCommand(getUserContractCMD())

	return userContractCmd
}

// createUserContractCMD create user contract command
// @return *cobra.Command
func createUserContractCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create user contract command",
		Long:  "create user contract command",
		RunE: func(_ *cobra.Command, _ []string) error {
			return createUserContract()
		},
	}

	attachFlags(cmd, []string{
		flagUserTlsKeyFilePath, flagUserTlsCrtFilePath, flagUserSignKeyFilePath, flagUserSignCrtFilePath,
		flagSdkConfPath, flagContractName, flagVersion, flagByteCodePath, flagOrgId, flagChainId, flagSendTimes,
		flagRuntimeType, flagTimeout, flagParams, flagSyncResult, flagEnableCertHash, flagAbiFilePath,
		flagAdminKeyFilePaths, flagAdminCrtFilePaths, flagAdminOrgIds, flagGasLimit,
	})

	cmd.MarkFlagRequired(flagSdkConfPath)
	cmd.MarkFlagRequired(flagContractName)
	cmd.MarkFlagRequired(flagVersion)
	cmd.MarkFlagRequired(flagByteCodePath)
	cmd.MarkFlagRequired(flagRuntimeType)

	return cmd
}

// invokeUserContractCMD invoke user contract command
// @return *cobra.Command
func invokeUserContractCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "invoke",
		Short: "invoke user contract command",
		Long:  "invoke user contract command",
		RunE: func(_ *cobra.Command, _ []string) error {
			return invokeUserContract()
		},
	}

	attachFlags(cmd, []string{
		flagUserSignKeyFilePath, flagUserSignCrtFilePath, flagUserTlsKeyFilePath, flagUserTlsCrtFilePath,
		flagConcurrency, flagTotalCountPerGoroutine, flagSdkConfPath, flagOrgId, flagChainId, flagSendTimes,
		flagEnableCertHash, flagContractName, flagMethod, flagParams, flagTimeout, flagSyncResult, flagAbiFilePath,
		flagGasLimit, flagTxId, flagContractAddress,
	})

	cmd.MarkFlagRequired(flagSdkConfPath)
	cmd.MarkFlagRequired(flagMethod)

	return cmd
}

// invokeContractTimesCMD invoke contract and set invoke times
//多次的并发调用指定合约方法
// @return *cobra.Command
func invokeContractTimesCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "invoke-times",
		Short: "invoke contract times command",
		Long:  "invoke contract times command",
		RunE: func(_ *cobra.Command, _ []string) error {
			return invokeContractTimes()
		},
	}

	attachFlags(cmd, []string{
		flagUserSignKeyFilePath, flagUserSignCrtFilePath, flagUserTlsKeyFilePath, flagUserTlsCrtFilePath,
		flagEnableCertHash, flagConcurrency, flagTotalCountPerGoroutine, flagSdkConfPath, flagOrgId, flagChainId,
		flagSendTimes, flagContractName, flagMethod, flagParams, flagTimeout, flagSyncResult, flagAbiFilePath,
		flagContractAddress,
	})

	cmd.MarkFlagRequired(flagSdkConfPath)
	cmd.MarkFlagRequired(flagMethod)

	return cmd
}

// getUserContractCMD query user contract command
// @return *cobra.Command
func getUserContractCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "get user contract command",
		Long:  "get user contract command",
		RunE: func(_ *cobra.Command, _ []string) error {
			return getUserContract()
		},
	}

	attachFlags(cmd, []string{
		flagUserSignKeyFilePath, flagUserSignCrtFilePath, flagUserTlsKeyFilePath, flagUserTlsCrtFilePath,
		flagEnableCertHash, flagConcurrency, flagTotalCountPerGoroutine, flagSdkConfPath, flagOrgId, flagChainId,
		flagSendTimes, flagContractName, flagMethod, flagParams, flagTimeout, flagContractAddress, flagAbiFilePath,
	})

	cmd.MarkFlagRequired(flagSdkConfPath)
	cmd.MarkFlagRequired(flagMethod)

	return cmd
}

// upgradeUserContractCMD upgrade user contract command
// @return *cobra.Command
func upgradeUserContractCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "upgrade user contract command",
		Long:  "upgrade user contract command",
		RunE: func(_ *cobra.Command, _ []string) error {
			return upgradeUserContract()
		},
	}

	attachFlags(cmd, []string{
		flagUserSignKeyFilePath, flagUserSignCrtFilePath, flagUserTlsKeyFilePath, flagUserTlsCrtFilePath,
		flagSdkConfPath, flagContractName, flagVersion, flagByteCodePath, flagOrgId, flagChainId, flagSendTimes,
		flagRuntimeType, flagTimeout, flagParams, flagSyncResult, flagEnableCertHash, flagContractAddress,
		flagAdminCrtFilePaths, flagAdminKeyFilePaths, flagAdminOrgIds, flagGasLimit,
	})

	cmd.MarkFlagRequired(flagSdkConfPath)
	cmd.MarkFlagRequired(flagVersion)
	cmd.MarkFlagRequired(flagByteCodePath)
	cmd.MarkFlagRequired(flagRuntimeType)

	return cmd
}

// freezeUserContractCMD freeze user contract command
// @return *cobra.Command
func freezeUserContractCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "freeze",
		Short: "freeze user contract command",
		Long:  "freeze user contract command",
		RunE: func(_ *cobra.Command, _ []string) error {
			return freezeOrUnfreezeOrRevokeUserContract(1)
		},
	}

	attachFlags(cmd, []string{
		flagUserSignKeyFilePath, flagUserSignCrtFilePath, flagUserTlsKeyFilePath, flagUserTlsCrtFilePath,
		flagSdkConfPath, flagContractName, flagOrgId, flagChainId, flagSendTimes, flagTimeout, flagParams,
		flagSyncResult, flagEnableCertHash, flagAdminCrtFilePaths, flagAdminKeyFilePaths, flagAdminOrgIds,
		flagContractAddress,
	})

	cmd.MarkFlagRequired(flagSdkConfPath)

	return cmd
}

// unfreezeUserContractCMD unfreeze user contract command
// @return *cobra.Command
func unfreezeUserContractCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unfreeze",
		Short: "unfreeze user contract command",
		Long:  "unfreeze user contract command",
		RunE: func(_ *cobra.Command, _ []string) error {
			return freezeOrUnfreezeOrRevokeUserContract(2)
		},
	}

	attachFlags(cmd, []string{
		flagUserSignKeyFilePath, flagUserSignCrtFilePath, flagUserTlsKeyFilePath, flagUserTlsCrtFilePath,
		flagSdkConfPath, flagContractName, flagOrgId, flagChainId, flagSendTimes, flagTimeout, flagParams,
		flagSyncResult, flagEnableCertHash, flagAdminCrtFilePaths, flagAdminKeyFilePaths, flagAdminOrgIds,
		flagContractAddress,
	})

	cmd.MarkFlagRequired(flagSdkConfPath)

	return cmd
}

// revokeUserContractCMD revoke user contract command
// @return *cobra.Command
func revokeUserContractCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke",
		Short: "revoke user contract command",
		Long:  "revoke user contract command",
		RunE: func(_ *cobra.Command, _ []string) error {
			return freezeOrUnfreezeOrRevokeUserContract(3)
		},
	}

	attachFlags(cmd, []string{
		flagUserSignKeyFilePath, flagUserSignCrtFilePath, flagUserTlsKeyFilePath, flagUserTlsCrtFilePath,
		flagSdkConfPath, flagContractName, flagOrgId, flagChainId, flagSendTimes, flagTimeout, flagParams,
		flagSyncResult, flagEnableCertHash, flagAdminCrtFilePaths, flagAdminKeyFilePaths, flagAdminOrgIds,
		flagContractAddress,
	})

	cmd.MarkFlagRequired(flagSdkConfPath)

	return cmd
}

func createUserContract() error {
	client, err := util.CreateChainClient(sdkConfPath, chainId, orgId, userTlsCrtFilePath, userTlsKeyFilePath,
		userSignCrtFilePath, userSignKeyFilePath)
	if err != nil {
		return err
	}
	defer client.Stop()

	adminKeys, adminCrts, adminOrgs, err := util.MakeAdminInfo(client, adminKeyFilePaths, adminCrtFilePaths, adminOrgIds)
	if err != nil {
		return err
	}

	rt, ok := common.RuntimeType_value[runtimeType]
	if !ok {
		return fmt.Errorf("unknown runtime type [%s]", runtimeType)
	}

	var kvs []*common.KeyValuePair

	if runtimeType != "EVM" {
		if params != "" {
			kvsMap := make(map[string]string)
			err := json.Unmarshal([]byte(params), &kvsMap)
			if err != nil {
				return err
			}
			kvs = util.ConvertParameters(kvsMap)
		}
	} else { // EVM contract deploy
		if abiFilePath == "" {
			return errors.New("required abi file path when deploy EVM contract")
		}
		abiBytes, err := ioutil.ReadFile(abiFilePath)
		if err != nil {
			return err
		}

		contractAbi, err := abi.JSON(bytes.NewReader(abiBytes))
		if err != nil {
			return err
		}

		inputData, err := util.Pack(contractAbi, "", params)
		if err != nil {
			if err.Error() != noConstructorErrMsg {
				return err
			}
		}

		inputDataHexStr := hex.EncodeToString(inputData)
		kvs = []*common.KeyValuePair{
			{
				Key:   "data",
				Value: []byte(inputDataHexStr),
			},
		}

		byteCode, err := ioutil.ReadFile(byteCodePath)
		if err != nil {
			return err
		}
		byteCodePath = string(byteCode)
	}

	payload, err := client.CreateContractCreatePayload(
		contractName,
		version,
		byteCodePath,
		common.RuntimeType(rt),
		kvs,
	)
	if err != nil {
		return err
	}

	if gasLimit > 0 {
		var limit = &common.Limit{GasLimit: gasLimit}
		payload = client.AttachGasLimit(payload, limit)
	}

	endorsementEntrys, err := util.MakeEndorsement(adminKeys, adminCrts, adminOrgs, client, payload)
	if err != nil {
		return err
	}
	resp, err := client.SendContractManageRequest(payload, endorsementEntrys, timeout, syncResult)
	if err != nil {
		return err
	}
	err = util.CheckProposalRequestResp(resp, false)
	if err != nil {
		return err
	}
	return createUpgradeUserContractOutput(resp)
}

func createUpgradeUserContractOutput(resp *common.TxResponse) error {
	if resp.ContractResult != nil && resp.ContractResult.Result != nil {
		var contract common.Contract
		err := contract.Unmarshal(resp.ContractResult.Result)
		if err != nil {
			return err
		}
		util.PrintPrettyJson(types.CreateUpgradeContractTxResponse{
			TxResponse: resp,
			ContractResult: &types.CreateUpgradeContractContractResult{
				ContractResult: resp.ContractResult,
				Result:         &contract,
			},
		})
	} else {
		util.PrintPrettyJson(resp)
	}
	return nil
}

func invokeUserContract() error {
	cc, err := sdk.NewChainClient(
		sdk.WithConfPath(sdkConfPath),
		sdk.WithChainClientChainId(chainId),
		sdk.WithChainClientOrgId(orgId),
		sdk.WithUserCrtFilePath(userTlsCrtFilePath),
		sdk.WithUserKeyFilePath(userTlsKeyFilePath),
		sdk.WithUserSignCrtFilePath(userSignCrtFilePath),
		sdk.WithUserSignKeyFilePath(userSignKeyFilePath),
	)
	if err != nil {
		return err
	}
	defer cc.Stop()
	if err := util.DealChainClientCertHash(cc, enableCertHash); err != nil {
		return err
	}

	if contractAddress != "" {
		contractName = contractAddress
	}
	if contractName == "" {
		return errors.New("either contract-name or contract-address must be set")
	}

	var kvs []*common.KeyValuePair
	var contractAbi *abi.ABI

	if abiFilePath != "" { // abi file path 非空 意味着调用的是EVM合约
		abiBytes, err := ioutil.ReadFile(abiFilePath)
		if err != nil {
			return err
		}

		contractAbi, err = abi.JSON(bytes.NewReader(abiBytes))
		if err != nil {
			return err
		}

		inputData, err := util.Pack(contractAbi, method, params)
		if err != nil {
			return err
		}

		inputDataHexStr := hex.EncodeToString(inputData)

		kvs = []*common.KeyValuePair{
			{
				Key:   "data",
				Value: []byte(inputDataHexStr),
			},
		}
	} else {
		if params != "" {
			kvsMap := make(map[string]string)
			err := json.Unmarshal([]byte(params), &kvsMap)
			if err != nil {
				return err
			}
			kvs = util.ConvertParameters(kvsMap)
		}
	}

	var limit *common.Limit
	if gasLimit > 0 {
		limit = &common.Limit{GasLimit: gasLimit}
	}

	if txId != "" {
		invokeContract(cc, contractName, method, txId, kvs, contractAbi, limit)
	} else {
		Dispatch(cc, contractName, method, kvs, contractAbi, limit)
	}
	return nil
}

func invokeContractTimes() error {
	client, err := util.CreateChainClient(sdkConfPath, chainId, orgId, userTlsCrtFilePath, userTlsKeyFilePath,
		userSignCrtFilePath, userSignKeyFilePath)
	if err != nil {
		return err
	}
	defer client.Stop()

	if contractAddress != "" {
		contractName = contractAddress
	}
	if contractName == "" {
		return errors.New("either contract-name or contract-address must be set")
	}

	var kvs []*common.KeyValuePair
	var evmMethod *abi.ABI

	if abiFilePath != "" { // abi file path 非空 意味着调用的是EVM合约
		abiBytes, err := ioutil.ReadFile(abiFilePath)
		if err != nil {
			return err
		}

		contractAbi, err := abi.JSON(bytes.NewReader(abiBytes))
		if err != nil {
			return err
		}

		//m, exist := contractAbi.Methods[method]
		//if !exist {
		//	return fmt.Errorf("method '%s' not found", method)
		//}
		//evmMethod = &m
		evmMethod = contractAbi

		inputData, err := util.Pack(contractAbi, method, params)
		if err != nil {
			return err
		}

		inputDataHexStr := hex.EncodeToString(inputData)
		method = inputDataHexStr[0:8]

		kvs = []*common.KeyValuePair{
			{
				Key:   "data",
				Value: []byte(inputDataHexStr),
			},
		}
	} else {
		if params != "" {
			kvsMap := make(map[string]string)
			err := json.Unmarshal([]byte(params), &kvsMap)
			if err != nil {
				return err
			}
			kvs = util.ConvertParameters(kvsMap)
		}
	}

	DispatchTimes(client, contractName, method, kvs, evmMethod)
	return nil
}

func getUserContract() error {
	client, err := util.CreateChainClient(sdkConfPath, chainId, orgId, userTlsCrtFilePath, userTlsKeyFilePath,
		userSignCrtFilePath, userSignKeyFilePath)
	if err != nil {
		return err
	}
	defer client.Stop()

	if contractAddress != "" {
		contractName = contractAddress
	}
	if contractName == "" {
		return errors.New("either contract-name or contract-address must be set")
	}

	var kvs []*common.KeyValuePair
	var contractAbi *abi.ABI

	if abiFilePath != "" { // abi file path 非空 意味着调用的是EVM合约
		abiBytes, err := ioutil.ReadFile(abiFilePath)
		if err != nil {
			return err
		}

		contractAbi, err = abi.JSON(bytes.NewReader(abiBytes))
		if err != nil {
			return err
		}

		inputData, err := util.Pack(contractAbi, method, params)
		if err != nil {
			return err
		}

		inputDataHexStr := hex.EncodeToString(inputData)

		kvs = []*common.KeyValuePair{
			{
				Key:   "data",
				Value: []byte(inputDataHexStr),
			},
		}
	} else {
		if params != "" {
			kvsMap := make(map[string]string)
			err := json.Unmarshal([]byte(params), &kvsMap)
			if err != nil {
				return err
			}
			kvs = util.ConvertParameters(kvsMap)
		}
	}

	resp, err := client.QueryContract(contractName, method, kvs, -1)
	if err != nil {
		return fmt.Errorf("query contract failed, %s", err.Error())
	}

	if resp.Code != common.TxStatusCode_SUCCESS {
		util.PrintPrettyJson(resp)
		return nil
	}

	if contractAbi != nil && resp.ContractResult != nil {
		output, err := contractAbi.Unpack(method, resp.ContractResult.Result)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		util.PrintPrettyJson(types.EvmTxResponse{
			TxResponse: resp,
			ContractResult: &types.EvmContractResult{
				ContractResult: resp.ContractResult,
				Result:         fmt.Sprintf("%v", output),
			},
		})
	} else {
		util.PrintPrettyJson(resp)
	}
	return nil
}

func upgradeUserContract() error {
	client, err := util.CreateChainClient(sdkConfPath, chainId, orgId, userTlsCrtFilePath, userTlsKeyFilePath,
		userSignCrtFilePath, userSignKeyFilePath)
	if err != nil {
		return err
	}
	defer client.Stop()

	if contractAddress != "" {
		contractName = contractAddress
	}
	if contractName == "" {
		return errors.New("either contract-name or contract-address must be set")
	}

	adminKeys, adminCrts, adminOrgs, err := util.MakeAdminInfo(client, adminKeyFilePaths, adminCrtFilePaths, adminOrgIds)
	if err != nil {
		return err
	}

	rt, ok := common.RuntimeType_value[runtimeType]
	if !ok {
		return fmt.Errorf("unknown runtime type [%s]", runtimeType)
	}

	pairs := make(map[string]string)
	if params != "" {
		err := json.Unmarshal([]byte(params), &pairs)
		if err != nil {
			return err
		}
	}
	pairsKv := util.ConvertParameters(pairs)
	payload, err := client.CreateContractUpgradePayload(contractName, version, byteCodePath, common.RuntimeType(rt),
		pairsKv)
	if err != nil {
		return err
	}

	if gasLimit > 0 {
		var limit = &common.Limit{GasLimit: gasLimit}
		payload = client.AttachGasLimit(payload, limit)
	}

	endorsementEntrys, err := util.MakeEndorsement(adminKeys, adminCrts, adminOrgs, client, payload)
	if err != nil {
		return err
	}
	// 发送更新合约请求
	resp, err := client.SendContractManageRequest(payload, endorsementEntrys, timeout, syncResult)
	if err != nil {
		return fmt.Errorf(SEND_CONTRACT_MANAGE_REQUEST_FAILED_FORMAT, err.Error())
	}

	err = util.CheckProposalRequestResp(resp, false)
	if err != nil {
		return fmt.Errorf(CHECK_PROPOSAL_RESPONSE_FAILED_FORMAT, err.Error())
	}
	return createUpgradeUserContractOutput(resp)
}

func freezeOrUnfreezeOrRevokeUserContract(which int) error {
	client, err := util.CreateChainClient(sdkConfPath, chainId, orgId, userTlsCrtFilePath, userTlsKeyFilePath,
		userSignCrtFilePath, userSignKeyFilePath)
	if err != nil {
		return err
	}
	defer client.Stop()

	if contractAddress != "" {
		contractName = contractAddress
	}
	if contractName == "" {
		return errors.New("either contract-name or contract-address must be set")
	}

	adminKeys, adminCrts, adminOrgs, err := util.MakeAdminInfo(client, adminKeyFilePaths, adminCrtFilePaths, adminOrgIds)
	if err != nil {
		return err
	}

	var (
		payload        *common.Payload
		whichOperation string
	)

	switch which {
	case 1:
		payload, err = client.CreateContractFreezePayload(contractName)
		whichOperation = "freeze"
	case 2:
		payload, err = client.CreateContractUnfreezePayload(contractName)
		whichOperation = "unfreeze"
	case 3:
		payload, err = client.CreateContractRevokePayload(contractName)
		whichOperation = "revoke"
	default:
		err = fmt.Errorf("wrong which param")
	}
	if err != nil {
		return fmt.Errorf("create cert manage %s payload failed, %s", whichOperation, err.Error())
	}

	endorsementEntrys, err := util.MakeEndorsement(adminKeys, adminCrts, adminOrgs, client, payload)
	if err != nil {
		return err
	}
	// 发送创建合约请求
	resp, err := client.SendContractManageRequest(payload, endorsementEntrys, timeout, syncResult)
	if err != nil {
		return fmt.Errorf(SEND_CONTRACT_MANAGE_REQUEST_FAILED_FORMAT, err.Error())
	}

	err = util.CheckProposalRequestResp(resp, false)
	if err != nil {
		return fmt.Errorf(CHECK_PROPOSAL_RESPONSE_FAILED_FORMAT, err.Error())
	}
	util.PrintPrettyJson(resp)
	return nil
}
