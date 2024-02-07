package entity

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

type Client struct {
	centralServ net.Conn
	localServ   *net.TCPListener
	readPeers   []net.Conn
	writePeers  []net.Conn
	isRunning   bool
	mut         sync.Mutex

	chain    []Node //copy of blockchain
	chainMut sync.Mutex
}

func (c *Client) SetCentralServ(cs net.Conn) {
	c.centralServ = cs
}

func (c *Client) SetLocalServ(ls *net.TCPListener) {
	c.localServ = ls
}

func (c *Client) CentralServ() net.Conn {
	return c.centralServ
}

func (c *Client) LocalServ() *net.TCPListener {
	return c.localServ
}

func (c *Client) Start() {
	c.ResetPeers()
	c.isRunning = true
}

func (c *Client) Stop() {
	c.isRunning = false
}

func (c *Client) IsRunning() bool {
	return c.isRunning
}

func (c *Client) AddReadPeer(peer net.Conn) {
	if c.readPeers == nil {
		fmt.Println("read peers slice is nil, pls call Start first! skipping...")
		return
	}
	c.mut.Lock()
	c.readPeers = append(c.readPeers, peer)
	c.mut.Unlock()
}

func (c *Client) AddWritePeer(peer net.Conn) {
	if c.writePeers == nil {
		fmt.Println("read peers slice is nil, pls call Start first! skipping...")
		return
	}
	c.mut.Lock()
	c.writePeers = append(c.writePeers, peer)
	c.mut.Unlock()
}

func (c *Client) ResetPeers() {
	c.readPeers = make([]net.Conn, 0)
	c.writePeers = make([]net.Conn, 0)
}

func (c *Client) WritePeers() []net.Conn {
	c.mut.Lock()
	toReturn := make([]net.Conn, len(c.writePeers))
	copy(toReturn, c.writePeers)
	c.mut.Unlock()
	return toReturn
}

func (c *Client) ReadPeers() []net.Conn {
	c.mut.Lock()
	toReturn := make([]net.Conn, len(c.readPeers))
	copy(toReturn, c.readPeers)
	c.mut.Unlock()
	return toReturn
}

func (c *Client) AddNode(n *Node) {
	c.chainMut.Lock()
	lastNode := c.chain[len(c.chain)-1]
	if strings.EqualFold(lastNode.Header, n.PHeader) {
		c.chain = append(c.chain, *n)
	}
	c.chainMut.Unlock()

}

func (c *Client) Chain() []Node {
	copychain := make([]Node, len(c.chain))
	copy(copychain, c.chain)
	return copychain
}

func (c *Client) ReceiveChain(chainRaw string) {
	chainSt := strings.Split(chainRaw, ";")
	chain := make([]Node, len(chainSt))

	for _, nodeSt := range chainSt {
		node, err := Unpack(nodeSt)
		if err != nil {
			return
		}
		chain = append(chain, *node)
	}
	c.chain = chain
}
