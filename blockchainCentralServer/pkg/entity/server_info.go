package entity

import (
	"fmt"
	"net"
	"sync"
)

type ServerInfo struct {
	isRunning   bool
	listener    *net.TCPListener
	connections map[*net.TCPConn]string
	mut         sync.Mutex
}

func (s *ServerInfo) Stop() {
	fmt.Println("stopping the server")
	s.isRunning = false
}

func (s *ServerInfo) Start() {
	fmt.Println("starting the server(press ctrl + C to stop)")
	s.isRunning = true
	s.connections = make(map[*net.TCPConn]string)
}

func (s *ServerInfo) Listener() *net.TCPListener {
	return s.listener
}

func (s *ServerInfo) SetListener(l *net.TCPListener) {
	s.listener = l
}

func (s *ServerInfo) IsRunning() bool {
	return s.isRunning
}

func (s *ServerInfo) AddConnection(conn *net.TCPConn) {
	if s.connections == nil {
		fmt.Println("call Start function before adding a connection, skipping...")
		return
	}

	buff := make([]byte, 25)
	conn.Read(buff)

	s.mut.Lock()
	s.connections[conn] = string(buff)
	s.mut.Unlock()
}

func (s *ServerInfo) Connections() []*net.TCPConn {
	dst := make([]*net.TCPConn, len(s.connections))
	i := 0
	s.mut.Lock()
	for connection := range s.connections {
		dst[i] = connection
		i++
	}
	s.mut.Unlock()
	return dst
}

func (s *ServerInfo) CloseConnection(conn *net.TCPConn) {
	s.mut.Lock()
	delete(s.connections, conn)
	s.mut.Unlock()
}

func (s *ServerInfo) Addrs() []string {
	var addrs []string
	s.mut.Lock()
	for _, value := range s.connections {
		addrs = append(addrs, value)
	}
	s.mut.Unlock()
	return addrs
}
