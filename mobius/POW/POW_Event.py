from Net import Network
from Transaction import Transaction
from Message import *
from Event import *
from POW.POW_Node import *
import random
import numpy as np
import copy

class  DoPOWEvent:

    @classmethod
    def handle_event(cls, event):
        if event.time >= config.over_time and not config.shutdown:
            config.shutdown = True
            EventQueue.add_event(Event("fork_resolution", 9999999999))
        if event.name == "start":
            for i in range(config.node_num):
                node = Network.node[i]
                cls.generate_next_block(node, event.time)
        elif event.name == "send_trans" and not config.shutdown:
            cls.send_tx(event)
        elif event.name == "recv_trans" and not config.shutdown:
            cls.recv_tx(event)
        elif event.name == "gen_block" and not config.shutdown:
            cls.generate_block(Network.node[event.node_id], event.time, event.message.block)
        elif event.name == "send_message":
            cls.send_message(event)
        elif event.name == "recv_message":
            cls.recv_message(event)
        elif event.name == "fork_resolution":
            cls.fork_resolution()

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

    @classmethod
    def recv_tx(cls, event):
        node = Network.node[event.node_id]
        if event.tx_id in node.tx_pool or event.tx_id in node.tx_done_pool:
            return
        node.log_info_tx(event)
        node.add_tx(event.tx_id)
        receiver = Network.gen_recv(node.id, event)
        for item in receiver:
            cls.add_recv_trans_event(node.id, item, event)

    @classmethod
    def send_message(cls, event):
        node = Network.node[event.node_id]
        receiver = Network.gen_recv(node.id, event)
        height = event.message.height
        if event.message.packet_type == "propagate":
            node.log_info_send(event)
            for item in receiver:
                cls.add_recv_message_event(node.id, item, event)

    @classmethod
    def add_recv_message_event(cls, sender_id, receiver_id, event):
        if receiver_id == event.sender:
            return
        sender = Network.node[sender_id]
        bandwidth = sender.bandwidth
        if event.time > sender.execute_time:
            time = event.time + event.message.size / bandwidth
        else:
            time = sender.execute_time + event.message.size / bandwidth
        sender.executor_time = time
        new_event = MessageEvent("recv_message", time, event.message, event.dis + 1)
        new_event.set_node_id(receiver_id)
        new_event.set_sender_id(sender_id)
        EventQueue.add_event(new_event)

    @classmethod
    def recv_message(cls, event):
        node = Network.node[event.node_id]
        miner = Network.node[event.sender]
        receiver = Network.gen_recv(node.id, event)
        if event.message.packet_type == "propagate":
            if event.message.height < node.height:
                return
            block = event.message.block
            blockPrev = block.previous
            lastBlockId = node.last_block().id

            for item in receiver:
                cls.add_recv_message_event(node.id, item, event)

            if blockPrev == lastBlockId:
                node.log_info_recv(event)
                node.blockchain.append(block)
                node.height += 1
                node.remove_tx(block)
                cls.generate_next_block(node, event.time)
            else:
                depth = node.height
                if (depth > len(node.blockchain)):
                    node.log_info_update_local(miner, depth, event.time)
                    cls.update_local_blockchain(node, miner)
                    cls.generate_next_block(node, event.time)
            

    @classmethod
    def generate_next_block(cls, node, time):
        if node.hash_power > 0 and not config.shutdown:
            height = node.height
            result = Transaction.select_tx(Block.block_size_limit, node.tx_pool)
            tx_list = result[0]
            size = result[1]
            if size == 0:
                print("There is no transaction in the pool.")
                tx_over_time = float(config.tx_num) / config.tx_rate
                if time >= tx_over_time:
                    return
            block = Block(len(node.blockchain), tx_list, size + Message.Message.hash_sig_size, node.last_block().id, node, random.randrange(100000000000), time)
            new_message = BlockMessage(node.id, "gen_block", height = height, block = block)
            blockTime = time + cls.cal_time(node)
            new_event = MessageEvent("gen_block", blockTime, new_message)
            new_event.set_sender_id(node.id)
            new_event.set_node_id(node.id)
            EventQueue.add_event(new_event)

    @classmethod
    def cal_time(cls, miner):
        TOTAL_hash_power = Network.hash_power
        hash_power = miner.hash_power / TOTAL_hash_power
        return random.expovariate(hash_power / config.block_interval)

    @classmethod
    def generate_block(cls, node, time, block):
        if config.shutdown:
            return
        blockPrev = block.previous
        if blockPrev == node.last_block().id:
            new_message = BlockMessage(node.id, "propagate", height = node.height, block = block)
            new_event = MessageEvent("send_message", time, new_message)
            new_event.set_node_id(node.id)
            new_event.set_sender_id(node.id)
            node.blockchain.append(block)
            node.remove_tx(block)
            node.height += 1
            cls.send_message(new_event)
            cls.generate_next_block(node, time)

    @classmethod
    def update_local_blockchain(cls, node, miner):
        node.tx_done_pool = copy.deepcopy(miner.tx_done_pool)
        node.tx_pool = copy.deepcopy(miner.tx_pool)
        node.blockchain = copy.deepcopy(miner.blockchain)
        # for i in range(height):
        #     if node.blockchain[i].id != miner.blockchain[i].id:
        #         node.blockchain[i] = miner.blockchain[i]

    @classmethod
    def fork_resolution(cls):
        Network.block = []

        a = []
        for i in Network.node:
            a += [i.blockchain_length()]
        x = max(a)

        b = []
        z = 0
        for i in Network.node:
            if i.blockchain_length() == x:
                b += [i.id]
                z = i.id
        
        if len(b) > 1:
            c = []
            for i in Network.node:
                if i.blockchain_length() == x:
                    c += [i.last_block().miner.id]
            z = np.bincount(c)
            z = np.argmax(z)
        
        for i in Network.node:
            if i.blockchain_length() == x and i.last_block().miner.id == z:
                for item in i.blockchain:
                    Network.block.append(item)
                break
