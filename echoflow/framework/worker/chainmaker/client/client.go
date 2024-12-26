// coding:utf-8
package client

import (
	"fmt"
	"framework/worker/chainmaker/config"
	"framework/worker/chainmaker/logger"

	sdk "chainmaker.org/chainmaker/sdk-go/v2"
)

// 创建链客户端
// @param sdkConfPath:客户端sdk路径
// @param SugaredLogger:客户端日志选项
func createClient(sdkConfPath string) (*sdk.ChainClient, error) {
	sdkLogger := logger.NewLogger("[SDK]", &config.LogConfig{LogPath: "./sdk.log"})
	cc, err := sdk.NewChainClient(
		sdk.WithConfPath(sdkConfPath),
		// 启用后使用异步订阅，而不是轮询方式
		sdk.WithChainClientLogger(sdkLogger),
		sdk.WithEnableTxResultDispatcher(true),
	)
	if err != nil {
		info := fmt.Sprintf("sdk.NewChainClient: %s", err)
		return nil, fmt.Errorf(info)
	}
	if cc.GetAuthType() == sdk.PermissionedWithCert {
		if err := cc.EnableCertHash(); err != nil {
			info := fmt.Sprintf("cc.EnableCertHash: %s", err)
			return nil, fmt.Errorf(info)
		}
	}
	return cc, nil
}

func NewClient(configPath string) (*sdk.ChainClient, error) {
	// 建立与链交互的客户端
	client, err := createClient(configPath)
	if err != nil {
		fmt.Println("configPath: " + configPath)
		info := fmt.Sprintf("chainmakerClient NewClient CreateClient: %s", err)
		return nil, fmt.Errorf(info)
	}
	return client, nil
}
