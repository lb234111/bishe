# -*- coding = utf-8 -*-
from PBFT.PBFT_Event import *
from HotStuff.HotStuff_Event import *
from TBFT.TBFT_Event import *
from POW.POW_Event import *
import random
import config
import matplotlib.pyplot as plt
from PStatistics import PStatistics
from Net import Network
import os
import json
import sys
import time

tx_time = 0

def random_generate_tx(node_num, tx_rate):
    global tx_time
    tx_list = []
    time_list = []
    for i in range(tx_rate):
        tx_time = tx_time + (1 / float(tx_rate))
        time_list.append(tx_time)
        sender = random.randint(0, node_num - 1)
        receiver = random.randint(0, node_num - 1)
        # size = random.randint(1, 10)
        # 单位字节，一般默认为512
        size = config.tx_size
        tx = Transaction(i, sender, receiver, size, tx_time)
        tx_list.append(tx)
    return tx_list, time_list


def add_node_event(network, node_id, event):
    network.node[node_id].add_event(event)


def main():
    config_json = False
    if config_json:
        with open('./config.json','r',encoding='utf8')as fp:
            json_data = json.load(fp)
            config.tx_size = json_data['tx_size']
            config.tx_pool_size = json_data['tx_pool_size']
            config.node_num = json_data['node_num']
            config.network_flag = json_data['network_flag']
            config.bandwidth = json_data['bandwidth']
            config.tx_flag = json_data['tx_flag']
            config.tx_num = json_data['tx_num']
            config.tx_rate = json_data['tx_rate']
            config.random_topology_flag = json_data['random_topology_flag']
            config.network_topology = json_data['network_topology']
            config.protocol = json_data['protocol']
            config.fan_out = json_data['fan_out']
            config.heartbeat_interval = json_data['heartbeat_interval']
            config.bucket_size = json_data['bucket_size']
            config.block_size_limit = json_data['block_size_limit']
            config.hash_sig_size = json_data['hash_sig_size']
            config.start_time = json_data['start_time']
            config.over_time = json_data['over_time']
            config.consensus = json_data['consensus']
            config.process_cost = json_data['process_cost']
            config.log_flag = json_data['log_flag']
            config.timeout_propose = json_data['timeout_propose']
            config.timeout_propose_delta = json_data['timeout_propose_delta']
            config.timeout_pre_vote = json_data['timeout_pre_vote']
            config.timeout_pre_vote_delta = json_data['timeout_pre_vote_delta']
            config.timeout_pre_commit = json_data['timeout_pre_commit']
            config.timeout_pre_commit_delta = json_data['timeout_pre_commit_delta']
            fp.close()
    if config.random_topology_flag == 0 and len(config.network_topology) != config.node_num:
        print("ERROR:network_topology: len(network_topology) != node_num")
        return -1,-1,-1
    if config.protocol == "gossip" and config.fan_out >=config.node_num:
        print("ERROR:fan_out: fan_out >= node_num")
    path = "./log/"
    files = os.listdir(path=path)
    for file in files:
        os.remove(path=path+file)
    network = Network()
    network.reset()
    Block.block_size_limit = config.block_size_limit
    if config.protocol == "gossip":
        network.init(config.node_num, config.consensus, config.protocol, fan_out=config.fan_out)
    elif config.protocol == "kad":
        network.init(config.node_num, config.consensus, config.protocol, bucket_size=config.bucket_size)
    queue = EventQueue()
    global tx_time
    tx_time = 0

    tx_list = []
    time_list = []
    tx_id = 0
    if config.tx_flag == 0:
        tx_list, time_list = random_generate_tx(network.node_num, config.tx_rate)
        network.transaction_list += tx_list
    for i in range(len(tx_list)):
        event = Event("send_trans", time_list[i])
        # if protocol == "kad":
        #     event.set_dis(0)
        event.set_node_id(tx_list[i].sender)
        event.set_tx_id(tx_id)
        tx_id += 1
        queue.add_event(event)
    tx_list.clear()
    time_list.clear()

    # 第一个区块的共识开始时间设定
    if Network.consensus == "TBFT" or Network.consensus == "POW":
        new_event = Event("start", config.start_time)
        queue.add_event(new_event)
    else:
        new_event = Event("gen_block", config.start_time)
        new_event.set_node_id(0)
        # add_node_event(network, 0, new_event)
        queue.add_event(new_event)
    gen_tx_flag = True
    while not queue.is_empty():
        event = queue.get_next_event()
        if gen_tx_flag:
            if tx_time > config.tx_num / config.tx_rate:
                gen_tx_flag = False
            elif event.time > tx_time:
                tx_list, time_list = random_generate_tx(network.node_num, config.tx_rate)
                network.transaction_list += tx_list
                for i in range(len(tx_list)):
                    new_event = Event("send_trans", time_list[i])
                    new_event.set_node_id(tx_list[i].sender)
                    new_event.set_tx_id(tx_id)
                    tx_id += 1
                    queue.add_event(new_event)
                tx_list.clear()
                time_list.clear()
            
        if network.consensus == "PBFT":
            DoPBFTEvent.handle_event(event)
        elif network.consensus == "HotStuff":
            DoHotStuffEvent.handle_event(event)
        elif network.consensus == "TBFT":
            DoTBFTEvent.handle_event(event)
        elif network.consensus == "POW":
            DoPOWEvent.handle_event(event)

    statics = PStatistics(network)
    return statics.getTPS(), statics.getTxLatency(), statics.getAvgLatency()

