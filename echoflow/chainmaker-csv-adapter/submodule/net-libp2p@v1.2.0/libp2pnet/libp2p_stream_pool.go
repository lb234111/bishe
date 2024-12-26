/*
Copyright (C) BABEC. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package libp2pnet

import (
	"errors"
	"math"
	"sync"
	"sync/atomic"

	"chainmaker.org/chainmaker/protocol/v2"
	"github.com/libp2p/go-libp2p-core/network"
)

const (
	// DefaultStreamPoolCap is the default cap of stream pool.
	DefaultStreamPoolCap int = 10000

	defaultIdlePercentage float64 = 0.25
	defaultInitSize       int32   = 10
)

var (
	// ErrStreamPoolClosed will be returned if stream pool closed.
	ErrStreamPoolClosed = errors.New("stream pool closed")
	// ErrNoStreamCanBeBorrowed will be returned if no stream can be borrowed.
	ErrNoStreamCanBeBorrowed = errors.New("no stream can be borrowed")
)

// StreamPool is a stream pool.
type StreamPool struct {
	once          sync.Once
	newStreamFunc func() (network.Stream, error)
	initSize      int32
	cap           int32
	idleSize      int32
	currentSize   int32
	pool          chan network.Stream

	expandingSignalChan chan struct{}
	idlePercentage      float64
	closeChan           chan struct{}

	log protocol.Logger
}

func newStreamPool(cap int, newStreamFunc func() (network.Stream, error), log protocol.Logger) *StreamPool {
	return &StreamPool{
		once:                sync.Once{},
		newStreamFunc:       newStreamFunc,
		initSize:            defaultInitSize,
		cap:                 int32(cap),
		idleSize:            0,
		currentSize:         0,
		pool:                make(chan network.Stream, cap),
		expandingSignalChan: make(chan struct{}, 1),
		idlePercentage:      defaultIdlePercentage,
		closeChan:           make(chan struct{}),
		log:                 log,
	}
}

// nolint: unused
func (sp *StreamPool) isClosing() bool {
	select {
	case <-sp.closeChan:
		return true
	default:
		return false
	}
}

func (sp *StreamPool) borrowStream() (network.Stream, error) {
	sp.initStreams()
	var (
		stream network.Stream
	)
	select {
	case <-sp.closeChan:
		return nil, ErrStreamPoolClosed
	case stream = <-sp.pool:
		atomic.AddInt32(&sp.idleSize, -1)
		sp.triggerExpand()
		return stream, nil
	default:
		if sp.getCurrentSize() >= sp.maxSize() {
			return nil, ErrNoStreamCanBeBorrowed
		}
		select {
		case <-sp.closeChan:
			return nil, ErrStreamPoolClosed
		case stream = <-sp.pool:
			atomic.AddInt32(&sp.idleSize, -1)
			sp.triggerExpand()
			return stream, nil
		}
	}
}

func (sp *StreamPool) initStreams() {
	sp.once.Do(func() {
		var pid string
		for i := 0; i < int(sp.initSize) && atomic.LoadInt32(&sp.currentSize) < sp.initSize; i++ {
			stream, e := sp.newStreamFunc()
			if e != nil {
				sp.log.Warnf("[StreamPool] new stream failed, %s", e.Error())
				return
			}
			sp.pool <- stream
			atomic.AddInt32(&sp.idleSize, 1)
			atomic.AddInt32(&sp.currentSize, 1)
			if pid == "" {
				pid = stream.Conn().RemotePeer().String()
			}
		}
		sp.log.Infof("[StreamPool] init streams finished.(init-size: %d, remote pid: %s)",
			sp.initSize, pid)
		go sp.expandLoop()
	})
}

func (sp *StreamPool) triggerExpand() {
	select {
	case <-sp.closeChan:
	case sp.expandingSignalChan <- struct{}{}:
	default:
	}
}

func (sp *StreamPool) expandLoop() {
	for {
		select {
		case <-sp.closeChan:
			return
		case <-sp.expandingSignalChan:
			if sp.needExpand() {
				expandSize := sp.initSize / 2
				if expandSize == 0 {
					expandSize = 1
				}
				sp.log.Debugf("[StreamPool] expanding started."+
					"(expand: %d, cap: %d, current: %d, idle: %d)",
					expandSize, sp.cap, atomic.LoadInt32(&sp.currentSize),
					atomic.LoadInt32(&sp.idleSize))
				if expandSize+atomic.LoadInt32(&sp.currentSize) > sp.cap {
					expandSize = sp.cap - atomic.LoadInt32(&sp.currentSize)
				}
				temp := 0
				var pid string
				for i := 0; i < int(expandSize); i++ {
					stream, e := sp.newStreamFunc()
					if e != nil {
						sp.log.Warnf("[StreamPool] create send stream failed, %s", e.Error())
						continue
					}
					sp.pool <- stream
					atomic.AddInt32(&sp.idleSize, 1)
					atomic.AddInt32(&sp.currentSize, 1)
					temp++
					if pid == "" {
						pid = stream.Conn().RemotePeer().String()
					}
				}
				sp.log.Infof("[StreamPool] expanding finished."+
					"(expand: %d, cap: %d, current: %d, idle: %d, pid: %s)",
					temp, sp.cap, atomic.LoadInt32(&sp.currentSize),
					atomic.LoadInt32(&sp.idleSize), pid)
				sp.triggerExpand()
			}
		}
	}
}

func (sp *StreamPool) needExpand() bool {
	threshold := sp.expandingThreshold()
	if threshold >= atomic.LoadInt32(&sp.idleSize) &&
		atomic.LoadInt32(&sp.currentSize) < sp.cap &&
		atomic.LoadInt32(&sp.currentSize) != 0 {
		return true
	}
	return false
}

func (sp *StreamPool) expandingThreshold() int32 {
	return int32(math.Round(sp.idlePercentage * float64(atomic.LoadInt32(&sp.currentSize))))
}

// nolint: unused
func (sp *StreamPool) addStream(stream network.Stream) {
	if stream == nil {
		return
	}
	select {
	case <-sp.closeChan:
		return
	case sp.pool <- stream:
		sp.log.Infof("[StreamPool] add stream success(pid:%s)", stream.Conn().RemotePeer())
		atomic.AddInt32(&sp.currentSize, 1)
		atomic.AddInt32(&sp.idleSize, 1)
	default:
		sp.log.Infof("[StreamPool] add stream failed[pool full], dropped(pid:%s)", stream.Conn().RemotePeer())
		_ = stream.Reset()
	}
}

func (sp *StreamPool) returnStream(stream network.Stream) {
	select {
	case <-sp.closeChan:
		return
	case sp.pool <- stream:
		atomic.AddInt32(&sp.idleSize, 1)
	default:
		_ = stream.Reset()
	}
}

func (sp *StreamPool) dropStream(stream network.Stream) {
	sp.log.Debugf("[StreamPool] before drop stream (currentSize:%d , pid:%s)",
		sp.currentSize, stream.Conn().RemotePeer())
	_ = stream.Reset()
	if atomic.AddInt32(&sp.currentSize, -1) < 0 {
		atomic.AddInt32(&sp.currentSize, 1)
	}
}

func (sp *StreamPool) getCurrentSize() int {
	return int(atomic.LoadInt32(&sp.currentSize))
}

func (sp *StreamPool) maxSize() int {
	return int(sp.cap)
}

func (sp *StreamPool) cleanAndDisable() {
	select {
	case <-sp.closeChan:
		return
	default:

	}
	close(sp.closeChan)
Loop:
	for {
		select {
		case stream := <-sp.pool:
			sp.dropStream(stream)
		default:
			break Loop
		}
	}
	atomic.StoreInt32(&sp.currentSize, 0)
	atomic.StoreInt32(&sp.idleSize, 0)
	sp.log.Infof("clean the stream pool successfully")
}
