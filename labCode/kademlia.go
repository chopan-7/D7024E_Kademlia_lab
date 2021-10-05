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

// ########################################################################### \\

// Given a hash from data, finds the closest node where the data is to be stored
func (kademlia *Kademlia) LookupData(hash string) []byte {
	net := &Network{kademlia}
	var wg sync.WaitGroup // gorutine waiting pool

	hashID := NewKademliaID(hash) // create kademlia ID from the hashed data
	/*
		listContact (below) is a LookupList which both contains the contacts
		that need to be traversed in order to find the data as well
		as data itself.
	*/

	shortlist := kademlia.NewLookupList(hashID)

	ch := make(chan []Contact)      // channel -> returns contacts
	targetData := make(chan []byte) // channel -> when the data is found it is communicated through this channel

	/* 	// Find the k closest node to target
	   	for _, insItem := range myClosest {
	   		lookupitem := &LookupListItems{insItem, false}
	   		listContact.Nodelist = append(listContact.Nodelist, *lookupitem)
	   	} */

	// sending RPCs to the alpha nodes async
	for i := 0; i < alpha; i++ {
		go asyncLookupData(hash, shortlist.Nodelist[i].Node, *net, ch, targetData)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		shortlist.updateLookupData(hash, ch, targetData, *net, wg)
	}()
	wg.Wait()

	// creating the result list
	return shortlist.Data
}

// runs SendFindDataMessage and loads response into two channels:
// ch -> contacts close to the data hash
// target -> the target data
func asyncLookupData(hash string, receiver Contact, net Network, ch chan []Contact, target chan []byte) {
	targetData, reslist, _ := net.SendFindDataMessage(&receiver, hash)
	ch <- reslist
	target <- targetData
}

func (lookuplist *LookupList) updateLookupData(hash string, ch chan []Contact, target chan []byte, net Network, wg sync.WaitGroup) {
	defer wg.Done()
	for {
		contacts := <-ch
		targetData := <-target

		// data not nil = correct data is found
		if targetData != nil {
			fmt.Printf("\nFind data done!\n")
			lookuplist.Data = targetData
			return
		}

		tempList := LookupList{}         // holds the response []Contact
		tempList2 := lookuplist.Nodelist // Copy of lookuplist
		for _, contact := range contacts {
			listItem := LookupListItems{contact, false}
			tempList.Nodelist = append(tempList.Nodelist, listItem)
		}

		// sorting/filtering list
		sortingList := LookupCandidates{}
		sortingList.Append(tempList2)         // right order??
		sortingList.Append(tempList.Nodelist) // right order??
		sortingList.Sort()

		// update the lookuplist
		if len(sortingList.Nodelist) < bucketSize {
			lookuplist.Nodelist = sortingList.GetContacts(len(sortingList.Nodelist))
		} else {
			lookuplist.Nodelist = sortingList.GetContacts(bucketSize)
		}

		nextContact, Done := lookuplist.findNextLookup()
		if Done {
			fmt.Printf("\nLookupdone!\n")
			return
		} else {
			go asyncLookupData(hash, nextContact, net, ch, target)
		}
	}
}

func findNextLookupData(lookuplist *LookupList) (Contact, bool) {
	var nextItem Contact
	done := true
	for i, item := range lookuplist.Nodelist {
		if item.Flag == false {
			nextItem = item.Node
			lookuplist.Nodelist[i].Flag = true
			done = false
			break
		}
	}
	return nextItem, done
}

// ########################################################################### \\
func (kademlia *Kademlia) Store(data []byte) {
	net := &Network{kademlia}
	hashFile := HashData(string(data))
	hashID := NewKademliaID(hashFile)

	fileDestinations := kademlia.Routingtable.FindClosestContacts(hashID, bucketSize)

	for _, target := range fileDestinations {
		go net.SendStoreMessage(&target, data)
		// net.SendStoreMessage(&target, data)
	}

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
