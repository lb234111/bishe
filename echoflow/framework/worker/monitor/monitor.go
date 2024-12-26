package monitor

import (
	"framework/utils"
	"framework/worker/clientList"
	"log"
	"math"
)

type Monitor struct {
	clients clientList.ClientList
}

var NodeId int
var logger *log.Logger

func init() {
	logger = utils.InitLogger("txloadtest")
}

func NewMonitor(clients clientList.ClientList) *Monitor {
	return &Monitor{clients: clients}
}

func (monitor *Monitor) CalPerformance(start, end int) (float64, error) {
	clients := monitor.clients
	logger.Println("监测并计算链上性能指标")
	//计算tps
	tps, err := clients.CalTps(start, end)
	//计算延迟
	//同步发送交易计算延迟
	if err != nil {
		return 0.0, err
	}
	logger.Println("链上性能指标计算完成")
	return math.Round(tps), nil
}
