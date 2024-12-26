# from Net import Network
import Log
import config
import Net


def process_block(block):
    return config.process_cost + len(block.tx_list) * config.tx_cost


# 此处定义了节点的不同状态
WAIT_PROPOSAL = 1
WAIT_PREVOTE = 2
WAIT_PRECOMMIT = 3
OVER = 4
COMMIT = 5


class TBFTNode:

    def __init__(self, id, threshold, tx_pool_size):
        self.id = id
        self.bandwidth = 0
        self.logger = Log.Log(id)
        self.logger.file_handle()
        self.tx_pool = []
        self.tx_done_pool = []
        self.tx_pool_size = tx_pool_size
        self.round = 0
        self.height = 0
        self.status = []
        # time意为计时器的时间，但是在主动添加超时计时器事件后，这个变量也没有什么用处了
        self.time = 0
        # timeout_flag用于记录当前height和round下，是否已经发生了计时器超时的事件
        # self.timeout_flag = False
        self.execute_time = 0
        # self.event_list = []
        self.blockchain_timestamp = []
        # 由于考虑到proposal在相同height下，有不同round。因此需要二维数组，后续add_vote,add_commit函数等都进行了对应的修改
        self.proposal = []
        self.pre_vote_list = []
        self.pre_commit_list = []
        # 这个列表用于记录各个高度的block_info的广播情况，避免重复转发
        self.blockinfo_list = []
        self.threshold = threshold
        # self.forward = []



    def set_bandwidth(self, bandwidth):
        self.bandwidth = bandwidth

    def add_tx(self, tx_id):
        self.tx_pool.append(tx_id)


    def add_proposal(self, event):
        height = event.message.height
        round = event.message.round
        while len(self.proposal) <= height:
            self.proposal.append([])
        while len(self.proposal[height]) <= round:
            self.proposal[height].append(None)
        if len(self.proposal[height]) > round:
            self.proposal[height][round] = event.message
            self.log_info_recv(event)

    def add_vote(self, event):
        height = event.message.height
        round = event.message.round
        node_id = event.message.idx
        while len(self.pre_vote_list) <= height:
            self.pre_vote_list.append([])
        while len(self.pre_vote_list[height]) <= round:
            self.pre_vote_list[height].append([])
        if node_id not in self.pre_vote_list[height][round]:
            self.pre_vote_list[height][round].append(node_id)
            self.log_info_recv(event)
        # if len(self.pre_vote_list) > height:
        #     if len(self.pre_vote_list[height]) > round:
        #         if node_id not in self.pre_vote_list[height][round]:
        #             self.pre_vote_list[height - 1].append(node_id)
        #             self.log_info_recv(event)
        #     else:  
        #         while len(self.pre_vote_list[height]) <= round:
        #             self.pre_vote_list[height].append([])
        #         if node_id not in self.pre_vote_list[height][round]:
        #             self.pre_vote_list[height][round].append(node_id)
        #             self.log_info_recv(event)
        # else:
        #     while len(self.pre_vote_list) <= height:
        #         self.pre_vote_list.append([])
        #     while len(self.pre_vote_list[height]) <= round:
        #         self.pre_vote_list[height].append([])
        #     if node_id not in self.pre_vote_list[height][round]:
        #         self.pre_vote_list[height][round].append(node_id)
        #         self.log_info_recv(event)

    def add_commit(self, event):
        height = event.message.height
        round = event.message.round
        node_id = event.message.idx
        while len(self.pre_commit_list) <= height:
            self.pre_commit_list.append([])
        while len(self.pre_commit_list[height]) <= round:
            self.pre_commit_list[height].append([])
        if node_id not in self.pre_commit_list[height][round]:
            self.pre_commit_list[height][round].append(node_id)
            self.log_info_recv(event)
        # if len(self.pre_commit_list) > height:
        #     if len(self.pre_commit_list[height]) > round:
        #         if node_id not in self.pre_commit_list[height][round]:
        #             self.pre_commit_list[height - 1].append(node_id)
        #             self.log_info_recv(event)
        #     else:  
        #         while len(self.pre_commit_list[height]) <= round:
        #             self.pre_commit_list[height].append([])
        #         if node_id not in self.pre_commit_list[height][round]:
        #             self.pre_commit_list[height][round].append(node_id)
        #             self.log_info_recv(event)
        # else:
        #     while len(self.pre_commit_list) <= height:
        #         self.pre_commit_list.append([])
        #     while len(self.pre_commit_list[height]) <= round:
        #         self.pre_commit_list[height].append([])
        #     if node_id not in self.pre_commit_list[height][round]:
        #         self.pre_commit_list[height][round].append(node_id)
        #         self.log_info_recv(event)

    def add_block_timestamp(self, event):
        height = event.message.height
        while len(self.blockchain_timestamp) <= height:
            self.blockchain_timestamp.append(None)
        self.blockchain_timestamp[height] = event.time
        self.logger.log_info("Node:%d adds the block:%d to its chain, time:%f" % (
            self.id, height, event.time))
        if config.single_log:
            Net.Network.logger.log_info("Node:%d adds the block:%d to its chain, time:%f" % (
            self.id, height, event.time))

    def add_blockinfo(self, event):
        height = event.message.height
        # round = event.message.round
        node_id = event.message.idx
        while len(self.blockinfo_list) <= height:
            self.blockinfo_list.append([])
        if node_id not in self.blockinfo_list[height]:
            self.blockinfo_list[height].append(node_id)
            self.log_info_recv(event)

    def is_leader(self):
        if (self.height + self.round) % config.node_num == self.id:
            return True
        else:
            return False

    def init_status(self, height, round):
        # 初始化在此height和round下的状态
        while len(self.status) <= height:
            self.status.append([])
            # self.forward.append(False)
        while len(self.status[height]) <= round:
            self.status[height].append(None)
        if self.status[height][round] == None:
            self.status[height][round] = WAIT_PROPOSAL
        # 初始化proposal列表
        while len(self.proposal) <= height:
            self.proposal.append([])
        while len(self.proposal[height]) <= round:
            self.proposal[height].append(None)
        # 初始化pre_vote_list列表
        while len(self.pre_vote_list) <= height:
            self.pre_vote_list.append([])
        while len(self.pre_vote_list[height]) <= round:
            self.pre_vote_list[height].append([])
        # 初始化pre_commit_list列表
        while len(self.pre_commit_list) <= height:
            self.pre_commit_list.append([])
        while len(self.pre_commit_list[height]) <= round:
            self.pre_commit_list[height].append([])
        # 初始化blockinfo_list列表
        while len(self.blockinfo_list) <= height:
            self.blockinfo_list.append([])


    

    def remove_tx(self, height, round):
        proposal_message = self.proposal[height][round]
        for tx in proposal_message.block.tx_list:
            self.tx_done_pool.append(tx)
            if tx in self.tx_pool:
                self.tx_pool.remove(tx)
    
    def check_num(self, height, round):
        # 这个函数用于统计在height和round下，是否收集了足够的票数
        # 本模拟器中，当一个节点超时后将会广播一个空票。例如：id=1的节点超时，进入prevote阶段
        # 后续进行广播prevote的空票，在prevote消息包中，message中的idx将为负数，表明该节点
        # 投的是空票，因此positive和negative分别对应了空票和非空票的数量，只要有任意一种票数达
        # 到Network.threshold后即可进入下个阶段
        count_positive = 0
        count_negative = 0
        if self.status[height][round] == WAIT_PREVOTE:
            for item in self.pre_vote_list[height][round]:
                if item > 0:
                    count_positive += 1
                else:
                    count_negative += 1
        if self.status[height][round] == WAIT_PRECOMMIT:
            for item in self.pre_commit_list[height][round]:
                if item > 0:
                    count_positive += 1
                else:
                    count_negative += 1
        if count_positive >= self.threshold:
                return True, True
        elif count_negative > config.node_num - self.threshold:
            return True, False
        else:
            return False, False
            
    def log_info_tx(self, event):
        if config.log_flag and config.single_log:
            message = "Transaction " + str(event.tx_id)
            if event.name == "send_trans":
                Net.Network.logger.log_info("[Node %d] [Send] [Broadcast] [%f] [%s], Size=%d" %
                                 (self.id, event.time, message, Net.Network.transaction_list[event.tx_id].size))
            else:
                Net.Network.logger.log_info("[Node %d] [Recv] [from %d] [%f] [%s], Size=%d" %
                                 (self.id, event.sender, event.time, message, Net.Network.transaction_list[event.tx_id].size))

    def log_info_send(self, event):
        if config.log_flag:
            if event.message.packet_type == "proposal":
                size = event.message.size
            elif event.message.packet_type == "commit_block":
                size = event.message.size
            else:
                size = event.message.hash_sig_size
            self.logger.log_info("[Send] [Broadcast] [%f] [%s], Height:%d, Size=%d, Round=%d, Id=%d" %
                                 (event.time, event.message.packet_type, event.message.height, size,
                                  self.round, self.id))
            if config.single_log:
                Net.Network.logger.log_info("[Node %d] [Send] [Broadcast] [%f] [%s], Height:%d, Size=%d, Round=%d, Id=%d" %
                                 (self.id, event.time, event.message.packet_type, event.message.height, size,
                                  self.round, self.id))

    def log_error(self, msg):
        if config.log_flag:
            self.logger.log_error(msg)

    def log_info_recv(self, event):
        if config.log_flag:
            p_type = event.message.packet_type
            if p_type == "proposal":
                if self.is_leader():
                    return
                size = event.message.size
            elif p_type == "commit_block":
                size = event.message.size
            else:
                size = event.message.hash_sig_size
            self.logger.log_info("[Recv] [from %d] [%f] [%s], Height=%d, Size=%d, Round=%d, Id=%d" %
                                 (event.message.idx, event.time, p_type, event.message.height,
                                  size, self.round, self.id))
            if config.single_log:
                Net.Network.logger.log_info("[Node %d] [Recv] [from %d] [%f] [%s], Height=%d, Size=%d, Round=%d, Id=%d" %
                                 (self.id, event.message.idx, event.time, p_type, event.message.height,
                                  size, self.round, self.id))

    def log_info_enter_nextphase(self, phase, time):
        if config.log_flag:
            self.logger.log_info("[Enter] [from %d] [%f] [%s], Height=%d, Round=%d" %
                                 (self.id, time, phase, self.height, self.round))
            if config.single_log:
                Net.Network.logger.log_info("[Node %d] [Enter] [from %d] [%f] [%s], Height=%d, Round=%d" %
                                 (self.id, self.id, time, phase, self.height, self.round))

    def log_info_set_timer(self, event, start_time):
        if config.log_flag:
            p_type = event.message.packet_type
            phase = ""
            if p_type == "commit":
                phase = "Commit"
            elif p_type == "pre_commit":
                phase = "Precommit"
            elif p_type == "pre_vote":
                phase = "Prevote"
            self.logger.log_info("[Timer] [from %d] [%f] [%s], Height=%d, Round=%d Start=%f, End=%f" %
                                 (self.id, start_time, phase, event.message.height, event.message.round, start_time, event.time))
            if config.single_log:
                Net.Network.logger.log_info("[Node %d] [Timer] [from %d] [%f] [%s], Height=%d, Round=%d Start=%f, End=%f" %
                                 (self.id, self.id, start_time, phase, event.message.height, event.message.round, start_time, event.time))

    def log_info_timeout(self, event):
        if config.log_flag:
            self.logger.log_info("[Timeout] [from %d] [%f] [%s], Height=%d, Round=%d" %
                                 (event.message.idx, event.time, event.message.packet_type, event.message.height,
                                  event.message.round))
            if config.single_log:
                Net.Network.logger.log_info("[Node %d] [Timeout] [from %d] [%f] [%s], Height=%d, Round=%d" %
                                 (self.id, event.message.idx, event.time, event.message.packet_type, event.message.height,
                                  event.message.round))
            
    def log_info_commit(self, event):
        if config.log_flag:
            height = event.message.height
            round = event.message.round
            commit_res = ""
            if self.status[height][round] == OVER:
                commit_res = "EMPTY"
            else:
                commit_res = "VALID"
            self.logger.log_info("[Commit] [from %d] [%f] [%s], Height=%d, Round=%d" %
                                 (self.id, event.time, commit_res, height, round))
            if config.single_log:
                Net.Network.logger.log_info("[Node %d] [Commit] [from %d] [%f] [%s], Height=%d, Round=%d" %
                                 (self.id, self.id, event.time, commit_res, height, round))
