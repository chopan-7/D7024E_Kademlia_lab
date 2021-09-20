package labCode

import "fmt"

// alpha parameter
var alpha int = 3

// Kademlia node definition
// store the routingtable
type Kademlia struct {
	Routingtable *RoutingTable
}

// NewKademliaNode returns a new instance of a Kademlianode
// containing a routingtable for now...
func NewKademliaNode(address string) *Kademlia {
	kademliaNode := &Kademlia{}
	nodeID := NewRandomKademliaID()   // Assign new KademliaID to this node
	me := NewContact(nodeID, address) // and store to contact object
	kademliaNode.Routingtable = NewRoutingTable(me)

	// print trace
	fmt.Println("New node created! :)")
	fmt.Printf("Me: %s \n", me.String())

	return kademliaNode
}

// LookupContact finds the bucketSize closest nodes and returns a list of contacts
func (kademlia *Kademlia) LookupContact(target *Contact) []Contact {
	return kademlia.Routingtable.FindClosestContacts(target.ID, bucketSize)
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
