package utils

import (
	"blockchainCentralServer/pkg/entity"
	"sync"
	"time"
)

func Listen(si *entity.ServerInfo, wg *sync.WaitGroup) {
	defer wg.Done()

	for si.IsRunning() {
		si.Listener().SetDeadline(time.Now().Add(30 * time.Second))
		conn, err := si.Listener().AcceptTCP()
		if err != nil {
			continue
		}
		si.AddConnection(conn)
	}
}