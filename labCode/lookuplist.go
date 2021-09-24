package labCode

import (
	"sort"
	"sync"
)

// ContactCandidates definition
// stores an array of Contacts
type LookupCandidates struct {
	Nodelist []LookupListItems
	Mux      sync.Mutex
}

// LookupList for temporary storing Nodeitems
type LookupList struct {
	Nodelist []LookupListItems
	Mux      sync.Mutex
}

type LookupListItems struct {
	Node Contact
	Flag bool
}

// NewLookupList retuns a LookupList with k-closest nodes from the nodes routingtable.
func (kademlia *Kademlia) NewLookupList(targetID *KademliaID) (ls *LookupList) {
	ls = &LookupList{}
	closestK := kademlia.Routingtable.FindClosestContacts(targetID, bucketSize)

	for _, item := range closestK {
		lsItem := &LookupListItems{item, false}
		ls.Nodelist = append(ls.Nodelist, *lsItem)
	}
	return
}

// Len returns the lenght of the LookupList
func (ls *LookupList) Len() int {
	return len(ls.Nodelist)
}

func (lookuplist *LookupList) updateLookupList(targetID KademliaID, ch chan []Contact, net Network, wg sync.WaitGroup) {
	defer wg.Done()
	for {
		contacts := <-ch
		tempList := LookupList{}         // holds the response []Contact
		tempList2 := lookuplist.Nodelist // Copy of lookuplist
		for _, contact := range contacts {
			listItem := LookupListItems{contact, false}
			tempList.Nodelist = append(tempList.Nodelist, listItem)
		}

		// sorting/filtering list
		sortingList := LookupCandidates{}
		sortingList.Append(tempList2)
		sortingList.Append(tempList.Nodelist)
		sortingList.Sort()

		// update the lookuplist
		if len(sortingList.Nodelist) < bucketSize {
			lookuplist.Nodelist = sortingList.GetContacts(len(sortingList.Nodelist))
		} else {
			lookuplist.Nodelist = sortingList.GetContacts(bucketSize)
		}

		nextContact, Done := lookuplist.findNextLookup()
		if Done {
			return
		} else {
			go AsyncLookup(targetID, nextContact, net, ch)
		}
	}
}

// findNextLookup returns the next contact to visit in the LookupList.
// It also returns true if all the contacts in the LookupList has been visisted.
func (lookuplist *LookupList) findNextLookup() (Contact, bool) {
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

// Append an array of Contacts to the ContactCandidates if not duplicate
func (candidates *LookupCandidates) Append(Contacts []LookupListItems) {
	for _, newCandidate := range Contacts {
		add := true
		for _, candidate := range candidates.Nodelist {
			if candidate.Node.ID.Equals(newCandidate.Node.ID) {
				add = false
				break
			}
		}
		if add {
			candidates.Nodelist = append(candidates.Nodelist, newCandidate)
		}
	}
}

// GetContacts returns the first count number of Contacts
func (candidates *LookupCandidates) GetContacts(count int) []LookupListItems {
	return candidates.Nodelist[:count]
}

// Sort the Contacts in ContactCandidates
func (candidates *LookupCandidates) Sort() {
	sort.Sort(candidates)
}

// Len returns the length of the ContactCandidates
func (candidates *LookupCandidates) Len() int {
	return len(candidates.Nodelist)
}

// Swap the position of the Contacts at i and j
// WARNING does not check if either i or j is within range
func (candidates *LookupCandidates) Swap(i, j int) {
	candidates.Nodelist[i], candidates.Nodelist[j] = candidates.Nodelist[j], candidates.Nodelist[i]
}

// Less returns true if the Contact at index i is smaller than
// the Contact at index j
func (candidates *LookupCandidates) Less(i, j int) bool {
	return candidates.Nodelist[i].Node.Less(&candidates.Nodelist[j].Node)
}
