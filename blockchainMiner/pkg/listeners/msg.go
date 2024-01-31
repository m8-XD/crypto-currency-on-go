package listeners

import (
	"blockchain/pkg/entity"
	"blockchain/pkg/utils"
	"strings"
	"sync"
	"time"
)

func MsgListen(wg *sync.WaitGroup, m *entity.Miner) {
	defer wg.Done()
	for m.IsRunning() {
		time.Sleep(5 * time.Second)
		msgs := utils.Read(m.Client())
		go handleMsgs(m, msgs)
	}
}

func handleMsgs(m *entity.Miner, msgs []string) {
	for _, msg := range msgs {
		if strings.EqualFold(msg, "") {
			continue
		}
		handleMsg(m, msg)
	}
}

func handleMsg(m *entity.Miner, msg string) {
	if strings.HasPrefix(msg, "valid") {
		//validate
		error
	} else {
		m.AddTX(msg)
	}
}
