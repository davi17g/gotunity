package main

import (
	"testing"
	"fmt"
)

func TestStartServer(t *testing.T) {
	host := "127.0.0.1"
	pin = 8888
	pout = 8089
	StartServer(host, pin, pout)
	fmt.Println("Here")

}

func SendUDPMessage(host *string, port *int) {
	addr := net.UDPAddr{Port: *port, IP: net.ParseIP(*host)}
	packet :=  make([]byte, 1024)
	address := *host + ":" + string(port)
	fmt.Println(address)
	conn, err := net.Dial("udp", "127.0.0.1:1234")


}
