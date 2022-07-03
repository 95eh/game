package main

import (
	"fmt"
	"github.com/95eh/eg"
	"github.com/95eh/eg/network"
	"strconv"
	"time"
)

var (
	dialer   eg.IDialer
	msgCount = 0
)

func main() {
	connectServer()
	<-(make(chan struct{}))
}

func connectServer() {
	//addr := "47.243.177.125:10001"
	addr := "127.0.0.1:10001"
	dialer = network.NewTcpDialer("test", addr, func(agent eg.IAgent, bytes []byte, err eg.FnErr) {
		eg.Info("receive", eg.M{
			"msg": string(bytes),
		})
		time.Sleep(time.Second)
		sendTestMsg()
	})

	dialer.Connect(
		eg.AgentConnected(func() {
			fmt.Printf("connected")
		}),
		eg.AgentClosed(func(agent eg.IAgent, err eg.IErr) {
			if err != nil {
				eg.Error(err)
			}
			reconnectServer()
		}),
	)
	sendTestMsg()
}

func sendTestMsg() {
	msgCount++
	bytes := []byte(strconv.Itoa(msgCount))
	err := dialer.Agent().Send(bytes)
	if err != nil {
		eg.Error(err)
		return
	}
	eg.Info("send", eg.M{
		"msg": string(bytes),
	})
}

func reconnectServer() {
	time.Sleep(time.Second)
	connectServer()
}

type msg struct {
}