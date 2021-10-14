package labCode

import (
	"github.com/pkg/errors"
)

/*
networkSend creates send messages for network
All of the RPC are sent using these message functions. The functions will create a Response
object with the data it has to send. Then the MessageHandler function will send and retreive the response
from the other contact
*/

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
func (network *Network) SendFindDataMessage(contact *Contact, hash string) ([]byte, []Contact, Contact, error) {
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

	// fmt.Printf("\nRes Data: %x\nRes Nodes: %x\n", res.Body.Data, res.Body.Nodes)

	return res.Body.Data, res.Body.Nodes, *res.SendingContact, nil
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
