package main

import (
	"fmt"
	lc "kademlia/labCode"
	"log"
	"net"
	"strings"
)

// Testing bootstrapnodes
func main() {

	port := "10001"

	localIP := GetOutboundIP()
	localIPstr := localIP.String() + ":" + port // currentNode IP
	bnIP := GetbnIP(localIP.String())           // bootstrapNode IP

	fmt.Println("Your IP is:", localIPstr)

	bnID := lc.NewKademliaID(lc.HashData(bnIP))
	bnContact := lc.NewContact(bnID, bnIP)

	me := lc.NewKademliaNode(localIPstr)

	network := &lc.Network{}
	network.Node = &me

	fmt.Printf("\nIP: %s\n", localIP.String())
	// Join network if not a BootstrapNode
	if localIPstr != bnIP {
		// Join network by sending LookupContact to bootstrapNode
		me.JoinNetwork(&bnContact)
	}

	// goroutine for data expiration runs every 20 seconds
	go network.Node.CheckDataExpired(20)

	go network.Listen()
	cliListener := &lc.CLI{&me, network}
	cliListener.CLIListen()
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

func GetbnIP(ipString string) string {
	stringList := strings.Split(ipString, ".")
	value := stringList[1]
	bnIP := "172." + value + ".0.2:10001"
	return bnIP
}
