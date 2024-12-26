# param_adapter

# 长安链参数动态调整服务

### 共识层 tbft 算法
propose_timeout  
propose_delta_timeout  
blocks_per_proposer  

### 共识层 dpos 算法
epoch_validator_num  
3.0.0底链才支持  

### 共识层 maxbft 算法
MaxBftViewNumsPerEpoch  
3.0.0底链才支持  

### （后续会继续更新其他算法，存储层，网络层的参数修改）

## 使用说明：
下载代码到本地，其中 main.go 的监听端口可以自行修改  
configs 文件夹下的 client_config.yml 和 client_config_dpos.yml 文件需要根据不同算法自行修改对应的密钥、证书路径  
struct.go 中的 users、permissionedPkUsers、pkUsers 也需要修改对应的密钥路径  
全部修改好后，go build，然后 ./param_adapter 启动服务  
可以在命令行体验效果，具体示例代码可以参考 testdata.txt
