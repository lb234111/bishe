# -*- coding = utf-8 -*-
import operator
import heapq

class Event:

    def __init__(self, name, time, dis=0):
        # 事件名分为：send_trans, recv_trans, gen_block
        self.name = name
        self.time = time
        # 此处的node_id意为处理该事件的节点id
        self.sender = 0
        self.node_id = 0
        self.tx_id = 0
        self.dis = dis
    
    def __lt__(self, other):
        return self.time < other.time

    def set_sender_id(self, node_id):
        self.sender = node_id

    def set_node_id(self, node_id):
        self.node_id = node_id

    def set_tx_id(self, tx_id):
        self.tx_id = tx_id

    def set_dis(self, dis):
        self.dis = dis


class MessageEvent(Event):

    def __init__(self, name, time, message, dis=0):
        # 事件名分为：send_message, recv_message
        super().__init__(name, time, dis)
        self.message = message


class EventQueue:
    event_list = []

    @classmethod
    def add_event(cls, event):
        heapq.heappush(cls.event_list, event)
        # cls.event_list.append(event)

    @classmethod
    def remove_event(cls):
        del cls.event_list[0]

    @classmethod
    def remove_all(cls):
        cls.event_list = []

    @classmethod
    def get_next_event(cls):
        if cls.get_size() == 0:
            return False
        event = heapq.heappop(cls.event_list)
        # cls.event_list.sort(key=operator.attrgetter('time'), reverse=False)
        # event = cls.event_list[0]
        # if event.time >= Network.block_interval and Network.block_num == 0:
        #     DoEvent.generate_block(Network.node[0], Network.block_interval)
        # cls.remove_event()
        return event

    @classmethod
    def get_size(cls):
        return len(cls.event_list)

    @classmethod
    def is_empty(cls):
        if len(cls.event_list) == 0:
            return True
        else:
            return False
