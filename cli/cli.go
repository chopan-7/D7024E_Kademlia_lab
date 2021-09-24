package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/pkg/errors"
)

type CLIApp struct {
	IP   string
	Port int
}

// Message body is used to stora any information that we want to send in an RPC
type Msgbody struct {
	IP   string
	Data []byte // Hashed key value
}

type Response struct {
	RPC  string // String representing what kind of rpc the message is
	Body Msgbody
}

/*
	CLI module is detached from the rest of the service hence it
	has its own main package and function

	The CLI will only communicate outwards to the running kademlia service
	over port 10002 and not recieve any incoming messages.
*/
func main() {
	app := &CLIApp{GetOutboundIP().String(), 10002}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(`Welcome to hecnet!: 
		Available commands:
		put <file> 	(coming soon)
		ger <hash> 	(coming soon)
		exit 		terminates program
	`)

	for {
		fmt.Print("Enter command:")
		text, _ := reader.ReadString('\n') //read from terminal
		text = strings.TrimRight(text, "\n")

		words := strings.Fields(text)

		err := app.parser(words)

		if err != nil {
			fmt.Println(errors.Wrap(err, "Failed to parse command"))
		}
	}
}

func (app *CLIApp) parser(cmd []string) error {
	if len(cmd) > 2 {
		return errors.New("Invalid command!")
	}
	switch command := cmd[0]; command {
	case "put":
		app.put(cmd[1])
	case "get":
		app.get(cmd[1])
	case "exit":
		app.exit()
	default:
		return errors.New(command + "is not a valid command...")
	}
	return nil
}

func (app *CLIApp) put(filePath string) error {
	dat, err := os.ReadFile(filePath)

	if err != nil {
		errors.Wrap(err, "Failed to read from file at:"+filePath)
	}

	hashFile := HashData(string(dat))

	req := Response{
		RPC: "put",
		Body: Msgbody{
			Data: []byte(hashFile),
		},
	}

	marshalledReq := marshallData(req)

	server := net.UDPAddr{
		Port: app.Port,
		IP:   net.ParseIP(app.IP),
	}

	Conn, err := net.DialUDP("udp", nil, &server)

	if err != nil {
		return errors.Wrap(err, "CLI: Failed to open connection to "+app.IP)
	}

	defer Conn.Close()

	Conn.Write([]byte(marshalledReq))
	buf := make([]byte, 1024)
	n, remoteaddr, _ := Conn.ReadFromUDP(buf)
	rec := unmarshallData(buf[0:n])

	fmt.Println("Received RPC: ", rec.RPC, "\nBody: ", rec.Body, "\nFrom ", remoteaddr)

	return nil
}

func (app *CLIApp) get(data string) error {

	hashFile := HashData(data)

	res := Response{
		RPC: "get",
		Body: Msgbody{
			Data: []byte(hashFile),
		},
	}

	marshalledReq := marshallData(res)

	server := net.UDPAddr{
		Port: app.Port,
		IP:   net.ParseIP(app.IP),
	}

	Conn, err := net.DialUDP("udp", nil, &server)

	if err != nil {
		return errors.Wrap(err, "CLI: Failed to open connection to local ip:"+app.IP)
	}

	defer Conn.Close()

	Conn.Write([]byte(marshalledReq))
	buf := make([]byte, 1024)
	n, remoteaddr, _ := Conn.ReadFromUDP(buf)
	rec := unmarshallData(buf[0:n])

	fmt.Println("Received RPC: ", rec.RPC, "\nBody: ", rec.Body, "\nFrom ", remoteaddr)

	return nil
}

func (app *CLIApp) exit() error {
	res := Response{
		RPC:  "exit",
		Body: Msgbody{},
	}

	marshalledReq := marshallData(res)

	fmt.Println(app.IP)

	server := net.UDPAddr{
		Port: app.Port,
		IP:   net.ParseIP(app.IP),
	}

	Conn, err := net.DialUDP("udp", nil, &server)

	if err != nil {
		return errors.Wrap(err, "CLI: Failed to open connection to local ip:"+app.IP)
	}

	defer Conn.Close()

	Conn.Write([]byte(marshalledReq))
	buf := make([]byte, 1024)
	n, remoteaddr, _ := Conn.ReadFromUDP(buf)
	rec := unmarshallData(buf[0:n])

	fmt.Println("Received RPC: ", rec.RPC, "\nBody: ", rec.Body, "\nFrom ", remoteaddr)
	os.Exit(1)
	return nil
}

// Will marshall the response object into json
func marshallData(data Response) []byte {
	marshalledData, _ := json.Marshal(data)
	return marshalledData
}

func unmarshallData(data []byte) Response {
	var unmarshalledData Response
	json.Unmarshal([]byte(data), &unmarshalledData)
	return unmarshalledData
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

// Helper function for hashing data returing hexstring
func HashData(data string) (hashString string) {
	newHash := sha1.New()
	newHash.Write([]byte(data))
	hashString = hex.EncodeToString(newHash.Sum(nil))
	return
}
