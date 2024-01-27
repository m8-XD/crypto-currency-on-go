package main_test

import (
	"net"
	"testing"
)

func TestConnectToServer(t *testing.T) {
	raddr := net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9000}
	_, err := net.DialTCP("tcp", nil, &raddr)

	if err != nil {
		t.Fail()
	}
	
	
}
