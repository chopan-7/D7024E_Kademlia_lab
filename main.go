package main

import (
	"fmt"
	"kademlia/labCode"

	"log"
	"net"
)

func main() {

	ip := GetOutboundIP()
	fmt.Println("Listening to port:", 10001)
	go labCode.Listen(ip.String(), 10001)
	labCode.CLIListen(ip.String(), 10002)

}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
