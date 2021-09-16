package main

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/pkg/errors"
	//"github.com/urfave/cli"
)

/*
	CLI module is detached from the rest of the service hence it
	has its own main package and function

	The CLI will only communicate outwards to the running kademlia service
	over port 10002 and not recieve any incoming messages.
*/
func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(`Welcome to hecnet!: 
Available commands:
ping <IP-address> (Almost)
put <file> 	(coming soon)
ger <hash> 	(coming soon)
exit 		(coming soon)
`)
	for {
		fmt.Print("Enter command:")
		text, _ := reader.ReadString('\n')
		text = strings.TrimRight(text, "\n")
		if text == "q" || text == "quit" || text == "exit" {
			fmt.Println("Bye!")
			os.Exit(1)
		}

		words := strings.Fields(text)

		err := parser(words)

		if err != nil {
			fmt.Println(errors.Wrap(err, "Failed to parse command"))
		}
	}
}

func parser(cmd []string) error {
	port := 10002
	if len(cmd) != 2 {
		return errors.New("Invalid command!")
	}
	switch command := cmd[0]; command {
	case "ping":
		fmt.Println("we shall run a", command, "command!")
		ipAddress := []byte(cmd[1])
		ping(string(ipAddress), port)
	case "put":
		fmt.Println("we shall run a", command, "command!")
		put(cmd[1], port)
	case "get":
		fmt.Println("we shall run a", command, "command!")
		// Perform store command
	case "exit":
		fmt.Println("we shall run a", command, "command!")
		// Perform store command
	default:
		return errors.New(command + "is not a valid commaand...")
	}
	return nil
}

func ping(ip string, port int) error {
	addr := net.ParseIP(ip)
	fmt.Println(addr)

	server := net.UDPAddr{
		Port: port,
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

func put(filePath string, port int) {
	dat, err := os.ReadFile(filePath)

	if err != nil {
		errors.Wrap(err, "Failed to read from file at :"+filePath)
	}
	fmt.Print(string(dat))
	hashFile := sha1.New()
	hashFile.Write([]byte(dat))
}

func get(data string) {
	// TODO
}

func exit(data []byte) {
	// TODO
}
