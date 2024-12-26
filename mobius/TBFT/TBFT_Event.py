# -*- coding = utf-8 -*-
import operator
from Net import Network
from Message import *
from Transaction import Transaction
from Event import *
from TBFT.TBFT_Node import *
# from .TBFT_Node import *


class DoTBFTEvent:

    @classmethod
    def handle_event(cls, event):
        if event.time >= config.over_time:
            config.shutdown = True
        if event.name == "start":
            for i in range(config.node_num):
                node = Network.node[i]
                cls.enter_new_height(node, 0, event.time)
        elif event.name == "send_trans" and not config.shutdown:
            cls.send_tx(event)
        elif event.name == "recv_trans" and not config.shutdown:
            cls.recv_tx(event)
        elif event.name == "send_message":
            cls.send_message(event)
        elif event.name == "recv_message":
            cls.recv_message(event)
        elif event.name == "gen_block" and not config.shutdown:
            cls.generate_block(Network.node[event.node_id], event.time)
        elif event.name == "timer":
            cls.check_time(event)

    @classmethod
    def check_time(cls, event):
        height = event.message.height
        round = event.message.round
        node = Network.node[event.node_id]
        if height != node.height or round != node.round:
            return
        if event.message.packet_type == "pre_vote":
            if node.status[height][round] == WAIT_PROPOSAL:
                node.log_info_timeout(event)
                cls.enter_prevote(node, event, timeout_flag=True)
        elif event.message.packet_type == "pre_commit":
            if node.status[height][round] == WAIT_PREVOTE:
                node.log_info_timeout(event)
                cls.enter_precommit(node, event, timeout_flag=True)
        elif event.message.packet_type == "commit":
            if node.status[height][round] == WAIT_PRECOMMIT:
                node.log_info_timeout(event)
                cls.enter_commit(node, event, timeout_flag=True)
                

    @classmethod
    def add_recv_trans_event(cls, sender_id, receiver_id, event):
        if receiver_id == event.sender:
            return
        sender = Network.node[sender_id]
        bandwidth = sender.bandwidth
        tx = Network.transaction_list[event.tx_id]
        if event.time > sender.execute_time:
            time = event.time + tx.size / bandwidth
        else:
            time = sender.execute_time + tx.size / bandwidth
        sender.execute_time = time
        new_event = Event("recv_trans", time, event.dis + 1)
        new_event.set_node_id(receiver_id)
        new_event.set_sender_id(sender_id)
        new_event.set_tx_id(event.tx_id)
        # Network.node[receiver_id].add_event(new_event)
        EventQueue.add_event(new_event)

    @classmethod
    def send_tx(cls, event):
        node = Network.node[event.node_id]
        if len(node.tx_pool) >= node.tx_pool_size:
            return
        node.log_info_tx(event)
        node.add_tx(event.tx_id)
        receiver = Network.gen_recv(node.id, event)
        for item in receiver:
            cls.add_recv_trans_event(node.id, item, event)
        # node.check_event_list()

    @classmethod
    def recv_tx(cls, event):
        node = Network.node[event.node_id]
        if event.tx_id in node.tx_pool or event.tx_id in node.tx_done_pool:
            return
        node.log_info_tx(event)
        node.add_tx(event.tx_id)
        receiver = Network.gen_recv(node.id, event)
        for item in receiver:
            # if event.tx_id in Network.node[item].tx_pool:
            #     continue
            cls.add_recv_trans_event(node.id, item, event)
        # node.check_event_list()

    @classmethod
    def add_recv_message_event(cls, sender_id, receiver_id, event):
        # 此函数是转发message，event的发起者id并不是sender转发者的id
        if receiver_id == event.sender:
            # 不会转发给这个事件的发送节点
            return
        sender = Network.node[sender_id]
        bandwidth = sender.bandwidth
        if event.message.packet_type == "proposal" or event.message.packet_type == "commit_block":
            # 这里还可以进行调整，因为proposal包和commit_block包中数据相比，commit_block中还包含了各个节点的有效签名
            if event.time > sender.execute_time:
                time = event.time + event.message.size / bandwidth
            else:
                time = sender.execute_time + event.message.size / bandwidth
            sender.execute_time = time
        else:
            if event.time > sender.execute_time:
                time = event.time + event.message.hash_sig_size / bandwidth
            else:
                time = sender.execute_time + event.message.hash_sig_size / bandwidth
            sender.execute_time = time
        new_event = MessageEvent("recv_message", time, event.message, event.dis + 1)
        new_event.set_node_id(receiver_id)
        new_event.set_sender_id(sender_id)
        # Network.node[receiver_id].add_event(new_event)
        EventQueue.add_event(new_event)

    @classmethod
    def send_message(cls, event):
        node = Network.node[event.node_id]
        receiver = Network.gen_recv(node.id, event)
        height = event.message.height
        round = event.message.round
        # node.init_status(event)
        # 由leader产生，向全网广播send_prep消息，开始共识
        if event.message.packet_type == "proposal":
            if node.id != (event.message.height + event.message.round) % Network.node_num:
                print("ERROR: ONLY LEADER CAN SEND PREPARE MESSAGE!")
                exit(-1)
            # 当Leader发起新的一轮共识时，Leader直接进入prevote阶段，不会
            # 调用enter_prevote函数，因此一些状态修改需要在这里进行
            node.log_info_send(event)
            node.add_proposal(event)
            node.status[height][round] = WAIT_PREVOTE
            # 进入prevote阶段后，需要设定一个precommit的计时器
            timer_message = Message(node.id, "pre_commit",height, round=round)
            timeoutClock = event.time + config.timeout_pre_vote + config.timeout_pre_vote_delta * round
            timeEvent = MessageEvent("timer", timeoutClock, timer_message)
            timeEvent.set_node_id(node.id)
            node.log_info_set_timer(timeEvent, event.time)
            EventQueue.add_event(timeEvent)
            node.add_vote(event)
            for item in receiver:
                cls.add_recv_message_event(node.id, item, event)
            # node.check_event_list()
        elif event.message.packet_type == "pre_vote":
            # node.time = event.time
            # node.status[height][round] = WAIT_PREVOTE
            node.log_info_send(event)
            # 冗余了，调用send_message函数前，add_vote过了
            # node.add_vote(event)
            for item in receiver:
                cls.add_recv_message_event(node.id, item, event)
            # node.check_event_list()
        elif event.message.packet_type == "pre_commit":
            # node.time = event.time
            # node.status[height][round] = WAIT_PRECOMMIT
            node.log_info_send(event)
            # node.add_commit(event)
            for item in receiver:
                cls.add_recv_message_event(node.id, item, event)
            # node.check_event_list()
        elif event.message.packet_type == "commit_block":
            # node.status[height][round] = COMMIT
            node.log_info_send(event)
            for item in receiver:
                cls.add_recv_message_event(node.id, item, event)

    @classmethod
    def recv_message(cls, event):
        node = Network.node[event.node_id]
        receiver = Network.gen_recv(node.id, event)
        height = event.message.height
        round = event.message.round
        # node.init_status(event)

        if event.message.packet_type == "commit_block":
            if event.message.idx in node.blockinfo_list[height]:
                return
            for item in receiver:
                cls.add_recv_message_event(node.id, item, event)
            node.add_blockinfo(event)
            if node.status[height][round] < COMMIT:
                node.status[height][round] = COMMIT
                node.blockchain_timestamp.append(event.time)
                # node.log_info_recv(event)
                # node.logger.log_info("node:%d receives a commit block and adds the block:%d to its chain, time:%f" % (
                #     node.id, event.message.height, event.time))
            if node.height <= height:
                # 如果收到了height更高的合法区块，那就直接进入下个高度
                # 同时收到这个区块，节点需要验证这个区块中的交易和签名信息，检查是否合法，因此还需要时间来检查
                time = event.time + process_block(event.message.block)
                cls.enter_new_height(node, height+1, time)
            return

        # 收到过时的包，那直接选择丢弃。两种情况判定为过时的包：1、height更小；2、height一样大，但是round更小
        if (event.message.height < node.height) or (event.message.height == node.height and event.message.round < node.round):
            return
        if height > node.height or (height == node.height and round > node.round):
            # 收到了后续共识轮次的消息，按照TBFT算法的流程，需要把这些东西进行缓存，只是把这些数据保存起来，并不会对共识状态造成任何影响
            node.init_status(height, round)
            if event.message.packet_type == "proposal":
                if node.proposal[height][round] != None:
                    # 这里表明已经收到这个proposal过了，那么就直接无视，下面也是同理
                    return
                else:
                    node.add_proposal(event)
            elif event.message.packet_type == "pre_vote":
                if event.message.idx in node.pre_vote_list[height][round]:
                    # 这里和下面都同理，如果这个消息已经接受后了，在本地已经有存储，那么就直接无视
                    return
                node.add_vote(event)
            elif event.message.packet_type == "pre_commit":
                if event.message.idx in node.pre_commit_list[height][round]:
                    return
                node.add_commit(event)
            # 收到这些包后，还需要进行转发，由于这些包仅需缓存起来即可，因此转发结束就return。
            for item in receiver:
                cls.add_recv_message_event(node.id, item, event)
            return
            
        if event.message.packet_type == "proposal":
            if node.status[height][round] > WAIT_PROPOSAL:
                # 如果节点在没有收到proposal就进入了后续阶段，表明计时器超时了，即使
                # 后续收到了proposal包也会将这个proposal视为非法的
                return
            else:
                # 普通节点在收到proposal包后，就会进入prevote阶段
                node.add_proposal(event)
                for item in receiver:
                    cls.add_recv_message_event(node.id, item, event)
                cls.enter_prevote(node, event)
        elif event.message.packet_type == "pre_vote":
            if event.message.idx in node.pre_vote_list[height][round]:
                return
            else:
                for item in receiver:
                    cls.add_recv_message_event(node.id, item, event)
                # node.check_event_list()
            node.add_vote(event)
            if node.status[height][round] != WAIT_PREVOTE:
                return
            enter_next, flag = node.check_num(height, round)
            if enter_next:
                if flag:
                    cls.enter_precommit(node, event)
                else:
                    cls.enter_precommit(node, event, hash_is_none=True)
        elif event.message.packet_type == "pre_commit":
            if event.message.idx in node.pre_commit_list[height][round]:
                return
            else:
                for item in receiver:
                    cls.add_recv_message_event(node.id, item, event)
                # node.check_event_list()
            node.add_commit(event)
            if node.status[height][round] != WAIT_PRECOMMIT:
                return
            enter_next, flag = node.check_num(height, round)
            if enter_next:
                if flag:
                    cls.enter_commit(node, event)
                else:
                    cls.enter_commit(node, event, hash_is_none=True)
            # if len(node.pre_commit_list[height][round]) >= Network.threshold:
                # 经过修改，将代码都放到enter_commit函数中
                # cls.enter_commit(node, event)

    

    @classmethod
    def enter_new_height(cls, node, height, time):
        node.height = height
        node.round = 0
        node.log_info_enter_nextphase("New Height", time)
        cls.enter_new_round(node, 0, time)

    @classmethod
    def enter_new_round(cls, node, round, time):
        node.round = round
        height = node.height
        node.log_info_enter_nextphase("New Round", time)
        # 在进入了新的Round后，需要对节点的各项信息都进行初始化
        node.init_status(height, round)
        # node.timeout_flag = False
        # 这里还需要判断一下，是否有未来的proposal已经被缓存了，如果有缓存，那么就直接进入下一轮的wait_provote阶段
        if node.proposal[height][round] != None:
            new_message = Message(node.id, "tempMessage", height, round=round)
            new_event = MessageEvent("tempEvent", time, new_message)
            cls.enter_prevote(node, new_event)
            return
        node.log_info_enter_nextphase("Proposal", time)
        if node.is_leader():
            # 如果该节点在该height和round下为leader的话，那么就应该发起一个新的proposal
            cls.generate_block(node, time)
        else:
            # 如果节点在该height和round下不为leader，那么就设立一个设定一个计时器等待proposal的到来
            timer_message = Message(node.id, "pre_vote", height, round=round)
            timeoutClock = time + config.timeout_propose + config.timeout_propose_delta * round
            timeEvent = MessageEvent("timer", timeoutClock, timer_message)
            timeEvent.set_node_id(node.id)
            node.log_info_set_timer(timeEvent, time)
            EventQueue.add_event(timeEvent)

    @classmethod
    def enter_prevote(cls, node, event, timeout_flag=False):
        # 根据event的类型可以判断是否发生超时，如果event的类型为timer，表明是计时器超时
        # 调用的enter_prevote，那么后续就需要广播空票
        node.log_info_enter_nextphase("Prevote", event.time)
        height = node.height
        round = node.round
        node.status[height][round] = WAIT_PREVOTE
        # 进入prevote阶段后，需要设定一个precommit的计时器
        timer_message = Message(node.id, "pre_commit",height, round=round)
        timeoutClock = event.time + config.timeout_pre_vote + config.timeout_pre_vote_delta * round
        timeEvent = MessageEvent("timer", timeoutClock, timer_message)
        timeEvent.set_node_id(node.id)
        node.log_info_set_timer(timeEvent, event.time)
        EventQueue.add_event(timeEvent)
        # 判断是否发生超时
        if timeout_flag:
            # 由计时器触发来到这里的，是超时的情况,Message的idx设定为负数，意味该节点投了空票
            new_message = Message(-node.id, "pre_vote", height, round=round)
            # 由于计时器压根没有收到proposal，发生超时来到这里的，因此没有process_block的过程
            new_event = MessageEvent("send_message", event.time, new_message)
        else:
            proposal = node.proposal[height][round]
            node.add_vote(event)
            time = event.time + process_block(proposal.block)
            new_message = Message(node.id, "pre_vote", height, round=round)
            new_event = MessageEvent("send_message", time, new_message)
        new_event.set_node_id(node.id)
        new_event.set_sender_id(node.id)
        # 这个函数内有两个add_vote，第一处是在else分支语句中，当收到来自leader的proposal，此时
        # leader已经进入WAIT_PREVOTE阶段，这个proposal包其实也代表了leader的pre_vote
        # 第二处就是下面这句，由于该节点enter_prevote后续就要广播自己的投票情况，同时需要将自己的投票
        # 结果保存到本地
        node.add_vote(new_event)
        cls.send_message(new_event)
        enter_next, flag = node.check_num(height, round)
        # 一般来说不会enter_prevote后立马就enter_pre_commit，只有以下这种特殊情况
        # 在收到proposal或者超时之前就已经收到了很多的prevote票，再enter_prevote后
        # 经过上述的两处或者一处的add_vote后，票数就达标了。
        if enter_next:
            if flag:
                cls.enter_precommit(node, event)
            else:
                # flag为false意为收齐的2f+1张票为空票，后续该节点也将会投空票
                cls.enter_precommit(node, event, hash_is_none=True)
    
    @classmethod
    def enter_precommit(cls, node, event, timeout_flag=False, hash_is_none=False):
        # 流程基本和ener_prevote一致，首先设立定时器，后续判断是否超时
        node.log_info_enter_nextphase("Precommit", event.time)
        height = event.message.height
        round = event.message.round
        node.status[height][round] = WAIT_PRECOMMIT

        timer_message = Message(node.id, "commit", height, round=round)
        timeoutClock = event.time + config.timeout_pre_commit + config.timeout_pre_commit_delta * round
        timeEvent = MessageEvent("timer", timeoutClock, timer_message)
        timeEvent.set_node_id(node.id)
        node.log_info_set_timer(timeEvent, event.time)
        EventQueue.add_event(timeEvent)

        if timeout_flag or node.proposal[height][round] == None or hash_is_none:
            # 有三种情况需要广播空包：1、发生超时进入precommit状态；2、在该height和round下
            # 节点没有收到proposal；3、prevote阶段共识的结果中的hash值为空，意为广播空包
            new_message = Message(-node.id, "pre_commit", height, round=round)
        else:
            new_message = Message(node.id, "pre_commit", height, round=round)
        new_event = MessageEvent("send_message", event.time, new_message)
        new_event.set_node_id(node.id)
        new_event.set_sender_id(node.id)
        node.add_commit(new_event)
        cls.send_message(new_event)
        enter_next, flag = node.check_num(height, round)
        if enter_next:
            if flag:
                cls.enter_commit(node, event)
            else:
                cls.enter_commit(node, event, hash_is_none=True)
        
    @classmethod
    def enter_commit(cls, node, event, timeout_flag=False, hash_is_none=False):
        # 这里有两个途径到达这里，1、在precommit阶段收齐了2/3的票数，后续就会调用entercommit；
        # 2、在precommit阶段的计时器发生超时，后续就会调用entercommit
        node.log_info_enter_nextphase("Commit", event.time)
        height = event.message.height
        round = event.message.round
        if timeout_flag or node.proposal[height][round] == None or hash_is_none:
            # 这种情况表明提交的是空块，需要新的一轮来完成一次合法的共识
            node.status[height][round] = OVER
            node.log_info_commit(event)
            cls.enter_new_round(node, round+1, event.time)
        else:
            # 反之就是通过正常共识流程来到这里的，将提交的区块保存到本地即可
            node.status[height][round] = COMMIT
            node.log_info_commit(event)
            node.add_block_timestamp(event)
            
            node.remove_tx(height, round)
            block_info = BlockInfoMessage(node.id, "commit_block", height, round=round,
                                                block=Network.block[height])
            new_event = MessageEvent("send_message", event.time, block_info)
            new_event.set_node_id(node.id)
            new_event.set_sender_id(node.id)
            cls.send_message(new_event)
            node.add_blockinfo(new_event)
            # 完成了这个高度的共识，后续进入新高度的共识
            cls.enter_new_height(node, height+1, event.time)

    @classmethod
    def generate_block(cls, node, time):
        if config.shutdown:
            return
        height = node.height
        round = node.round
        result = Transaction.select_tx(Block.block_size_limit, node.tx_pool)
        tx_list = result[0]
        size = result[1]
        if size == 0:
            print("There is no transaction in the pool.")
            tx_over_time = float(config.tx_num) / config.tx_rate
            if time >= tx_over_time:
                config.shutdown = True
                return
            if config.log_flag:
                node.logger.log_info("Node:%d becomes leader and will launch a empty block." % node.id)
        else:
            if config.log_flag:
                node.logger.log_info("Node:%d becomes leader and will launch a new block." % node.id)
        block = Block(len(node.blockchain_timestamp), tx_list, size + Message.hash_sig_size)
        Network.add_block(block, height)
        new_message = ProposalMessage(node.id, "proposal", height=height, round=round, block=block)
        new_event = MessageEvent("send_message", time, new_message)
        new_event.set_node_id(node.id)
        new_event.set_sender_id(node.id)
        cls.send_message(new_event)
