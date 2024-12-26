// 算法模块
package algo

import (
	"fmt"
	"framework/common"
	"framework/utils"
	"log"
	"math"
	"math/rand"
	"os"
	"time"
)

// 日志配置
var logger *log.Logger
var initWeight float64

// var initTps int

func init() {
	logger = utils.InitLogger("algo")
}

// 算法模块
type Algo struct {
	nodeWeightMap *NodeWeightMap
	weightCache   WeightCache
	nodeMap       NodeMap
	chainType     string
	round         int //当前轮次
	nodeId        int //当前最新的节点ID
	// lastRoundNodeId []int
	// lastRoundTps    []float64
	batchSize     int
	curPruneRound int
	pruneRound    int

	bestPerformance     float64
	continuousNoImprove int
}

//节点权重缓存，用于在测试后，给新节点赋予权重使用
type WeightCache map[int]common.ConfigWeight

// 节点ID和节点映射
type NodeMap map[int]*Node

// 节点和节点权重映射
type NodeWeightMap struct {
	weightMap map[*Node]float64
	weightSum float64
}

func NewAlgo(chainType string, batchSize int, pruneRound int) *Algo {
	logger.Printf("算法模块初始化, chainType:%s, batchSize:%d, pruneRound:%d\n", chainType, batchSize, pruneRound)
	return &Algo{
		nodeWeightMap:   &NodeWeightMap{weightMap: make(map[*Node]float64), weightSum: 0},
		weightCache:     make(map[int]common.ConfigWeight),
		nodeMap:         make(map[int]*Node),
		chainType:       chainType,
		batchSize:       batchSize,
		pruneRound:      pruneRound,
		bestPerformance: 0,
	}
}
func (algo *Algo) selectNodeAndConfig() map[int][]string {

	set := make(map[int]map[string]struct{}) //记录选取的节点及配置名
	res := make(map[int][]string)
	aliveNode := 0
	file, _ := os.OpenFile(common.AliveNodePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	defer file.Close()
	for _, w := range algo.nodeWeightMap.weightMap {
		if w != 0 {
			aliveNode += 1
		}
	}
	file.WriteString(fmt.Sprintf("%d\t%d\n", algo.round, aliveNode))
	logger.Printf("开始为轮次%d选择下一轮要调整的节点及配置\n", algo.round)
	for i := 0; i < algo.batchSize; {
		if algo.nodeWeightMap.weightSum == 0 {
			break
		}
		node := algo.selectNode()
		if _, ok := set[node.NodeId]; !ok {
			configChoice := make(map[string]struct{})
			set[node.NodeId] = configChoice
		}
		configName := algo.selectConfig(node)
		if _, ok := set[node.NodeId][configName]; !ok {
			set[node.NodeId][configName] = struct{}{}
			node.NodeWeight -= node.ConfigWeight[configName]
			algo.nodeWeightMap.weightMap[node] -= node.ConfigWeight[configName]
			algo.nodeWeightMap.weightSum -= node.ConfigWeight[configName]
			node.updateConfigWeightVal(configName, 0)
			i++
		}

		// var sum1 float64
		// randNum4Node := rand.Float64() * algo.nodeWeightMap.weightSum

		// //根据weight选取节点
		// for node, weight := range algo.nodeWeightMap.weightMap {
		// 	sum1 += weight
		// 	if sum1 >= randNum4Node {
		// 		if _, ok := set[node.NodeId]; !ok {
		// 			configChoice := make(map[string]struct{})
		// 			set[node.NodeId] = configChoice
		// 		}
		// 		//根据weight选取修改参数
		// 		var sum2 float64
		// 		randNum4Config := rand.Float64() * weight
		// 		for k, v := range node.ConfigWeight {
		// 			sum2 += v
		// 			if sum2 >= randNum4Config {
		// 				if _, ok := set[node.NodeId][k]; !ok {
		// 					set[node.NodeId][k] = struct{}{}
		// 					node.UpdateConfigWeightVal(k, 0)
		// 					node.NodeWeight -= v
		// 					algo.nodeWeightMap.weightMap[node] -= v
		// 					algo.nodeWeightMap.weightSum -= v
		// 					i++
		// 					break
		// 				}
		// 			}
		// 		}
		// 		break
		// 	}
		// }
	}

	for nodeId, tmp := range set {
		configNames := []string{}
		for configName := range tmp {
			configNames = append(configNames, configName)
		}
		res[nodeId] = configNames
	}
	return res
}

func (algo *Algo) selectNode() *Node {
	var sum float64
	randNum4Node := rand.Float64() * algo.nodeWeightMap.weightSum
	for node, weight := range algo.nodeWeightMap.weightMap {
		sum += weight
		if sum >= randNum4Node {
			return node
		}
	}
	return nil
}
func (algo *Algo) selectConfig(node *Node) string {
	var sum float64
	randNum4Config := rand.Float64() * node.NodeWeight
	for k, v := range node.ConfigWeight {
		sum += v
		if sum >= randNum4Config {
			return k
		}
	}
	return ""
}

func (algo *Algo) createRootNode() (*Node, error) {
	logger.Printf("为轮次1生成根节点")
	node := new(Node)
	node.NodeId = 1
	node.Round = 1
	config, err := common.CreateInitConfig(algo.chainType)
	if err != nil {
		return nil, err
	}
	node.Config = config
	step, err := common.GetStep(algo.chainType)
	if err != nil {
		return nil, err
	}
	node.Step = step
	node.ConfigWeight = common.ConfigWeight{}
	return node, nil
}

func (algo *Algo) createNewNodes(config2update map[int][]string) []*Node {
	newNodes := []*Node{}
	for nodeId, configNames := range config2update {
		oldNode := algo.nodeMap[nodeId]
		for _, configName := range configNames {
			algo.nodeId++
			newNode := oldNode.createNewNode(algo.nodeId, algo.round, algo.nodeMap)
			newNode.ChangeConfigName = configName
			newNode.Config = oldNode.createNewConfig()
			newNode.updateConfig(configName)
			oldNode.ChildNode = append(oldNode.ChildNode, newNode)
			newNodes = append(newNodes, newNode)
			algo.nodeMap[algo.nodeId] = newNode
		}
	}
	return newNodes
}
func (algo *Algo) algorun() ([]*Node, error) {
	//1. 如果是初始化
	algo.round++
	logger.Printf("轮次%d:算法模块开始运行", algo.round)
	if algo.round == 1 {
		algo.nodeId++
		newNode, err := algo.createRootNode()
		if err != nil {
			utils.PrintError("创建根节点失败")
			return nil, err
		}
		algo.nodeMap[1] = newNode
		return []*Node{newNode}, nil
	}
	//2. 选取节点
	config2update := algo.selectNodeAndConfig()
	logger.Printf("轮次%d要更新的节点及配置:%+v", algo.round, config2update)
	//3. 生成新节点
	newNodes := algo.createNewNodes(config2update)
	return newNodes, nil
}

func (algo *Algo) SendToController() ([]int, []common.Config) {
	algoStart := time.Now().UnixMilli()
	newNode, err := algo.algorun()
	algoEnd := time.Now().UnixMilli()
	logger.Printf("round%d生成配置用时%dms", algo.round, algoEnd-algoStart)
	if err != nil {
		utils.PrintError("生成节点失败")
		return nil, nil
	}
	nodeIds := []int{}
	configs := []common.Config{}
	for _, node := range newNode {
		nodeIds = append(nodeIds, node.NodeId)
		configs = append(configs, node.Config)
	}
	logger.Printf("轮次%d,发送节点%+v\n", algo.round, newNode)
	return nodeIds, configs
}

// 从controlller接受到nodeId和tps后进行操作
func (algo *Algo) RecvFromController(nodeIds []int, performance []float64) {
	//第一轮次处理方法
	if algo.round == 1 {
		//根据id获取节点
		node := algo.nodeMap[1]
		//更新节点TPS
		node.updatePerformance(performance[0])
		//更新节点状态
		node.State = ACTIVE
		//更新节点权重
		// initTps = tps[0]
		initWeight = math.Round(2 * common.THRESHOLD * performance[0])
		for k := range node.Config {
			node.ConfigWeight[k] = initWeight
			node.NodeWeight += initWeight
		}
		//更新节点及权重map
		algo.insertNodeWeightMap(node)
		// fmt.Printf("%+v", node.ConfigWeight)
		tmp := make(common.ConfigWeight)
		common.CopyConfigWeight(tmp, node.ConfigWeight)
		algo.weightCache[1] = tmp
	} else {
		//后续轮次处理
		for i := 0; i < len(nodeIds); i++ {
			//根据id获取节点
			node := algo.nodeMap[nodeIds[i]]
			//更新节点TPS
			node.updatePerformance(performance[i])
			//更新节点状态
			node.updateState()
			//更新节点权重及映射map
			algo.updateWeight(node)
			tmp := make(common.ConfigWeight)
			common.CopyConfigWeight(tmp, node.ConfigWeight)
			algo.weightCache[nodeIds[i]] = tmp
		}
	}
	//结果写入文件
	for _, nodeId := range nodeIds {
		writeNodeToFile(algo.nodeMap[nodeId])
	}
	algo.updateContinuousNoImprove(performance)
	algo.updatePruneRound(nodeIds)
}

func (algo *Algo) insertNodeWeightMap(node *Node) {
	algo.nodeWeightMap.weightMap[node] = node.NodeWeight
	algo.nodeWeightMap.weightSum += node.NodeWeight
}

func (algo *Algo) updateWeight(node *Node) {
	if node.State == PRUNE {
		node.ConfigWeight = nil
		node.NodeWeight = 0
		logger.Println("裁剪该节点,nodeId =", node.NodeId)
		return
	}
	node.ConfigWeight = make(common.ConfigWeight)
	common.CopyConfigWeight(node.ConfigWeight, algo.weightCache[node.PreNode.NodeId])
	algo.weightCache[node.PreNode.NodeId][node.ChangeConfigName] = 0
	var newWeight float64
	if node.State == ACTIVE {
		newWeight = math.Round(math.Max(node.Performance-node.PreNode.Performance, initWeight)) * 2
	} else {
		if node.Performance > node.PreNode.Performance {
			newWeight = math.Round(math.Max((node.Performance - node.PreNode.Performance), initWeight))
		} else {
			// newWeight = math.Round(math.Min(math.Abs(node.TPS-node.PreNode.TPS), initWeight))
			newWeight = math.Round(math.Min(node.PreNode.Performance*common.THRESHOLD-(node.PreNode.Performance-node.Performance)/2, initWeight))
		}
	}
	logger.Printf("节点%d的%s配置赋予新权重%.1f,节点状态为%d", node.NodeId, node.ChangeConfigName, newWeight, node.State)
	node.updateConfigWeightVal(node.ChangeConfigName, newWeight)
	var nodeWeight float64
	for _, v := range node.ConfigWeight {
		nodeWeight += v
	}
	node.NodeWeight = nodeWeight
	algo.insertNodeWeightMap(node)
}

func (algo *Algo) updatePruneRound(nodeIds []int) {
	for _, nodeId := range nodeIds {
		if algo.nodeMap[nodeId].State != PRUNE {
			algo.curPruneRound = 0
			return
		} else {
			continue
		}
	}
	algo.curPruneRound += 1
}

func (algo *Algo) updateContinuousNoImprove(performance []float64) {
	maxtmp := utils.GetMaxFromAry(performance)
	// logger.Println("tps,maxtmp", tps, maxtmp)
	if maxtmp <= algo.bestPerformance {
		logger.Printf("轮次%d，该轮性能无增长", algo.round)
		algo.continuousNoImprove += 1
	} else {
		algo.continuousNoImprove = 0
		algo.bestPerformance = maxtmp
		logger.Printf("轮次%d，该轮性能有增长,当前最大Tps：%f", algo.round, algo.bestPerformance)
	}
}
func (algo *Algo) IsFinish() bool {
	logger.Printf("当前轮次：%d,连续性能无增长：%d", algo.round, algo.continuousNoImprove)
	if algo.continuousNoImprove == 15 {
		return true
	}
	if algo.curPruneRound == algo.pruneRound {
		return true
	} else {
		return false
	}
}
