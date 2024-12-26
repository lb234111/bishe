import config
class PStatistics:
    def __init__(self, network):
        self.network = network

    def getTPS(self):
        tx_num = 0
        time = 0
        if config.consensus == "POW":
            for block in self.network.block:
                tx_num += len(block.tx_list)
            time = self.network.block[-1].timestamp
        else:
            for node in self.network.node:
                time += node.blockchain_timestamp[-1]
            for i in range(self.network.block_num):
                block = self.network.block[i]
                tx_num += len(block.tx_list)
            time = time / self.network.node_num
        time = time - config.start_time
        return tx_num / time
    
    def getTxLatency(self):
        # 平均每笔交易的打包时
        time = 0
        timestamp = []
        tx_num = 0
        if config.consensus == "POW":
            for block in self.network.block:
                tx_num += len(block.tx_list)
            time = self.network.block[-1].timestamp
        else:
            for index in range(self.network.node_num):
                if index == 0:
                    timestamp = self.network.node[index].blockchain_timestamp
                    continue
                timestamp = [i+j for i,j in zip(timestamp, self.network.node[index].blockchain_timestamp)]
            timestamp = [item / self.network.node_num for item in timestamp]
            for i in range(self.network.block_num):
                block = self.network.block[i]
                tx_num += len(block.tx_list)
                for tx_id in block.tx_list:
                    time += timestamp[i] - self.network.transaction_list[tx_id].time
        return time / tx_num

    def getAvgLatency(self):
        # 平均出块延迟
        time = 0
        if config.consensus == "POW":
            block_list = self.network.block
            return block_list[-1].timestamp / len(block_list)
        else:
            timestamp = []
            for index in range(self.network.node_num):
                if index == 0:
                    timestamp = self.network.node[index].blockchain_timestamp
                    continue
                timestamp = [i + j for i, j in zip(timestamp, self.network.node[index].blockchain_timestamp)]
            timestamp = [item / self.network.node_num for item in timestamp]
            for i in range(len(timestamp)):
                if i == 0:
                    continue
                time += timestamp[i] - timestamp[i - 1]
            time = time - config.start_time
            return time / (self.network.block_num - 1)
    



