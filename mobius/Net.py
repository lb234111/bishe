# -*- coding = utf-8 -*-
from PBFT.PBFT_Node import PBFTNode
from HotStuff.HotStuff_Node import HotStuffNode
from TBFT.TBFT_Node import TBFTNode
from POW.POW_Node import POWNode
import random
from Transaction import *
import math
import config
import Log


class Gossip:
    def __init__(self, fan_out, node_num):
        # gossip协议中需要给每个节点都定义fan_out个固定的相邻节点，后续消息转发都将固定转发给这fan_out个节点
        self.name = "gossip"
        self.node_num = node_num
        self.fan_out = fan_out
        if config.random_topology_flag == 1:
            self.neighbor = self.init_neighbor()
        else:
            self.neighbor = config.network_topology

    # 根据gossip协议随机生成fan_out个节点的id，这些节点将作为node_id节点的相邻节点
    def generate_neighbor(self, node_id):
        receiver = []
        fan_out = config.fan_out
        if fan_out >= config.node_num:
            fan_out = config.node_num - 1
        for i in range(fan_out):
            receiver.append((node_id + i + 1) % config.node_num)
        return receiver
        # receiver = []
        # while len(receiver) != self.fan_out:
        #     item = random.randint(0, self.node_num - 1)
        #     if item not in receiver and item != node_id:
        #         receiver.append(item)
        # return receiver

    def init_neighbor(self):
        neighbor = []
        for i in range(self.node_num):
            neighbor.append(self.generate_neighbor(i))
        return neighbor


class KadCast:
    def __init__(self, node_num, size):
        self.name = "kad"
        self.id_len = 0
        self.node_num = node_num
        self.node_id = self.init_node_id()
        self.bucket_size = size
        self.forward_table = self.init_forward_table()

    def init_node_id(self):
        # 根据节点的数量来确定id的长度
        up_limit = self.node_num - 1
        up_str = format(up_limit, "b")
        self.id_len = len(up_str)
        str_list = []
        # 每个节点的id就是其十进制的id号转化为二进制
        for i in range(up_limit):
            str_temp = format(i, "b")
            while len(str_temp) < self.id_len:
                str_temp = "0" + str_temp
            str_list.append(str_temp)
        str_list.append(up_str)
        return str_list

    def init_forward_table(self):
        # 生成转发表，该表是一个三维数组，第一维表示的第几个节点的转发表，第二维表示这个节点的第几个桶，第三维中就是数据，即这个桶中的节点
        forward_table = []
        for i in range(self.node_num):
            forward_table.append(self.get_bucket(i, self.id_len))
        return forward_table

    def get_bucket(self, d, size):
        base = 0
        forward = []
        for dis in range(size):
            low_low = base
            low_high = base + 2 ** (size - dis - 1) - 1
            if low_high >= self.node_num - 1:
                continue
            high_low = low_high + 1
            high_high = base + 2 ** (size - dis) - 1
            if high_high > self.node_num - 1:
                high_high = self.node_num - 1
            temp = []
            if low_low <= d <= low_high:
                base = low_low
                item = high_low
                up_limit = high_high
            else:
                base = high_low
                item = low_low
                up_limit = low_high
            while item <= up_limit:
                temp.append(item)
                if len(temp) == self.bucket_size:
                    break
                item += 1
            forward.append(temp)
        return forward

    def gen_recv(self, node_id, dis):
        # 节点收到dis的消息，转而进行转发，则需要转发给自身转发表中distance>=dis的节点
        num = dis
        recv_list = []
        while num < len(self.forward_table[node_id]):
            recv_list += self.forward_table[node_id][num]
            num += 1
        return recv_list


# 网络中定义了节点的个数，gossip协议的fan_out的大小
# 初始化网络中没有节点，leader从0号node开始轮流担任
class Network:
    # 节点数量
    node_num = 0
    # 全网的节点列表
    node = []
    # 网络协议的选择
    protocol = None
    # 共识算法的选择
    consensus = None

    # 网络中所有的交易
    transaction_list = []
    # 区块数量
    block_num = 0
    # 网络中的所有区块
    block = []
    # 阈值2f+1
    threshold = 0
    # 全网算力总和
    hash_power = 0
    # 出块间隔
    # block_interval = 10

    @classmethod
    def reset(cls):
        cls.node_num = 0
        cls.node = []
        cls.protocol = None
        cls.consensus = None
        cls.transaction_list = []
        cls.block_num = 0
        cls.block = []
        cls.threshold = 0
        cls.hash_power = 0
        config.shutdown = False

    @classmethod
    def init(cls, node_num, consensus, protocol, fan_out=0, bucket_size=0):
        cls.consensus = consensus
        cls.init_node(node_num)
        cls.logger = Log.Log(-1)
        cls.logger.file_handle()
        if config.consensus == "POW":
            cls.hash_power = sum([miner.hash_power for miner in Network.node])
        if protocol == "gossip":
            cls.protocol = Gossip(fan_out, node_num)
        elif protocol == "kad":
            cls.protocol = KadCast(node_num, bucket_size)
        # cls.block_interval = block_interval

    @classmethod
    def init_node(cls, num):
        cls.node_num = num
        cls.init_threshold()
        size = config.tx_pool_size
        for i in range(num):
            if cls.consensus == "PBFT":
                cls.node.append(PBFTNode(i, size))
            elif cls.consensus == "HotStuff":
                cls.node.append(HotStuffNode(i, size))
            elif cls.consensus == "TBFT":
                cls.node.append(TBFTNode(i, cls.threshold, size))
            elif cls.consensus == "POW":
                cls.node.append(POWNode(i, config.tx_pool_size, random.randint(30, 50)))
                # block = Block(block_id = 0, previous = -1)
                # cls.node.blockchain.append(block)
            # 为了测试就给所有节点都设置随机的带宽大小
            # cls.node[i].set_bandwidth(random.randint(500, 1000))
            cls.node[i].set_bandwidth(config.bandwidth)

    @classmethod
    def add_transaction(cls, tx_id, sender, receiver, size):
        new_transaction = Transaction(tx_id, sender, receiver, size)
        cls.transaction_list.append(new_transaction)

    @classmethod
    def add_block(cls, block, height=-1):
        if height != -1:
            if cls.block_num == height+1:
                cls.block[height] = block
            else:
                cls.block.append(block)
                cls.block_num = len(cls.block)
        else:
            cls.block.append(block)
            cls.block_num = len(cls.block)

    @classmethod
    def init_threshold(cls):
        cls.threshold = math.ceil(((cls.node_num - 1) / 3) * 2) + 1

    @classmethod
    def gen_recv(cls, node_id, event):
        if cls.protocol.name == "gossip":
            return cls.protocol.neighbor[node_id]
        elif cls.protocol.name == "kad":
            return cls.protocol.gen_recv(node_id, event.dis)
