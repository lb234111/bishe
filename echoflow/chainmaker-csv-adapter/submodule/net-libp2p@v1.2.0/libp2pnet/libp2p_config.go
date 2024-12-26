/*
Copyright (C) BABEC. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package libp2pnet

var Config *netConfig

type netConfig struct {
	Provider                string            `mapstructure:"provider"`
	ListenAddr              string            `mapstructure:"listen_addr"`
	PeerStreamPoolSize      int               `mapstructure:"peer_stream_pool_size"`
	MaxPeerCountAllow       int               `mapstructure:"max_peer_count_allow"`
	PeerEliminationStrategy int               `mapstructure:"peer_elimination_strategy"`
	Seeds                   []string          `mapstructure:"seeds"`
	TLSConfig               netTlsConfig      `mapstructure:"tls"`
	BlackList               blackList         `mapstructure:"blacklist"`
	CustomChainTrustRoots   []chainTrustRoots `mapstructure:"custom_chain_trust_roots"`
}

type netTlsConfig struct {
	Enabled     bool   `mapstructure:"enabled"`
	PrivKeyFile string `mapstructure:"priv_key_file"`
	CertFile    string `mapstructure:"cert_file"`
}

type blackList struct {
	Addresses []string `mapstructure:"addresses"`
	NodeIds   []string `mapstructure:"node_ids"`
}

type chainTrustRoots struct {
	ChainId    string       `mapstructure:"chain_id"`
	TrustRoots []trustRoots `mapstructure:"trust_roots"`
}

type trustRoots struct {
	OrgId string `mapstructure:"org_id"`
	Root  string `mapstructure:"root"`
}

/*
// LibP2pNetPrepare prepare the config options.
type LibP2pNetPrepare struct {
	listenAddr              string              // listenAddr
	bootstrapsPeers         map[string]struct{} // bootstrapsPeers
	pubSubMaxMsgSize        int                 // pubSubMaxMsgSize
	peerStreamPoolSize      int                 // peerStreamPoolSize
	maxPeerCountAllow       int                 // maxPeerCountAllow
	peerEliminationStrategy int                 // peerEliminationStrategy

	keyBytes      []byte // keyBytes
	certBytes     []byte // certBytes
	signKeyBytes  []byte // signKeyBytes used when use dual certificate mode
	signCertBytes []byte // signCertBytes used when use dual certificate mode

	chainTrustRootCertsBytes map[string][][]byte // chainTrustRootCertsBytes

	blackAddresses map[string]struct{} // blackAddresses
	blackPeerIds   map[string]struct{} // blackPeerIds

	isInsecurity       bool
	pktEnable          bool
	priorityCtrlEnable bool

	lock sync.Mutex
}
*/
