package utils_test

import (
	"blockchain/pkg/entity"
	"blockchain/pkg/utils"
	"fmt"
	"net"
	"strings"
	"sync"
	"testing"
	"time"
)

func prepare(port int) (*entity.Client, error) {

	clientEnt := entity.Client{}

	serverConn, err := net.Dial("tcp", "127.0.0.1:9000")
	if err != nil {
		return nil, err
	}
	serverConn.Write([]byte(fmt.Sprint(port)))

	clientEnt.SetCentralServ(serverConn)

	tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:"+fmt.Sprint(port))
	if err != nil {
		fmt.Println(err)
		return &clientEnt, err
	}

	localServer, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println(err)
		return &clientEnt, err
	}

	clientEnt.SetLocalServ(localServer)

	clientEnt.Start()
	return &clientEnt, nil
}

func TestConnection(t *testing.T) {
	c, err := prepare(9001)
	if err != nil {
		t.Fatal(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(2)
	go utils.ListenForPeers(c, &wg)
	go utils.ServerListen(c, &wg)

	time.Sleep(40 * time.Second)

	c.Stop()
	wg.Wait()
	c.LocalServ().Close()
	c.CentralServ().Close()
}

func TestListeners(t *testing.T) {
	c, err := prepare(9002)
	if err != nil {
		t.Fatal(err)
	}
	c1, err := prepare(9003)
	if err != nil {
		t.Fatal(err)
	}
	defer c.LocalServ().Close()
	defer c.CentralServ().Close()
	defer c1.LocalServ().Close()
	defer c1.CentralServ().Close()

	wg := sync.WaitGroup{}
	wg.Add(2)
	go utils.ListenForPeers(c, &wg)
	go utils.ServerListen(c, &wg)

	wg.Add(2)
	go utils.ListenForPeers(c1, &wg)
	go utils.ServerListen(c1, &wg)

	time.Sleep(40 * time.Second)

	if len(c.WritePeers()) == 0 {
		t.Log(c.WritePeers())
		t.Fatal("no write peers for client 1")
	}

	if len(c1.WritePeers()) == 0 {
		t.Log(c1.WritePeers())
		t.Fatal("no write peers for client 2")
	}

	if len(c.ReadPeers()) == 0 {
		t.Log(c.ReadPeers())
		t.Fatal("no read peers for client 1")
	}

	if len(c1.ReadPeers()) == 0 {
		t.Log(c1.ReadPeers())
		t.Fatal("no read peers for client 2")
	}

	msg1 := "test"
	msg2 := "test1"

	utils.Write(msg1, c)
	msgs1 := utils.Read(c1)

	if strings.Compare(msg1, msgs1[0]) != 0 {
		t.Log(len(msgs1))
		t.Log(msgs1)
		t.Log(len(c.ReadPeers()))
		t.Fatal("first messages arent the same")
	}

	utils.Write(msg2, c1)
	msgs2 := utils.Read(c)

	if strings.Compare(msg2, msgs2[0]) != 0 {
		t.Log(len(msgs2))
		t.Log(msgs2)
		t.Fatal("second messages arent the same")
	}

	c.Stop()
	c1.Stop()
	wg.Wait()
}
