# 由于模拟器的数值准确性难以实现，最终选择的变化趋势性上的预测，这里选择了形状距离来判断模拟器的模拟是否准确
import numpy as np
import json
import re
import matplotlib.pyplot as plt
import sys

def judge_m(val):
    if val > 0:
        return 1
    elif val == 0:
        return 0
    else:
        return -1

def judge_trend(old, new):
    rate = old / new
    if rate > 0.95 or rate < 1.05:
        return 2
    if new > old:
        return 3
    else:
        return 1

def cal_accuracy_block_size():
    sim_data = []
    real_data = []
    x = []
    for i in range(100, 601, 10):
        x.append(i)

    #选择读取当前路径下的指定文件获取两组数据，并对这两组数据进行相关性分析
    with open('./data/SIM_BLOCKSIZE.json', 'r', encoding='utf-8') as fp1:
        sim_data = json.load(fp1)
        # sim_data = sim_data[:]
        fp1.close()
    with open('./data/REAL_BLOCKSIZE.json', 'r', encoding='utf-8') as fp2:
        real_data = fp2.read()
        real_data = re.sub(r'\t|\n','',real_data)
        real_data = real_data.split("}")
        real_data = [d.strip() + "}" for d in real_data]
        # real_data = list(filter(("}").__ne__), real_data)
        del real_data[-1]
        real_data = [json.loads(d) for d in real_data]
        # real_data = real_data[:51]
        # real_data = json.load(fp2)
        fp2.close()

    sim_arr = []
    real_arr = []
    for item in sim_data:
        sim_arr.append(item["tps"])
    for item in real_data:
        real_arr.append(item["tps"])

    z1 = np.polyfit(x, sim_arr, 4)
    p1 = np.poly1d(z1)
    yvals1 = p1(x)
    z2 = np.polyfit(x, real_arr, 4)
    p2 = np.poly1d(z2)
    yvals2 = p2(x)

    plt.figure(1)
    plt.subplot(1, 2, 1)
    plt.title('sim_polyfitting')
    plt.plot(x, sim_arr, '*', label = 'original values')
    plt.plot(x, yvals1, 'r', label = 'polyfit values')
    plt.xlabel('blockSize')
    plt.ylabel('tps')


    plt.subplot(1, 2, 2)
    plt.title('real_polyfitting')
    plt.plot(x, real_arr, '*', label = 'original values')
    plt.plot(x, yvals2, 'r', label = 'polyfit values')
    plt.xlabel('blockSize')
    plt.ylabel('tps')

    plt.show()

    total = 0
    hit_rate = 0.
    for x in range(100, 601, 5):
        total += 1
        diff = abs(p1(x) - p2(x))
        hit_rate += float(diff) / p2(x)

    hit_rate /= total
    print("平均准确率：", 1 - hit_rate)

    # shape_dis = 0
    # count = 0
    # total = 0

    # for x in range(100, 600, 5):
    #     total += 1
    #     diff_sim = p1(x) - p1(x - 10)
    #     diff_real = p2(x) - p2(x - 10)
    #     m_sim = judge_m(diff_sim)
    #     m_real = judge_m(diff_real)
    #     if m_sim == m_real:
    #         count += 1
    #     shape_dis += abs(m_sim - m_real)

    # print(shape_dis / 50)
    # print(count / total)

    # diff_sim = []
    # diff_real = []

    # size = len(real_data)

    # for i in range(size - 1):
    #     diff_sim.append(sim_data[i + 1]["tps"] - sim_data[i]["tps"])
    #     diff_real.append(real_data[i + 1]["tps"] - real_data[i]["tps"])

    # shape_dis = 0.0
    # count = 0
    # for i in range(len(diff_real)):
    #     m_sim = judge_m(diff_sim[i])
    #     m_real = judge_m(diff_real[i])
    #     if m_real == m_sim:
    #         count += 1
    #     # if i != 0:
    #     #     m_sim *= judge_trend(sim_data[i - 1], sim_data[i])
    #     #     m_real *= judge_trend(real_data[i - 1], real_data[i])
    #     shape_dis += abs(m_sim - m_real)
    
    # shape_dis /= 50
    # print(shape_dis)
    # print(count / len(diff_real))

######################################################
    # tmp_arr1 = []
    # tmp_arr2 = []
    # for item in sim_data:
    #     tmp_arr1.append(item["tps"])
    # for item in real_data:
    #     tmp_arr2.append(item["tps"])


    # v1_numpy = np.array(tmp_arr1)
    # v2_numpy = np.array(tmp_arr2)

    # cos_sim = cosine_similarity(v1_numpy.reshape(1, -1), v2_numpy.reshape(1 -1))
    # print(cos_sim[0][0])

def cal_accuracy_bandwidth():
    sim_data = []
    real_data = []
    x = []
    for i in range(30, 61, 2):
        x.append(i)

    #选择读取当前路径下的指定文件获取两组数据，并对这两组数据进行相关性分析
    with open('./data/SIM_BANDWIDTH.json', 'r', encoding='utf-8') as fp1:
        sim_data = json.load(fp1)
        # sim_data = sim_data[:]
        fp1.close()
    with open('./data/REAL_BANDWIDTH.json', 'r', encoding='utf-8') as fp2:
        real_data = fp2.read()
        real_data = re.sub(r'\t|\n','',real_data)
        real_data = real_data.split("}")
        real_data = [d.strip() + "}" for d in real_data]
        # real_data = list(filter(("}").__ne__), real_data)
        del real_data[-1]
        real_data = [json.loads(d) for d in real_data]
        # real_data = real_data[:51]
        # real_data = json.load(fp2)
        fp2.close()

    sim_arr = []
    real_arr = []
    for item in sim_data:
        sim_arr.append(item["tps"])
    for item in real_data:
        real_arr.append(item["tps"])

    z1 = np.polyfit(x, sim_arr, 4)
    p1 = np.poly1d(z1)
    yvals1 = p1(x)
    z2 = np.polyfit(x, real_arr, 4)
    p2 = np.poly1d(z2)
    yvals2 = p2(x)

    plt.figure(1)
    plt.subplot(1, 2, 1)
    plt.title('sim_polyfitting')
    plt.plot(x, sim_arr, '*', label = 'original values')
    plt.plot(x, yvals1, 'r', label = 'polyfit values')
    plt.xlabel('bandwidth')
    plt.ylabel('tps')


    plt.subplot(1, 2, 2)
    plt.title('real_polyfitting')
    plt.plot(x, real_arr, '*', label = 'original values')
    plt.plot(x, yvals2, 'r', label = 'polyfit values')
    plt.xlabel('bandwidth')
    plt.ylabel('tps')

    plt.show()

    total = 0
    hit_rate = 0.
    for x in range(30, 61, 1):
        total += 1
        diff = abs(p1(x) - p2(x))
        hit_rate += float(diff) / p2(x)

    hit_rate /= total
    print("平均准确率：", 1 - hit_rate)


if __name__ == '__main__':
    test_type = sys.argv[1]
    if test_type == "blocksize":
        cal_accuracy_block_size()
    elif test_type == "bandwidth":
        cal_accuracy_bandwidth()