package common

//Master和server之间通信的格式
type Message struct {
	Config Config  `json:"configs"`
	NodeId int     `json:"nodeId"`
	TPS    float64 `json:"tps"`
}
