package listeners

import (
	"blockchain/pkg/entity"
	"blockchain/pkg/utils"
	"bytes"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

func ServerListen(c *entity.Client, wg *sync.WaitGroup) {
	defer wg.Done()
	centralServ := c.CentralServ()
	buff := make([]byte, utils.BUFFER_SIZE)
	for c.IsRunning() {
		centralServ.SetReadDeadline(time.Now().Add(1 * time.Minute))
		_, err := centralServ.Read(buff)

		if err != nil {
			panic(err)
		}
		go createWriteConnections(buff, c)
	}
}

func createWriteConnections(buff []byte, c *entity.Client) {
	addrs := bytes.Split(buff, []byte{','})
	for _, addr := range addrs {
		addrSt := utils.TrimAndCast(addr)
		// ip: n.n.n.n:port, so minimum length of ip is 9
		if utf8.RuneCountInString(addrSt) < 9 {
			addrSt = "127.0.0.1:" + addrSt
		}

		if strings.Compare(addrSt, c.LocalServ().Addr().String()) == 0 {
			continue
		}
		fmt.Println("creating connection from: " + c.LocalServ().Addr().String() + " to: " + addrSt)
		go createWriteConnection(addrSt, c)
	}
}

func createWriteConnection(addr string, c *entity.Client) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("error during conecting: " + err.Error())
		return
	}
	c.AddWritePeer(conn)
}
