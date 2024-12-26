# -*- coding = utf-8 -*-
import math
import operator
from Net import Network
from Message import *
from Transaction import Transaction
from HotStuff.HotStuff_Node import process_block
from Event import *
import config


def sync_tx_from_neighbor(node_id, lack_tx_list):
    size = len(lack_tx_list) * 512
    time = size / Network.node[node_id].bandwidth
    return time


def view_to_leader(view):
    # 该函数用于决定下一轮view的leader应该是谁
    # 目前是和PBFT一样的采用轮换的形式
    leader_id = view % Network.node_num
    return leader_id


class DoHotStuffEvent:

    @classmethod
    def handle_event(cls, event):
        if event.time >= config.over_time:
            config.shutdown = True
        if event.name == "send_trans" and not config.shutdown:
            cls.send_tx(event)
        if event.name == "recv_trans" and not config.shutdown:
            cls.recv_tx(event)
        elif event.name == "send_message":
            cls.send_message(event)
        elif event.name == "recv_message":
            cls.recv_message(event)
        elif event.name == "gen_block" and not config.shutdown:
            cls.generate_block(Network.node[event.node_id], event.time)

    @classmethod
    def add_recv_trans_event(cls, sender_id, receiver_id, event):
        sender = Network.node[sender_id]
        receiver = Network.node[receiver_id]
        bandwidth = sender.bandwidth
        tx = Network.transaction_list[event.tx_id]
        if event.time >= sender.execute_time:
            time = event.time + tx.size / bandwidth
        else:
            time = sender.execute_time + tx.size / bandwidth
        sender.execute_time = time
        new_event = Event("recv_trans", time, event.dis + 1)
        new_event.set_node_id(receiver_id)
        new_event.set_tx_id(event.tx_id)
        Network.node[receiver_id].add_event(new_event)
        EventQueue.add_event(new_event)
        # print("trans---> sender:%d, receiver:%d, time:%f" % (sender_id, receiver_id, new_event.time))
        # logging.info('<Event Name> | {}'.format())

    @classmethod
    # 发送交易仅被交易发起者调用，开始让该笔交易在网络上传播
    def send_tx(cls, event):
        node = Network.node[event.node_id]
        if len(node.tx_pool) >= node.tx_pool_size:
            return
        node.add_tx(event.tx_id)
        receiver = Network.gen_recv(node.id, event)
        for item in receiver:
            cls.add_recv_trans_event(node.id, item, event)
        node.check_event_list()

    @classmethod
    # 接受交易仅被非交易发起者调用，交易的接受和传播，每笔交易每个节点仅执行一次，再次收到该交易不会做什么
    def recv_tx(cls, event):
        node = Network.node[event.node_id]
        if event.tx_id in node.tx_pool or event.tx_id in node.tx_done_pool:
            # print("node:%d drop the duplicate/done event" % (node.id))
            return
        node.add_tx(event.tx_id)
        receiver = Network.gen_recv(node.id, event)

        for item in receiver:
            # 如果随机挑选的节点中已经接受过该笔交易的，直接跳过
            if event.tx_id in Network.node[item].tx_pool:
                continue
            cls.add_recv_trans_event(node.id, item, event)
        node.check_event_list()

    @classmethod
    def add_recv_message_event(cls, sender_id, receiver_id, event):
        # 此函数是转发message，event的发起者id并不是sender转发者的id
        sender = Network.node[sender_id]
        receiver = Network.node[receiver_id]
        bandwidth = sender.bandwidth
        if event.message.packet_type == "prepare":
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
        Network.node[receiver_id].add_event(new_event)
        EventQueue.add_event(new_event)

    @classmethod
    def send_message(cls, event):
        node = Network.node[event.node_id]
        receiver = Network.gen_recv(node.id, event)
        node.init_status(event)
        leader_id = view_to_leader(event.message.view)
        # 由leader产生，向全网广播send_prep消息，开始共识
        if event.message.packet_type == "prepare":
            if node.id != leader_id:
                node.log_error("ONLY LEADER CAN SEND PREPARE MESSAGE!")
                # print("ERROR: ONLY LEADER CAN SEND PREPARE MESSAGE!")
                exit(-1)
            node.add_message(event)
            # node.add_time(event)
            node.status[event.message.height - 1] = "wait_prepare"
            node.log_info_leader(event)
            # print("leader%d send the prepare message:%d, time:%f" % (node.id, event.message.height, event.time))
            node.add_prepare_vote(event, node.id)
            for item in receiver:
                cls.add_recv_message_event(node.id, item, event)
            node.check_event_list()
        elif event.message.packet_type == "prepare_vote":
            node.status[event.message.height - 1] = "prepare"
            node.log_info_follower(leader_id, event)
            # print("node%d send prepare_vote:%d to leader%d, time:%f" % (
            #     node.id, event.message.height, leader_id, event.time))
            cls.add_recv_message_event(node.id, leader_id, event)
            node.check_event_list()
        elif event.message.packet_type == "pre_commit":
            if node.id != leader_id:
                node.log_error("ONLY LEADER CAN SEND PREPARE MESSAGE!")
                # print("ERROR: ONLY LEADER CAN SEND PRE_COMMIT MESSAGE!")
                exit(-1)
            node.status[event.message.height - 1] = "wait_pre_commit"
            node.log_info_leader(event)
            # print("leader%d send the pre_commit message:%d, time:%f" % (node.id, event.message.height, event.time))
            node.add_pre_commit_vote(event, node.id)
            for item in receiver:
                cls.add_recv_message_event(node.id, item, event)
            node.check_event_list()
        elif event.message.packet_type == "pre_commit_vote":
            node.status[event.message.height - 1] = "pre_commit"
            node.log_info_follower(leader_id, event)
            # print("node%d send pre_commit_vote:%d to leader%d, time:%f" % (
            #     node.id, event.message.height, leader_id, event.time))
            cls.add_recv_message_event(node.id, leader_id, event)
            node.check_event_list()
        elif event.message.packet_type == "commit":
            if node.id != leader_id:
                node.log_error("ONLY LEADER CAN SEND PREPARE MESSAGE!")
                # print("ERROR: ONLY LEADER CAN SEND COMMIT MESSAGE!")
                exit(-1)
            node.status[event.message.height - 1] = "wait_commit"
            node.log_info_leader(event)
            # print("leader%d send the commit message:%d, time:%f" % (node.id, event.message.height, event.time))
            node.add_commit_vote(event, node.id)
            for item in receiver:
                cls.add_recv_message_event(node.id, item, event)
            node.check_event_list()
        elif event.message.packet_type == "commit_vote":
            node.status[event.message.height - 1] = "commit"
            node.log_info_follower(leader_id, event)
            # print("node%d send commit_vote:%d to leader%d, time:%f" % (
            #     node.id, event.message.height, leader_id, event.time))
            cls.add_recv_message_event(node.id, leader_id, event)
            node.check_event_list()
        elif event.message.packet_type == "decide":
            if node.id != leader_id:
                node.log_error("ERROR: ONLY LEADER CAN SEND DECIDE MESSAGE!")
                # print("ERROR: ONLY LEADER CAN SEND DECIDE MESSAGE!")
                exit(-1)
            node.status[event.message.height - 1] = "over"
            node.log_info_leader(event)
            # print("leader%d send the decide message:%d, time:%f" % (node.id, event.message.height, event.time))
            for item in receiver:
                cls.add_recv_message_event(node.id, item, event)
            node.check_event_list()
        elif event.message.packet_type == "new_view":
            node.status[event.message.height - 2] = "over"
            node.log_info_follower(leader_id, event)
            # print("node%d send new_view:%d to next leader%d, time:%f" % (
            #     node.id, event.message.height, leader_id, event.time))
            cls.add_recv_message_event(node.id, leader_id, event)
            node.check_event_list()

    @classmethod
    def recv_message(cls, event):
        node = Network.node[event.node_id]
        receiver = Network.gen_recv(node.id, event)
        node.init_status(event)
        # if len(node.prepare_message) >= event.message.height:
        #     prepare = node.prepare_message[event.message.height - 1]
        if event.message.packet_type == "prepare":
            if node.status[event.message.height - 1] == "normal":
                node.add_message(event)
                # node.add_time(event)
                prepare = node.prepare_message[event.message.height - 1]
                node.log_info_recv(event)
                # print("node%d receive the prepare message:%d from leader:%d, time:%f" % (
                #     node.id, event.message.height, event.message.idx, event.time))
                time = event.time + process_block(prepare.block)
                node.status[event.message.height - 1] = "prepare"
                lack_tx_list = list(set(event.message.block.tx_list) - set(node.tx_pool))
                if len(lack_tx_list) > 0:
                    # 如果节点发现区块中的部分交易，本节点尚未同步到区块中的全部交易，则需要向邻居节点索要缺失的交易信息
                    # 仅需计算出这个向邻居节点索要交易信息的时间，然后把这段时间加上，就当做交易信息同步成功
                    time += sync_tx_from_neighbor(node.id, lack_tx_list)
                new_message = Message(node.id, "prepare_vote", prepare.height, prepare.view)
                new_event = MessageEvent("send_message", time, new_message)
                new_event.set_node_id(node.id)
                cls.send_message(new_event)
                # 同时还需要继续转发这个包给其他节点
                for item in receiver:
                    cls.add_recv_message_event(node.id, item, event)
                node.check_event_list()
            else:
                return
        elif event.message.packet_type == "prepare_vote":
            node.add_prepare_vote(event)
            prepare = node.prepare_message[event.message.height - 1]
            if node.status[event.message.height - 1] != "wait_prepare":
                return
            if len(node.prepare_vote[event.message.height - 1]) >= Network.threshold:
                new_message = Message(node.id, "pre_commit", prepare.height, prepare.view)
                new_event = MessageEvent("send_message", event.time, new_message)
                new_event.set_node_id(node.id)
                cls.send_message(new_event)
        elif event.message.packet_type == "pre_commit":
            if node.status[event.message.height - 1] not in ["normal", "prepare"]:
                return
            prepare = node.prepare_message[event.message.height - 1]
            node.log_info_recv(event)
            # print("node%d receive the pre_commit message:%d from leader:%d, time:%f" % (
            #     node.id, event.message.height, event.message.idx, event.time))
            new_message = Message(node.id, "pre_commit_vote", prepare.height, prepare.view)
            new_event = MessageEvent("send_message", event.time, new_message)
            new_event.set_node_id(node.id)
            cls.send_message(new_event)
            for item in receiver:
                cls.add_recv_message_event(node.id, item, event)
            node.check_event_list()
        elif event.message.packet_type == "pre_commit_vote":
            node.add_pre_commit_vote(event)
            if node.status[event.message.height - 1] != "wait_pre_commit":
                return
            prepare = node.prepare_message[event.message.height - 1]
            if len(node.pre_commit_vote[event.message.height - 1]) >= Network.threshold:
                new_message = Message(node.id, "commit", prepare.height, prepare.view)
                new_event = MessageEvent("send_message", event.time, new_message)
                new_event.set_node_id(node.id)
                cls.send_message(new_event)
        elif event.message.packet_type == "commit":
            if node.status[event.message.height - 1] not in ["normal", "prepare", "pre_commit"]:
                return
            prepare = node.prepare_message[event.message.height - 1]
            node.log_info_recv(event)
            # print("node%d receive the commit message:%d from leader:%d, time:%f" % (
            #     node.id, event.message.height, event.message.idx, event.time))
            new_message = Message(node.id, "commit_vote", prepare.height, prepare.view)
            new_event = MessageEvent("send_message", event.time, new_message)
            new_event.set_node_id(node.id)
            cls.send_message(new_event)
            for item in receiver:
                cls.add_recv_message_event(node.id, item, event)
            node.check_event_list()
        elif event.message.packet_type == "commit_vote":
            node.add_commit_vote(event)
            if node.status[event.message.height - 1] != "wait_commit":
                return
            prepare = node.prepare_message[event.message.height - 1]
            if len(node.commit_vote[event.message.height - 1]) >= Network.threshold:
                node.blockchain_timestamp.append(event.time)
                node.remove_tx(event.message.height)
                node.view += 1
                new_message = Message(node.id, "decide", prepare.height, prepare.view)
                new_event = MessageEvent("send_message", event.time, new_message)
                new_event.set_node_id(node.id)
                cls.send_message(new_event)
        elif event.message.packet_type == "decide":
            if node.status[event.message.height - 1] == "over":
                return
            else:
                prepare = node.prepare_message[event.message.height - 1]
                node.log_info_recv(event)
                # print("node%d receive the decide message:%d from leader:%d, time:%f" % (
                #     node.id, event.message.height, event.message.idx, event.time))
                node.blockchain_timestamp.append(event.time)
                node.remove_tx(event.message.height)
                if prepare.view + 1 > node.view:
                    node.view = prepare.view + 1
                new_message = Message(node.id, "new_view", prepare.height + 1, prepare.view + 1)
                new_event = MessageEvent("send_message", event.time, new_message)
                new_event.set_node_id(node.id)
                cls.send_message(new_event)
                for item in receiver:
                    cls.add_recv_message_event(node.id, item, event)
                node.check_event_list()
        elif event.message.packet_type == "new_view":
            node.add_view_vote(event)
            if node.status[event.message.height - 1] != "normal":
                return
            # if node.status[event.message.height - 1] != "normal":
            #     return
            if len(node.view_vote[event.message.height - 1]) >= Network.threshold:
                if node.status[event.message.height - 2] != "over":
                    node.status[event.message.height - 2] = "over"
                    node.blockchain_timestamp.append(event.time)
                    node.remove_tx(event.message.height - 1)
                    node.view = event.message.view
                cls.generate_block(node, event.time)

    @classmethod
    def generate_block(cls, node, time):
        result = Transaction.select_tx(Block.block_size_limit, node.tx_pool)
        tx_list = result[0]
        size = result[1]
        if size == 0:
            # print("there is no transaction in the pool.")
            tx_over_time = float(config.tx_num) / config.tx_rate
            if time >= tx_over_time:
                return
            if config.log_flag:
                node.logger.log_info("Node%d becomes leader and will launch a empty block." % node.id)
            # print("Node%d becomes leader and will launch a empty block." % node.id)
        else:
            if config.log_flag:
                node.logger.log_info("Node%d becomes leader and will launch a new block." % node.id)
            # print("Node%d becomes leader and will launch a new block." % node.id)
        block = Block(len(node.blockchain_timestamp), tx_list, size + Message.hash_sig_size)
        Network.add_block(block)
        new_message = PrepareMessage(node.id, "prepare", len(node.blockchain_timestamp) + 1, node.view, block)
        new_event = MessageEvent("send_message", time, new_message)
        new_event.set_node_id(node.id)
        cls.send_message(new_event)
