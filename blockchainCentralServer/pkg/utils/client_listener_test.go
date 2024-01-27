package utils_test

import (
	"blockchainCentralServer/pkg/entity"
	"blockchainCentralServer/pkg/utils"
	"net"
	"sync"
	"testing"
	"time"
)

func prepare(servAddr string) *entity.ServerInfo {
	addr, _ := net.ResolveTCPAddr("tcp", servAddr)

	listener, _ := net.ListenTCP("tcp", addr)

	si := entity.ServerInfo{}
	si.Start()
	si.SetListener(listener)

	return &si
}

func TestConnectionListener(t *testing.T) {
	servAddr := "127.0.0.1:9001"
	si := prepare(servAddr)
	wg := sync.WaitGroup{}

	wg.Add(1)
	go utils.Listen(si, &wg)
	time.Sleep(2 * time.Second)
	_, err := net.Dial("tcp", servAddr)

	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(5 * time.Second)

	connections := si.Connections()

	if len(connections) == 0 {
		t.Fatal("there are no connections")
	}

	si.Stop()
	wg.Wait()
	si.Listener().Close()
}
