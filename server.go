package main

import (
	"net"
	"fmt"
	"github.com/akamensky/argparse"
	"os"
)


const (
	MESSAGE_SIZE = 1024
	CONNECT_MSG = "CONNECT\r\n"
	DISCONNECT_MSG = "DISCONNECT\r\n"
)
func StartServer(host *string, pin, pout *int){
	fmt.Println("Started Server execution")
	pin_addr := net.UDPAddr{Port: *pin, IP: net.ParseIP(*host)}
	pout_addr :=net.UDPAddr{Port: *pout, IP: net.ParseIP(*host)}
	conn_in, err := net.ListenUDP("udp", &pin_addr)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
	}
	conn_out, err := net.ListenUDP("udp", &pout_addr)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
	}
	go StartServerListener(conn_in, conn_out)
}


func StartServerListener(conn_in *net.UDPConn, conn_out *net.UDPConn){

	p := make([]byte, MESSAGE_SIZE)
	comm := make(chan []byte)
	go outputListener(conn_out, comm)
	for {
		_, _, err := conn_in.ReadFromUDP(p)
		if err != nil {
			fmt.Printf("Couldn't send response %v", err)
			continue
		}
		comm <- p
	}
}

func outputListener(ser *net.UDPConn, comm <- chan []byte) {
	packet := make([]byte, MESSAGE_SIZE)
	sessions := make(map[string]chan bool)
	for {
		n, addr, err  := ser.ReadFromUDP(packet)
		if err != nil {
			fmt.Printf("Couldn't get a message %v\n", err)
		}
		msg := string(packet[0:n])
		ipaddr := addr.String()
		switch msg {
		case CONNECT_MSG:
			fmt.Println("Connect", ipaddr)
			sessions[ipaddr] = make(chan bool)
			go broadcast(ser, addr, comm, sessions[ipaddr])

		case DISCONNECT_MSG:
			fmt.Println("Disconect ", ipaddr)
			sessions[addr.String()] <- false
			delete(sessions, ipaddr)
		}


	}

}
func broadcast(ser *net.UDPConn, addr *net.UDPAddr, comm <- chan []byte, quit <- chan bool)  {
	for {
		select {
		case msg:= <-comm:
			_,err := ser.WriteToUDP(msg, addr)
			if err != nil {
				fmt.Printf("Couldn't send response %v", err)
				}
		case <- quit:
			break
		}
	}
}



func cliOutput() (host* string, pout* int, pin* int){
	parser := argparse.NewParser("server", "Simple UDP server written in GO")
	host = parser.String("H", "host", &argparse.Options{Required: true, Help: "Specify ip address"})
	pout = parser.Int("o", "output", &argparse.Options{Required: true, Help: "Specify output port"})
	pin = parser.Int("i", "input", &argparse.Options{Required: true, Help: "Specify input port"})
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}
	return
}

func main() {
	host, pout, pin := cliOutput()
	switch {
	case *host == "":
		fmt.Println("Specify Servers ip address")
		os.Exit(1)
	case *pout == 0:
		fmt.Println("Specify output port")
		os.Exit(1)
	case *pin == 0:
		fmt.Println("Specify input port")
		os.Exit(1)
	default:
		// start the server
		quit := make(chan bool)
		go StartServer(host, pin, pout)
		for {
			select {
			case <-quit:
				break
			}
		}

		fmt.Println("Server Execution Finished")

	}
}
