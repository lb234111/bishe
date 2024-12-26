#!/bin/bash
#对本地四个节点之间的通信进行设置

show_help() {
  echo "Usage: ./network_script.sh [option]"
  echo "Options:"
  echo "  show          Show the delay between ports"
  echo "  set           Set the delay between ports"
  echo "  setWD         Set the bandwidth and delay between ports"
  echo "  cancel        Cancel the setting between ports"
  echo "  --help        Display this help message"
}

#监测实时端口带宽流量的命令 sudo iftop -i lo -f "port 11301"
set_delay() {
  # 对本地四个节点（端口）添加延迟
  sudo tc qdisc add dev lo root handle 1: prio bands 4
  sudo tc qdisc add dev lo parent 1:4 handle 40: netem delay 200ms
  for port in {11301..11304}; do
    sudo tc filter add dev lo protocol ip parent 1:0 prio 4 u32 match ip dport $port 0xffff flowid 1:4
  done

  # sudo tc filter add dev lo protocol ip parent 1:0 prio 4 u32 match ip dport 11301 0xffff flowid 1:4
  # sudo tc filter add dev lo protocol ip parent 1:0 prio 4 u32 match ip dport 11302 0xffff flowid 1:4
  # sudo tc filter add dev lo protocol ip parent 1:0 prio 4 u32 match ip dport 11303 0xffff flowid 1:4
  # sudo tc filter add dev lo protocol ip parent 1:0 prio 4 u32 match ip dport 11304 0xffff flowid 1:4
  # # sudo tc filter add dev lo protocol ip parent 1:0 prio 4 u32 match ip sport 11301 0xffff flowid 1:4
  # # sudo tc filter add dev lo protocol ip parent 1:0 prio 4 u32 match ip sport 11302 0xffff flowid 1:4
  # # sudo tc filter add dev lo protocol ip parent 1:0 prio 4 u32 match ip sport 11303 0xffff flowid 1:4
  # # sudo tc filter add dev lo protocol ip parent 1:0 prio 4 u32 match ip sport 11304 0xffff flowid 1:4

  # 显示当前延迟设置
  echo "The latest delay setting:"
  sudo tc -s qdisc show dev lo

  echo "Delay has been set between the ports."
}

set_delay_bandwidth() {
  # 对本地四个节点（端口）限制带宽
  sudo tc qdisc add dev lo root handle 1:0 htb default 1 
  sudo tc class add dev lo parent 1:0 classid 1:1 htb rate 1000Mbit burst 15k
  #对每个端口设置10Mbit带宽
  sudo tc class add dev lo parent 1:1 classid 1:10 htb rate 100Mbit burst 15k
  # sudo tc qdisc add dev lo parent 1:10 handle 20: netem delay 200ms
  for port in {11301..11310}; do
    sudo tc filter add dev lo protocol ip parent 1:0 prio 1 u32 match ip dport $port 0xffff flowid 1:10
  done

  # 显示当前延迟设置
  echo "The latest delay setting:"
  sudo tc -s qdisc show dev lo

  echo "Delay has been set between the ports."
}

cancel_delay() {
  # 取消延迟设置
  sudo tc qdisc del dev lo root

  echo "Delay has been canceled between the ports."
}

show_delay(){
  # 显示当前延迟设置
  echo "The latest qdisc setting:"
  sudo tc -s qdisc show dev eth0
  echo "The latest filter setting:"
  sudo tc -s filter show dev lo
}
# 处理命令行参数
if [ "$1" == "--help" ]; then
  show_help
elif [ "$1" == "set" ]; then
  set_delay
elif [ "$1" == "setWD" ]; then
  set_delay_bandwidth
elif [ "$1" == "cancel" ]; then
  cancel_delay
elif [ "$1" == "show" ]; then
  show_delay
else
  echo "Invalid option. Use --help to see the available options."
fi
