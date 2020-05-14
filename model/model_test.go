package model_test

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/vsliouniaev/packet-loss/model"
	"testing"
)

func TestRoundTrip(t *testing.T) {
	send, err := proto.Marshal(&model.Packet{
		Number:    2890123,
		Timestamp: ptypes.TimestampNow(),
	})
	if err != nil {
		t.Error(err)
	}

	receive := &model.Packet{}
	err = proto.Unmarshal(send, receive)
	if err != nil {
		t.Error(err)
	}
	if receive.Timestamp == nil {
		t.Error("failed unmarshal timestamp")
	}
}
