package listeners

import (
	"blockchain/pkg/entity"
	"blockchain/pkg/utils"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

func Listen(c *entity.Client, wg *sync.WaitGroup) {
	defer wg.Done()

	for c.IsRunning() {
		time.Sleep(5 * time.Second)
		msgs := utils.Read(c)
		for _, msg := range msgs {
			if utf8.RuneCountInString(msg) == 0 {
				continue
			}

			if strings.Contains(msg, ";") {
				c.ReceiveChain(msg)
				continue
			}

			node, err := entity.Unpack(msg)
			if err != nil {
				continue
			}

			c.AddNode(node)
		}
	}
}
