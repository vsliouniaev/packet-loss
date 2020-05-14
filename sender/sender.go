package sender

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/vsliouniaev/packet-loss/model"
	"net"
	"time"
)

type Sender struct {
	logger     log.Logger
	connection *net.UDPConn
	stop       chan interface{}
	stopped    bool
	ticker     *time.Ticker
}

func (s *Sender) Stop() {
	if s.stopped {
		return
	}
	s.stopped = true
	close(s.stop)
	s.ticker.Stop()
	if err := s.connection.Close(); err != nil {
		level.Error(s.logger).Log("msg", "unable to close connection", "err", err)
	}
	level.Info(s.logger).Log("msg", "stopped")
}

func New(logger log.Logger, address string, interval time.Duration) (*Sender, error) {
	s, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		return nil, err
	}

	connection, err := net.DialUDP("udp4", nil, s)
	if err != nil {
		return nil, err
	}

	return &Sender{
		logger:     logger,
		connection: connection,
		stop:       make(chan interface{}),
		stopped:    false,
		ticker:     time.NewTicker(interval),
	}, nil
}

func (s *Sender) Send() {
	number := uint64(0)
	level.Info(s.logger).Log("msg", "started")
	for {
		select {
		case <-s.stop:
			return
		case <-s.ticker.C:
			number++
			m, err := proto.Marshal(&model.Packet{
				Number:    number,
				Timestamp: ptypes.TimestampNow(),
			})
			if err != nil {
				level.Error(s.logger).Log("msg", "failed marshal", "err", err, "num", number)
			}
			_, err = s.connection.Write(m)
			if err != nil {
				level.Error(s.logger).Log("msg", "failed send", "err", err, "num", number)
			}

		default:
		}
	}
}
