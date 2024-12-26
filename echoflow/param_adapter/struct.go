package main

var (
	// TBFT_propose_timeout_key
	TBFT_propose_timeout_key = "TBFT_propose_timeout"
	// TBFT_propose_delta_timeout_key
	TBFT_propose_delta_timeout_key = "TBFT_propose_delta_timeout"
	// TBFT_blocks_per_proposer
	TBFT_blocks_per_proposer = "TBFT_blocks_per_proposer"
	// DPOS
	DPOS = "DPOS"
	// MAXBFT_view_nums_per_epoch
	MAXBFT_view_nums_per_epoch = "MaxBftViewNumsPerEpoch"
	// network
	max_peer_count_allow = "MaxPeerCountAllow"
	block_tx_capacity    = "block_tx_capacity"

	batch_max_size       = "BatchMaxSize"
	batch_create_timeout = "BatchCreateTimeout"

	connectors              = "Connectors"
	heartbeatInterval       = "HeartbeatInterval"
	gossipRetransmission    = "GossipRetransmission"
	opportunisticGraftPeers = "OpportunisticGraftPeers"
	dout                    = "Dout"
	dlazy                   = "Dlazy"
	gossipFactor            = "GossipFactor"
	opportunisticGraftTicks = "OpportunisticGraftTicks"
	floodPublish            = "FloodPublish"
)

// ParamAdapter 请求参数结构
type ParamAdapter struct {
	Module     int32  `json:"module"`      // 需要调整的模块 1:共识层
	Type       int32  `json:"type"`        // 调整类别 例:共识层 1:tbft 2:raft 3:dpos 4:hotstuff(maxbft)
	ParamName  string `json:"param_name"`  // 调整的参数名字
	ParamValue int32  `json:"param_value"` // 调整后的数值
}

const (
	OrgId1  = "wx-org1.chainmaker.org"
	OrgId2  = "wx-org2.chainmaker.org"
	OrgId3  = "wx-org3.chainmaker.org"
	OrgId4  = "wx-org4.chainmaker.org"
	OrgId5  = "wx-org5.chainmaker.org"
	OrgId6  = "wx-org6.chainmaker.org"
	OrgId7  = "wx-org7.chainmaker.org"
	OrgId8  = "wx-org8.chainmaker.org"
	OrgId9  = "wx-org9.chainmaker.org"
	OrgId10 = "wx-org10.chainmaker.org"

	UserNameOrg1Admin1  = "org1admin1"
	UserNameOrg2Admin1  = "org2admin1"
	UserNameOrg3Admin1  = "org3admin1"
	UserNameOrg4Admin1  = "org4admin1"
	UserNameOrg5Admin1  = "org5admin1"
	UserNameOrg6Admin1  = "org6admin1"
	UserNameOrg7Admin1  = "org7admin1"
	UserNameOrg8Admin1  = "org8admin1"
	UserNameOrg9Admin1  = "org9admin1"
	UserNameOrg10Admin1 = "org10admin1"
)

var Chain1Admins = []string{UserNameOrg1Admin1, UserNameOrg2Admin1, UserNameOrg3Admin1, UserNameOrg4Admin1}

// PkUsers ...
type PkUsers struct {
	SignKeyPath string
}

// PermissionedPkUsers ...
type PermissionedPkUsers struct {
	SignKeyPath string
	OrgId       string
}

// User ...
type User struct {
	TlsKeyPath  string
	TlsCrtPath  string
	SignKeyPath string
	SignCrtPath string
}

