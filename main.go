package main

import (
	"kademlia/labCode"
	"log"
	"net"
	"os"

	"github.com/urfave/cli"
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
	con := labCode.Contact{
		ID:      labCode.NewRandomKademliaID(),
		Address: "192.168.10.114:10002",
	}
	network := labCode.Network{}

	app := cli.NewApp()
	app.Name = "Network CLI"
	app.Usage = "Lets you send command to the distributed network"

	app.Commands = []cli.Command{
		{
			Name:  "ping",
			Usage: "Will ping another node in the network given its IP adress",
			Action: func(c *cli.Context) error {
				network.SendPingMessage(&con)
				return nil

			},
		}, {
			Name:  "start",
			Usage: "Will start a listener on this node",
			Action: func(c *cli.Context) error {
				ip := GetOutboundIP()
				labCode.Listen(ip.String(), 10001, *newNode)
				return nil

			},
		},
		{
			Name:  "lookup",
			Usage: "Uses the find_node rpc",
			Action: func(c *cli.Context) error {
				network.SendFindContactMessage(&con, labCode.NewRandomKademliaID())
				return nil

			},
		},
		{
			Name:  "find",
			Usage: "Uses the find_data rpc to find data in any node with matching hash key",
			Action: func(c *cli.Context) error {
				network.SendFindDataMessage(&con, "1123421342342343515315511525151415445235422525435243245345232534245324534523453224534523452324534523")
				return nil

			},
		},
		{
			Name:  "store",
			Usage: "Uses the store_data rpc",
			Action: func(c *cli.Context) error {
				network.SendStoreMessage(&con, []byte("123123"))
				return nil

			},
		},
	}

	// start our application
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
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
