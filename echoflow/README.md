## 目录结构描述

```shell
.
|-- README.md
|-- chainmaker-cryptogen	#长安链证书生成工具
|-- chainmaker-csv-adapter	#针对指标采修改过的长安链
|   |-- ...
|   |-- submodule			#修改过的模块
|   |   |--  consensus-maxbft
|   |   |--  consensus-tbft-v2.4.2
|   |   |-- libp2p-pubsub@v1.1.3
|   |   |-- localconf
|   |   |-- net-libp2p@v1.2.0
|   |   |-- net-liquid@v1.1.0
|   |   |-- protocol
|   |   |-- recorderfile	#指标采集的工具包，其他模块通过调用该包实现指标采集
|   |   |-- store
|   |   |-- txpool-batch
|   |   |-- txpool-normal
|   |   |-- txpool-single
|   |-- ...
|-- deploytool			#长安链部署工具
|-- framework			#性能边界框架
|   |-- README.md
|   |-- common			#通用数据结构以及配置
|   |   |-- config.go	#配置文件
|   |   |-- nodeConfigs.go
|   |   |-- struct.go
|   |   |-- worker.go
|   |-- go.mod
|   |-- go.sum
|   |-- main.go
|   |-- master			#master层代码
|   |   |-- algo		#参数优化模块
|   |   |-- controller	#控制模块
|   |-- utils			#工具模块
|   |   |-- files.go
|   |   |-- logger.go
|   |   |-- print.go
|   |   |-- utils.go
|   |-- worker			#Worker层代码
|       |-- adapter		#参数调整模块，结合下面param_adapter使用
|       |-- chainmaker	#长安链的sdk的调用
|       |-- clientList	#客户端类
|       |-- monitor		#性能监测模块
|       |-- stresstest	#压力测试模块
|-- network_script.sh   #端口带宽限制脚本
|-- param_adapter		#针对长安链的链上参数调整的工具
```

## 链部署

安装好长安链的部署环境

`cd deploytool`，进入deploy文件夹下，`make`编译长安链并部署，部署完毕后，进入`deploytool/scripts`下`./chain_start.sh`启动

### 框架启动

### 1. 配置

进行worker和master对应的配置

```go
var (
	//for master，针对master的配置，在对应机器修改后启动
	WorkerIP      []string = []string{"10.15.9.14", "10.15.9.15", "10.15.9.16", "10.15.9.17"}	//worker的IP
	DataFilePath  string   = "./data_0413/nodeinfo" + time.Now().Format("2006010215") + ".json"	//配置以及测得性能的存储文件地址
	AliveNodePath string   = "./data_0413/aliveNode" + time.Now().Format("2006010215") + ".dat" //算法中记录还有多少个活跃节点数量的，不重要，我用来画图的
	BatchSize     int      = 10	//每次迭代生成多少组参数配置

	//for worker 针对worker的配置在对应机器修改后启动
	MasterIP string = "10.15.9.37"
	//for tx，worker中发送交易的特征
	TransactionType        string = "store"  //交易类型，即合约名称
	TransactionCount       int    = 40000	//发送交易总数
	TransactionRate        int    = 4000	//发送交易速率，即每秒发多少笔
	Connections            int    = 40		//并发数量
	ChainmakerContractPath string = "./worker/chainmaker/contracts"	//合约文件的地址

	//for all
	ChainType string = "chainmaker"	//链类型
)
```

### 2. 启动前准备

进入param_adapter，go build编译

### 3. 启动

1. 在每个要运行worker的机器上部署长安链并启动
2. 进入param_adapter目录下，启动param_adapter，可能会出现配置文件问题，将长安链的crypto_config文件夹软链接到该目录下。
3. 再打开一个终端，再framework目录下，`go run main.go -mode worker`，启动worker，第一次比较慢
4. worker启动完成后，再要运行master的机器上，在framework目录下运行`go run main.go -mode master`，系统开始迭代。
