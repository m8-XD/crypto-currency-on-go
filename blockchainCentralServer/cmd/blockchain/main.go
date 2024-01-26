package main

import (
	"net"
	"os"
	"os/signal"
	"sync"

	"blockchainCentralServer/pkg/entity"
	"blockchainCentralServer/pkg/utils"
)

var sInfo = entity.ServerInfo{}
var wg = sync.WaitGroup{}

func main() {
	sInfo.Start()

	listener, _ := net.ListenTCP("tcp",
		&net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9000})

	sInfo.SetListener(listener)
	defer listener.Close()

	wg.Add(2)
	go utils.Listen(&sInfo, &wg)
	go utils.BroadcastClientsData(&sInfo, &wg)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		sInfo.Stop()
	}()

	wg.Wait()
}
