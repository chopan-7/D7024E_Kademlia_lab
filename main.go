package main

import (
	"fmt"
	lc "kademlia/labCode"
	"log"
	"net"
)

// Testing bootstrapnodes
func main() {

	port := "10001"
	cliPort := 10002

	localIP := GetOutboundIP()
	localIPstr := localIP.String() + ":" + port // currentNode IP
	bnIP := "172.20.0.2:10001"                  // bootstrapNode IP


	fmt.Println("Your IP is:", localIPstr)

	bnIP := "172.18.0.2:10001" // bootstrapNode IP
	bnID := lc.NewKademliaID(lc.HashData(bnIP))
	bnContact := lc.NewContact(bnID, bnIP)

	me := lc.NewKademliaNode(localIPstr)
	me.JoinNetwork(&bnContact)

	network := &lc.Network{}
	network.Node = &me

	fmt.Printf("\nIP: %s\n", localIP.String())
	// Join network if not a BootstrapNode
	if localIPstr != bnIP {
		// Join network by sending LookupContact to bootstrapNode
		bnContact := lc.NewContact(lc.NewKademliaID(lc.HashData(bnIP)), bnIP)
		nn.JoinNetwork(&bnContact)
		fmt.Printf("\nRoutingtable: %x\n", nn.Routingtable.FindClosestContacts(nn.Me.ID, 2))
	}

	go network.Listen()
	lc.CLIListen(localIPstr, cliPort)
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
