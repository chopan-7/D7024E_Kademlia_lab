package labCode

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

// alpha parameter
var alpha int = 3

// Kademlia node definition
// store the routingtable
type Kademlia struct {
	Routingtable *RoutingTable
}

// LookupList for temporary storing nodeitems
type LookupList struct {
	nodelist []Contact
}

// NewKademliaNode returns a new instance of a Kademlianode
// containing a routingtable for now...
func NewKademliaNode(address string) (node Kademlia) {
	nodeID := NewKademliaID(hashData(address)) // Assign a KademliaID to this node
	me := NewContact(nodeID, address)          // and store to contact object
	node.Routingtable = NewRoutingTable(me)

	// print trace, remove later
	fmt.Println("New node created! :)")
	fmt.Printf("Me: %s \n", me.String())

	return
}

// LookupContact finds the bucketSize closest nodes and returns a list of contacts
func (kademlia *Kademlia) LookupContact(target *Contact) []Contact {
	net := &Network{}
	// lookuplist := &LookupList{}
	listItems := make([]Contact, alpha)

	// Find the k closest node to target
	closest := kademlia.Routingtable.FindClosestContacts(target.ID, bucketSize)

	// add the the alpha closest to the LookupList
	for i := 0; i < alpha; i++ {
		listItems[i] = closest[i]

		// print distance to target during test
		// fmt.Printf("Distance to target: %X \n", *nodeItem.contact.distance)

	}

	nodelist := &listItems

	// sending RPCs to the alpha nodes
	for _, node := range *nodelist {
		go net.SendFindContactMessage(&node)
		fmt.Printf("Sending RPC to: %s \n", string(node.String()))
	}

	return closest
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}

// UpdateBuckets adds/update the appropriate buckets in the routingtable.
func (kademlia *Kademlia) updateBuckets(clist []Contact) {
	rt := kademlia.Routingtable

	// insert all contacts from clist to appropriate bucket
	for _, contact := range clist {
		rt.AddContact(contact)
	}
}

// Helper function for hashing data returing hexstring
func hashData(data string) (hashed string) {
	newHash := sha1.New()
	newHash.Write([]byte(data))
	hashed = hex.EncodeToString(newHash.Sum(nil))
	return
}
