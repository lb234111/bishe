# -*- coding = utf-8 -*-
import operator
from Net import Network
from Message import *
from Transaction import Transaction
from PBFT.PBFT_Node import process_block
from Event import *


def sync_tx_from_neighbor(node_id, lack_tx_list):
    size = len(lack_tx_list) * 512
    time = size / Network.node[node_id].bandwidth
    return time


class DoPBFTEvent:

    @classmethod
    def handle_event(cls, event):
        if event.time >= config.over_time:
            config.shutdown = True
        if event.name == "send_trans" and not config.shutdown:
            cls.send_tx(event)
        elif event.name == "recv_trans" and not config.shutdown:
            cls.recv_tx(event)
        elif event.name == "send_message":
            cls.send_message(event)
        elif event.name == "recv_message":
            cls.recv_message(event)
        elif event.name == "gen_block" and not config.shutdown:
            cls.generate_block(Network.node[event.node_id], event.time)

    @classmethod
    def add_recv_trans_event(cls, sender_id, receiver_id, event):
        if receiver_id == event.sender:
            return
        sender = Network.node[sender_id]
        receiver = Network.node[receiver_id]
        bandwidth = sender.bandwidth
        tx = Network.transaction_list[event.tx_id]
        if event.time >= sender.execute_time:
            time = event.time + tx.size / bandwidth
        else:
            time = sender.execute_time + tx.size / bandwidth
        sender.execute_time = time
        # if Network.protocol.name == "gossip":
        #     new_event = Event("recv_trans", time)
        # elif Network.protocol.name == "kad":
        #     new_event = Event("recv_trans", time, event.dis + 1)
        new_event = Event("recv_trans", time, event.dis + 1)
        new_event.set_node_id(receiver_id)
        new_event.set_sender_id(sender_id)
        new_event.set_tx_id(event.tx_id)
        # Network.node[receiver_id].add_event(new_event)
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
        # tx = Network.transaction_list[event.tx_id]
        # print("node:%d send the tx:%d, time:%f" % (node.id, event.tx_id, event.time))

        receiver = Network.gen_recv(node.id, event)

        # if Network.protocol.name == "gossip":
        #     # 随机生成fan_out个节点作为交易接受对象：创建fan_out个接受交易的事件，接受交易的事件会让该交易继续传播开来
        #     receiver = Network.protocol.neighbor[node.id]
        # elif Network.protocol.name == "kad":
        #     event.set_dis(0)
        #     # 根据该节点的转发表生成接受节点，接受节点的选择要从distance=0开始
        #     receiver = Network.protocol.gen_recv(node.id, 0)

        for item in receiver:
            cls.add_recv_trans_event(node.id, item, event)
        node.check_event_list()

    @classmethod
    # 接受交易仅被非交易发起者调用，交易的接受和传播，每笔交易每个节点仅执行一次，再次收到该交易不会做什么
    # 同时需要考虑特殊情况：如果节点收到了一个pre_req但是块中的交易本地尚未同步到，当节点收到的交易后，发现块中的交易全部合法，随后需要进行签名广播sign_req
    def recv_tx(cls, event):
        node = Network.node[event.node_id]
        if event.tx_id in node.tx_pool or event.tx_id in node.tx_done_pool:
            # print("node:%d drop the duplicate/done event" % (node.id))
            return
        node.add_tx(event.tx_id)
        # print("node:%d receive the tx:%d, time:%f" % (node.id, event.tx_id, event.time))
        # tx = Network.transaction_list[event.tx_id]

        receiver = Network.gen_recv(node.id, event)

        # if Network.protocol.name == "gossip":
        #     receiver = Network.protocol.neighbor[node.id]
        # elif Network.protocol.name == "kad":
        #     receiver = Network.protocol.gen_recv(node.id, event.dis)
        for item in receiver:
            # 如果随机挑选的节点中已经接受过该笔交易的，直接跳过
            # if event.tx_id in Network.node[item].tx_pool:
            #     continue
            cls.add_recv_trans_event(node.id, item, event)
        node.check_event_list()

        # for index in range(len(node.status)):
        #     if node.status[index] == "wait_tx":
        #         prepare_message = node.prepare_message[index]
        #         if event.tx_id in prepare_message.block.tx_list and node.tx_pool >= prepare_message.block.tx_list:
        #             # 发现一个prepare包中的所有交易都已经集齐并合法，即可进入发送签名，进入下一阶段
        #
        #             # node.status[index] = "wait_sig"
        #             # msg_receiver = Network.generate_receiver(node.id)
        #             # 验证通过，将自身加入签名通过列表中，当该列表中节点数量达到2f+1即可发送commit消息
        #             node.add_sig_id(prepare_message.height, node.id)
        #             # 时间分为将prepare消息中的区块进行处理，后续即可进行发送sig消息
        #             time = event.time + node.process_block(prepare_message.block)
        #             new_message = Message(node.id, "sig", prepare_message.height, prepare_message.view)
        #             new_event = MessageEvent("send_message", time, new_message)
        #             new_event.set_node_id(node.id)
        #             cls.send_message(new_event)
        #             if len(node.sig_list[prepare_message.height - 1]) >= Network.threshold:
        #                 # 节点在进入sig阶段后，发现加上自己的签名后，数量达到2f+1了，立即进入wait_commit阶段，并发送commit包
        #                 # node.status[index] = "wait_commit"
        #                 new_message = Message(node.id, "commit", prepare_message.height, prepare_message.view)
        #                 new_event = MessageEvent("send_message", time, new_message)
        #                 new_event.set_node_id(node.id)
        #                 cls.send_message(new_event)

    @classmethod
    def add_recv_message_event(cls, sender_id, receiver_id, event):
        # 此函数是转发message，event的发起者id并不是sender转发者的id
        if receiver_id == event.sender:
            return
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
        # if Network.protocol.name == "gossip":
        #     new_event = MessageEvent("recv_message", time, event.message)
        # elif Network.protocol.name == "kad":
        #     new_event = MessageEvent("recv_message", time, event.message, event.dis + 1)
        new_event = MessageEvent("recv_message", time, event.message, event.dis + 1)
        new_event.set_node_id(receiver_id)
        new_event.set_sender_id(sender_id)
        # Network.node[receiver_id].add_event(new_event)
        EventQueue.add_event(new_event)

    @classmethod
    def send_message(cls, event):
        node = Network.node[event.node_id]
        # if Network.protocol.name == "gossip":
        #     receiver = Network.protocol.neighbor[node.id]
        # elif Network.protocol.name == "kad":
        #     event.set_dis(0)
        #     receiver = Network.protocol.gen_recv(node.id, 0)
        receiver = Network.gen_recv(node.id, event)
        node.init_status(event)
        # 由leader产生，向全网广播send_prep消息，开始共识
        if event.message.packet_type == "prepare":
            if node.id != (event.message.height - 1) % Network.node_num:
                print("ERROR: ONLY LEADER CAN SEND PREPARE MESSAGE!")
                exit(-1)
            node.add_message(event)
            # node.add_time(event)
            node.status[event.message.height - 1] = "wait_sig"
            node.log_info_send(event)
            # print("node:%d send the prepare message:%d, time:%f" % (node.id, event.message.height, event.time))
            node.add_sig_id(event)
            for item in receiver:
                cls.add_recv_message_event(node.id, item, event)
            # node.check_event_list()
        elif event.message.packet_type == "sig":
            node.status[event.message.height - 1] = "wait_sig"
            node.log_info_send(event)
            node.add_sig_id(event)
            # print("node:%d send the sig message:%d, time:%f" % (node.id, event.message.height, event.time))
            for item in receiver:
                cls.add_recv_message_event(node.id, item, event)
            # node.check_event_list()
        elif event.message.packet_type == "commit":
            node.status[event.message.height - 1] = "wait_commit"
            node.log_info_send(event)
            node.add_commit_id(event)
            # print("node:%d send the commit message:%d, time:%f" % (node.id, event.message.height, event.time))
            for item in receiver:
                cls.add_recv_message_event(node.id, item, event)
            # node.check_event_list()

    @classmethod
    def recv_message(cls, event):
        node = Network.node[event.node_id]
        # if Network.protocol.name == "gossip":
        #     receiver = Network.protocol.neighbor[node.id]
        # elif Network.protocol.name == "kad":
        #     receiver = Network.protocol.gen_recv(node.id, event.dis)
        receiver = Network.gen_recv(node.id, event)
        node.init_status(event)
        if event.message.packet_type == "prepare":
            if node.status[event.message.height - 1] == "normal":
                # 注意，节点如果处于非normal状态，表明节点收到这个prepare包过迟，他都已经收集到了足够了sig包或者是commit包
                node.add_message(event)
                # 每个节点需要记录每个共识开始的时间
                # node.add_time(event)
                prepare = node.prepare_message[event.message.height - 1]
                # 只有leader节点才会发送prepare消息，该消息等同于leader的sig消息，leader节点并不会发送sig消息，节点收到来自leader节点的prepare消息直接在sig_list中加入leader节点
                node.add_sig_id(event)
                # print("node:%d receive the prepare message:%d from node:%d, time:%f" % (node.id, event.message.height,
                #                                                                         event.message.idx, event.time))
                time = event.time + process_block(prepare.block)
                lack_tx_list = list(set(event.message.block.tx_list) - set(node.tx_pool))
                if len(lack_tx_list) > 0:
                    # 如果发现本地少了区块中的某些交易，则向其他节点索取相关交易的信息
                    time += sync_tx_from_neighbor(node.id, lack_tx_list)
                # node.status[event.message.height - 1] = "wait_sig"


                # 节点处于wait_sig阶段则表明区块中的所有交易都已到齐，可以进行区块的处理，并进行消息的转发，同时还要进行sig消息的发送
                new_message = Message(node.id, "sig", prepare.height, prepare.view)
                new_event = MessageEvent("send_message", time, new_message)
                new_event.set_node_id(node.id)
                # 发送给其他节点，自己签名了
                cls.send_message(new_event)
                # else:
                #     node.status[event.message.height - 1] = "wait_tx"
                if len(node.sig_list[event.message.height - 1]) >= Network.threshold:
                    new_message = Message(node.id, "commit", prepare.height, prepare.view)
                    new_event = MessageEvent("send_message", time, new_message)
                    new_event.set_node_id(node.id)
                    # 发送给其他节点，自己commit了
                    cls.send_message(new_event)
                for item in receiver:
                    # 这里是转发这个消息给其他节点
                    cls.add_recv_message_event(node.id, item, event)
                # node.check_event_list()
            else:
                # 注意只有leader节点才会发送prepare消息，当节点第一次收到prepare消息后就会改变自身的status，后续收到其他节点转发来的重复的prepare消息可以直接无视
                return
        elif event.message.packet_type == "sig":
            if len(node.sig_list) >= event.message.height and event.message.idx \
                    in node.sig_list[event.message.height - 1]:
                # 已经收到过了
                return
            else:
                for item in receiver:
                    cls.add_recv_message_event(node.id, item, event)
                # node.check_event_list()
            # node.log_info_recv(event)
            # print("node:%d receive the sig message:%d from node:%d, time:%f" % (node.id, event.message.height,
            #                                                                     event.message.idx, event.time))
            node.add_sig_id(event)
            if node.status[event.message.height - 1] == "wait_commit" or \
                    node.status[event.message.height - 1] == "over":
                # 已经进入commit状态或这个区块的共识已经结束了，则对这个sig消息没有兴趣
                return
            if len(node.sig_list[event.message.height - 1]) >= Network.threshold:
                if node.status[event.message.height - 1] == "normal":
                    # node.consensus_block_id = Network.block_num - 1
                    node.prepare_message[event.message.height - 1] = \
                        Network.node[event.message.idx].prepare_message[event.message.height - 1]
                # node.status[event.message.height - 1] = "wait_commit"
                new_message = Message(node.id, "commit", event.message.height, event.message.view)
                new_event = MessageEvent("send_message", event.time, new_message)
                new_event.set_node_id(node.id)
                cls.send_message(new_event)
        elif event.message.packet_type == "commit":
            if len(node.commit_list) >= event.message.height and event.message.idx in \
                    node.commit_list[event.message.height - 1]:
                return
            else:
                for item in receiver:
                    cls.add_recv_message_event(node.id, item, event)
                # node.check_event_list()
                node.add_commit_id(event)
                # node.log_info_recv(event)
                # print("node:%d receive the commit message:%d from node:%d, time:%f" % (node.id, event.message.height,
                #                                                                        event.message.idx, event.time))
                if node.status[event.message.height - 1] == "over":
                    # 表明该节点已经对该区块达成共识，无需处理该区块的commit消息
                    return
                else:
                    if len(node.commit_list[event.message.height - 1]) >= Network.threshold:
                        node.blockchain_timestamp.append(event.time)
                        node.logger.log_info("node:%d adds the block:%d to its chain, time:%f" % (
                            node.id, event.message.height, event.time))
                        # print("node:%d adds the block:%d to its chain, time:%f" % (
                        #     node.id, event.message.height, event.time))
                        node.status[event.message.height - 1] = "over"
                        node.remove_tx(event.message.height)
                        if node.is_leader():
                            cls.generate_block(node, event.time)
                            # time_start = node.start_time[event.message.height - 1]
                            # time_end = event.time
                            # if time_end - time_start >= Network.block_interval:
                            #     cls.generate_block(node, event.time)
                            # else:
                            #     new_event = Event("gen_block", time_start + Network.block_interval)
                            #     new_event.set_node_id(node.id)
                            #     EventQueue.add_event(new_event)

    @classmethod
    def generate_block(cls, node, time):
        result = Transaction.select_tx(Block.block_size_limit, node.tx_pool)
        tx_list = result[0]
        size = result[1]
        if size == 0:
            print("There is no transaction in the pool.")
            tx_over_time = float(config.tx_num) / config.tx_rate
            # if EventQueue.is_empty():
            #     print("Over")
            #     return

            if time >= tx_over_time:
                return
            if config.log_flag:
                node.logger.log_info("Node%d becomes leader and will launch a empty block." % node.id)
            # print("leader will send a empty block.")
            # else:
            #     print("After block_interval try again")
            #     time = time + Network.block_interval
            #     new_event = Event("gen_block", time)
            #     new_event.set_node_id(node.id)
            #     EventQueue.add_event(new_event)
            #     return
        else:
            if config.log_flag:
                node.logger.log_info("Node%d becomes leader and will launch a new block." % node.id)
            # print("leader%d will launch next block." % node.id)
        block = Block(len(node.blockchain_timestamp), tx_list, size + Message.hash_sig_size)
        Network.add_block(block)
        new_message = PrepareMessage(node.id, "prepare", len(node.blockchain_timestamp) + 1, node.view, block)
        new_event = MessageEvent("send_message", time, new_message)
        new_event.set_node_id(node.id)
        new_event.set_sender_id(node.id)
        cls.send_message(new_event)
