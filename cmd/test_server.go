package main

import (
	"fmt"
	"github.com/95eh/eg"
	"github.com/95eh/eg/network"
	"net"
)

func main() {
	listener := network.NewTcpListener(":10001", func(conn net.Conn) {
		onAddConn(conn, network.NewTcpAgent)
	})
	listener.Start()
	<-(make(chan struct{}))
}

func onAddConn(conn net.Conn, newAgent eg.ToNewAgent) {
	addr := conn.RemoteAddr().String()
	agent := newAgent(addr, receiver,
		eg.AgentMaxBadPacketCount(5),
		eg.AgentMaxBadPacketInterval(5),
		eg.AgentClosed(onAgentClosed))
	agent.Start(conn)
}

func onAgentClosed(agent eg.IAgent, err eg.IErr) {
	err.AddParam("addr", agent.Id())
	eg.Error(err)
}

func receiver(agent eg.IAgent, bytes []byte, err eg.FnErr) {
	fmt.Printf("receive:%s\n", string(bytes))
	agent.Send(bytes)
}
