package utils_test

import (
	"blockchainCentralServer/pkg/utils"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"
)

func TestBroadcastData(t *testing.T) {
	servAddr := "127.0.0.1:9000"
	si := prepare(servAddr)
	wg := sync.WaitGroup{}

	wg.Add(2)
	go utils.Listen(si, &wg)
	go utils.BroadcastClientsData(si, &wg)
	time.Sleep(2 * time.Second)
	clientConn, err := net.Dial("tcp", servAddr)

	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(30 * time.Second)

	buff := make([]byte, 25)

	clientConn.SetDeadline(time.Now().Add(10 * time.Second))
	n, err := clientConn.Read(buff)
	if err != nil {
		t.Fatal(err)
	}

	if n == 0 {
		t.Fatal("length of readed buffer is: 0")
	}

	fmt.Println("\ncheck it pls: " + string(buff) + "\n")

	clientConn.Close()

	time.Sleep(61 * time.Second)

	connections := si.Connections()

	if len(connections) != 0 {
		t.Fatal("connections slice ain't empty")
	}

	si.Stop()
	wg.Wait()
	si.Listener().Close()
}
