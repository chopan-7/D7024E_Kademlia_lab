package labCode

import (
	"fmt"
	"net"
	//"github.com/urfave/cli"
)

func CLIListen(ip string, port int) {

	addr := net.ParseIP(ip)
	server := net.UDPAddr{
		Port: port,
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

/*
	old CLI saved if needed later
	--------------------------------------
	app := cli.NewApp()
	app.Name = "Network CLI"
	app.Usage = "Lets you send command to the distributed network"

	app.Commands = []cli.Command{
		{
			Name:  "ping",
			Usage: "Will ping another node in the network given its IP adress",
			Action: func(c *cli.Context) error {
				TestPing(c.Args()[0])
				return nil

			},
		}, {
			Name:  "start",
			Usage: "Will start a listener on this node",
			Action: func(c *cli.Context) error {
				ip := GetOutboundIP()
				Listen(ip.String(), 10001)
				return nil

			},
		},
	} */
