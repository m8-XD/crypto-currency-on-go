package entity

import (
	"fmt"
	"net"
	"sync"
)

type ServerInfo struct {
	isRunning   bool
	listener    *net.TCPListener
	connections map[string]*net.TCPConn
	mut         sync.Mutex
}

func (s *ServerInfo) Stop() {
	fmt.Println("stopping the server")
	s.isRunning = false
}

func (s *ServerInfo) Start() {
	fmt.Println("starting the server(press ctrl + C to stop)")
	s.isRunning = true
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
	s.mut.Lock()
	s.connections[conn.LocalAddr().String()] = conn
	s.mut.Unlock()
}

func (s *ServerInfo) Connections() []*net.TCPConn {
	dst := make([]*net.TCPConn, len(s.connections))
	i := 0
	s.mut.Lock()
	for _, connection := range s.connections {
		dst[i] = connection
		i++
	}
	s.mut.Unlock()
	return dst
}

func (s *ServerInfo) CloseConnection(conn string) {
	s.mut.Lock()
	delete(s.connections, conn)
	s.mut.Unlock()
}

func (s *ServerInfo) Ports() []string {
	var ports []string
	s.mut.Lock()
	for key := range s.connections {
		ports = append(ports, key)
	}
	s.mut.Unlock()
	return ports
}
