package labCode

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"

	"github.com/pkg/errors"
)

type Network struct {
	Node *Kademlia
}

// Message body is used to stora any information that we want to send in an RPC
type Msgbody struct {
	Nodes []Contact   // List of contact nodes
	Data  []byte      // Hashed key value
	KadID *KademliaID // For find_node RPCs the id we are looking for
	Hash  string      // Used for find_value rpcs to find the node that has stored data for the hashed key value
}

type Response struct {
	RPC            string      // String representing what kind of rpc the message is
	ID             *KademliaID // A randomly generated kademlia id to identify the ping
	SendingContact *Contact
	Body           Msgbody
}

// Will open up a UDP listener on itself with a given port.
// If a message is received on the listener it will use the response handler to do the
// correct operations
func (network *Network) Listen() {
	server := GetUDPAddrFromContact(&network.Node.Me)
	fmt.Printf("Listening on: %s:%d\n", server.IP, server.Port)

	ServerConn, _ := net.ListenUDP("udp", &server)
	defer ServerConn.Close()
	buf := make([]byte, 1024)
	for {
		n, remoteaddr, _ := ServerConn.ReadFromUDP(buf)
		res := unmarshallData(buf[0:n])

		//fmt.Println("Received RPC: ", res.RPC, "\nWith RPC ID: ", res.ID, "\nBody: ", res.Body, "\nFrom ", remoteaddr)
		network.Node.Routingtable.AddContact(*res.SendingContact)
		responseMsg := network.responseHandler(res, *network.Node)
		marshalledMsg := marshallData(responseMsg)

		sendResponse(ServerConn, remoteaddr, marshalledMsg)
	}
}

// All of the RPC are sent using these message functions. The functions will create a Response
// object with the data it has to send. Then the MessageHandler function will send and retreive the response
// from the other contact

// Creates correct message for a ping
func (network *Network) SendPingMessage(contact *Contact) error {

	msg := Response{
		RPC:            "ping",
		ID:             NewRandomKademliaID(),
		SendingContact: &network.Node.Me,
	}

	_, err := network.MessageHandler(contact, msg)
	if err != nil {
		errors.Wrap(err, "Something went wrong")
	}
	return nil
}

// Creates correct message object for find_node RPC
func (network *Network) SendFindContactMessage(contact *Contact, kadID *KademliaID) ([]Contact, error) {

	body := Msgbody{
		KadID: kadID, // Puts the nodes ID we are looking for in the body
	}

	msg := Response{
		RPC:            "find_node",
		ID:             NewRandomKademliaID(),
		SendingContact: &network.Node.Me,
		Body:           body,
	}

	res, err := network.MessageHandler(contact, msg)
	if err != nil {
		errors.Wrap(err, "Something went wrong")
	}

	return res.Body.Nodes, nil

}

// Creates correct message object for find_data RPC
func (network *Network) SendFindDataMessage(contact *Contact, hash string) ([]byte, []Contact, error) {
	body := Msgbody{
		Hash: hash, // Hashed id is put in the body
	}

	msg := Response{
		RPC:            "find_data",
		ID:             NewRandomKademliaID(),
		SendingContact: &network.Node.Me,
		Body:           body,
	}

	res, err := network.MessageHandler(contact, msg)
	if err != nil {
		errors.Wrap(err, "Something went wrong")
	}

	return res.Body.Data, res.Body.Nodes, nil
}

// Creates the correct message for a store_data RPC
func (network *Network) SendStoreMessage(contact *Contact, data []byte) error {
	body := Msgbody{
		Data: data, // Data to store is put in the body
	}

	msg := Response{
		RPC:            "store_data",
		ID:             NewRandomKademliaID(),
		SendingContact: &network.Node.Me,
		Body:           body,
	}

	_, err := network.MessageHandler(contact, msg)
	if err != nil {
		errors.Wrap(err, "Something went wrong")
	}
	return nil
}

// Handles the UDP connection dial up.
// Will fetch address of the desired contact and will then open a connection to that adress
// Will marshal the data and then send it over the connection
// Waits for a response from the connection and then returns the response object if it is validated.
func (network *Network) MessageHandler(contact *Contact, msg Response) (Response, error) {
	udpAddr := GetUDPAddrFromContact(contact)

	marshalledMsg := marshallData(msg)

	Conn, err := net.DialUDP("udp", nil, &udpAddr)

	if err != nil {
		return Response{}, errors.Wrap(err, "Client: Failed to open connection to "+udpAddr.IP.String())
	}

	defer Conn.Close()
	Conn.Write([]byte(marshalledMsg))
	buf := make([]byte, 2048)

	//Conn.SetDeadline(time.Now().Add(deadline)) TODODODO

	n, _, _ := Conn.ReadFromUDP(buf)
	res := unmarshallData(buf[0:n])

	// fmt.Println("msg id: ", msg.ID, "\nrec id: ", res.ID, "\nrec rpc : ", res.RPC)

	if Validate(msg, res) {
		network.Node.Routingtable.AddContact(*res.SendingContact)
		// fmt.Println("Received RPC: ", res.RPC, "\nWith RPC ID: ", res.ID, "\nBody: ", res.Body, "\nFrom ", remoteaddr)
	}
	return res, nil
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
func (network *Network) responseHandler(res Response, node Kademlia) Response {
	switch res.RPC {
	case "ping":
		return network.createPingResponse(res)
	case "find_node":
		return network.createFindNodeResponse(res, node)
	case "find_data":
		return network.createFindDataResponse(res, node)
	case "store_data":
		return network.createStoreResponse(res)
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
func (network *Network) createPingResponse(res Response) Response {
	responseMessage := Response{
		RPC:            "ping",
		ID:             res.ID,
		SendingContact: &network.Node.Me,
	}
	return responseMessage
}

// Creates a find_node RPC response containing the nodes k closest contacts to a given ID
func (network *Network) createFindNodeResponse(res Response, node Kademlia) Response {

	contacts := node.Routingtable.FindClosestContacts(res.Body.KadID, bucketSize)

	resBody := Msgbody{
		Nodes: contacts,
	}

	responseMessage := Response{
		RPC:            "find_node",
		ID:             res.ID,
		SendingContact: &network.Node.Me,
		Body:           resBody,
	}

	return responseMessage
}

// Creates a find_data RPC response containing only the data requested if it is stored in the node
// Or will return the 20 closest contacts to the hashed value ID
func (network *Network) createFindDataResponse(res Response, node Kademlia) Response {
	// Function for finding data in node
	// if data in node:

	// else

	contacts := node.Routingtable.FindClosestContacts(NewKademliaID(res.Body.Hash), 20)
	// fmt.Println(contacts)

	resBody := Msgbody{
		Nodes: contacts,
	}

	responseMessage := Response{
		RPC:  "find_data",
		ID:   res.ID,
		Body: resBody,
	}

	// fmt.Printf("%+v \n", responseMessage)

	return responseMessage

}

// Creates a simple store_data RPC response to confirm that the data has been stored on the node
func (network *Network) createStoreResponse(res Response) Response {
	//Stores data in the node

	responseMessage := Response{
		RPC: "store_data",
		ID:  res.ID,
	}
	return responseMessage
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

// Creates a UDPAddr from a contacts ip address.
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
