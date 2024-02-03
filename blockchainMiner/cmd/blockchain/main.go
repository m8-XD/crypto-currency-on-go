package main

import (
	"blockchain/pkg/entity"
	"blockchain/pkg/entity/mining"
	"blockchain/pkg/listeners"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
)

var wg = sync.WaitGroup{}

func main() {

	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("please pass valid arguments, there should be only one (port number)")
	}

	clientEnt := entity.Client{}
	minerEnt := mining.Miner{}

	serverConn, err := net.Dial("tcp", "127.0.0.1:9000")
	defer serverConn.Close()
	serverConn.Write([]byte(fmt.Sprint(port)))

	clientEnt.SetCentralServ(serverConn)

	tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:"+fmt.Sprint(port))
	if err != nil {
		fmt.Println(err)
		return
	}

	localServer, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer localServer.Close()

	clientEnt.SetLocalServ(localServer)

	minerEnt.SetClient(&clientEnt)

	minerEnt.Start()

	wg.Add(3)
	go listeners.ServerListen(&clientEnt, &wg)
	go listeners.ListenForPeers(&clientEnt, &wg)
	go listeners.MsgListen(&minerEnt, &wg)

	wg.Wait()
}
