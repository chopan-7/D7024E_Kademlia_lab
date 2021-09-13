package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func client(ip string) error {

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

func server() error {
	IPAddr := GetOutboundIP()
	addr := net.ParseIP(IPAddr.String())
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

func sendResponse(conn *net.UDPConn, addr *net.UDPAddr) {
	_, err := conn.WriteToUDP([]byte("HELLOOOOOOOOOO"), addr)
	if err != nil {
		fmt.Printf("Couldn't send response %v", err)
	}
}

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func main() {
	app := cli.NewApp()
	app.Name = "Network CLI"
	app.Usage = "Lets you send command to the distributed network"

	app.Commands = []cli.Command{
		{
			Name:  "ping",
			Usage: "Will ping another node in the network given its IP adress",
			Action: func(c *cli.Context) error {
				client(c.Args()[0])
				return nil

			},
		}, {
			Name:  "start",
			Usage: "Will start a listener on this node",
			Action: func(c *cli.Context) error {
				server()
				return nil

			},
		},
	}

	// start our application
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	// connect := flag.String("connect", "", "IP address of process to join. If empty, go into listen mode.")
	// flag.Parse()

	// // If the connect flag is set, go into client mode.
	// if *connect != "" {
	// 	err := client(*connect)
	// 	if err != nil {
	// 		log.Println("Error:", errors.WithStack(err))
	// 	}
	// 	log.Println("Client done.")
	// 	return
	// }

	// // Else go into server mode.
	// err := server()
	// if err != nil {
	// 	log.Println("Error:", errors.WithStack(err))
	// }

	// log.Println("Server done.")
}
