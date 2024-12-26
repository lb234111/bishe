package config

// User用户结构定义了用户基础信息
type User struct {
	TlsKeyPath, TlsCrtPath   string
	SignKeyPath, SignCrtPath string
}

// 所有链1上的管理员用户信息
// var Chain1Admins = []string{Chain1Org1Admin1, Chain1Org2Admin1, Chain1Org3Admin1, Chain1Org4Admin1, Chain1Org5Admin1, Chain1Org6Admin1, Chain1Org7Admin1, Chain1Org8Admin1, Chain1Org9Admin1, Chain1Org10Admin1}
var Chain1Admins = []string{Chain1Org1Admin1, Chain1Org2Admin1, Chain1Org3Admin1, Chain1Org4Admin1}

var OrgNum int = 4
var ClientNum int = 10

const (
	//组织管理员信息
	Chain1Org1Admin1  = "chain1org1admin1"  //链1组织1的管理员名称
	Chain1Org2Admin1  = "chain1org2admin1"  //链1组织2的管理员名称
	Chain1Org3Admin1  = "chain1org3admin1"  //链1组织3的管理员名称
	Chain1Org4Admin1  = "chain1org4admin1"  //链1组织4的管理员名称
	Chain1Org5Admin1  = "chain1org5admin1"  //链1组织5的管理员名称
	Chain1Org6Admin1  = "chain1org6admin1"  //链1组织6的管理员名称
	Chain1Org7Admin1  = "chain1org7admin1"  //链1组织7的管理员名称
	Chain1Org8Admin1  = "chain1org8admin1"  //链1组织8的管理员名称
	Chain1Org9Admin1  = "chain1org9admin1"  //链1组织9的管理员名称
	Chain1Org10Admin1 = "chain1org10admin1" //链1组织10的管理员名称
)
const (
	ChainmakerConfigPath = "./worker/chainmaker/config_files"
)

type LogConfig struct {
	LogInConsole bool   `mapstructure:"log_in_console"`
	ShowColor    bool   `mapstructure:"show_color"`
	LogLevel     string `mapstructure:"log_level"`
	LogPath      string `mapstructure:"log_path"`
}

var Admins = map[string]*User{
	Chain1Org1Admin1: {
		ChainmakerConfigPath + "/chain1/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.tls.key",
		ChainmakerConfigPath + "/chain1/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.tls.crt",
		ChainmakerConfigPath + "/chain1/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key",
		ChainmakerConfigPath + "/chain1/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt",
	},
	Chain1Org2Admin1: {
		ChainmakerConfigPath + "/chain1/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.tls.key",
		ChainmakerConfigPath + "/chain1/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.tls.crt",
		ChainmakerConfigPath + "/chain1/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key",
		ChainmakerConfigPath + "/chain1/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt",
	},
	Chain1Org3Admin1: {
		ChainmakerConfigPath + "/chain1/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.tls.key",
		ChainmakerConfigPath + "/chain1/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.tls.crt",
		ChainmakerConfigPath + "/chain1/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key",
		ChainmakerConfigPath + "/chain1/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt",
	},
	Chain1Org4Admin1: {
		ChainmakerConfigPath + "/chain1/crypto-config/wx-org4.chainmaker.org/user/admin1/admin1.tls.key",
		ChainmakerConfigPath + "/chain1/crypto-config/wx-org4.chainmaker.org/user/admin1/admin1.tls.crt",
		ChainmakerConfigPath + "/chain1/crypto-config/wx-org4.chainmaker.org/user/admin1/admin1.sign.key",
		ChainmakerConfigPath + "/chain1/crypto-config/wx-org4.chainmaker.org/user/admin1/admin1.sign.crt",
	},
	Chain1Org5Admin1: {
		ChainmakerConfigPath + "/chain1/crypto-config/wx-org5.chainmaker.org/user/admin1/admin1.tls.key",
		ChainmakerConfigPath + "/chain1/crypto-config/wx-org5.chainmaker.org/user/admin1/admin1.tls.crt",
		ChainmakerConfigPath + "/chain1/crypto-config/wx-org5.chainmaker.org/user/admin1/admin1.sign.key",
		ChainmakerConfigPath + "/chain1/crypto-config/wx-org5.chainmaker.org/user/admin1/admin1.sign.crt",
	},
	Chain1Org6Admin1: {
		ChainmakerConfigPath + "/chain1/crypto-config/wx-org6.chainmaker.org/user/admin1/admin1.tls.key",
		ChainmakerConfigPath + "/chain1/crypto-config/wx-org6.chainmaker.org/user/admin1/admin1.tls.crt",
		ChainmakerConfigPath + "/chain1/crypto-config/wx-org6.chainmaker.org/user/admin1/admin1.sign.key",
		ChainmakerConfigPath + "/chain1/crypto-config/wx-org6.chainmaker.org/user/admin1/admin1.sign.crt",
	},
	Chain1Org7Admin1: {
		ChainmakerConfigPath + "/chain1/crypto-config/wx-org7.chainmaker.org/user/admin1/admin1.tls.key",
		ChainmakerConfigPath + "/chain1/crypto-config/wx-org7.chainmaker.org/user/admin1/admin1.tls.crt",
		ChainmakerConfigPath + "/chain1/crypto-config/wx-org7.chainmaker.org/user/admin1/admin1.sign.key",
		ChainmakerConfigPath + "/chain1/crypto-config/wx-org7.chainmaker.org/user/admin1/admin1.sign.crt",
	},
	Chain1Org8Admin1: {
		ChainmakerConfigPath + "/chain1/crypto-config/wx-org8.chainmaker.org/user/admin1/admin1.tls.key",
		ChainmakerConfigPath + "/chain1/crypto-config/wx-org8.chainmaker.org/user/admin1/admin1.tls.crt",
		ChainmakerConfigPath + "/chain1/crypto-config/wx-org8.chainmaker.org/user/admin1/admin1.sign.key",
		ChainmakerConfigPath + "/chain1/crypto-config/wx-org8.chainmaker.org/user/admin1/admin1.sign.crt",
	},
	Chain1Org9Admin1: {
		ChainmakerConfigPath + "/chain1/crypto-config/wx-org9.chainmaker.org/user/admin1/admin1.tls.key",
		ChainmakerConfigPath + "/chain1/crypto-config/wx-org9.chainmaker.org/user/admin1/admin1.tls.crt",
		ChainmakerConfigPath + "/chain1/crypto-config/wx-org9.chainmaker.org/user/admin1/admin1.sign.key",
		ChainmakerConfigPath + "/chain1/crypto-config/wx-org9.chainmaker.org/user/admin1/admin1.sign.crt",
	},
	Chain1Org10Admin1: {
		ChainmakerConfigPath + "/chain1/crypto-config/wx-org10.chainmaker.org/user/admin1/admin1.tls.key",
		ChainmakerConfigPath + "/chain1/crypto-config/wx-org10.chainmaker.org/user/admin1/admin1.tls.crt",
		ChainmakerConfigPath + "/chain1/crypto-config/wx-org10.chainmaker.org/user/admin1/admin1.sign.key",
		ChainmakerConfigPath + "/chain1/crypto-config/wx-org10.chainmaker.org/user/admin1/admin1.sign.crt",
	},
}
