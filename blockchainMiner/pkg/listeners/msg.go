package listeners

import (
	"blockchain/pkg/entity/mining"
	"blockchain/pkg/utils"
	"strings"
	"sync"
	"time"
)

func MsgListen(wg *sync.WaitGroup, m *mining.Miner) {
	defer wg.Done()
	for m.IsRunning() {
		time.Sleep(5 * time.Second)
		msgs := utils.Read(m.Client())
		go handleMsgs(m, msgs)
	}
}

func handleMsgs(m *mining.Miner, msgs []string) {
	for _, msg := range msgs {
		if strings.EqualFold(msg, "") {
			continue
		}
		if strings.HasPrefix(msg, "block") {
			continue
		}
		m.AddTX(msg)
	}
}
