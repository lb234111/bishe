package libp2pnet

import (
	"sync/atomic"

	"chainmaker.org/chainmaker/protocol/v2"
	"github.com/libp2p/go-libp2p-core/network"
)

// singleStreamPool
type singleStreamPool struct {
	// newStreamFunc  a stream
	newStreamFunc func() (network.Stream, error)
	// stream
	stream network.Stream
	// cap
	cap int32
	// idleSize
	idleSize int32
	// log
	log protocol.Logger
}

// Close stream
func (s *singleStreamPool) cleanAndDisable() error {
	if s.stream == nil {
		return nil
	}
	return s.stream.Reset()
}

// InitStreams create a stream
func (s *singleStreamPool) InitStreams() error {
	stream, err := s.newStreamFunc()
	s.stream = newSendStream(stream, s.log)
	return err
}

// BorrowStream return the stream
func (s *singleStreamPool) borrowStream() (network.Stream, error) {
	atomic.AddInt32(&s.idleSize, -1)
	return s.stream, nil
}

// ReturnStream no need do anything
func (s *singleStreamPool) returnStream(_ network.Stream) error {
	atomic.AddInt32(&s.idleSize, 1)
	return nil
}

// DropStream
// single stream pool, do nothing here
func (s *singleStreamPool) dropStream(_ network.Stream) {
	if s.stream != nil && s.stream.Conn() != nil {
		s.log.Warnf("DropStream %v,%v", s.stream.Conn().ID(), s.stream.Conn().RemoteMultiaddr().String())
	}
}

// nolint: unused
func (s *singleStreamPool) addStream(_ network.Stream) {

}

// NewSingleStreamPool create a new mgr.StreamPool .
// SingleStreamPool have only one stream
// SingleStreamPool is meat to remove stream-pool and dont change any other code
// close stream will not influence conn but only stream
func NewSingleStreamPool(cap int, newStreamFunc func() (network.Stream, error), log protocol.Logger) *singleStreamPool {
	log.Debug("NewSingleStreamPool")
	return &singleStreamPool{
		idleSize:      0,
		cap:           int32(cap),
		newStreamFunc: newStreamFunc,
		log:           log,
	}
}
