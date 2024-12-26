package libp2pnet

import (
	"context"
	"errors"
	"sync"
	"time"

	protocol2 "chainmaker.org/chainmaker/protocol/v2"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/protocol"
)

var _ network.Stream = (*ySendStream)(nil)

// ySendStream
type ySendStream struct {
	// network.Stream
	stream network.Stream
	// send chan
	sendCh chan []byte
	// log
	log protocol2.Logger
	// ctx
	ctx context.Context
	// writeLock lock when Write()
	writeLock sync.Mutex
	// for reset
	once sync.Once
}

// newSendStream ySendStream
func newSendStream(stream network.Stream, log protocol2.Logger) *ySendStream {
	ss := &ySendStream{
		stream: stream,
		log:    log,
		sendCh: make(chan []byte, 32),
	}
	//go ss.sendLoop()
	return ss
}

// for send until write return err
//func (y *ySendStream) sendLoop() {
//	for {
//		select {
//		case <-y.ctx.Done():
//			return
//		case buf := <-y.sendCh:
//			_, err := y.stream.Write(buf)
//			if err != nil {
//				y.log.Warnf("sendLoop err:", err.Error())
//				return
//			}
//		}
//	}
//}

// ySendStream Read
func (y *ySendStream) Read(p []byte) (n int, err error) {
	if y.stream == nil {
		return 0, nil
	}
	return y.stream.Read(p)
}

// ySendStream Write
func (y *ySendStream) Write(p []byte) (int, error) {
	if y.stream == nil {
		return 0, nil
	}
	y.writeLock.Lock()
	defer y.writeLock.Unlock()
	nAlreadySend := 0
	n := len(p)
	for {
		nSend, err := y.stream.Write(p[nAlreadySend:])
		if err != nil {
			return 0, err
		}
		if nSend <= 0 {
			return nSend, errors.New("EOF")
		}
		nAlreadySend += nSend
		if nAlreadySend == n {
			return n, nil
		} else if nAlreadySend > n {
			//never happened before ,just for case
			return nAlreadySend, nil
		}
	}
}

// ySendStream write
// for write until all []byte write ok or time out
//func (y *ySendStream) write(p []byte) (n int, err error) {
//select {
//case y.sendCh <- p:
//	return len(p), err
//case <-time.After(20 * time.Second):
//	return 0, errors.New("SendStream write timeout")
//}
//}

// ySendStream Close
func (y *ySendStream) Close() error {
	if y.stream == nil {
		return nil
	}
	return y.stream.Close()
}

// ySendStream Reset
func (y *ySendStream) Reset() error {
	if y.stream == nil {
		return nil
	}
	y.once.Do(
		func() {
			y.stream.Reset()
		})
	return nil
}

// ySendStream SetDeadline
func (y *ySendStream) SetDeadline(time time.Time) error {
	if y.stream == nil {
		return nil
	}
	return y.stream.SetDeadline(time)
}

// ySendStream SetReadDeadline
func (y *ySendStream) SetReadDeadline(time time.Time) error {
	if y.stream == nil {
		return nil
	}
	return y.stream.SetReadDeadline(time)
}

// ySendStream SetWriteDeadline
func (y *ySendStream) SetWriteDeadline(time time.Time) error {
	if y.stream == nil {
		return nil
	}
	return y.stream.SetWriteDeadline(time)
}

// ySendStream ID
func (y *ySendStream) ID() string {
	if y.stream == nil {
		return ""
	}
	return y.stream.ID()
}

// ySendStream Protocol
func (y *ySendStream) Protocol() protocol.ID {
	if y.stream == nil {
		return ""
	}
	return y.stream.Protocol()
}

// ySendStream SetProtocol
func (y *ySendStream) SetProtocol(id protocol.ID) {
	if y.stream == nil {
		return
	}
	y.stream.SetProtocol(id)
}

// ySendStream Stat
func (y *ySendStream) Stat() network.Stat {
	if y.stream == nil {
		return network.Stat{}
	}
	return y.stream.Stat()
}

// ySendStream Conn
func (y *ySendStream) Conn() network.Conn {
	if y.stream == nil {
		return nil
	}
	return y.stream.Conn()
}
