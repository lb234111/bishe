# -*- coding = utf-8 -*-
import config
import Net

class Block:
    # 区块大小上限
    block_size_limit = config.block_size_limit

    # 出块间隔，每interval秒进行出块，随后进行leader的轮换(该参数根据自定义的出块规则选择是否需要)
    # interval = 5

    def __init__(self, block_id, tx_list=[], block_size=0, previous=-1, miner=None, id=0, timestamp=0):
        self.block_id = block_id
        self.tx_list = tx_list
        self.block_size = block_size
        self.previous = previous
        self.miner = miner
        self.id = id
        self.timestamp = timestamp

    def add_tx(self, tx_id):
        self.tx_list.append(tx_id)


class Message:
    hash_sig_size = config.hash_sig_size

    def __init__(self, idx, packet_type, height=0, view=0, round=0):
        # idx为节点的索引值，意为node_id=idx的节点发出的消息
        # packet_type分为prepare,sig,commit
        # packet_type在TBFT中分为proposal,pre_vote,pre_commit,commit_block
        # packet_type在POW中为gen_block,propagate_block
        # height为该消息对应的区块高度
        # view为该消息对应的view值
        # round为TBFT共识中的view值
        self.idx = idx
        self.packet_type = packet_type
        self.height = height
        self.view = view
        self.round = round


class PrepareMessage(Message):

    def __init__(self, idx, packet_type, height=0, view=0, block=None):
        super(PrepareMessage, self).__init__(idx, packet_type, height, view)
        self.block = block
        self.size = block.block_size + super().hash_sig_size

class ProposalMessage(Message):

    def __init__(self, idx, packet_type, height=0, view=0, round=0, block=None):
        super().__init__(idx, packet_type, height, view, round)
        self.block = block
        self.size = block.block_size + super().hash_sig_size


class BlockInfoMessage(Message):
    # 在TBFT共识算法中，每个节点提交了一个区块后，都会把这个提交的区块进行广播
    # 这个完整的区块中不仅包含交易信息，还包含达成这个区块共识的其他节点的签名信息
    def __init__(self, idx, packet_type, height=0, view=0, round=0, block=None):
        super().__init__(idx, packet_type, height, view, round)
        self.block = block
        self.size = block.block_size + super().hash_sig_size * Net.Network.threshold
        
class BlockMessage(Message):

    def __init__(self, idx, packet_type, height=0, view=0, round=0, block=None):
        super().__init__(idx, packet_type, height, view, round)
        self.block = block
        self.size = block.block_size + super().hash_sig_size
