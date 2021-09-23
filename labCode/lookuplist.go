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
