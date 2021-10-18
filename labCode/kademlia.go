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
	DataStore    map[string][]byte
}

// NewKademliaNode returns a new instance of a Kademlianode
func NewKademliaNode(address string) (node Kademlia) {
	nodeID := NewKademliaID(HashData(address)) // Assign a KademliaID to this node
	node.Me = NewContact(nodeID, address)      // and store to contact object
	node.Routingtable = NewRoutingTable(node.Me)
	node.DataStore = make(map[string][]byte)

	// print trace, remove later
	fmt.Printf("Node %s created on address %s \n", node.Me.ID.String(), node.Me.Address)
	return
}

// LookupContact finds the bucketSize closest nodes and returns a list of contacts
func (kademlia *Kademlia) LookupContact(targetID *KademliaID) (resultlist []Contact) {
	net := &Network{} // network object
	net.Node = kademlia
	ch := make(chan []Contact)  // channel for response
	conCh := make(chan Contact) // Channel for response contact

	// shortlist of k-closest nodes
	shortlist := kademlia.NewLookupList(targetID)

	// if LookupContact on JoinNetwork
	if shortlist.Len() < alpha {
		go AsyncLookup(*targetID, shortlist.Nodelist[0].Node, *net, ch, conCh)
	} else {
		// sending RPCs to the alpha nodes async
		for i := 0; i < alpha; i++ {
			go AsyncLookup(*targetID, shortlist.Nodelist[i].Node, *net, ch, conCh)
		}
	}

	shortlist.updateLookupList(*targetID, ch, conCh, *net)

	// creating the result list
	for _, insItem := range shortlist.Nodelist {
		resultlist = append(resultlist, insItem.Node)
	}
	return
}

// AsyncLookup sends a FindContactMessage to the receiver and writes the response to a channel.
func AsyncLookup(targetID KademliaID, receiver Contact, net Network, ch chan []Contact, conCh chan Contact) {
	reslist, err := net.SendFindContactMessage(&receiver, &targetID)
	if err != nil {
		ch <- reslist
		conCh <- receiver
	} else {
		ch <- reslist
		conCh <- receiver
	}
}

// ########################################################################### \\

// Given a hash from data, finds the closest node where the data is to be stored
func (kademlia *Kademlia) LookupData(hash string) ([]byte, Contact) {
	net := &Network{}
	net.Node = kademlia
	var wg sync.WaitGroup // gorutine waiting pool

	hashID := NewKademliaID(hash) // create kademlia ID from the hashed data
	/*
		shortlist (below) is a LookupList which both contains the contacts
		that need to be traversed in order to find the data as well
		as data itself.
	*/

	shortlist := kademlia.NewLookupList(hashID)

	ch := make(chan []Contact)          // channel -> returns contacts
	targetData := make(chan []byte)     // channel -> when the data is found it is communicated through this channel
	dataContactCh := make(chan Contact) // channel that only takes the contact that returned the data

	if shortlist.Len() < alpha {
		go asyncLookupData(hash, shortlist.Nodelist[0].Node, *net, ch, targetData, dataContactCh)
	} else {
		// sending RPCs to the alpha nodes async
		for i := 0; i < alpha; i++ {
			go asyncLookupData(hash, shortlist.Nodelist[i].Node, *net, ch, targetData, dataContactCh)
		}
	}

	data, con := shortlist.updateLookupData(hash, ch, targetData, dataContactCh, *net, wg)

	// creating the resultdata, con :=shortlist.updateLook list
	return data, con
}

// runs SendFindDataMessage and loads response into two channels:
// ch -> contacts close to the data hash
// target -> the target data
func asyncLookupData(hash string, receiver Contact, net Network, ch chan []Contact, target chan []byte, dataContactCh chan Contact) {
	targetData, reslist, dataContact, _ := net.SendFindDataMessage(&receiver, hash)
	ch <- reslist
	target <- targetData
	dataContactCh <- dataContact
}

// ########################################################################### \\
func (kademlia *Kademlia) Store(data []byte) {
	net := &Network{}
	net.Node = kademlia
	hashFile := HashData(string(data))
	hashID := NewKademliaID(hashFile)

	fileDestinations := kademlia.LookupContact(hashID)
	for _, target := range fileDestinations {
		net.SendStoreMessage(&target, data)
	}

}

// JoinNetwork takes knownpeer or bootstrapNode as input and joins the network.
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
