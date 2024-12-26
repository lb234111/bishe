package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"framework/common"
	"framework/master/algo"
	"framework/master/controller"
	"framework/utils"
	"framework/worker/adapter"
	"framework/worker/clientList"
	"framework/worker/monitor"
	"framework/worker/stresstest"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	mode       string
	pruneRound int
	logger     *log.Logger
)

func init() {
	flag.StringVar(&mode, "mode", "", "启动模式")
	flag.IntVar(&pruneRound, "finish", 3, "连续多少轮次没有性能增长就结束迭代")
	logger = utils.InitLogger("main")
	flag.Parse()
}
func main() {
	fmt.Println(mode)
	if mode == "master" {
		start := time.Now()
		//创建controlller，algo
		logger.Printf("初始化master,链类型:%s,每轮次生成配置节点数量:%d", common.ChainType, common.BatchSize)
		controllerObj := controller.NewController(common.WorkerIP)
		algoObj := algo.NewAlgo(common.ChainType, common.BatchSize, pruneRound)
		//创建master server
		var wg sync.WaitGroup
		var messageList []*common.Message = []*common.Message{}
		router := gin.New()
		router.POST("/post", func(ctx *gin.Context) {
			var body common.Message
			err := ctx.ShouldBindJSON(&body)
			if err != nil {
				fmt.Println("解析post错误：", err)
			}
			fmt.Println("收到worker回复")
			messageList = append(messageList, &body)
			ip := ctx.ClientIP()
			controllerObj.PutWoker(ip)
			wg.Done()
		})
		go router.Run(":11110")
		logger.Println("master初始化完成")
		logger.Printf("链类型:%s,worker数量:%d,每轮次生成配置节点数量:%d,调整参数数量:%d,调整参数列表:%v,优化性能指标:Tps",
			common.ChainType, len(common.WorkerIP), common.BatchSize, len(common.ValidConfigName), common.ValidConfigName)
		for {
			logger.Println("从算法模块处获取下一轮测试配置")
			sendNodeIds, sendConfigs := controllerObj.RecvFromAlgo(algoObj)
			fmt.Printf("下一轮测试节点id%v,节点配置%v\n", sendNodeIds, sendConfigs)
			wg.Add(len(sendNodeIds))
			logger.Println("controller发送配置给worker")
			controllerObj.SendtoWorker(sendNodeIds, sendConfigs)

			logger.Println("controller发送配置成功")

			logger.Println("等待worker回复测试结果")
			wg.Wait()
			recvNodeIds, recvPerformance := controllerObj.RecvFromWorker(messageList)
			logger.Printf("收到worker运行测试结果:nodeId%v,Tps%v\n", recvNodeIds, recvPerformance)
			messageList = []*common.Message{}
			logger.Println("controller发送结果给Algo模块")
			controllerObj.SendToAlgo(algoObj, recvNodeIds, recvPerformance)
			finish := algoObj.IsFinish()
			if finish {
				utils.PrintInfo("迭代结束")
				logger.Println("迭代结束")
				timeElapsed := time.Since(start).Minutes()
				logger.Printf("迭代共计用时%f分", timeElapsed)
				return
			}
		}

	} else if mode == "worker" {
		//初始化
		//1. 脚本启动长安链
		//2. 创建客户端clients
		//3. 创建adapter，monitor，server，txloadtest
		//1.
		//todo
		logger.Println("部署区块链")

		//2.
		logger.Println("创建客户端")
		clients, err := clientList.NewClients(common.ChainType)
		if err != nil {
			utils.PrintError("创建客户端失败")
			logger.Fatal(err)
		}
		//3.
		logger.Println("创建adapter")
		adapterObj := adapter.NewAdapter(clients, common.ChainType)
		logger.Println("创建monitor")
		monitorObj := monitor.NewMonitor(clients)
		logger.Println("创建压测模块")
		stresstestObj := stresstest.NewStresstest(clients)
		logger.Println("创建worker的server")
		engine := gin.New()
		engine.POST("/run", func(ctx *gin.Context) {
			var body common.Message
			err := ctx.ShouldBindJSON(&body)
			if err != nil {
				fmt.Println("解析错误：", err)
				logger.Fatal()
			}
			utils.PrintInfo(fmt.Sprintf("接收到信息%+v\n", body))
			adapterObj.AdapterRun(body.Config)
			start, end, err := stresstestObj.SendTx()
			if err != nil {
				utils.PrintError("发送交易出错")
				logger.Fatal(err)
			}
			tps, err := monitorObj.CalPerformance(start, end)
			if err != nil {
				utils.PrintError("计算性能指标出错")
				logger.Fatal(err)
			}
			sendMessage := &common.Message{Config: body.Config, NodeId: body.NodeId, TPS: tps}
			messageJson, err := json.Marshal(sendMessage)
			if err != nil {
				utils.PrintError("配置序列化错误")
				logger.Fatal(err)
			}
			utils.PrintInfo("给master发送回复")
			hostname := fmt.Sprintf("http://%s:11110", common.MasterIP)
			http.Post(hostname+"/post", "application/json", bytes.NewBuffer(messageJson))
		})
		engine.Run(":11110")

		//后续操作
	} else {
		fmt.Println("输入正确的MODE")
		return
	}
}
