package labCode

import "testing"

func TestNewLookupList(t *testing.T) {
	// Create new node object
	kademlia := NewKademliaNode("127.0.0.1")

	// Populate routing table
	kademlia.Routingtable.AddContact(NewContact(NewKademliaID("2111111190000000000000000000000000000000"), "localhost:8002"))
	for i := 0; i < 19; i++ {
		kademlia.Routingtable.AddContact(NewContact(NewRandomKademliaID(), "localhost:8002"))
	}

	// Create new lookupList
	lookup := kademlia.NewLookupList(NewKademliaID("2111111400000000000000000000000000000000"))

	// Check if len of lookup is = 20
	want_len := 20
	got_len := lookup.Len()
	if got_len != want_len {
		t.Errorf("Failed: Wrong length. Want: %d, Got: %d", want_len, got_len)
	}

}

func TestRefresh(t *testing.T) {
	// Create new node objects
	kademlia := NewKademliaNode("127.0.0.1")
	alpha := NewKademliaNode("255.255.255.255") // alpha node

	// Populate routing table
	alpha.Routingtable.AddContact(NewContact(NewKademliaID("2111111190000000000000000000000000000000"), "localhost:8002"))
	for i := 0; i < 50; i++ {
		kademlia.Routingtable.AddContact(NewContact(NewRandomKademliaID(), "localhost:8002"))
		alpha.Routingtable.AddContact(NewContact(NewRandomKademliaID(), "localhost:8002"))
	}

	// Create new lookupList
	lookup := kademlia.NewLookupList(NewKademliaID("2111111400000000000000000000000000000000"))
	alphasClosest := alpha.Routingtable.FindClosestContacts(NewKademliaID("2111111400000000000000000000000000000000"), 20)

	// refresh lookupList
	nextLookupContact, _ := lookup.refresh(alphasClosest)

	want_nextContact := lookup.Nodelist[0]

	// check if refresh returns correct next contact
	if !nextLookupContact.ID.Equals(want_nextContact.Node.ID) {
		t.Errorf("Failed: Next lookup contact doesn't match. Want: %s, Got: %s", want_nextContact.Node.String(), nextLookupContact.String())
	}

}
