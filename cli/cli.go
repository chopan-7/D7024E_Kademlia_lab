package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	lc "kademlia/labCode"
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

// Message body is used to store any information that we want to send in an RPC
type CLIMsgBody struct {
	Data        []byte // Hashed key value
	DataContact lc.Contact
}

type CLIResponse struct {
	RPC  string // String representing what kind of rpc the message is
	Body CLIMsgBody
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
	fmt.Printf("Welcome to kademlia network!\nAvailable commands:\nput <string> (stores data on k closest nodes to hash)	\nget <hash> (fetches data object with this hash if it is stored in the network)	\nexit (terminates this node)\n")

	for {
		fmt.Printf("\nEnter command:")
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
		return errors.New("Too many arguments")
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

func (app *CLIApp) put(data string) error {

	req := CLIResponse{
		RPC: "put",
		Body: CLIMsgBody{
			Data: []byte(data),
		},
	}

	marshalledReq := marshallCLIData(req)

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
	n, _, _ := Conn.ReadFromUDP(buf)
	rec := unmarshallCLIData(buf[0:n])

	hashData := string(rec.Body.Data)

	if rec.RPC == req.RPC {
		fmt.Printf("\nSuccesfully stored data with hash: %s\n", hashData)
	} else {
		fmt.Println("Something went wrong")
	}

	return nil
}

func (app *CLIApp) get(hash string) error {

	req := CLIResponse{
		RPC: "get",
		Body: CLIMsgBody{
			Data: []byte(hash),
		},
	}

	marshalledReq := marshallCLIData(req)

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
	n, _, _ := Conn.ReadFromUDP(buf)
	rec := unmarshallCLIData(buf[0:n])

	receivedData := string(rec.Body.Data)
	dataContact := rec.Body.DataContact

	if rec.RPC == req.RPC {
		if rec.Body.Data == nil {
			fmt.Printf("\nNo data found in the network with that hash\n")
		} else {
			fmt.Printf("\nFound data: %s\nFrom contact: %s\n", receivedData, &dataContact)
		}
	} else {
		fmt.Println("Something went wrong")
	}

	return nil
}

func (app *CLIApp) exit() error {
	res := CLIResponse{
		RPC:  "exit",
		Body: CLIMsgBody{},
	}

	marshalledReq := marshallCLIData(res)

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
	n, _, _ := Conn.ReadFromUDP(buf)
	rec := unmarshallCLIData(buf[0:n])
	if rec.RPC == "exit" {
		fmt.Printf("\nTerminating node with ip: %s\n", app.IP)
		os.Exit(1)
	}
	return nil
}

// Will marshall the response object into json
func marshallCLIData(data CLIResponse) []byte {
	marshalledData, _ := json.Marshal(data)
	return marshalledData
}

// Will unmarshall the byte stream into json
func unmarshallCLIData(data []byte) CLIResponse {
	var unmarshalledData CLIResponse
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
