package labCode

// alpha parameter
var alpha int = 3

// Kademlia node definition
// store the routingtable
type Kademlia struct {
	routingtable *RoutingTable
}

// NewKademliaNode returns a new instance of a Kademlianode
// containing a routingtable for now...
func NewKademliaNode(address string) *Kademlia {
	kademliaNode := &Kademlia{}
	nodeID := NewKademliaID(address)	// Assign new KademliaID to this node
	me := NewContact(nodeID, address)	// and store to contact object
	kademliaNode.routingtable = NewRoutingTable(me)

	return kademliaNode
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	// TODO
	// 1. Choose alpha nodes from appropriate bucket
	// 2. Calculate the distances for each node
	closest := FindClosestContacts(target, alpha)
	// 3. Send 'SendFindContactMessage' to the closest node
	// 4. Update bucket with new contact?
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
