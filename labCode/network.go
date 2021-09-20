package labCode

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"

	"github.com/pkg/errors"
)

type Network struct {
	node Kademlia
}

// Message body is used to stora any information that we want to send in an RPC
type Msgbody struct {
	Nodes []Contact   // List of contact nodes
	Data  []byte      // Hashed key value
	KadID *KademliaID // For find_node RPCs the id we are looking for
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

func (network *Network) SendPingMessage(contact *Contact) error {

	udpAddr := GetUDPAddrFromContact(contact)

	msg := Response{
		RPC: "ping",
		ID:  NewRandomKademliaID(),
	}

	marshalledMsg := marshallData(msg)

	Conn, err := net.DialUDP("udp", nil, &udpAddr)

	if err != nil {
		return errors.Wrap(err, "Client: Failed to open connection to "+udpAddr.IP.String())
	}

	defer Conn.Close()
	Conn.Write([]byte(marshalledMsg))
	buf := make([]byte, 1024)

	for {
		n, remoteaddr, _ := Conn.ReadFromUDP(buf)
		rec := unmarshallData(buf[0:n])

		if Validate(msg, rec) {
			fmt.Println("Received RPC: ", rec.RPC, "\nWith RPC ID: ", rec.ID, "\nBody: ", rec.Body, "\nFrom ", remoteaddr)
			break
		}
	}
	return nil
}

func (network *Network) SendFindContactMessage(contact *Contact, kadID *KademliaID) error {
	// TODO

	udpAddr := GetUDPAddrFromContact(contact)

	body := Msgbody{
		KadID: kadID,
	}

	msg := Response{
		RPC:  "find_node",
		ID:   NewRandomKademliaID(),
		Body: body,
	}

	marshalledMsg := marshallData(msg)

	Conn, err := net.DialUDP("udp", nil, &udpAddr)

	if err != nil {
		return errors.Wrap(err, "Client: Failed to open connection to "+udpAddr.IP.String())
	}

	defer Conn.Close()
	Conn.Write([]byte(marshalledMsg))
	buf := make([]byte, 1024)

	n, remoteaddr, _ := Conn.ReadFromUDP(buf)
	rec := unmarshallData(buf[0:n])

	fmt.Println("msg id: ", msg.ID, "\nrec id: ", rec.ID, "\nrec rpc : ", rec.RPC)

	if Validate(msg, rec) {
		fmt.Println("Received RPC: ", rec.RPC, "\nWith RPC ID: ", rec.ID, "\nBody: ", rec.Body, "\nFrom ", remoteaddr)
	}

	return nil

}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(contact *Contact, data []byte) error {
	// TODO
	udpAddr := GetUDPAddrFromContact(contact)

	body := Msgbody{
		Data: data,
	}

	msg := Response{
		RPC:  "store_data",
		ID:   NewRandomKademliaID(),
		Body: body,
	}

	marshalledMsg := marshallData(msg)

	Conn, err := net.DialUDP("udp", nil, &udpAddr)

	if err != nil {
		return errors.Wrap(err, "Client: Failed to open connection to "+udpAddr.IP.String())
	}

	defer Conn.Close()
	Conn.Write([]byte(marshalledMsg))
	buf := make([]byte, 1024)

	n, remoteaddr, _ := Conn.ReadFromUDP(buf)
	rec := unmarshallData(buf[0:n])

	fmt.Println("msg id: ", msg.ID, "\nrec id: ", rec.ID, "\nrec rpc : ", rec.RPC)

	if Validate(msg, rec) {
		fmt.Println("Received RPC: ", rec.RPC, "\nWith RPC ID: ", rec.ID, "\nBody: ", rec.Body, "\nFrom ", remoteaddr)
	}

	return nil
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

	fmt.Println("msg id: ", msg.ID, "\nrec id: ", rec.ID, "\nrec rpc : ", rec.RPC)

	if rec.RPC == "find_node" && rec.ID == msg.ID {
		fmt.Println("Received RPC: ", rec.RPC, "\nWith RPC ID: ", rec.ID, "\nBody: ", rec.Body, "\nFrom ", remoteaddr)
	}

	return nil

}

func TestPing(ip string) error {

	// Dummy data for test sending udp messages
	msg := Response{
		RPC: "ping",
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

	for {
		n, remoteaddr, _ := Conn.ReadFromUDP(buf)
		rec := unmarshallData(buf[0:n])

		if Validate(msg, rec) {
			fmt.Println("Received RPC: ", rec.RPC, "\nWith RPC ID: ", rec.ID, "\nBody: ", rec.Body, "\nFrom ", remoteaddr)
			break
		}
	}

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
	case "find_data":
	case "store_data":
	}

	return Response{}

}

// Will validate if a response to an RPC has correct RPC name and correct RPC ID
func Validate(msg Response, res Response) bool {
	if (msg.RPC == res.RPC) && msg.ID.Equals(res.ID) {
		return true
	} else {
		return false
	}
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

	// Here is just dummy data that is the contact list

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

func GetUDPAddrFromContact(contact *Contact) net.UDPAddr {
	addr, port, _ := net.SplitHostPort(contact.Address)
	netAddr := net.ParseIP(addr)
	intPort, _ := strconv.Atoi(port)
	receiver := net.UDPAddr{
		IP:   netAddr,
		Port: intPort,
	}
	return receiver
}
