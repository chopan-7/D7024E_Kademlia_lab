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

	localIP := GetOutboundIP()
	localIPstr := localIP.String() + ":" + port // currentNode IP
	bnIP := "172.19.0.2:10001"                  // bootstrapNode IP

	fmt.Println("Your IP is:", localIPstr)

	bnID := lc.NewKademliaID(lc.HashData(bnIP))
	bnContact := lc.NewContact(bnID, bnIP)

	me := lc.NewKademliaNode(localIPstr)

	network := &lc.Network{}
	network.Node = &me
	network.Store = make(map[string][]byte)

	fmt.Printf("\nIP: %s\n", localIP.String())
	// Join network if not a BootstrapNode
	if localIPstr != bnIP {
		// Join network by sending LookupContact to bootstrapNode
		me.JoinNetwork(&bnContact)
		fmt.Printf("\nRoutingtable: %s\n", me.Routingtable.FindClosestContacts(me.Me.ID, 1))
	}

	go network.Listen()

	// test store function in bootstrapnode
	/* testData := []byte("hej hej thomas!")
	if localIPstr != bnIP {
		me.Store(testData)
	}
	time.Sleep(10 * time.Second)
	// test lookup data from other nodes
	if localIPstr != bnIP {
		lookupdata := me.LookupData(lc.HashData(string(testData)))
		fmt.Printf("\nLookup from %s and found %s\n", localIPstr, lookupdata)
	} */

	sch := lc.Scheduler{}
	sch.RunningRoutine()

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
