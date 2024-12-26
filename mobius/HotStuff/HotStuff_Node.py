# -*- coding = utf-8 -*-
import Log
import config
import operator
import re

def process_block(block):
    # 该函数为节点收到区块后进行验证的时间，如果仅仅只是转发prepare消息无需计算该部分时间
    return config.process_cost


class HotStuffNode:
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
        # 由于HotStuff共识算法下不会出现分叉现象，因此不同链上的区块的差别仅是共识达成的时间而已，其余交易信息都统一存放在Network中即可
        self.blockchain_timestamp = []
        # 节点的每个区块的状态：normal, prepare, pre_commit, commit, decide
        self.status = []
        # self.start_time = []
        # self.end_time = []
        self.prepare_message = []
        self.prepare_vote = []
        self.pre_commit_vote = []
        self.commit_vote = []
        self.view_vote = []

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

    def add_tx(self, tx_id):
        self.tx_pool.append(tx_id)

    def init_status(self, event):
        while len(self.status) < event.message.height:
            self.status.append("normal")

    def add_message(self, event):
        while len(self.prepare_message) < event.message.height:
            self.prepare_message.append(None)
        if len(self.prepare_message) >= event.message.height:
            self.prepare_message[event.message.height - 1] = event.message
        # else:
        #     self.prepare_message.append(event.message)

    # def add_time(self, event):
    #     while len(self.start_time) < event.message.height:
    #         self.start_time.append(None)
    #     self.start_time.append(event.time)
    #
    # def add_end_time(self, event):
    #     while len(self.end_time) < event.message.height:
    #         self.end_time.append(None)
    #     self.end_time.append(event.time)

    def add_prepare_vote(self, event, node_id=-1):
        # from HotStuff.HotStuff_Event import view_to_leader
        height = event.message.height
        if node_id == -1:
            node_id = event.message.idx
        if len(self.prepare_vote) >= height:
            if node_id not in self.prepare_vote[height - 1]:
                self.prepare_vote[height - 1].append(node_id)
                # if self.id != view_to_leader(self.view):
                self.log_info_recv(event)
        else:
            while len(self.prepare_vote) < height:
                self.prepare_vote.append([])
            self.prepare_vote[height - 1].append(node_id)
            # if self.id != view_to_leader(self.view):
            self.log_info_recv(event)

    def add_pre_commit_vote(self, event, node_id=-1):
        # from HotStuff.HotStuff_Event import view_to_leader
        height = event.message.height
        if node_id == -1:
            node_id = event.message.idx
        if len(self.pre_commit_vote) >= height:
            if node_id not in self.pre_commit_vote[height - 1]:
                self.pre_commit_vote[height - 1].append(node_id)
                # if self.id != view_to_leader(self.view):
                self.log_info_recv(event)
        else:
            while len(self.pre_commit_vote) < height:
                self.pre_commit_vote.append([])
            self.pre_commit_vote[height - 1].append(node_id)
            # if self.id != view_to_leader(self.view):
            self.log_info_recv(event)

    def add_commit_vote(self, event, node_id=-1):
        # from HotStuff.HotStuff_Event import view_to_leader
        height = event.message.height
        if node_id == -1:
            node_id = event.message.idx
        if len(self.commit_vote) >= height:
            if node_id not in self.commit_vote[height - 1]:
                self.commit_vote[height - 1].append(node_id)
                # if self.id != view_to_leader(self.view):
                self.log_info_recv(event)
        else:
            while len(self.commit_vote) < height:
                self.commit_vote.append([])
            self.commit_vote[height - 1].append(node_id)
            # if self.id != view_to_leader(self.view):
            self.log_info_recv(event)

    def add_view_vote(self, event, node_id=-1):
        height = event.message.height
        if node_id == -1:
            node_id = event.message.idx
        if len(self.view_vote) >= height:
            if node_id not in self.view_vote[height - 1]:
                self.view_vote[height - 1].append(node_id)
                # if self.id != HotStuff.HotStuff_Event.view_to_leader(self.view):
                self.log_info_recv(event)
        else:
            while len(self.view_vote) < height:
                self.view_vote.append([])
            self.view_vote[height - 1].append(node_id)
            # if self.id != HotStuff.HotStuff_Event.view_to_leader(self.view):
            self.log_info_recv(event)

    def remove_tx(self, height):
        prepare_message = self.prepare_message[height - 1]
        for tx in prepare_message.block.tx_list:
            self.tx_done_pool.append(tx)
            if tx in self.tx_pool:
                self.tx_pool.remove(tx)

    def log_info_follower(self, receiver, event):
        if config.log_flag:
            if event.message.packet_type == "prepare":
                size = event.message.size
            else:
                size = event.message.hash_sig_size
            self.logger.log_info("[Send] [%d to %d] [%f] [%s], Height:%d, Size=%d, View=%d, Id=%d" %
                                 (self.id, receiver, event.time, event.message.packet_type, event.message.height,
                                  size, self.view, self.id))

    def log_info_leader(self, event):
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
            if event.message.idx == self.id:
                if p_type == "prepare":
                    p_type = "prepare_vote"
                elif p_type == "pre_commit":
                    p_type = "pre_commit_vote"
                elif p_type == "commit":
                    p_type = "commit_vote"
            if p_type == "prepare":
                size = event.message.size
            else:
                size = event.message.hash_sig_size
            self.logger.log_info("[Recv] [from %d] [%f] [%s], Height=%d, Size=%d, View=%d, Id=%d" %
                                 (event.message.idx, event.time, p_type, event.message.height,
                                  size, self.view, self.id))

    # def log_info_leader_recv(self):
    #     if config.log_flag:
    #         self.logger.log_info("[Recv] [from %d] [%f] Type=%s, Height=%d, View=%d, Id=%d" %
    #                              ())
