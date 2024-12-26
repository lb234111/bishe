# -*- coding = utf-8 -*-
import Net


class Transaction:

    def __init__(self, tx_id, sender, receiver, size, time):
        self.tx_id = tx_id
        self.sender = sender
        self.receiver = receiver
        self.size = size
        self.time = time

    @classmethod
    def select_tx(cls, block_size, tx_list):
        block_list = []
        size = 0
        for tx_id in tx_list:
            if block_size < Net.Network.transaction_list[tx_id].size:
                break
            block_size -= Net.Network.transaction_list[tx_id].size
            size += Net.Network.transaction_list[tx_id].size
            block_list.append(tx_id)
        result = ()
        result += (block_list,)
        result += (size,)
        return result
