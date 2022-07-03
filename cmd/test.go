package main

import (
	"game/pkg/proto/pb"
	"github.com/95eh/eg"
	"net"
	"time"
)

var (
	conn  net.Conn
	count int32
	//addr           = "10.6.2.26:10011"
	addr           = "47.243.177.125:10011"
	sndCode uint16 = 249
	rcvCode uint16 = 250
)

func main() {
	connect()

	c := make(chan struct{})
	<-c
	eg.Info("exit", nil)
}

func reconnect() {
	time.Sleep(time.Second * 2)
	connect()
}

func connect() {
	count = 0
	eg.Info("connect", eg.M{
		"addr": addr,
	})
	tcpAddr, e := net.ResolveTCPAddr("tcp", addr)
	if e != nil {
		eg.Error(eg.NewErr(eg.EcConnectErr, eg.M{
			"addr":  addr,
			"error": e.Error(),
		}))
		reconnect()
		return
	}

	conn, e = net.DialTCP("tcp", nil, tcpAddr)
	if e != nil {
		eg.Error(eg.NewErr(eg.EcConnectErr, eg.M{
			"addr":  addr,
			"error": e.Error(),
		}))
		reconnect()
		return
	}

	eg.Info("connected", eg.M{
		"addr": addr,
	})
	go read()

	writePkt()
}

func read() {
	defer func() {
		reconnect()
	}()

	const (
		Cap = 64
	)
	var (
		buffer = make([]byte, Cap)
		ring   = eg.NewRing[byte](
			eg.RingMaxCap[byte](Cap),
			eg.RingMinCap[byte](Cap<<4),
			eg.RingResize[byte](func(c uint32) {
				eg.Debug("resize", eg.M{
					"size": c,
					"addr": addr,
				})
				buffer = make([]byte, c)
			}),
		)
		headLen uint16
	)
	for {
		if e := conn.SetReadDeadline(time.Now().Add(time.Second * 10)); e != nil {
			eg.Error(eg.WrapErr(eg.EcTimeout, e))
			return
		}
		newLen, e := conn.Read(buffer)
		if e != nil {
			eg.Error(eg.NewErr(eg.EcConnectErr, eg.M{
				"addr":  addr,
				"error": e.Error(),
			}))
			return
		}
		err := ring.Write(buffer[:newLen]...)
		if err != nil {
			eg.Error(err)
			return
		}
		for {
			if headLen == 0 {
				if ring.Available() < 2 {
					break
				}
				ring.Read(buffer, 2)
				headLen = readUint16(buffer)
				if headLen == 0 {
					err = eg.NewErr(eg.EcBadHead, nil)
					return
				}
			}
			if headLen == 2 {
				headLen = 0
				continue
			}
			l := uint32(headLen)
			if ring.Available() < l {
				break
			}
			headLen = 0
			ring.Read(buffer, l)
			receivePkt(buffer[:l])
		}
	}
}

func receivePkt(bytes []byte) {
	code := readUint16(bytes)
	if code != rcvCode {
		eg.Error(eg.NewErr(eg.EcBadPacket, eg.M{
			"error": "wrong code",
		}))
		return
	}
	var sd pb.G2C_ServerDate
	e := sd.Unmarshal(bytes[2:])
	if e != nil {
		eg.Error(eg.WrapErr(eg.EcUnmarshallErr, e))
		return
	}
	eg.Info("receive", eg.M{
		"pkg": sd,
		"hex": eg.Hex(bytes),
	})
	time.Sleep(time.Second * 1)
	writePkt()
}

func readUint16(bytes []byte) uint16 {
	return uint16(bytes[0]) | uint16(bytes[1])<<8
}

func writeUint16(bytes []byte, v uint16, offset int) {
	bytes[offset] = byte(v)
	bytes[offset+1] = byte(v) >> 8
}

func writePkt() {
	count++
	pkt := &pb.C2G_ServerDate{
		RpcId: count,
	}
	bytes, _ := pkt.Marshal()
	l := uint16(len(bytes) + 4)
	all := make([]byte, l)
	writeUint16(all, l-2, 0)
	writeUint16(all, sndCode, 2)
	copy(all[4:], bytes)
	_, e := conn.Write(all)
	if e != nil {
		eg.Error(eg.WrapErr(eg.EcSendErr, e))
		return
	}
	eg.Info("send", eg.M{
		"hex": eg.Hex(all[2:]),
		"pkg": pkt,
	})
}
