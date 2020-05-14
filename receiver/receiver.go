package receiver

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/vsliouniaev/packet-loss/model"
	"net"
	"time"
)

type Receiver struct {
	logger     log.Logger
	connection *net.UDPConn
	stop       chan interface{}
	stopped    bool
}

func (r *Receiver) Stop() {
	if r.stopped {
		return
	}
	r.stopped = true
	close(r.stop)
	if err := r.connection.Close(); err != nil {
		level.Error(r.logger).Log("msg", "unable to close connection", "err", err)
	}
	level.Info(r.logger).Log("msg", "stopped")
}

func New(logger log.Logger, address string) (*Receiver, error) {
	s, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		return nil, err
	}

	connection, err := net.ListenUDP("udp4", s)
	if err != nil {
		return nil, err
	}

	return &Receiver{
		logger:     logger,
		connection: connection,
		stop:       make(chan interface{}),
	}, nil
}

func (r *Receiver) Receive() {
	buffer := make([]byte, 4096)
	current := uint64(0)
	level.Info(r.logger).Log("msg", "started")
	for {
		select {
		case <-r.stop:
			return
		default:
			numRead, _, err := r.connection.ReadFromUDP(buffer)
			if r.stopped {
				continue
			}
			if err != nil {
				level.Error(r.logger).Log("msg", "unable to read", "err", err)
				continue
			}
			packet := &model.Packet{}
			if err := proto.Unmarshal(buffer[0:numRead], packet); err != nil {
				level.Error(r.logger).Log("msg", "failed unmarshal", "err", err)
			}

			diff := packet.Number - current
			if current != 0 && diff != 1 {
				// TODO: This whole thing
				level.Info(r.logger).Log("msg", "out of oder", "diff", diff)
			}
			if diff < 0 {
				// TODO: and this
				current = packet.Number
			}

			ts, err := ptypes.Timestamp(packet.Timestamp)
			if err != nil {
				level.Error(r.logger).Log("msg", "failed unmarshal timestamp", "err", err)
			}
			level.Info(r.logger).Log("msg", "received", "diff", time.Since(ts), "num", packet.Number)
		}
	}
}
