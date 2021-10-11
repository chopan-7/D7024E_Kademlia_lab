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

// Will open up a UDP listener on itself with a given port.
// If a message is received on the listener it will use the response handler to do the
// correct operations
func (network *Network) Listen() {
	server := GetUDPAddrFromContact(&network.Node.Me)
	fmt.Printf("Listening on: %s:%d\n", server.IP, server.Port)

	ServerConn, _ := net.ListenUDP("udp", &server)
	defer ServerConn.Close()
	buf := make([]byte, 2048)
	for {
		n, remoteaddr, _ := ServerConn.ReadFromUDP(buf)
		res := unmarshallData(buf[0:n])

		network.Node.Routingtable.AddContact(*res.SendingContact)
		responseMsg := network.responseHandler(res, *network.Node)

		marshalledMsg := marshallData(responseMsg)

		sendResponse(ServerConn, remoteaddr, marshalledMsg)
	}
}

// Handles the UDP connection dial up. Will fetch address of the desired contact and then open a connection to that adress.
// The data will be marshalled and then send it over the connection, then waits for a response from the connection and then
// returns the response object if it is validated.
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

	// Conn.SetDeadline(time.Now().Add(deadline)) TODODODO

	n, _, _ := Conn.ReadFromUDP(buf)
	res := unmarshallData(buf[0:n])

	if Validate(msg, res) {
		network.Node.Routingtable.AddContact(*res.SendingContact)
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
		return network.CreatePingResponse(res)
	case "find_node":
		return network.CreateFindNodeResponse(res)
	case "find_data":
		return network.CreateFindDataResponse(res)
	case "store_data":
		return network.CreateStoreResponse(res)
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
