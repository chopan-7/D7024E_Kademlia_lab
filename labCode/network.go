package labCode

import (
	"fmt"
	"net"

	"github.com/pkg/errors"
)

type Network struct {
}

func Listen(ip string, port int) {
	addr := net.ParseIP(ip)
	fmt.Println(addr)
	server := net.UDPAddr{
		Port: 10001,
		IP:   addr,
	}
	ServerConn, _ := net.ListenUDP("udp", &server)
	defer ServerConn.Close()
	buf := make([]byte, 1024)
	for {
		n, remoteaddr, _ := ServerConn.ReadFromUDP(buf)
		fmt.Println("Received ", string(buf[0:n]), " from ", remoteaddr)
		sendResponse(ServerConn, remoteaddr)
	}
}

func (network *Network) SendPingMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}

func TestPing(ip string) error {
	addr := net.ParseIP(ip)
	fmt.Println(addr)

	server := net.UDPAddr{
		Port: 10001,
		IP:   addr,
	}

	Conn, err := net.DialUDP("udp", nil, &server)

	if err != nil {
		return errors.Wrap(err, "Client: Failed to open connection to "+ip)
	}

	defer Conn.Close()
	Conn.Write([]byte("hello"))
	buf := make([]byte, 1024)
	n, remoteaddr, _ := Conn.ReadFromUDP(buf)
	fmt.Println("Received ", string(buf[0:n]), " from ", remoteaddr)

	return nil
}

func sendResponse(conn *net.UDPConn, addr *net.UDPAddr) {
	_, err := conn.WriteToUDP([]byte("HELLOOOOOOOOOO"), addr)
	if err != nil {
		fmt.Printf("Couldn't send response %v", err)
	}
}
