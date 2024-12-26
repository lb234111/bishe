package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

// tbftParam 记录 tbft 算法中的参数是否已更新过
var tbftParam map[string]interface{}
var mutex sync.Mutex

// paramAdapter 解析参数
func paramAdapter(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var param ParamAdapter
	err = json.Unmarshal(body, &param)
	if err != nil {
		panic(err)
	}
	log.Println(param)
	// 判断修改参数逻辑
	// 修改共识层参数
	var block_height uint64
	if param.Module == 1 {
		block_height, err = chooseConsensusType(param)
		if err != nil {
			fmt.Fprintln(w, err)
		}
	}
	// 修改存储层参数 todo
	if param.Module == 2 {
		err = chooseStorageType(param)
		if err != nil {
			fmt.Fprintln(w, err)
		}
	}
	// 修改网络层参数 todo
	if param.Module == 3 {
		block_height, err = chooseNetworkType(param)
		if err != nil {
			fmt.Fprintln(w, err)
		}
	}
	fmt.Fprintln(w, "success!")
	fmt.Fprintln(w, block_height)
}

// StartSvr 启动并监听
func StartSvr() {
	http.HandleFunc("/", paramAdapter)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func main() {
	tbftParam = make(map[string]interface{})
	log.Println("server start!")
	StartSvr()
}
