package entity

import (
	"fmt"
	"net"
	"sync"
)

type Client struct {
	centralServ net.Conn
	localServ   *net.TCPListener
	readPeers   []net.Conn
	writePeers  []net.Conn
	isRunning   bool
	mut         sync.Mutex
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
