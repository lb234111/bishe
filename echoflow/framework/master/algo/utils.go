package algo

import (
	"encoding/json"
	"fmt"
	"framework/common"
	"log"
	"os"
)

func writeNodeToFile(node *Node) {
	file, err := os.OpenFile(common.DataFilePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", *node)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)
	var pre int
	if node.PreNode == nil {
		pre = 0
	} else {
		pre = node.PreNode.NodeId
	}
	nodeInfo := struct {
		NodeId int
		Round  int
		Pre    int
		State  int
		Tps    float64
		Config common.Config
	}{node.NodeId, node.Round, pre, node.State, node.Performance, node.Config}
	data, err := json.MarshalIndent(nodeInfo, "", "\t")
	if err != nil {
		fmt.Println("转换失败")
		return
	}
	_, err = file.Write(data)
	if err != nil {
		log.Fatal(err)
	}

}
