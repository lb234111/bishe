package adapter

import (
	"errors"
	"fmt"
	"framework/common"
	"framework/utils"
	"framework/worker/clientList"
	"log"
	"os/exec"
)

type Adapter struct {
	clients   clientList.ClientList
	chainType string
}

var Config common.Config //该适配器当前要调整的参数值
var logger *log.Logger

func init() {
	logger = utils.InitLogger("adapter")
}

func NewAdapter(clients clientList.ClientList, chainType string) *Adapter {
	return &Adapter{
		clients:   clients,
		chainType: chainType,
	}
}

func (adapter *Adapter) AdapterRun(config common.Config) {
	logger.Println("开始调整参数")
	for configName, configVal := range config {
		err := adapter.adaptConfig(configName, configVal)
		if err != nil {
			utils.PrintInfo("调整参数错误", err)
			logger.Fatal(err)
		}
	}
	logger.Println("调整参数完成")

}

func (adapter *Adapter) adaptConfig(configName string, configVal interface{}) error {
	switch adapter.chainType {
	case "chainmaker":
		var command string
		switch configVal.(type) {
		case int:
			configVal = configVal.(int)
			command = fmt.Sprintf("curl -d '{\"module\":%d,\"type\":%d,\"param_name\":\"%s\",\"param_value\":%d}' -H \"Content-Type:application/json\" -X POST http://localhost:8000", 3, 1, configName, configVal)

		case float64:
			configVal = configVal.(float64)
			command = fmt.Sprintf("curl -d '{\"module\":%d,\"type\":%d,\"param_name\":\"%s\",\"param_value\":%f}' -H \"Content-Type:application/json\" -X POST http://localhost:8000", 3, 1, configName, configVal)
		default:
			return errors.New("参数值类型异常")
		}
		cmd := exec.Command("bash", "-c", command)
		if err := cmd.Run(); err != nil {
			logger.Fatal("调整链上参数异常", err)
		}
	default:
		return errors.New("未知的链类型")
	}

	return nil
}
