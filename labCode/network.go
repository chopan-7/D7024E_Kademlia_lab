package labCode

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/pkg/errors"
)

type Network struct {
}

// Message body is used to stora any information that we want to send in an RPC
type Msgbody struct {
	Nodes []Contact // List of contact nodes
	Data  []byte    // Hashed key value
}

type Response struct {
	RPC  string      // String representing what kind of rpc the message is
	ID   *KademliaID // A randomly generated kademlia id to identify the ping
	Body Msgbody
}

func Listen(ip string, port int) {
	addr := net.ParseIP(ip)
	fmt.Println(addr)
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

		fmt.Println("Received RPC: ", res.RPC, "\nWith RPC ID: ", res.ID, "\nBody: ", res.Body, "\nFrom ", remoteaddr)

		responseMsg := responseHandler(res)
		marshalledMsg := marshallData(responseMsg)
		sendResponse(ServerConn, remoteaddr, marshalledMsg)
	}
}

func (network *Network) SendPingMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}

func TestFindNode(ip string) error {
	msg := Response{
		RPC: "find_node",
		ID:  NewRandomKademliaID(),
	}

	marshalledMsg := marshallData(msg)

	addr := net.ParseIP(ip)
	fmt.Println(addr)

	server := net.UDPAddr{
		Port: 10002,
		IP:   addr,
	}

	Conn, err := net.DialUDP("udp", nil, &server)

	if err != nil {
		return errors.Wrap(err, "Client: Failed to open connection to "+ip)
	}

	defer Conn.Close()
	Conn.Write([]byte(marshalledMsg))
	buf := make([]byte, 1024)

	n, remoteaddr, _ := Conn.ReadFromUDP(buf)
	rec := unmarshallData(buf[0:n])

	fmt.Println("Received RPC: ", rec.RPC, "\nWith RPC ID: ", rec.ID, "\nBody: ", rec.Body, "\nFrom ", remoteaddr)

	return nil

}

func TestPing(ip string) error {

	// Dummy data for test sending udp messages
	con := NewContact(NewRandomKademliaID(), "123,123,123,123")
	conList := make([]Contact, 1)
	conList[0] = con

	testbody := Msgbody{
		Nodes: conList,
		Data:  nil,
	}

	msg := Response{
		RPC:  "ping",
		ID:   NewRandomKademliaID(),
		Body: testbody,
	}

	marshalledMsg := marshallData(msg)

	addr := net.ParseIP(ip)
	fmt.Println(addr)

	server := net.UDPAddr{
		Port: 10002,
		IP:   addr,
	}

	Conn, err := net.DialUDP("udp", nil, &server)

	if err != nil {
		return errors.Wrap(err, "Client: Failed to open connection to "+ip)
	}

	defer Conn.Close()
	Conn.Write([]byte(marshalledMsg))
	buf := make([]byte, 1024)

	n, remoteaddr, _ := Conn.ReadFromUDP(buf)
	rec := unmarshallData(buf[0:n])

	fmt.Println("Received RPC: ", rec.RPC, "\nWith RPC ID: ", rec.ID, "\nBody: ", rec.Body, "\nFrom ", remoteaddr)

	return nil
}

// Will marshall the response object into json
func marshallData(data Response) []byte {
	marshalledData, _ := json.Marshal(data)
	return marshalledData
}

// Will unmarshall the byte stream into json
func unmarshallData(data []byte) Response {
	var unmarshalledData Response
	json.Unmarshal([]byte(data), &unmarshalledData)
	return unmarshalledData
}

// Given a connection channel, ip address of the node that sent a message
// and a response object sends back the queried data to the initiator
func sendResponse(conn *net.UDPConn, addr *net.UDPAddr, responseMsg []byte) {
	_, err := conn.WriteToUDP([]byte(responseMsg), addr)
	if err != nil {
		fmt.Printf("Couldn't send response %v", err)
	}
}

// The response handler will return the correct response based on which RPC it received
func responseHandler(res Response) Response {
	switch res.RPC {
	case "ping":
		return createPingResponse(res.ID)
	case "find_node":
		return createFindNodeResponse(res.ID)
	}

	return Response{}

}

// Will create a simple ping RPC response object
func createPingResponse(resID *KademliaID) Response {
	responseMessage := Response{
		RPC: "ping",
		ID:  resID,
	}
	return responseMessage
}

func createFindNodeResponse(resID *KademliaID) Response {
	// Need own node in order to reach own routingtable to find K closest nodes to the target.
	// Something like contacts := ownNode.ownRoutingtable.FindClosestContacts

	// Here is just dummy data that is the contact list obtained from above

	con := NewContact(NewRandomKademliaID(), "111.111.111.111")
	con2 := NewContact(NewRandomKademliaID(), "222.222.222.222")
	con3 := NewContact(NewRandomKademliaID(), "333.333.333.333")
	con4 := NewContact(NewRandomKademliaID(), "444.444.444.444")
	con5 := NewContact(NewRandomKademliaID(), "555.555.555.555")
	conList := make([]Contact, 5)
	conList[0] = con
	conList[1] = con2
	conList[2] = con3
	conList[3] = con4
	conList[4] = con5

	resBody := Msgbody{
		Nodes: conList,
	}

	responseMessage := Response{
		RPC:  "find_node",
		ID:   resID,
		Body: resBody,
	}
	return responseMessage
}
