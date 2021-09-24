package labCode

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"sync"
)

// alpha parameter
const alpha int = 3

// Kademlia node definition
// store the routingtable
type Kademlia struct {
	Me           Contact
	Routingtable *RoutingTable
}

// NewKademliaNode returns a new instance of a Kademlianode
func NewKademliaNode(address string) (node Kademlia) {
	nodeID := NewKademliaID(HashData(address)) // Assign a KademliaID to this node
	node.Me = NewContact(nodeID, address)      // and store to contact object
	node.Routingtable = NewRoutingTable(node.Me)

	// print trace, remove later
	fmt.Printf("Node %s created on address %s \n", node.Me.ID.String(), node.Me.Address)
	return
}

// LookupContact finds the bucketSize closest nodes and returns a list of contacts
func (kademlia *Kademlia) LookupContact(targetID *KademliaID) (resultlist []Contact) {
	net := &Network{kademlia}  // network object
	var wg sync.WaitGroup      // gorutine waiting pool
	ch := make(chan []Contact) // channel for response

	// shortlist of k-closest nodes
	shortlist := kademlia.NewLookupList(targetID)

	// if LookupContact on JoinNetwork
	if shortlist.Len() < alpha {
		go AsyncLookup(*targetID, shortlist.Nodelist[0].Node, *net, ch)
	} else {
		// sending RPCs to the alpha nodes async
		for i := 0; i < alpha; i++ {
			go AsyncLookup(*targetID, shortlist.Nodelist[i].Node, *net, ch)
		}
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		shortlist.updateLookupList(*targetID, ch, *net, wg)
	}()
	wg.Wait()

	// creating the result list
	for _, insItem := range shortlist.Nodelist {
		resultlist = append(resultlist, insItem.Node)
	}
	return
}

// AsyncLookup sends a FindContactMessage to the receiver and writes the response to a channel.
func AsyncLookup(targetID KademliaID, receiver Contact, net Network, ch chan []Contact) {
	reslist, _ := net.SendFindContactMessage(&receiver, &targetID)
	ch <- reslist
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}

// JoinNetwork takes knownpeer or bootstrapNode
func (kademlia *Kademlia) JoinNetwork(knownpeer *Contact) {
	kademlia.Routingtable.AddContact(*knownpeer)
	kademlia.LookupContact(kademlia.Me.ID)
	fmt.Printf("Joining network")
}

// Helper function for hashing data returing hexstring
func HashData(data string) (hashString string) {
	newHash := sha1.New()
	newHash.Write([]byte(data))
	hashString = hex.EncodeToString(newHash.Sum(nil))
	return
}
