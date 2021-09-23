package main

import (
	"fmt"
	lc "kademlia/labCode"
	"log"
	"net"
)

// Testing bootstrapnodes
// func main() {
// 	port := "10001"
// 	localIP := GetOutboundIP()
// 	localIPstr := localIP.String() + ":" + port // currentNode IP
// 	bnIP := "172.18.0.2:10001"                  // bootstrapNode IP

// 	// create a new node and init network with current node
// 	nn := lc.NewKademliaNode(localIPstr)
// 	network := &lc.Network{}
// 	network.Node = &nn

// 	fmt.Printf("\nIP: %s\n", localIP.String())
// 	// Join network if not a BootstrapNode
// 	if localIPstr != bnIP {
// 		// Join network by sending LookupContact to bootstrapNode
// 		bnContact := lc.NewContact(lc.NewKademliaID(lc.HashData(bnIP)), bnIP)
// 		nn.JoinNetwork(&bnContact, localIPstr)
// 		fmt.Printf("\nRoutingtable: %x\n", nn.Routingtable.FindClosestContacts(nn.Me.ID, 2))
// 	}

// 	go network.Listen()
// 	for {
// 	}

// }

func main() {
	ip := "192.168.10.165"
	//testnodes
	n0 := lc.NewKademliaNode(ip + ":10000")
	n1 := lc.NewKademliaNode(ip + ":10001")
	n2 := lc.NewKademliaNode(ip + ":10002")
	n3 := lc.NewKademliaNode(ip + ":10003")
	n4 := lc.NewKademliaNode(ip + ":10004")
	n5 := lc.NewKademliaNode(ip + ":10005")

	// adding contact to n0 routingtable
	n0.Routingtable.AddContact(n1.Me)
	n0.Routingtable.AddContact(n2.Me)
	n0.Routingtable.AddContact(n3.Me)
	n0.Routingtable.AddContact(n4.Me)
	n0.Routingtable.AddContact(n5.Me)

	n1.Routingtable.AddContact(n3.Me)
	n1.Routingtable.AddContact(n5.Me)
	n2.Routingtable.AddContact(n1.Me)
	n2.Routingtable.AddContact(n3.Me)
	n3.Routingtable.AddContact(n4.Me)
	n3.Routingtable.AddContact(n0.Me)
	n4.Routingtable.AddContact(n2.Me)
	n4.Routingtable.AddContact(n3.Me)
	n5.Routingtable.AddContact(n0.Me)
	n5.Routingtable.AddContact(n3.Me)
	n5.Routingtable.AddContact(n4.Me)

	// open Listener for all nodes
	net0 := lc.Network{&n0}
	net1 := lc.Network{&n1}
	net2 := lc.Network{&n2}
	net3 := lc.Network{&n3}
	net4 := lc.Network{&n4}
	net5 := lc.Network{&n5}

	go net0.Listen()
	go net1.Listen()
	go net2.Listen()
	go net3.Listen()
	go net4.Listen()
	go net5.Listen()

	closest := n0.LookupContact(n5.Me.ID)
	fmt.Printf("Closest from n0 to n5: %s\n", closest)
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
