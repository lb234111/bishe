# -*- coding = utf-8 -*-
# 每个交易的大小固定设置为200-512字节
tx_size = 200
# 交易池大小
tx_pool_size = 10000
# 模拟节点的数量
node_num = 4
# network_flag = 0，模拟中的所有节点使用固定自定义的带宽；network_flag=1，模拟中所有节点使用真实的节点的带宽数据
network_flag = 0
# 自定义的固定的节点带宽,单位是字节/秒,目前设置为100Mbps
bandwidth = 13107200
# tx_flag = 0 ，使用自动生成的交易大小固定的交易来模拟；tx_flag=1，使用真实的链上交易信息来模拟
tx_flag = 0
# 自动生成的交易数
tx_num = 40000
# 自动生成中每秒生成的交易数
tx_rate = 2000
# random_topology_flag = 1，表示gossip算法中，节点之间网络拓扑关系是随机生成的;random_topology_flag = 0，使用json文件中自定义的拓扑结构
random_topology_flag = 1
#  gossips协议的网络拓扑结构
network_topology = []
# p2p的网络协议
protocol = "kad"
# 这个参数用于gossip协议
fan_out = 3
# 这个参数用于gossip协议，单位毫秒
heartbeat_interval = 100
# 这个参数用于kad协议
bucket_size = 2
# 每个区块的大小上线
block_size_limit = 100 * tx_size
# 区块哈希和签名大小之和均为一个固定值，sha256 32byte each
hash_sig_size = 64
# 第一区块的共识开始时间
start_time = 5
# 结束的时间  
over_time = 20
# 共识算法选择
consensus = "TBFT"
# 区块的基本处理耗时，总耗时为process_cost + tx_count * tx_cost
process_cost = 1 / 30
# 区块中的每笔交易的处理耗时
tx_cost = 1 / 13000
# 是否开启log功能
log_flag = True
# log是否输出到同一个文件中
single_log = True
# TBFT中的propose阶段的超时时间设置。如果超时则直接进入PreVote阶段
timeout_propose = 30
# TBFT中的delta值，每轮超时时间叠加
timeout_propose_delta = 1
# 下面同理，对应着TBFT共识算法中的PreVote、PreCommit、Commit
timeout_pre_vote = 30
timeout_pre_vote_delta = 1
timeout_pre_commit = 30
timeout_pre_commit_delta = 1
# 平均出一个块的时间,默认为以太坊的
block_interval = 2
# 是否关闭模拟器
shutdown = False
