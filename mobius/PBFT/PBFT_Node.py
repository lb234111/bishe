# -*- coding = utf-8 -*-
import Log
import config
import operator
import re


def process_block(block):
    # 该函数为节点收到区块后进行验证的时间，如果仅仅只是转发prepare消息无需计算该部分时间
    return config.process_cost


class PBFTNode:

    def __init__(self, id, tx_pool_size):
        self.id = id
        self.bandwidth = 0
        # self.upload_rate = self.bandwidth / 2
        # self.download_rate = self.bandwidth / 2
        self.logger = Log.Log(id)
        self.logger.file_handle()
        self.tx_pool = []
        self.tx_done_pool = []
        self.tx_pool_size = tx_pool_size
        self.view = 0
        self.execute_time = 0
        # self.execute_time_down = 0
        self.event_list = []
        # 由于PBFT共识算法下不会出现分叉现象，因此不同链上的区块的差别仅是共识达成的时间而已，其余交易信息都统一存放在Network中即可
        self.blockchain_timestamp = []
        # 节点的每个区块的状态用于共识的阶段:normal, wait_tx, wait_sig, wait_commit
        # normal:节点处于正常状态，和全节点进行交易信息的同步，若收到prep_req将进入wait_tx或wait_sig状态
        # wait_tx:节点处于共识状态，收到了prep_req但是区块中的交易尚未全部收到因此需要等待交易到达，收到全部交易后并执行验证缓存后，进入wait_sig状态
        # wait_sig:节点处于共识状态，此时节点等待其他节点的验证，收到2f+1个sig后进入wait_commit状态
        # wait_commit:节点处于共识状态，此时节点等待其他节点的确认提交，收到2f+1个commit后，将完成缓存中区块的上链操作，随后进去normal状态
        # over:节点已经完成了该区块的共识
        self.status = []
        # 收到prep_req后，将该消息中的区块信息进行缓存到本地，prepare_message应该为PrepareMessage类的对象
        self.prepare_message = []
        # 收到prep_req后，记录下共识开始的时间
        # self.start_time = []
        # 节点接下来要上链的区块的索引
        self.latest_id = 0
        # # 节点收到pre_req后发现其中的块中的交易本地尚未全部收到，需要等待块中交易全部同步到本地后才能进行验证，执行，签名，转发
        # # pre_req参数如果不为-1，表明需要等待指定区块中的交易到来。
        # self.pre_req = -1

        # 已经收到了来自列表中的节点的sig包，列表中的节点数量达到一定后广播commit包，切换状态
        self.sig_list = []
        # 已经收到了来自列表中的节点的commit包，列表中的节点数量达到一定后将该接受区块，切换状态
        self.commit_list = []

    def add_event(self, event):
        self.event_list.append(event)

    def get_next_event(self):
        if len(self.event_list) == 0:
            return False
        self.event_list.sort(key=operator.attrgetter('time'), reverse=False)
        event = self.event_list[0]
        del self.event_list[0]
        return event

    def check_event_list(self):
        self.event_list.sort(key=operator.attrgetter('time'), reverse=False)
        for event in self.event_list:
            if event.name == "gen_block":
                continue
            flag = 0
            # if re.search('recv', event.name) is not None:
            #     if event.time < self.execute_time_down:
            #         flag = 1
            #         event.time = self.execute_time_down
            if re.search('send', event.name) is not None:
                if event.time < self.execute_time:
                    flag = 1
                    event.time = self.execute_time
            if flag == 0:
                return

    def set_id(self, id):
        self.id = id

    def set_bandwidth(self, bandwidth):
        self.bandwidth = bandwidth

    def add_sig_id(self, event):
        height = event.message.height
        node_id = event.message.idx
        if len(self.sig_list) >= height:
            if node_id not in self.sig_list[height - 1]:
                self.sig_list[height - 1].append(node_id)
                self.log_info_recv(event)
        else:
            while len(self.sig_list) < height:
                self.sig_list.append([])
            self.sig_list[height - 1].append(node_id)
            self.log_info_recv(event)

    def add_commit_id(self, event):
        height = event.message.height
        node_id = event.message.idx
        if len(self.commit_list) >= height:
            if node_id not in self.commit_list[height - 1]:
                self.commit_list[height - 1].append(node_id)
                self.log_info_recv(event)
        else:
            while len(self.commit_list) < height:
                self.commit_list.append([])
            self.commit_list[height - 1].append(node_id)
            self.log_info_recv(event)

    def add_tx(self, tx_id):
        self.tx_pool.append(tx_id)

    def is_leader(self):
        return self.id == len(self.blockchain_timestamp) % config.node_num

    def init_status(self, event):
        while len(self.status) < event.message.height:
            self.status.append("normal")

    def add_message(self, event):
        while len(self.prepare_message) < event.message.height:
            self.prepare_message.append(None)
        if len(self.prepare_message) >= event.message.height:
            self.prepare_message[event.message.height - 1] = event.message

    # def add_time(self, event):
    #     while len(self.start_time) < event.message.height:
    #         self.start_time.append(None)
    #     self.start_time.append(event.time)

    def remove_tx(self, height):
        prepare_message = self.prepare_message[height - 1]
        for tx in prepare_message.block.tx_list:
            self.tx_done_pool.append(tx)
            if tx in self.tx_pool:
                self.tx_pool.remove(tx)

    # def log_info_follower(self, receiver, event):
    #     if config.log_flag:
    #         if event.message.packet_type == "prepare":
    #             size = event.message.size
    #         else:
    #             size = event.message.hash_sig_size
    #         self.logger.log_info("[Send] [%d to %d] [%f] [%s], Height:%d, Size=%d, View=%d, Id=%d" %
    #                              (self.id, receiver, event.time, event.message.packet_type, event.message.height,
    #                               size, self.view, self.id))

    def log_info_send(self, event):
        if config.log_flag:
            if event.message.packet_type == "prepare":
                size = event.message.size
            else:
                size = event.message.hash_sig_size
            self.logger.log_info("[Send] [Broadcast] [%f] [%s], Height:%d, Size=%d, View=%d, Id=%d" %
                                 (event.time, event.message.packet_type, event.message.height, size,
                                  self.view, self.id))

    def log_error(self, msg):
        if config.log_flag:
            self.logger.log_error(msg)

    def log_info_recv(self, event):
        if config.log_flag:
            p_type = event.message.packet_type
            # if event.message.idx == self.id:
            #     if p_type == "prepare":
            #         p_type = "sig"
            #     elif p_type == "sig":
            #         p_type = "commit"
            if p_type == "prepare":
                if self.is_leader():
                    return
                size = event.message.size
            else:
                size = event.message.hash_sig_size
            self.logger.log_info("[Recv] [from %d] [%f] [%s], Height=%d, Size=%d, View=%d, Id=%d" %
                                 (event.message.idx, event.time, p_type, event.message.height,
                                  size, self.view, self.id))
