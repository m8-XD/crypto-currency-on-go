package listeners

import (
	"blockchain/pkg/entity/mining"
	"fmt"
	"sync"
	"time"
)

func ListenForPeers(m *mining.Miner, wg *sync.WaitGroup) {
	defer wg.Done()
	c := m.Client()
	for c.IsRunning() {
		c.LocalServ().SetDeadline(time.Now().Add(30 * time.Second))
		conn, err := c.LocalServ().Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go c.AddReadPeer(conn)
		go m.SendChain()
	}
}
