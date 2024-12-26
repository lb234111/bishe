/*
Copyright (C) BABEC. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package net

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"chainmaker.org/chainmaker/common/v2/msgbus"
	rootLog "chainmaker.org/chainmaker/logger/v2"
	"chainmaker.org/chainmaker/net-common/common/priorityblocker"
	configPb "chainmaker.org/chainmaker/pb-go/v2/config"
	netPb "chainmaker.org/chainmaker/pb-go/v2/net"
	"chainmaker.org/chainmaker/protocol/v2"
	"chainmaker.org/chainmaker/recorderfile"
	"github.com/gogo/protobuf/proto"
)

var (
	ErrorChainMsgBusBeenBound = errors.New("chain msg bus has been bound")
	ErrorNetNotRunning        = errors.New("net instance is not running")
)

const (
	topicNamePrefix            = "topic"
	consensusTopicNamePrefix   = "consensus_topic"
	msgBusTopicPrefix          = "msgbus_topic"
	msgBusConsensusTopicPrefix = "msgbus_consensus_topic"
	msgBusMsgFlagPrefix        = "msgbus"

	topicSeparator = "::"
)

// CreateFlagWithPrefixAndMsgType will join prefix with msg type string.
func CreateFlagWithPrefixAndMsgType(prefix string, msgType netPb.NetMsg_MsgType) string {
	var builder strings.Builder
	builder.WriteString(prefix)
	builder.WriteString(topicSeparator)
	builder.WriteString(msgType.String())
	return builder.String()
}

var _ protocol.NetService = (*NetService)(nil)

// NetService provide a net service for modules.
type NetService struct {
	chainId                   string
	localNet                  protocol.Net
	msgBus                    msgbus.MessageBus
	logger                    *rootLog.CMLogger
	configWatcher             *ConfigWatcher
	netContractEventSubscribe *NetContractEventSubscribe
	consensusNodeIds          map[string]struct{}
	consensusNodeIdsLock      sync.RWMutex

	ac            protocol.AccessControlProvider
	revokeNodeIds sync.Map // nolint: structcheck,unused // node id of node cert revoked , map[string]struct{}
	vmWatcher     *VmWatcher
}

// NewNetService create a new net service instance.
func NewNetService(chainId string, localNet protocol.Net, ac protocol.AccessControlProvider) *NetService {
	logger := rootLog.GetLoggerByChain(rootLog.MODULE_NET, chainId)
	ns := &NetService{
		chainId:          chainId,
		localNet:         localNet,
		consensusNodeIds: make(map[string]struct{}),
		ac:               ac,
		logger:           logger,
	}
	return ns
}

// BroadcastMsg broadcast a net msg to other nodes belongs to the same chain.
func (ns *NetService) BroadcastMsg(msg []byte, msgType netPb.NetMsg_MsgType) error {
	err := ns.broadcastMsg(msg, CreateFlagWithPrefixAndMsgType(topicNamePrefix, msgType))
	if err != nil {
		return err
	}
	return nil
}

func (ns *NetService) broadcastMsg(msg []byte, topic string) error {
	return ns.localNet.BroadcastWithChainId(ns.chainId, topic, msg)
}

// Subscribe a pub-sub topic for receiving the msg that be broadcast by the other node.
func (ns *NetService) Subscribe(msgType netPb.NetMsg_MsgType, handler protocol.MsgHandler) error {
	err := ns.subscribe(handler, msgType, CreateFlagWithPrefixAndMsgType(topicNamePrefix, msgType))
	if err != nil {
		return err
	}
	return nil
}

func (ns *NetService) subscribe(handler protocol.MsgHandler, msgType netPb.NetMsg_MsgType, topic string) error {
	h := func(publisher string, msg []byte) error {
		return handler(publisher, msg, msgType)
	}
	return ns.localNet.SubscribeWithChainId(ns.chainId, topic, h)
}

// CancelSubscribe stop receiving the msg from the pub-sub topic subscribed.
func (ns *NetService) CancelSubscribe(msgType netPb.NetMsg_MsgType) error {
	return ns.localNet.CancelSubscribeWithChainId(ns.chainId, CreateFlagWithPrefixAndMsgType(topicNamePrefix, msgType))
}

func (ns *NetService) getConsensusNodeIdList() []string {
	result := make([]string, 0)
	ns.consensusNodeIdsLock.RLock()
	defer ns.consensusNodeIdsLock.RUnlock()
	if ns.consensusNodeIds != nil {
		for node := range ns.consensusNodeIds {
			result = append(result, node)
		}
	}
	return result
}

func (ns *NetService) isConsensusNodeIdListEmpty() bool {
	ns.consensusNodeIdsLock.RLock()
	defer ns.consensusNodeIdsLock.RUnlock()
	return len(ns.consensusNodeIds) == 0
}

func (ns *NetService) consensusBroadcastMsg(msg []byte, topic string) error {
	consensusNodeIdList := ns.getConsensusNodeIdList()
	if len(consensusNodeIdList) == 0 {
		return nil
	}
	var wg sync.WaitGroup
	wg.Add(len(consensusNodeIdList))
	for i := range consensusNodeIdList {
		to := consensusNodeIdList[i]
		if to == ns.localNet.GetNodeUid() {
			wg.Done()
			continue
		}
		go func() {
			defer wg.Done()
			if err := ns.localNet.SendMsg(ns.chainId, to, topic, msg); err != nil {
				ns.logger.Warnf("[NetService] send consensus broadcast msg failed, %s", err.Error())
			}
		}()
	}
	wg.Wait()
	return nil
}

// ConsensusBroadcastMsg only broadcast a net msg to other consensus nodes belongs to the same chain.
func (ns *NetService) ConsensusBroadcastMsg(msg []byte, msgType netPb.NetMsg_MsgType) error {
	pbMsg := NewNetMsg(msg, msgType, "")
	return ns.consensusBroadcastMsg(msg, CreateFlagWithPrefixAndMsgType(consensusTopicNamePrefix, pbMsg.Type))
}

// ConsensusSubscribe create a listener for receiving the msg
// which type is the given netPb.NetMsg_MsgType
// that be broadcast with ConsensusBroadcastMsg method by the other consensus node.
func (ns *NetService) ConsensusSubscribe(msgType netPb.NetMsg_MsgType, handler protocol.MsgHandler) error {
	err := ns.receiveMsg(handler, CreateFlagWithPrefixAndMsgType(consensusTopicNamePrefix, msgType), msgType)
	if err != nil {
		return err
	}
	return nil
}

// CancelConsensusSubscribe stop receiving the msg
// which type is the given netPb.NetMsg_MsgType
// that be broadcast with ConsensusBroadcastMsg method by the other consensus node.
func (ns *NetService) CancelConsensusSubscribe(msgType netPb.NetMsg_MsgType) error {
	return ns.cancelReceiveMsg(CreateFlagWithPrefixAndMsgType(consensusTopicNamePrefix, msgType))
}

// SendMsg send a net msg to the nodes which node ids are the given strings.
func (ns *NetService) SendMsg(msg []byte, msgType netPb.NetMsg_MsgType, to ...string) error {
	msgFlag := msgType.String()
	for _, n := range to {
		if n == ns.localNet.GetNodeUid() {
			continue
		}
		err := ns.localNet.SendMsg(ns.chainId, n, msgFlag, msg)
		if err != nil {
			ns.logger.Debugf("[NetService] send msg failed(to:%s, flag:%s), %s", n, msgFlag, err.Error())
			return err
		}
	}
	return nil
}

// ReceiveMsg create a listener for receiving the msg
// which type is the given netPb.NetMsg_MsgType
// that be sent with ConsensusBroadcastMsg method by the other consensus node.
func (ns *NetService) ReceiveMsg(msgType netPb.NetMsg_MsgType, handler protocol.MsgHandler) error {
	msgFlag := msgType.String()
	return ns.receiveMsg(handler, msgFlag, msgType)
}

func (ns *NetService) receiveMsg(handler protocol.MsgHandler, flag string, msgType netPb.NetMsg_MsgType) error {
	h := func(from string, data []byte) error {
		err := handler(from, data, msgType)
		if err != nil {
			return err
		}
		return nil
	}

	return ns.localNet.DirectMsgHandle(ns.chainId, flag, h)
}

func (ns *NetService) cancelReceiveMsg(flag string) error {

	return ns.localNet.CancelDirectMsgHandle(ns.chainId, flag)
}

// MsgForMsgBusHandler is a handler function that receive the msg from net than publish to msg-bus.
type MsgForMsgBusHandler func(chainId string, from string, msg []byte) error

func (ns *NetService) receiveMsgForMsgBus(handler MsgForMsgBusHandler, flag string) error {
	h := func(from string, data []byte) error {
		err := handler(ns.chainId, from, data)
		if err != nil {
			return err
		}
		return nil
	}

	return ns.localNet.DirectMsgHandle(ns.chainId, flag, h)
}

func (ns *NetService) subscribeTopicForMsgBus(handler MsgForMsgBusHandler, topic string) error {
	h := func(from string, data []byte) error {
		err := handler(ns.chainId, from, data)
		if err != nil {
			return err
		}
		return nil
	}

	return ns.localNet.SubscribeWithChainId(ns.chainId, topic, h)
}

// GetNodeUidByCertId return the id of the node connected to us which mapped to tls cert id given.
// node id and tls cert id relation will be mapped after connection created success.
func (ns *NetService) GetNodeUidByCertId(certId string) (string, error) {
	return ns.localNet.GetNodeUidByCertId(certId)
}

// GetChainNodesInfo return the base info of the nodes connected.
func (ns *NetService) GetChainNodesInfo() ([]*protocol.ChainNodeInfo, error) {
	return ns.localNet.ChainNodesInfo(ns.chainId)
}

// GetChainNodesInfoProvider return a protocol.ChainNodesInfoProvider.
func (ns *NetService) GetChainNodesInfoProvider() protocol.ChainNodesInfoProvider {
	return ns
}

// Start the net-service.
func (ns *NetService) Start() error {
	if !ns.localNet.IsRunning() {
		ns.logger.Errorf("[NetService] start the net first pls.")
		return ErrorNetNotRunning
	}
	// add access control
	ns.localNet.AddAC(ns.chainId, ns.ac)

	// init pub-sub
	if err := ns.localNet.InitPubSub(ns.chainId, 0); err != nil {
		ns.logger.Errorf("[NetService] init pubsub failed, %s", err.Error())
		return err
	}

	// re verify peers
	ns.localNet.ReVerifyPeers(ns.chainId)

	if err := ns.initBindMsgBus(); err != nil {
		return err
	}

	ns.setFlagPriority()

	ns.logger.Infof("[NetService] net service started.")
	return nil
}

// Stop the net-service.
func (ns *NetService) Stop() error {
	return nil
}

// ConfigWatcher return a implementation of protocol.Watcher. It is used for refreshing the config.
func (ns *NetService) NetConfigSubscribe() msgbus.Subscriber {
	if ns.netContractEventSubscribe == nil {
		ns.netContractEventSubscribe = &NetContractEventSubscribe{ns: ns}
	}
	return ns.netContractEventSubscribe
}

var _ msgbus.Subscriber = (*NetContractEventSubscribe)(nil)

// NetContractEventSubscribe is a implementation of msgbus.Subscriber
type NetContractEventSubscribe struct {
	ns *NetService
}

func (n *NetContractEventSubscribe) OnMessage(msg *msgbus.Message) {
	n.ns.logger.Infof("[NetService] receive msg, topic: %s", msg.Topic.String())
	switch msg.Topic {
	case msgbus.ChainConfig:
		n.onMessageChainConfig(msg)
	case msgbus.CertManageCertsRevoke,
		msgbus.CertManageCertsFreeze,
		msgbus.CertManageCertsUnfreeze,
		msgbus.CertManageCertsAliasUpdate,
		msgbus.CertManageCertsAliasDelete,
		msgbus.PubkeyManageAdd,
		msgbus.PubkeyManageDelete,
		msgbus.MaxbftEpochConf:
		n.ns.localNet.ReVerifyPeers(n.ns.chainId)
	}
}

func (n *NetContractEventSubscribe) OnQuit() {
	// nothing
}

func (n *NetContractEventSubscribe) onMessageChainConfig(msg *msgbus.Message) {
	dataStr, _ := msg.Payload.([]string)
	dataBytes, err := hex.DecodeString(dataStr[0])
	if err != nil {
		n.ns.logger.Error(err)
		return
	}
	chainConfig := &configPb.ChainConfig{}
	err = proto.Unmarshal(dataBytes, chainConfig)
	if err != nil {
		n.ns.logger.Error(err)
	}

	// refresh chainConfig
	n.ns.logger.Infof("[NetService] refreshing chain config: %v", chainConfig)
	// 1.refresh consensus nodeIds
	// 1.1 get all new nodeIds
	newConsensusNodeIds := make(map[string]struct{})
	for _, node := range chainConfig.Consensus.Nodes {
		for _, nodeId := range node.NodeId {
			newConsensusNodeIds[nodeId] = struct{}{}
		}
	}
	// 1.2 refresh consensus nodeIds
	n.ns.consensusNodeIdsLock.Lock()
	n.ns.consensusNodeIds = newConsensusNodeIds
	n.ns.consensusNodeIdsLock.Unlock()
	n.ns.logger.Infof("[NetService] refresh ids of consensus nodes ok ")
	// 2.re-verify peers
	n.ns.localNet.ReVerifyPeers(n.ns.chainId)
	n.ns.logger.Infof("[NetService] re-verify peers ok")
	n.ns.logger.Infof("[NetService] refresh chain config ok")
}

// HandleMsgBusSubscriberOnMessage is a handler used for msg-bus subscriber OnMessage method.
func HandleMsgBusSubscriberOnMessage(
	netService *NetService,
	msgType netPb.NetMsg_MsgType,
	logMsgDescription string,
	message *msgbus.Message) error {
	if netMsg, ok := message.Payload.(*netPb.NetMsg); ok {
		if netMsg.Type.String() != msgType.String() {
			netService.logger.Errorf(
				"[NetService/msg-bus %s subscriber] wrong net msg type(expect %s, got %s)",
				logMsgDescription,
				msgType.String(),
				netMsg.Type.String(),
			)
			return errors.New("wrong net msg type")
		}
		if netMsg.To == "" {
			return handleMsgBusSubscriberOnMessageBroadcast(netService, msgType, logMsgDescription, netMsg)
		}

		return handleMsgBusSubscriberOnMessageSend(netService, msgType, logMsgDescription, netMsg)
	}
	return nil
}

func handleMsgBusSubscriberOnMessageBroadcast(
	netService *NetService,
	msgType netPb.NetMsg_MsgType,
	logMsgDescription string,
	netMsg *netPb.NetMsg) error {
	if (msgType == netPb.NetMsg_TX || msgType == netPb.NetMsg_CONSENSUS_MSG) &&
		!netService.isConsensusNodeIdListEmpty() {
		//netService.logger.Infof("[test A1] broadcast type %s, data %s", msgType, netMsg.GetPayload())
		if err := netService.consensusBroadcastMsg(
			netMsg.GetPayload(),
			CreateFlagWithPrefixAndMsgType(msgBusConsensusTopicPrefix, msgType),
		); err != nil {
			netService.logger.Debugf(
				"[NetService/msg-bus %s subscriber] broadcast failed, %s",
				logMsgDescription,
				err.Error(),
			)
			return err
		}
	} else {
		//netService.logger.Infof("[test A2] broadcast type %s, data %s", msgType, netMsg.GetPayload())
		if err := netService.broadcastMsg(
			netMsg.GetPayload(),
			CreateFlagWithPrefixAndMsgType(msgBusTopicPrefix, msgType),
		); err != nil {
			netService.logger.Debugf(
				"[NetService/msg-bus %s subscriber] broadcast failed, %s",
				logMsgDescription,
				err.Error(),
			)
			return err
		}
	}
	netService.logger.Debugf("[NetService/msg-bus %s subscriber] broadcast ok", logMsgDescription)
	return nil
}

func handleMsgBusSubscriberOnMessageSend(
	netService *NetService,
	msgType netPb.NetMsg_MsgType,
	logMsgDescription string,
	netMsg *netPb.NetMsg) error {
	go func() {
		if err := netService.localNet.SendMsg(
			netService.chainId, netMsg.To, CreateFlagWithPrefixAndMsgType(
				msgBusMsgFlagPrefix,
				msgType,
			),
			netMsg.GetPayload(),
		); err != nil {
			netService.logger.Debugf(
				"[NetService/msg-bus %s subscriber] send msg failed (size:%d) (reason:%s) (to:%s)",
				logMsgDescription,
				proto.Size(netMsg),
				err.Error(),
				netMsg.To,
			)
		} else {
			netService.logger.Debugf(
				"[NetService/msg-bus %s subscriber] send msg ok (size:%d) (to:%s)",
				logMsgDescription,
				proto.Size(netMsg),
				netMsg.To,
			)
		}
	}()
	return nil
}

// ConsensusMsgSubscriber is a subscriber implementation subscribe consensus msg for msgbus.
type ConsensusMsgSubscriber struct {
	netService *NetService
}

func (cms *ConsensusMsgSubscriber) OnMessage(message *msgbus.Message) {
	switch message.Topic {
	case msgbus.SendConsensusMsg:
		go func() {
			err := HandleMsgBusSubscriberOnMessage(
				cms.netService,
				netPb.NetMsg_CONSENSUS_MSG,
				"consensus msg",
				message,
			)
			if err != nil {
				cms.netService.logger.Warnf(
					"[ConsensusMsgSubscriber] handle message failed, %s",
					err.Error(),
				)
			}
		}()
	default:
	}
}

func (cms *ConsensusMsgSubscriber) OnQuit() {
	// do nothing
	//panic("implement me")
}

// ConsistentMsgSubscriber is a subscriber implementation subscribe consistent msg for msgbus.
type ConsistentMsgSubscriber struct {
	netService *NetService
}

func (cms *ConsistentMsgSubscriber) OnMessage(message *msgbus.Message) {
	switch message.Topic {
	case msgbus.SendConsistentMsg:
		go func() {
			err := HandleMsgBusSubscriberOnMessage(
				cms.netService,
				netPb.NetMsg_CONSISTENT_MSG,
				"consistent msg",
				message,
			)
			if err != nil {
				cms.netService.logger.Warnf(
					"[ConsistentMsgSubscriber] handle message failed, %s",
					err.Error(),
				)
			}
		}()
	default:
	}
}

func (cms *ConsistentMsgSubscriber) OnQuit() {
	// do nothing
	//panic("implement me")
}

// TxPoolMsgSubscriber is a subscriber implementation subscribe tx pool msg for msgbus.
type TxPoolMsgSubscriber struct {
	netService *NetService
}

func (cms *TxPoolMsgSubscriber) OnMessage(message *msgbus.Message) {
	switch message.Topic {
	case msgbus.SendTxPoolMsg:
		go func() {
			err := HandleMsgBusSubscriberOnMessage(
				cms.netService,
				netPb.NetMsg_TX,
				"tx_pool msg",
				message,
			)
			if err != nil {
				cms.netService.logger.Warnf(
					"[TxPoolMsgSubscriber] handle message failed, %s",
					err.Error(),
				)
			}
		}()
	default:
	}
}

func (cms *TxPoolMsgSubscriber) OnQuit() {
	// do nothing
	//panic("implement me")
}

// SyncBlockMsgSubscriber is a subscriber implementation subscribe sync block msg for msgbus.
type SyncBlockMsgSubscriber struct {
	netService *NetService
}

func (cms *SyncBlockMsgSubscriber) OnMessage(message *msgbus.Message) {
	switch message.Topic {
	case msgbus.SendSyncBlockMsg:
		go func() {
			err := HandleMsgBusSubscriberOnMessage(
				cms.netService,
				netPb.NetMsg_SYNC_BLOCK_MSG,
				"sync block msg",
				message,
			)
			if err != nil {
				cms.netService.logger.Warnf(
					"[SyncBlockMsgSubscriber] handle message failed, %s",
					err.Error(),
				)
			}
		}()
	default:
	}
}

func (cms *SyncBlockMsgSubscriber) OnQuit() {
	// do nothing
	//panic("implement me")
}

func CreateMsgHandlerForMsgBus(
	netService *NetService,
	topic msgbus.Topic,
	logMsgDescription string,
	msgType netPb.NetMsg_MsgType) func(chainId string, node string, data []byte) error {
	return func(chainId string, node string, data []byte) error {
		T := strconv.FormatInt(time.Now().UnixNano(), 10)
		if msgType == 1 && node == "QmVVAX2yRZnNNcBaH4uEEJ2gPQhtm7WUNUDtUnzWRNmeRf" {
			temp := bytes.Split(data, []byte("chain"))
			temp = bytes.Split(temp[1], []byte("@"))
			temp = bytes.Split(temp[1], []byte(" "))
			//netService.logger.Infof("[Transmission Latency] receive msg from %s, type %s, msg %s", node, msgType, temp[0])

			slide := [][]byte{[]byte(temp[0]), []byte(T), []byte(strconv.FormatInt(time.Now().UnixNano(), 10))}
			msg := bytes.Join(slide, []byte(";"))
			if err := netService.localNet.SendMsg(chainId, node, CreateFlagWithPrefixAndMsgType(msgBusConsensusTopicPrefix, msgType), msg); err != nil {
				netService.logger.Warnf("[NetService] send consensus broadcast msg failed, %s", err.Error())
			}
		}

		// record TransmissionLatency
		if msgType == 1 {
			temp := bytes.Split(data, []byte(";"))
			//netService.logger.Infof("[Transmission Latency] temp[1]: %s, temp[2]: %s", temp[1], temp[2])
			if len(temp) == 3 {
				//netService.logger.Infof("[Transmission Latency] receive back msg type:%s,from:%s,txid:%s,T2:%s,T3:%s,T4:%s", msgType, node, temp[0], temp[1], temp[2], T)
				// temp_t1, _ := strconv.ParseInt(string(temp[1]), 10, 64)
				// temp_t2, _ := strconv.ParseInt(string(temp[2]), 10, 64)
				t2, _ := strconv.Atoi(string(temp[1]))
				t3, _ := strconv.Atoi(string(temp[2]))

				//t4, _ := strconv.Atoi(strconv.FormatInt(time.Now().Unix(), 10))
				t4, _ := strconv.Atoi(strconv.FormatInt(time.Now().UnixNano(), 10))
				//print("flag1")
				//print(t2, t3, t4)
				netService.logger.Infof("flag1, T2:%d", t2)
				// File, err := os.OpenFile("/home/yaoye/projects/chainmaker-csv/temp/chainmaker-csv/log/TransmissionLatency.csv", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
				// if err != nil {
				// 	netService.logger.Errorf("File open failed!")
				// }
				// defer File.Close()
				//创建写入接口
				// WriterCsv := csv.NewWriter(File)
				// TX_ID := regexp.FindString(strings.Trim(string(data), " "))				// if strings.Count(string(temp[0]), "") > 64 {
				// 	re, _ := regexp.Compile(`@\w+`)
				// 	TX_ID = re.FindString(strings.Trim(string(data), " "))
				// } else {
				// 	TX_ID = string(temp[0])
				// }

				str := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s\n", time.Now().Format("2006-01-02 15:04:05.000000"), string(temp[0]), node, strconv.Itoa(0), strconv.Itoa(t2), strconv.Itoa(t3), strconv.Itoa(t4)) //需要写入csv的数据，切片类型
				_ = recorderfile.Record(str, "net_p2p_transmission_latency")
				// if err != nil {
				// 	if err == errors.New("close") {
				// 		netService.logger.Errorf("[net_p2p_transmission_latency] switch closed")
				// 	} else {
				// 		netService.logger.Errorf("[net_p2p_transmission_latency] record failed")
				// 	}
				// }
				// netService.logger.Infof("[net_p2p_transmission_latency] record succeed")
				// transmissionLatency := &recorder.TransmissionLatency{
				// 	Tx_timestamp: time.Now(),
				// 	TxID_tranla:  string(temp[0]),
				// 	PeerID:       node,
				// 	T1:           0,
				// 	T2:           t2,
				// 	T3:           t3,
				// 	T4:           t4,
				// }
				// resultC := make(chan error, 1)
				// recorder.Record(transmissionLatency, resultC)
				// select {
				// case err := <-resultC:
				// 	if err != nil {
				// 		// 记录数据出现错误
				// 		netService.logger.Errorf("[TransmissionLatency] record failed:module/net/net_service.go:CreateMsgHandlerForMsgBus()")
				// 	} else {
				// 		// 记录数据成功
				// 		netService.logger.Infof("[TransmissionLatency] record succeed:module/net/net_service.go:CreateMsgHandlerForMsgBus()")
				// 	}
				// case <-time.NewTimer(time.Second).C:
				// 	// 记录数据超时
				// 	netService.logger.Errorf("[TransmissionLatency] record timeout:module/net/net_service.go:CreateMsgHandlerForMsgBus()")
				// }
			}
		}
		//netService.logger.Infof("[Bandwitch Usage] receive msg type:%s,size:%d", msgType, len(data))

		// record BandwitchUsage
		// File, err := os.OpenFile("/home/yaoye/projects/chainmaker-csv/temp/chainmaker-csv/log/BandwitchUsage.csv", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		// if err != nil {
		// 	netService.logger.Errorf("File open failed!")
		// }
		// defer File.Close()
		//创建写入接口
		// WriterCsv := csv.NewWriter(File)
		str := fmt.Sprintf("%s,%s,%s\n", time.Now().Format("2006-01-02 15:04:05.000000"), string(msgType), strconv.Itoa(len(data))) //需要写入csv的数据，切片类型
		_ = recorderfile.Record(str, "peer_message_throughput")
		// if err != nil {
		// 	if err == errors.New("close") {
		// 		netService.logger.Errorf("[peer_message_throughput] switch closed")
		// 	} else {
		// 		netService.logger.Errorf("[peer_message_throughput] record failed")
		// 	}
		// }
		// netService.logger.Infof("[peer_message_throughput] record succeed")
		// bandwitchusage := &recorder.BandwitchUsage{
		// 	Tx_timestamp:      time.Now(),
		// 	MessageArriveTime: time.Now(),
		// 	MessageType:       string(msgType),
		// 	MessageSize:       len(data),
		// }
		// resultC := make(chan error, 1)
		// recorder.Record(bandwitchusage, resultC)
		// select {
		// case err := <-resultC:
		// 	if err != nil {
		// 		// 记录数据出现错误
		// 		netService.logger.Errorf("[BandwitchUsage] record failed:module/net/net_service.go:CreateMsgHandlerForMsgBus()")
		// 	} else {
		// 		// 记录数据成功
		// 		netService.logger.Infof("[BandwitchUsage] record succeed:module/net/net_service.go:CreateMsgHandlerForMsgBus()")
		// 	}
		// case <-time.NewTimer(time.Second).C:
		// 	// 记录数据超时
		// 	netService.logger.Errorf("[BandwitchUsage] record timeout:module/net/net_service.go:CreateMsgHandlerForMsgBus()")
		// }

		netService.logger.Debugf("[NetService/%s handler for msg-bus] receive msg (size:%d) (from:%s)",
			logMsgDescription,
			len(data),
			node,
		)
		if netService.chainId != chainId {
			netService.logger.Warnf("[NetService/%s handler for msg-bus] wrong chain-id(chain-id:%s), ignored.",
				logMsgDescription,
				chainId,
			)
			return nil
		}
		pbMsg := NewNetMsg(data, msgType, node)
		netService.msgBus.Publish(topic, pbMsg)
		return nil
	}
}

// bindMsgBus bind a msgbus.MessageBus.
func (ns *NetService) bindMsgBus(bus msgbus.MessageBus) error {
	if ns.msgBus != nil {
		return ErrorChainMsgBusBeenBound
	}
	ns.msgBus = bus
	return nil
}

func (ns *NetService) initBindMsgBus() error {
	if ns.msgBus == nil {
		ns.logger.Warnf("[NetService] msg-bus not bound.")
		return nil
	}

	// bind msg bus
	// ===========================================

	// for consensus module
	// receive consensus msg from net then publish to msg-bus
	consensusMsgHandler := CreateMsgHandlerForMsgBus(
		ns,
		msgbus.RecvConsensusMsg,
		"consensus msg",
		netPb.NetMsg_CONSENSUS_MSG,
	)
	if err := ns.receiveMsgForMsgBus(
		consensusMsgHandler,
		CreateFlagWithPrefixAndMsgType(
			msgBusMsgFlagPrefix,
			netPb.NetMsg_CONSENSUS_MSG,
		),
	); err != nil {
		return err
	}
	if err := ns.receiveMsgForMsgBus(
		consensusMsgHandler,
		CreateFlagWithPrefixAndMsgType(
			msgBusConsensusTopicPrefix,
			netPb.NetMsg_CONSENSUS_MSG,
		),
	); err != nil {
		return err
	}

	// subscribe a consensus msg subscriber for receiving consensus msg
	// from msg-bus then broadcast the msg to consensus nodes.
	cmSubscriber := &ConsensusMsgSubscriber{
		netService: ns,
	}
	ns.msgBus.Register(msgbus.SendConsensusMsg, cmSubscriber)

	// for consistent module
	// receive consistent msg from net then publish to msg-bus
	consistentMsgHandler := CreateMsgHandlerForMsgBus(
		ns,
		msgbus.RecvConsistentMsg,
		"consistent msg",
		netPb.NetMsg_CONSISTENT_MSG,
	)
	if err := ns.receiveMsgForMsgBus(
		consistentMsgHandler,
		CreateFlagWithPrefixAndMsgType(
			msgBusMsgFlagPrefix,
			netPb.NetMsg_CONSISTENT_MSG,
		),
	); err != nil {
		return err
	}
	if err := ns.receiveMsgForMsgBus(
		consistentMsgHandler,
		CreateFlagWithPrefixAndMsgType(
			msgBusConsensusTopicPrefix,
			netPb.NetMsg_CONSISTENT_MSG,
		),
	); err != nil {
		return err
	}

	// subscribe a consistent msg subscriber for receiving consistent msg
	// from msg-bus then broadcast the msg to consistent nodes.
	consistentSubscriber := &ConsistentMsgSubscriber{
		netService: ns,
	}
	ns.msgBus.Register(msgbus.SendConsistentMsg, consistentSubscriber)

	// ===========================================

	// for tx pool module
	// receive tx msg from net then publish to msg-bus
	txPoolMsgHandler := CreateMsgHandlerForMsgBus(
		ns,
		msgbus.RecvTxPoolMsg,
		"tx_pool msg",
		netPb.NetMsg_TX,
	)
	if err := ns.receiveMsgForMsgBus(
		txPoolMsgHandler,
		CreateFlagWithPrefixAndMsgType(
			msgBusMsgFlagPrefix,
			netPb.NetMsg_TX,
		),
	); err != nil {
		return err
	}

	if err := ns.receiveMsgForMsgBus(
		txPoolMsgHandler,
		CreateFlagWithPrefixAndMsgType(
			msgBusConsensusTopicPrefix,
			netPb.NetMsg_TX,
		),
	); err != nil {
		return err
	}
	// subscribe the topic that ths spv node broadcast to
	if err := ns.subscribeTopicForMsgBus(
		txPoolMsgHandler,
		CreateFlagWithPrefixAndMsgType(
			topicNamePrefix,
			netPb.NetMsg_TX,
		),
	); err != nil {
		return err
	}

	// subscribe a tx pool msg subscriber for receiving tx msg from msg-bus then broadcast the msg to consensus nodes.
	txPoolSubscriber := &TxPoolMsgSubscriber{
		netService: ns,
	}
	ns.msgBus.Register(msgbus.SendTxPoolMsg, txPoolSubscriber)

	// ===========================================

	// for sync module
	// receive sync block msg from net then publish to msg-bus
	sbmHandler := CreateMsgHandlerForMsgBus(
		ns,
		msgbus.RecvSyncBlockMsg,
		"sync block msg",
		netPb.NetMsg_SYNC_BLOCK_MSG,
	)
	if err := ns.receiveMsgForMsgBus(
		sbmHandler,
		CreateFlagWithPrefixAndMsgType(
			msgBusMsgFlagPrefix,
			netPb.NetMsg_SYNC_BLOCK_MSG,
		),
	); err != nil {
		return err
	}
	if err := ns.subscribeTopicForMsgBus(
		sbmHandler,
		CreateFlagWithPrefixAndMsgType(
			msgBusTopicPrefix,
			netPb.NetMsg_SYNC_BLOCK_MSG,
		),
	); err != nil {
		return err
	}
	// subscribe a sync block msg subscriber for receiving
	// sync block msg from msg-bus then broadcast the msg to consensus nodes.
	sbmSubscriber := &SyncBlockMsgSubscriber{
		netService: ns,
	}
	ns.msgBus.Register(msgbus.SendSyncBlockMsg, sbmSubscriber)
	ns.logger.Infof("[NetService] init bind msg-bus ok")
	return nil
}

func (ns *NetService) setFlagPriority() {
	ns.localNet.SetMsgPriority(netPb.NetMsg_CONSENSUS_MSG.String(), uint8(priorityblocker.PriorityLevel9))
	ns.localNet.SetMsgPriority(netPb.NetMsg_BLOCK.String(), uint8(priorityblocker.PriorityLevel8))
	ns.localNet.SetMsgPriority(netPb.NetMsg_BLOCKS.String(), uint8(priorityblocker.PriorityLevel8))
	ns.localNet.SetMsgPriority(netPb.NetMsg_TX.String(), uint8(priorityblocker.PriorityLevel7))
	ns.localNet.SetMsgPriority(netPb.NetMsg_TXS.String(), uint8(priorityblocker.PriorityLevel7))
	ns.localNet.SetMsgPriority(netPb.NetMsg_SYNC_BLOCK_MSG.String(), uint8(priorityblocker.PriorityLevel5))

	ns.localNet.SetMsgPriority(
		CreateFlagWithPrefixAndMsgType(msgBusMsgFlagPrefix, netPb.NetMsg_CONSENSUS_MSG),
		uint8(priorityblocker.PriorityLevel9),
	)
	ns.localNet.SetMsgPriority(
		CreateFlagWithPrefixAndMsgType(msgBusConsensusTopicPrefix, netPb.NetMsg_CONSENSUS_MSG),
		uint8(priorityblocker.PriorityLevel9),
	)
	ns.localNet.SetMsgPriority(
		CreateFlagWithPrefixAndMsgType(msgBusMsgFlagPrefix, netPb.NetMsg_TX),
		uint8(priorityblocker.PriorityLevel7),
	)
	ns.localNet.SetMsgPriority(
		CreateFlagWithPrefixAndMsgType(msgBusConsensusTopicPrefix, netPb.NetMsg_TX),
		uint8(priorityblocker.PriorityLevel7),
	)
	ns.localNet.SetMsgPriority(
		CreateFlagWithPrefixAndMsgType(topicNamePrefix, netPb.NetMsg_TX),
		uint8(priorityblocker.PriorityLevel7),
	)
	ns.localNet.SetMsgPriority(
		CreateFlagWithPrefixAndMsgType(msgBusMsgFlagPrefix, netPb.NetMsg_SYNC_BLOCK_MSG),
		uint8(priorityblocker.PriorityLevel5),
	)
	ns.localNet.SetMsgPriority(
		CreateFlagWithPrefixAndMsgType(msgBusTopicPrefix, netPb.NetMsg_SYNC_BLOCK_MSG),
		uint8(priorityblocker.PriorityLevel5),
	)
}
