import Log
import config
import Message
import Net

def process_block(block):
    return config.process_cost


class POWNode:

    def __init__(self, id, tx_pool_size, hash_power):
        self.id = id
        self.bandwidth = 0
        self.logger = Log.Log(id)
        self.logger.file_handle()
        self.blockchain = []
        self.tx_pool = []
        self.tx_done_pool = []
        self.tx_pool_size = tx_pool_size
        self.height = 0
        self.hash_power = hash_power
        self.execute_time = 0
        self.generate_gensis_block()

    def generate_gensis_block(self):
        block = Message.Block(block_id = 0, previous = -1)
        self.blockchain.append(block)

    def set_bandwidth(self, bandwidth):
        self.bandwidth = bandwidth

    def add_tx(self, tx_id):
        self.tx_pool.append(tx_id)

    def last_block(self):
        return self.blockchain[-1]

    def add_block(self, block):
        self.blockchain.append(block)

    def blockchain_length(self):
        return len(self.blockchain) - 1
    
    def remove_tx(self, block):
        for tx in block.tx_list:
            self.tx_done_pool.append(tx)
            if tx in self.tx_pool:
                self.tx_pool.remove(tx)
    
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
            if event.message.packet_type == "propagate":
                size = event.message.size
            self.logger.log_info("[Send] [Broadcast] [%f] [%s], Height:%d, Size=%d, BlockId=%d ,Id=%d" %
                                 (event.time, event.message.packet_type, event.message.height, size,
                                   event.message.block.block_id, self.id))
            if config.single_log:
                Net.Network.logger.log_info("[Node %d] [Send] [Broadcast] [%f] [%s], Height:%d, Size=%d, BlockId=%d, Id=%d" %
                                 (self.id, event.time, event.message.packet_type, event.message.height, size,
                                  event.message.block.block_id, self.id))

    def log_error(self, msg):
        if config.log_flag:
            self.logger.log_error(msg)

    def log_info_recv(self, event):
        if config.log_flag:
            p_type = event.message.packet_type
            if p_type == "propagate":
                size = event.message.size
            self.logger.log_info("[Recv] [from %d] [%f] [%s], Height=%d, Size=%d, BlockId=%d, Id=%d" %
                                 (event.message.idx, event.time, p_type, event.message.height,
                                  size, event.message.block.block_id, self.id))
            if config.single_log:
                Net.Network.logger.log_info("[Node %d] [Recv] [from %d] [%f] [%s], Height=%d, Size=%d, BlockId=%d, Id=%d" %
                                 (self.id, event.message.idx, event.time, p_type, event.message.height,
                                  size, event.message.block.block_id, self.id))

    def log_info_update_local(self, miner, depth, time):
        if config.log_flag:
            self.logger.log_info("[Merge] [from %d] [%f], Height=%d, Id=%d" %
                                 (miner.id, time, depth, self.id))