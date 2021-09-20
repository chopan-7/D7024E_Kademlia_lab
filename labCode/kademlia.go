package labCode

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"sync"
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
	nodelist [][]Contact
	mux      sync.Mutex
}

type visited struct {
	nodes []Contact
	mux   sync.Mutex
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
	lupls := &LookupList{}
	v := &visited{}
	alphaNodes := make([]Contact, alpha)

	// Find the k closest node to target
	closest := kademlia.Routingtable.FindClosestContacts(target.ID, bucketSize)

	// add the the alpha closest to the alphaNodes
	for i := 0; i < alpha; i++ {
		alphaNodes[i] = closest[i]

		// print distance to target during test
		// fmt.Printf("Distance to target: %X \n", *nodeItem.contact.distance)
	}

	// sending RPCs to the alpha nodes async
	for i, node := range alphaNodes {
		go func() {
			// add node to visisted
			v.mux.Lock()
			v.nodes = append(v.nodes, node)

			fmt.Printf("Sending RPC to: %s \n", string(node.String()))
			_, contactList, _ := net.SendFindContactMessage(&node)

			// add contactList to appropriate nodelist
			lupls.mux.Lock()
			for j, node := range contactList {
				lupls.nodelist[i][j] = node
			}
			lupls.mux.Unlock()
			v.mux.Unlock()

		}()
	}

	// updating the appropriate bucket with contacts in nodelist when the lookup routine is done
	for !gotResponseFrom(*v) {
		lupls.mux.Lock()
		for _, contacts := range lupls.nodelist {
			kademlia.updateBuckets(contacts)
		}
		lupls.mux.Unlock()
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

// gotResponseFrom returns true if all alphanodes has responded
func gotResponseFrom(v visited) bool {
	v.mux.Lock()
	if len(v.nodes) < alpha-1 {
		return false
	}
	v.mux.Unlock()
	return true
}

// Helper function for hashing data returing hexstring
func hashData(data string) (hashed string) {
	newHash := sha1.New()
	newHash.Write([]byte(data))
	hashed = hex.EncodeToString(newHash.Sum(nil))
	return
}
