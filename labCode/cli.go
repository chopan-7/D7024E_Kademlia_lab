package labCode

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/pkg/errors"
	//"github.com/urfave/cli"
)

func CLI() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(`Welcome to hecnet!: 
Available commands:
ping <IP-address> (Almost)`)
	text, _ := reader.ReadString('\n')

	words := strings.Fields(text)

	err := parser(words)

	if err != nil {
		fmt.Println(errors.Wrap(err, "Failed to parse command"))
	}

	/* app := cli.NewApp()
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

func parser(cmd []string) error {
	fmt.Println(len(cmd))
	if len(cmd) != 2 {
		return errors.New("Invalid command!")
	}
	switch command := cmd[0]; command {
	case "ping":
		fmt.Println("we shall run a", command, "command!")
		// Perform ping-command
	case "lookup":
		fmt.Println("we shall run a", command, "command!")
		// Perform lookup-command
	default:

		return errors.New(command + "is not a valid commaand...")
	}
	return nil
}