def test_block_size():
    tps_list = []
    # tx_latency_list = []
    # avg_latency_list = []
    block_size_list = []
    write_data = []
    config.bandwidth = 13107200
    config.start_time = 0
    config.over_time = 2
    config.tx_rate = 4000
    config.tx_num = 8000
    for i in range(51):
        # 100 - 600
        # if i != 0 and i != 50:
        #     continue
        tmp_size = 10 * (i + 10)
        block_size_list.append(tmp_size)
        config.block_size_limit = tmp_size * config.tx_size
        tps, tx_latency, avg_latency = main()
        tps_list.append(tps)
        # tx_latency_list.append(tx_latency)
        # avg_latency_list.append(avg_latency)
        print("block_size:%d, tps:%f" % (tmp_size, tps))
        write_data.append({"blockSize" : tmp_size, "tps" : tps})

    filename = time.asctime() + '_SIM_BLOCKSIZE.json'
    filename = "./output/" + filename
    with open(filename, "w") as fp:
        json.dump(write_data, fp)
        fp.close()
    plt.figure(1)
    plt.subplot(1,3,1)
    plt.title("tps--blockSize")
    plt.xlabel("blockSize")
    plt.ylabel("tps")
    plt.grid(True,c="k",axis="both",linestyle="--",linewidth=0.3)
    plt.plot(block_size_list,tps_list,c="r",lw=1,marker="o",ms=4,mfc="r")

    # plt.subplot(1,3,2)
    # plt.title("txLatency---blockSize")
    # plt.xlabel("blockSize")
    # plt.ylabel("txLatency")
    # plt.grid(True,c="k",axis="both",linestyle="--",linewidth=0.3)
    # plt.plot(block_size_list,tx_latency_list,c="r",lw=1,marker="o",ms=4,mfc="r")

    # plt.subplot(1,3,3)
    # plt.title("blockLatency---blockSize")
    # plt.xlabel("blockSize")
    # plt.ylabel("blockLatency")
    # plt.grid(True, c="k", axis="both", linestyle="--", linewidth=0.3)
    # plt.plot(block_size_list, avg_latency_list, c="r", lw=1, marker="o", ms=4, mfc="r")
    plt.show()

    # print("blockSize:" + block_size_list)
    # print("tps_list:%" + tps_list)
    # print("avg_latency_list:" + avg_latency_list)

def test_bandwidth():
    tps_list = []
    # tx_latency_list = []
    # avg_latency_list = []
    bandwidth_list = []
    write_data = []
    config.start_time = 0
    config.over_time = 2
    config.tx_rate = 2000
    config.tx_num = 4000
    config.block_size_limit = 100 * config.tx_size
    # config.process_cost = 1 / 20
    for i in range(16):
        # 30 - 60
        # if i != 0 and i != 16:
        #     continue
        tmp_speed = 2 * (i + 15)
        config.bandwidth = tmp_speed * 1024 * 1024 / 8
        tps, tx_latency, avg_latency = main()
        tps_list.append(tps)
        bandwidth_list.append(tmp_speed)
        # tx_latency_list.append(tx_latency)
        # avg_latency_list.append(avg_latency)
        print("bandwidth:%d, tps:%f" % (tmp_speed, tps))
        write_data.append({"bandwidth" : tmp_speed, "tps" : tps})

    filename = time.asctime() + '_SIM_BANDWIDTH.json'
    filename = "./output/" + filename
    with open(filename, "w") as fp:
        json.dump(write_data, fp)
        fp.close()
    plt.figure(1)
    plt.subplot(1,3,1)
    plt.title("tps--bandwidth")
    plt.xlabel("bandwidth")
    plt.ylabel("tps")
    plt.grid(True,c="k",axis="both",linestyle="--",linewidth=0.3)
    plt.plot(bandwidth_list,tps_list,c="r",lw=1,marker="o",ms=4,mfc="r")
    
    plt.show()

    print("bandwidth_list:" + bandwidth_list)
    print("tps_list:%" + tps_list)


if __name__ == '__main__':
    if len(sys.argv) == 1:
        tps, tx_latency, avg_latency = main()
        print("tps:", tps)
    else:
        test_type = sys.argv[1]
        if test_type == "blocksize":
            print("Start test BlockSize")
            test_block_size()
        elif test_type == "bandwidth":
            print("Start test bandwidth")
            test_bandwidth()
        elif test_type == "time":
            print("Start test time")
            time_start = time.time()
            main()
            time_end = time.time()
            print("耗时:", time_end - time_start)
        else:
            print("Please input test type")
