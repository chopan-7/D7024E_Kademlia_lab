package labCode

import (
	"testing"
)

func TestNewKademliaNode(t *testing.T) {
	testAddr := "172.19.0.3:10101"
	testID := NewKademliaID("4bc578bd59ddbb005fc9ec86e3f44b9d60cf3f70")
	nn := NewKademliaNode(testAddr) // new node object
	contact := nn.Routingtable.me   // contact info

	if !contact.ID.Equals(testID) {
		t.Errorf("ID = %d; want: %d", contact.ID, testID)
	}
}

func TestStoreAndLookup(t *testing.T) {
	// Create network and kademlia object
	testIP := GetOutboundIP()
	port := "10012"
	testIPStr := testIP.String() + ":" + port
	node := NewKademliaNode(testIPStr)
	net := &Network{}
	net.Node = &node
	go net.Listen()

	// test JoinNetwork
	net.Node.JoinNetwork(&net.Node.Me)

	// test Store
	testStoreData := []byte("karlsson_på_taket")
	storeAt := net.Node.Store(testStoreData)

	if len(storeAt) < 1 {
		t.Error("Fail: couldn't store data in the network.")
	}

	// test LookupData
	foundData, _ := net.Node.LookupData(HashData(string(testStoreData)))

	if string(foundData) != string(testStoreData) {
		t.Error("Fail: couldn't find stored data in the network.")
	}
}
