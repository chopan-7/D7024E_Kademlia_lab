package labCode

import (
	"fmt"
	"net"
	"os"
)

// Message body is used to stora any information that we want to send in an RPC
type CLIMsgbody struct {
	Data []byte // Hashed key value
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
func CLIListen(ip string, port int) error {
	addr := net.ParseIP(ip)
	fmt.Println("CLI listener started")
	server := net.UDPAddr{
		Port: 10002,
		IP:   addr,
	}
	ServerConn, _ := net.ListenUDP("udp", &server)
	defer ServerConn.Close()
	buf := make([]byte, 1024)
	for {
		n, remoteaddr, _ := ServerConn.ReadFromUDP(buf)
		res := unmarshallData(buf[0:n])

		fmt.Println("Received RPC: ", res.RPC, "\nBody: ", res.Body, "\nFrom ", remoteaddr)

		responseMsg := cliresponseHandler(res)

		marshalledMsg := marshallData(responseMsg)
		sendResponse(ServerConn, remoteaddr, marshalledMsg)
		fmt.Println(responseMsg.RPC)
		if responseMsg.RPC == "exit" {
			fmt.Println("Exiting program...")
			os.Exit(1)
		}
	}
}

// The response handler will return the correct response based on which RPC it received
func cliresponseHandler(res Response) Response {
	switch res.RPC {
	case "ping":
		return cliPing(res.Body.Data)
	case "put":
		return cliPut(res.Body.Data)
	case "get":
		return cliGet(res.Body.Data)

		// TODO: if exit then shut down program
	case "exit":
		return exit()

	default:
		return Response{
			RPC: "Somethig went wrong: Invalid command",
		}
	}
}

// Will create a simple ping RPC response object
func cliPing(data []byte) Response {
	TestPing(string(data))
	responseMessage := Response{
		RPC: "ping",
		//ID:  resID,
	}
	return responseMessage
}

// Will create a simple ping RPC response object
func cliPut(data []byte) Response {
	responseMessage := Response{
		RPC: "put",
		//ID:  data,
	}
	return responseMessage
}

// Will create a simple ping RPC response object
func cliGet(data []byte) Response {
	responseMessage := Response{
		RPC: "get",
		// ID:  data,
	}
	return responseMessage
}

func exit() Response {

	responseMessage := Response{
		RPC: "exit",
		// ID:  data,
	}
	return responseMessage
}
