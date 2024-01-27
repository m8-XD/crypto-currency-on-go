package utils

import (
	"blockchainCentralServer/pkg/entity"
	"fmt"
	"strings"
	"sync"
	"time"
)

func BroadcastClientsData(s *entity.ServerInfo, wg *sync.WaitGroup) {
	defer wg.Done()
	for s.IsRunning() {
		time.Sleep(30 * time.Second)
		fmt.Println("broadcasting data to " + fmt.Sprint(len(s.Connections())) + " peers")
		for _, peer := range s.Connections() {
			addresses := strings.Join(s.Addrs(), ",")
			_, err := peer.Write([]byte(addresses))

			if err != nil {
				fmt.Println("client disconnected: " + peer.RemoteAddr().String())
				s.CloseConnection(peer)
			}
		}
	}
}
