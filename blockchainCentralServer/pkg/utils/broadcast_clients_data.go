package utils

import (
	"blockchainCentralServer/pkg/entity"
	"sync"
	"time"
)

func BroadcastClientsData(s *entity.ServerInfo, wg *sync.WaitGroup) {
	defer wg.Done()
	for s.IsRunning() {
		time.Sleep(30 * time.Second)

		for _, peer := range s.Connections() {
			_, err := peer.Write([]byte(s.Addrs()[0]))

			if err != nil {
				s.CloseConnection(peer.RemoteAddr().String())
			}
		}
	}
}