var users = map[string]*User{
	"org1admin1": {
		"./crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.tls.key",
		"./crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.tls.crt",
		"./crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key",
		"./crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt",
	},
	"org2admin1": {
		"./crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.tls.key",
		"./crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.tls.crt",
		"./crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key",
		"./crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt",
	},
	"org3admin1": {
		"./crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.tls.key",
		"./crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.tls.crt",
		"./crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key",
		"./crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt",
	},
	"org4admin1": {
		"./crypto-config/wx-org4.chainmaker.org/user/admin1/admin1.tls.key",
		"./crypto-config/wx-org4.chainmaker.org/user/admin1/admin1.tls.crt",
		"./crypto-config/wx-org4.chainmaker.org/user/admin1/admin1.sign.key",
		"./crypto-config/wx-org4.chainmaker.org/user/admin1/admin1.sign.crt",
	},
	"org5admin1": {
		"./crypto-config/wx-org5.chainmaker.org/user/admin1/admin1.tls.key",
		"./crypto-config/wx-org5.chainmaker.org/user/admin1/admin1.tls.crt",
		"./crypto-config/wx-org5.chainmaker.org/user/admin1/admin1.sign.key",
		"./crypto-config/wx-org5.chainmaker.org/user/admin1/admin1.sign.crt",
	},
	"org6admin1": {
		"./crypto-config/wx-org6.chainmaker.org/user/admin1/admin1.tls.key",
		"./crypto-config/wx-org6.chainmaker.org/user/admin1/admin1.tls.crt",
		"./crypto-config/wx-org6.chainmaker.org/user/admin1/admin1.sign.key",
		"./crypto-config/wx-org6.chainmaker.org/user/admin1/admin1.sign.crt",
	},
	"org7admin1": {
		"./crypto-config/wx-org7.chainmaker.org/user/admin1/admin1.tls.key",
		"./crypto-config/wx-org7.chainmaker.org/user/admin1/admin1.tls.crt",
		"./crypto-config/wx-org7.chainmaker.org/user/admin1/admin1.sign.key",
		"./crypto-config/wx-org7.chainmaker.org/user/admin1/admin1.sign.crt",
	},
	"org8admin1": {
		"./crypto-config/wx-org8.chainmaker.org/user/admin1/admin1.tls.key",
		"./crypto-config/wx-org8.chainmaker.org/user/admin1/admin1.tls.crt",
		"./crypto-config/wx-org8.chainmaker.org/user/admin1/admin1.sign.key",
		"./crypto-config/wx-org8.chainmaker.org/user/admin1/admin1.sign.crt",
	},
	"org9admin1": {
		"./crypto-config/wx-org9.chainmaker.org/user/admin1/admin1.tls.key",
		"./crypto-config/wx-org9.chainmaker.org/user/admin1/admin1.tls.crt",
		"./crypto-config/wx-org9.chainmaker.org/user/admin1/admin1.sign.key",
		"./crypto-config/wx-org9.chainmaker.org/user/admin1/admin1.sign.crt",
	},
	"org10admin1": {
		"./crypto-config/wx-org10.chainmaker.org/user/admin1/admin1.tls.key",
		"./crypto-config/wx-org10.chainmaker.org/user/admin1/admin1.tls.crt",
		"./crypto-config/wx-org10.chainmaker.org/user/admin1/admin1.sign.key",
		"./crypto-config/wx-org10.chainmaker.org/user/admin1/admin1.sign.crt",
	},
}

var permissionedPkUsers = map[string]*PermissionedPkUsers{
	"org1admin1": {
		"../../config-pk/permissioned-with-key/wx-org1/user/admin1/admin1.key",
		OrgId1,
	},
	"org2admin1": {
		"../../config-pk/permissioned-with-key/wx-org2/user/admin1/admin1.key",
		OrgId2,
	},
	"org3admin1": {
		"../../config-pk/permissioned-with-key/wx-org3/user/admin1/admin1.key",
		OrgId3,
	},
	"org4admin1": {
		"../../config-pk/permissioned-with-key/wx-org4/user/admin1/admin1.key",
		OrgId4,
	},
	// "org5admin1": {
	// 	"../../config-pk/permissioned-with-key/wx-org5/user/admin1/admin1.key",
	// 	OrgId5,
	// },
	// "org6admin1": {
	// 	"../../config-pk/permissioned-with-key/wx-org6/user/admin1/admin1.key",
	// 	OrgId6,
	// },
	// "org7admin1": {
	// 	"../../config-pk/permissioned-with-key/wx-org7/user/admin1/admin1.key",
	// 	OrgId7,
	// },
	// "org8admin1": {
	// 	"../../config-pk/permissioned-with-key/wx-org8/user/admin1/admin1.key",
	// 	OrgId8,
	// },
	// "org9admin1": {
	// 	"../../config-pk/permissioned-with-key/wx-org9/user/admin1/admin1.key",
	// 	OrgId9,
	// },
	// "org10admin1": {
	// 	"../../config-pk/permissioned-with-key/wx-org10/user/admin1/admin1.key",
	// 	OrgId10,
	// },
}

var pkUsers = map[string]*PkUsers{
	"org1admin1": {
		"./crypto-config/node1/admin/admin1/admin1.key",
	},
	"org2admin1": {
		"./crypto-config/node2/admin/admin1/admin1.key",
	},
	"org3admin1": {
		"./crypto-config/node3/admin/admin1/admin1.key",
	},
	"org4admin1": {
		"./crypto-config/node4/admin/admin1/admin1.key",
	},
	// "org5admin1": {
	// 	"./crypto-config/node5/admin/admin1/admin1.key",
	// },
	// "org6admin1": {
	// 	"./crypto-config/node6/admin/admin1/admin1.key",
	// },
	// "org7admin1": {
	// 	"./crypto-config/node7/admin/admin1/admin1.key",
	// },
	// "org8admin1": {
	// 	"./crypto-config/node8/admin/admin1/admin1.key",
	// },
	// "org9admin1": {
	// 	"./crypto-config/node9/admin/admin1/admin1.key",
	// },
	// "org10admin1": {
	// 	"./crypto-config/node10/admin/admin1/admin1.key",
	// },
}
