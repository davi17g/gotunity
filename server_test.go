package main

import (
	"testing"
	"fmt"
	"net"
	"strconv"
	"time"
)

func TestStartServer(t *testing.T) {
	host := "127.0.0.1"
	pin := 8888
	pout := 8089
	ans := make(chan string)
	msg := "TEST\r\n"
	go StartServer(&host, &pin, &pout)
	go SendUDPMessage(&host, &pin, &msg)
	go ReciveUDPmessage(&host, &pout, ans)
	select {
	case result := <-ans:
		if result == "TEST" {
			fmt.Println("Test Succeded")
			} else {
				t.Fail()
				break
			}
			case <- time.After(7 * time.Second):
				fmt.Println("Test timouted")
		}
	}

func SendUDPMessage(host *string, port *int, msg *string) {
	packet := make([]byte, 1024)
	address := *host + ":" + strconv.Itoa(*port)
	conn, err := net.Dial("udp", address)
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
	fmt.Fprintf(conn, *msg)
	conn.Write(packet)
	if err != nil {
		fmt.Printf("Some error %v\n", err)

	}
	conn.Close()

}
func ReciveUDPmessage(host *string, port *int, comm chan <- string){
	connect := "CONNECT\r\n"
	packet := make([]byte, 1024)
	p := make([]byte, 1024)
	address := *host + ":" + strconv.Itoa(*port)
	conn, err := net.Dial("udp", address)
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
	fmt.Fprintf(conn, connect)
	conn.Write(packet)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
	}
	_ ,err = conn.Read(p)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
	}
	result := string(p[0:4])
	comm <- result

}
