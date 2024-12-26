package main

import (
	sdk "chainmaker.org/chainmaker/sdk-go/v3"
)

// createClientWithConfig 通过配置创建客户端
func createClientWithConfig(consensusType int32) (*sdk.ChainClient, error) {
	var chainClient *sdk.ChainClient
	var err error
	// dpos
	if consensusType == 3 {
		chainClient, err = sdk.NewChainClient(sdk.WithConfPath("./configs/client_config_dpos.yml"))
		if err != nil {
			return nil, err
		}
		return chainClient, nil
	} else {
		chainClient, err = sdk.NewChainClient(sdk.WithConfPath("./configs/client_config.yml"))
		if err != nil {
			return nil, err
		}
	}
	// 启用证书压缩（开启证书压缩可以减小交易包大小，提升处理性能）
	err = chainClient.EnableCertHash()
	if err != nil {
		return nil, err
	}
	return chainClient, nil
}

// // createNode 创建节点
// func createNode(nodeAddr string, connCnt int) *sdk.NodeConfig {
// 	node := sdk.NewNodeConfig(
// 		// 节点地址，格式：127.0.0.1:12301
// 		sdk.WithNodeAddr(nodeAddr),
// 		// 节点连接数
// 		sdk.WithNodeConnCnt(connCnt),
// 		// 节点是否启用TLS认证
// 		sdk.WithNodeUseTLS(true),
// 		// 根证书路径，支持多个
// 		sdk.WithNodeCAPaths(caPaths),
// 		// TLS Hostname
// 		sdk.WithNodeTLSHostName(tlsHostName),
// 	)
// 	return node
// }

// // createClient 通过参数创建客户端
// func createClient() (*sdk.ChainClient, error) {
// 	var node1, node2 *sdk.NodeConfig
// 	if node1 == nil {
// 		// 创建节点1
// 		node1 = createNode(nodeAddr1, connCnt1)
// 	}
// 	if node2 == nil {
// 		// 创建节点2
// 		node2 = createNode(nodeAddr2, connCnt2)
// 	}

// 	chainClient, err := sdk.NewChainClient(
// 		// 设置归属组织
// 		sdk.WithChainClientOrgId(chainOrgId),
// 		// 设置链ID
// 		sdk.WithChainClientChainId(chainId),
// 		// 设置logger句柄，若不设置，将采用默认日志文件输出日志
// 		sdk.WithChainClientLogger(getDefaultLogger()),
// 		// 设置客户端用户私钥路径
// 		sdk.WithUserKeyFilePath(userKeyPath),
// 		// 设置客户端用户证书
// 		sdk.WithUserCrtFilePath(userCrtPath),
// 		// 添加节点1
// 		sdk.AddChainClientNodeConfig(node1),
// 		// 添加节点2
// 		sdk.AddChainClientNodeConfig(node2),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// 启用证书压缩（开启证书压缩可以减小交易包大小，提升处理性能）
// 	err = chainClient.EnableCertHash()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return chainClient, nil
// }
