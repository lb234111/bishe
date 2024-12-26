package controller

import (
	"bytes"
	"encoding/json"
	"framework/common"
	"framework/master/algo"
	"framework/utils"
	"log"
	"net/http"
)

type Controller struct {
	Workers   chan *common.Worker
	WorkerMap map[string]*common.Worker
}

var logger *log.Logger

func init() {
	logger = utils.InitLogger("controller")
}
func NewController(workersIP []string) *Controller {
	logger.Printf("初始化controller,worker数量%d", len(workersIP))
	utils.PrintInfo("初始化controller,worker数量：", len(workersIP), "workerIP:", workersIP)
	workers := make(chan *common.Worker, len(workersIP))
	workerMap := make(map[string]*common.Worker, len(workersIP))
	for _, IP := range workersIP {
		worker := common.NewWorker(IP)
		workers <- worker
		workerMap[IP] = worker
	}
	return &Controller{Workers: workers, WorkerMap: workerMap}
}

func (con *Controller) getWorker() common.Worker {
	return *<-con.Workers
}

func (con *Controller) PutWoker(ip string) {
	con.Workers <- con.WorkerMap[ip]
}

func (con *Controller) SendtoWorker(nodeIds []int, configs []common.Config) {
	for i := 0; i < len(nodeIds); i++ {
		worker := con.getWorker()
		sendMessage := &common.Message{Config: configs[i], NodeId: nodeIds[i]}
		messageJson, err := json.Marshal(sendMessage)
		if err != nil {
			utils.PrintError("配置反序列化错误", err)
			logger.Fatal("配置反序列化错误", err)
		}

		go func() {
			_, err := http.Post(worker.HostName+"/run", "application/json", bytes.NewBuffer(messageJson))
			if err != nil {
				utils.PrintError("发送配置错误", err)
				logger.Fatal("发送配置错误", err)
			}
		}()

	}
}

func (con *Controller) RecvFromWorker(messageList []*common.Message) ([]int, []float64) {
	var nodeIds []int
	var tps []float64
	for _, message := range messageList {
		nodeIds = append(nodeIds, message.NodeId)
		tps = append(tps, message.TPS)
	}
	return nodeIds, tps
}

func (con *Controller) SendToAlgo(algoMod *algo.Algo, nodeIds []int, tps []float64) {
	algoMod.RecvFromController(nodeIds, tps)
}

func (con *Controller) RecvFromAlgo(algoMod *algo.Algo) ([]int, []common.Config) {
	nodeIds, configs := algoMod.SendToController()
	return nodeIds, configs
}
