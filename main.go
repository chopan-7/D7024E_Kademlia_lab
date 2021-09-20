package main

import (
	"fmt"
	"kademlia/labCode"

	"log"
	"net"
)

func main() {

	// init node
	newNode := labCode.NewKademliaNode("192.168.0.1")

	// create 5 random ids
	randId1 := labCode.NewRandomKademliaID()
	randId2 := labCode.NewRandomKademliaID()
	randId3 := labCode.NewRandomKademliaID()
	randId4 := labCode.NewRandomKademliaID()
	randId5 := labCode.NewRandomKademliaID()
	// fmt.Println(randId1.String())
	// fmt.Println(randId2.String())
	// fmt.Println(randId3.String())
	// fmt.Println(randId4.String())
	// fmt.Println(randId5.String())

	// add 5 random nodes to the nodes to the routing table
	newNode.Routingtable.AddContact(labCode.NewContact(randId1, "192.168.0.2"))
	newNode.Routingtable.AddContact(labCode.NewContact(randId2, "192.168.0.3"))
	newNode.Routingtable.AddContact(labCode.NewContact(randId3, "192.168.0.4"))
	newNode.Routingtable.AddContact(labCode.NewContact(randId4, "192.168.0.5"))
	newNode.Routingtable.AddContact(labCode.NewContact(randId5, "192.168.0.6"))

	//
	//var targetpointer *labCode.KademliaID
	//targetpointer = randId5

	//fmt.Printf("Target: %x\n", targetpointer)

	// test FindClosestContact in Routingtable
	//closest := newNode.Routingtable.FindClosestContacts(targetpointer, 5)

	// test LookupContact

	//fmt.Printf("Closest: %x", closest)
	ip := GetOutboundIP()
	fmt.Print(ip)
	go labCode.Listen(ip.String(), 10001, *newNode)
	labCode.CLIListen(ip.String(), 10002)

}

// func main() {
// 	app := cli.NewApp()
// 	app.Name = "Network CLI"
// 	app.Usage = "Lets you send command to the distributed network"

// 	app.Commands = []cli.Command{
// 		{
// 			Name:  "ping",
// 			Usage: "Will ping another node in the network given its IP adress",
// 			Action: func(c *cli.Context) error {
// 				labCode.TestPing(c.Args()[0])
// 				return nil

// 			},
// 		}, {
// 			Name:  "start",
// 			Usage: "Will start a listener on this node",
// 			Action: func(c *cli.Context) error {
// 				ip := GetOutboundIP()
// 				labCode.Listen(ip.String(), 10001)
// 				return nil

// 			},
// 		},
// 	}

// 	// start our application
// 	err := app.Run(os.Args)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}
