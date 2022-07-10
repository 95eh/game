package service

import (
	"game/pkg/cmn"
	"game/pkg/proto"
	"game/pkg/proto/pb"
	"github.com/95eh/eg"
)

func PushToUser(tid, id int64, service, code int, pkt interface{}) {
	bytes, err := eg.Codec().Marshal(pkt)
	if err != nil {
		eg.Log().TraceErr(tid, err)
		return
	}
	buffer := eg.NewByteBufferWithLen(uint32(4 + len(bytes)))
	buffer.WUint16(uint16(service))
	buffer.WUint16(uint16(code))
	buffer.Write(bytes)
	eg.Req().Request(tid, id, cmn.SvcGate, proto.CdGatePushToUser, &pb.PushToUser{
		Id:    id,
		Bytes: buffer.All(),
	}, nil, nil)
}

func PushToUsers(tid int64, ids []int64, serve, code int, pkt interface{}) {
	bytes, err := eg.Codec().Marshal(pkt)
	if err != nil {
		eg.Log().TraceErr(tid, err)
		return
	}
	buffer := eg.NewByteBufferWithLen(uint32(4 + len(bytes)))
	buffer.WUint16(uint16(serve))
	buffer.WUint16(uint16(code))
	buffer.Write(bytes)
	allBytes := buffer.All()

	for _, id := range ids {
		eg.Req().Request(tid, id, cmn.SvcGate, proto.CdGatePushToUsers, &pb.PushToUsers{
			Ids:   ids,
			Bytes: eg.CopyBytes(allBytes),
		}, nil, nil)
	}
}

func PushToAllUsers(tid int64, serve, code int, pkt interface{}) {
	bytes, err := eg.Codec().Marshal(pkt)
	if err != nil {
		eg.Log().TraceErr(tid, err)
		return
	}
	buffer := eg.NewByteBufferWithLen(uint32(4 + len(bytes)))
	buffer.WUint16(uint16(serve))
	buffer.WUint16(uint16(code))
	buffer.Write(bytes)
	eg.Req().Request(tid, 0, cmn.SvcGate, proto.CdGatePushToAllUsers, &pb.PushToUsers{
		Bytes: buffer.All(),
	}, nil, nil)
}
