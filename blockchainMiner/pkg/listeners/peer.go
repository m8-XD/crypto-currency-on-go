package listeners

import (
	"blockchain/pkg/entity"
	"fmt"
	"sync"
	"time"
)

func ListenForPeers(c *entity.Client, wg *sync.WaitGroup) {
	defer wg.Done()
	for c.IsRunning() {
		c.LocalServ().SetDeadline(time.Now().Add(30 * time.Second))
		conn, err := c.LocalServ().Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go c.AddReadPeer(conn)
	}
}
