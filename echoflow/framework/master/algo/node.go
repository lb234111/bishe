package algo

import (
	"framework/common"
)

// 节点状态
const (
	TBD      int = 0
	ACTIVE   int = 1
	DECLINE  int = 2
	INACTIVE int = 3
	PRUNE    int = 4
)

// 节点数据结构
type Node struct {
	//节点相关信息
	NodeId           int     //节点编号
	PreNode          *Node   //该节点的父节点
	ChildNode        []*Node //该节点的子节点
	Round            int     //节点所处的迭代轮次
	State            int
	Performance      float64
	ChangeConfigName string //相较父节点修改的配置名称

	//配置相关信息
	Config       common.Config       //节点的的控制参数值
	NodeWeight   float64             //节点的的控制参数值
	ConfigWeight common.ConfigWeight //该节点各控制参数对应的权重
	Step         common.Step         //调整步长
}

func (node *Node) createNewNode(nodeId int, round int, nodeMap NodeMap) *Node {
	newNode := new(Node)
	newNode.NodeId = nodeId
	newNode.PreNode = node
	newNode.Round = round
	newNode.Step = make(common.Step)
	newNode.State = node.State
	common.CopyStep(newNode.Step, node.Step)
	nodeMap[nodeId] = newNode
	return newNode
}

func (node *Node) createNewConfig() common.Config {
	newConfig := make(common.Config)
	common.CopyConfig(newConfig, node.Config)
	return newConfig
}

func (node *Node) updateConfig(configName string) {
	switch configName {
	case common.CHAINMAKER_BATCH_MAX_SIZE_KEY:
		node.Config[configName] += node.Step[configName]
		node.Step[common.CHAINMAKER_BLOCK_TX_CAPACITY_KEY] = node.Config[configName]
	default:
		node.Config[configName] += node.Step[configName]
	}
}

func (node *Node) updatePerformance(performance float64) {
	node.Performance = performance
}

func (node *Node) updateState() {
	preNode := node.PreNode
	if node.Performance >= preNode.Performance {
		if node.Performance-preNode.Performance >= preNode.Performance*common.THRESHOLD*2 {
			if node.State > 1 {
				logger.Printf("%d节点性能有显著增长,设置节点State-1", node.NodeId)
				node.State -= 1
			} else {
				logger.Printf("%d节点性能有显著增长,但由于父节点为Active,设置节点State不变", node.NodeId)

			}
		} else if node.Performance-preNode.Performance >= preNode.Performance*common.THRESHOLD {
			logger.Printf("%d节点性能增长,节点State不变", node.NodeId)
		} else {
			node.State += 1
			logger.Printf("%d节点性能无显著增长,节点State+1", node.NodeId)
		}
	} else {
		if preNode.Performance-node.Performance >= preNode.Performance*common.THRESHOLD*2 {
			logger.Printf("%d节点性能显著降低,裁剪该节点分支", node.NodeId)
			node.State = PRUNE
		} else {
			node.State += 1
			logger.Printf("%d节点性能下降不显著,节点state+1", node.NodeId)
		}
	}
}
func (node *Node) updateConfigWeightVal(configName string, val float64) {
	node.ConfigWeight[configName] = val
}
