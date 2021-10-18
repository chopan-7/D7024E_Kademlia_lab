package labCode

import (
	"sort"
	"sync"
)

// ContactCandidates definition
// stores an array of Contacts
type LookupCandidates struct {
	Nodelist []LookupListItems
}

// LookupList for temporary storing Nodeitems
type LookupList struct {
	Nodelist []LookupListItems
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

func (lookuplist *LookupList) refresh(contacts []Contact, notConsidered []LookupListItems) {
	candidateList := LookupCandidates{} // holds the response []Contact
	tempList := lookuplist.Nodelist     // Copy of lookuplist
	for _, contact := range contacts {
		listItem := LookupListItems{contact, false}
		candidateList.Nodelist = append(candidateList.Nodelist, listItem)
	}
	sortingList := LookupCandidates{}
	candidateList.Remove(notConsidered)
	sortingList.Append(tempList)
	sortingList.Remove(notConsidered)
	sortingList.Append(candidateList.Nodelist)
	sortingList.Sort()

	if len(sortingList.Nodelist) < bucketSize {
		lookuplist.Nodelist = sortingList.GetContacts(len(sortingList.Nodelist))
	} else {
		lookuplist.Nodelist = sortingList.GetContacts(bucketSize)
	}

	return
}

func (lookuplist *LookupList) updateLookupList(targetID KademliaID, ch chan []Contact, conCh chan Contact, net Network) {
	notConsidered := LookupCandidates{}
	for {
		contacts := <-ch
		responder := <-conCh
		if len(contacts) > 0 {
			lookuplist.refresh(contacts, notConsidered.Nodelist)
		} else {
			resItem := LookupListItems{responder, true}
			notConList := []LookupListItems{resItem}
			notConsidered.Append(notConList)
			lookuplist.refresh([]Contact{}, notConsidered.Nodelist)
		}
		nextContact, Done := lookuplist.findNextLookup()
		if Done {
			return
		} else {
			go AsyncLookup(targetID, nextContact, net, ch, conCh)
		}
	}
}

// ########################################################################### \\
func (lookuplist *LookupList) updateLookupData(hash string, ch chan []Contact, target chan []byte, dataContactCh chan Contact, net Network, wg sync.WaitGroup) ([]byte, Contact) {
	for {
		contacts := <-ch
		targetData := <-target
		dataContact := <-dataContactCh

		// data not nil = correct data is found
		if targetData != nil {
			return targetData, dataContact
		}

		lookuplist.refresh(contacts, []LookupListItems{})
		nextContact, Done := lookuplist.findNextLookup()
		if Done {
			return nil, Contact{}
		} else {
			go asyncLookupData(hash, nextContact, net, ch, target, dataContactCh)
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

func (candidates *LookupCandidates) Remove(Contacts []LookupListItems) {
	for _, newCandidate := range Contacts {
		for i, candidate := range candidates.Nodelist {
			if candidate.Node.ID.Equals(newCandidate.Node.ID) {
				candidates.remove(i)
				break
			}
		}
	}
}

func (lc *LookupCandidates) remove(n int) {
	lc.Nodelist = append(lc.Nodelist[:n], lc.Nodelist[n+1:]...)
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
