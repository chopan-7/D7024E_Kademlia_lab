package labCode

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type CLI struct {
	Node *Kademlia
	Net  *Network
}

// Message body is used to stora any information that we want to send in an RPC
type CLIMsgbody struct {
	Data        []byte // Hashed key value
	DataContact Contact
}

type CLIResponse struct {
	RPC  string // String representing what kind of rpc the message is
	Body CLIMsgbody
}

/*
	CLIListener will listen for cli input on port 10002

	The cli input is interpreted and thereafter relevant code is executed in the network module,
	a response is generated and sent back to the CLI.

*/
func (this *CLI) CLIListen() error {
	addr := GetUDPAddrFromContact(&this.Node.Me)
	port := 10002
	fmt.Printf("\nCLI listener started on port %d\n", port)
	server := net.UDPAddr{
		Port: port,
		IP:   addr.IP,
	}

	ServerConn, _ := net.ListenUDP("udp", &server)
	defer ServerConn.Close()
	buf := make([]byte, 1024)
	for {
		n, remoteaddr, _ := ServerConn.ReadFromUDP(buf)
		res := unmarshallCLIData(buf[0:n])

		responseMsg := this.cliresponseHandler(res)

		marshalledMsg := marshallCLIData(responseMsg)
		sendResponse(ServerConn, remoteaddr, marshalledMsg)
		if responseMsg.RPC == "exit" {
			fmt.Println("Exiting program...")
			os.Exit(0)
		}
	}
}

// The response handler will return the correct response based on which RPC it received
func (this *CLI) cliresponseHandler(res CLIResponse) CLIResponse {
	switch res.RPC {
	case "put":
		return this.cliPut(res.Body.Data)
	case "get":
		return this.cliGet(res.Body.Data)
	case "exit":
		return this.exit()

	default:
		return CLIResponse{
			RPC: "Something went wrong: Invalid command",
		}
	}
}

// Will create a simple ping RPC response object
func (this *CLI) cliPut(data []byte) CLIResponse {
	this.Node.Store(data)
	hash := HashData(string(data))

	responseBody := CLIMsgbody{
		Data: []byte(hash),
	}
	responseMessage := CLIResponse{
		RPC:  "put",
		Body: responseBody,
	}
	return responseMessage
}

// Will create a simple ping RPC response object
func (this *CLI) cliGet(data []byte) CLIResponse {
	hash := string(data)
	receivedData, contact := this.Node.LookupData(hash)

	body := CLIMsgbody{
		Data:        receivedData,
		DataContact: contact,
	}

	responseMessage := CLIResponse{
		RPC:  "get",
		Body: body,
	}
	return responseMessage
}

func (this *CLI) exit() CLIResponse {

	responseMessage := CLIResponse{
		RPC: "exit",
	}
	return responseMessage
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
