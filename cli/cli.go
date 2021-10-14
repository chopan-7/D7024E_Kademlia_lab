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
	"time"

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
		fmt.Printf("\nEnter command: ")
		text, _ := reader.ReadString('\n') //read from terminal
		text = strings.TrimRight(text, "\n")

		splitStr := strings.SplitN(text, " ", 2) //Split string at first space

		err := app.parser(splitStr)

		if err != nil {
			fmt.Printf("%s", err)
		}
	}
}

func (app *CLIApp) parser(cmd []string) error {
	var err error
	switch command := cmd[0]; command {
	case "put":
		err = app.put(cmd[1])
	case "get":
		args := strings.Fields(cmd[1]) // Checks so that the get command only receives one argument
		if len(args) > 1 {
			return errors.New("Too many arguments in command")
		}
		if len(cmd[1]) != 40 {
			return errors.New("Hash should be 40 characters long")
		}
		if !validHex(cmd[1]) {
			return errors.New("Hash contains illegal characters")
		}
		err = app.get(cmd[1])
	case "exit":
		err = app.exit()
	default:
		return errors.New("(" + command + ")" + " is not a valid command")
	}
	if err != nil {
		return err
	}
	return nil
}

func (app *CLIApp) put(data string) error {

	msg := CLIResponse{
		RPC: "put",
		Body: CLIMsgBody{
			Data: []byte(data),
		},
	}

	res, err := app.CLIMessageHandler(msg)

	if err != nil {
		return err
	}

	hashData := string(res.Body.Data)

	if res.RPC == msg.RPC {
		fmt.Printf("\nSuccesfully stored data with hash: %s\n", hashData)
	} else {
		fmt.Printf("\nWrong RPC in response message\n")
	}

	return nil
}

func (app *CLIApp) get(hash string) error {

	msg := CLIResponse{
		RPC: "get",
		Body: CLIMsgBody{
			Data: []byte(hash),
		},
	}

	res, err := app.CLIMessageHandler(msg)

	if err != nil {
		return err
	}

	receivedData := string(res.Body.Data)
	dataContact := res.Body.DataContact

	if res.RPC == msg.RPC {
		if res.Body.Data == nil {
			fmt.Printf("\nNo data found in the network with that hash\n")
		} else {
			fmt.Printf("\nFound data: %s\nFrom contact: %s\n", receivedData, &dataContact)
		}
	} else {
		fmt.Printf("\nSomething went wrong\n")
	}

	return nil
}

func (app *CLIApp) exit() error {
	msg := CLIResponse{
		RPC:  "exit",
		Body: CLIMsgBody{},
	}

	res, err := app.CLIMessageHandler(msg)

	if err != nil {
		return err
	}

	if res.RPC == "exit" {
		fmt.Printf("\nTerminating node with ip: %s\n", app.IP)
		os.Exit(0)
	} else {
		fmt.Printf("\nWrong RPC in response\n")
	}
	return nil
}

func (app *CLIApp) CLIMessageHandler(msg CLIResponse) (CLIResponse, error) {
	server := net.UDPAddr{
		Port: app.Port,
		IP:   net.ParseIP(app.IP),
	}

	marshalledMsg := marshallCLIData(msg)

	Conn, err := net.DialUDP("udp", nil, &server)

	if err != nil {
		return CLIResponse{}, errors.New("CLI: Failed to open connection to local ip:" + app.IP)
	}

	defer Conn.Close()

	timeDeadline := time.Now().Add(20 * time.Second)

	Conn.SetDeadline(timeDeadline)

	Conn.Write([]byte(marshalledMsg))
	buf := make([]byte, 1024)
	n, _, err := Conn.ReadFromUDP(buf)

	if err != nil {
		return CLIResponse{}, errors.New("Connection to cli listener timer has expired")
	}

	res := unmarshallCLIData(buf[0:n])

	return res, nil

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

func validHex(str string) bool {
	_, err := hex.DecodeString(str)
	if err != nil {
		return false
	}
	return true
}
