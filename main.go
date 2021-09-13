package main

import (
	"kademlia/labCode"
	"log"
	"net"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Network CLI"
	app.Usage = "Lets you send command to the distributed network"

	app.Commands = []cli.Command{
		{
			Name:  "ping",
			Usage: "Will ping another node in the network given its IP adress",
			Action: func(c *cli.Context) error {
				labCode.TestPing(c.Args()[0])
				return nil

			},
		}, {
			Name:  "start",
			Usage: "Will start a listener on this node",
			Action: func(c *cli.Context) error {
				ip := GetOutboundIP()
				labCode.Listen(ip.String(), 10001)
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

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
