package labCode

import (
	"testing"
)

func TestCreatePingResponse(t *testing.T) {
	localIPstr := "172.19.0.3" + ":" + "10001" // currentNode IP
	me := NewKademliaNode(localIPstr)

	network := &Network{}
	network.Node = &me
	testID := NewKademliaID("4bc578bd59ddbb005fc9ec86e3f44b9d60cf3f70")

	testRes := Response{
		ID: testID,
	}
	res := network.CreatePingResponse(testRes)
	// check RPC
	if res.RPC != "ping" {
		t.Errorf("RPC = %s; want: ping", res.RPC)
	}
	// check destination ID
	if !res.ID.Equals(testRes.ID) {
		t.Errorf("ID = %d; want: %d", res.ID, testRes.ID)
	}
	// check sender's ID
	if !res.SendingContact.ID.Equals(network.Node.Me.ID) {
		t.Errorf("ID = %d; want: %d", res.SendingContact.ID, network.Node.Me.ID)
	}
}

func TestCreateFindNodeResponse(t *testing.T) {
	localIPstr := "172.19.0.3" + ":" + "10001"
	contactIPstr := "172.19.0.5" + ":" + "10001"

	me := NewKademliaNode(localIPstr)

	contactID := NewKademliaID("4bc578bd593f44b9d60cf3f7ddbb005fc9ec86e0")
	contact := NewContact(contactID, contactIPstr)
	network := &Network{}
	network.Node = &me

	me.Routingtable.AddContact(contact)
	testID := NewKademliaID("4bc578bd59ddbb005fc9ec86e3f44b9d60cf3f70")
	targetID := NewKademliaID("c9ec86e3f44b9d60cf3f704bc578bd59ddbb005f")

	testRes := Response{
		ID: testID,
		Body: Msgbody{
			KadID: targetID,
		},
	}

	res := network.CreateFindNodeResponse(testRes)

	// check RPC
	if res.RPC != "find_node" {
		t.Errorf("RPC = %s; want: find_node", res.RPC)
	}
	// check destination ID
	if !res.ID.Equals(testRes.ID) {
		t.Errorf("ID = %d; want: %d", res.ID, testRes.ID)
	}
	// check sender's ID
	if !res.SendingContact.ID.Equals(network.Node.Me.ID) {
		t.Errorf("ID = %d; want: %d", res.SendingContact.ID, network.Node.Me.ID)
	}
	// check body of RPC
	if len(res.Body.Nodes) > 20 || len(res.Body.Nodes) == 0 || res.Body.Nodes == nil {
		t.Errorf("Nodes length:%x | Expected 20", len(res.Body.Nodes))
	}
}

// does not contain the target data
func TestCreateFindDataResponse_notFound(t *testing.T) {
	localIPstr := "172.19.0.3" + ":" + "10001"
	contactIPstr := "172.19.0.5" + ":" + "10001"
	data := "This_is_the_data"
	hash := HashData(data)

	me := NewKademliaNode(localIPstr)

	contactID := NewKademliaID("4bc578bd593f44b9d60cf3f7ddbb005fc9ec86e0")
	contact := NewContact(contactID, contactIPstr)
	network := &Network{}
	network.Node = &me

	me.Routingtable.AddContact(contact)
	testID := NewKademliaID("4bc578bd59ddbb005fc9ec86e3f44b9d60cf3f70")

	testRes := Response{
		ID: testID,
		Body: Msgbody{
			Hash: hash,
		},
	}

	res := network.CreateFindDataResponse(testRes)

	// check RPC
	if res.RPC != "find_data" {
		t.Errorf("RPC = %s; want: find_data", res.RPC)
	}
	// check destination ID
	if !res.ID.Equals(testRes.ID) {
		t.Errorf("ID = %d; want: %d", res.ID, testRes.ID)
	}
	// check sender's ID
	if !res.SendingContact.ID.Equals(network.Node.Me.ID) {
		t.Errorf("ID = %d; want: %d", res.SendingContact.ID, network.Node.Me.ID)
	}
	// check body of RPC
	if len(res.Body.Nodes) > 20 || len(res.Body.Nodes) == 0 || res.Body.Nodes == nil {
		t.Errorf("Nodes length:%x | Expected 20", len(res.Body.Nodes))
	}
	if !res.Body.Nodes[0].ID.Equals(contactID) {
		t.Errorf("Expected ID: %s | Actual ID: %s", contact.ID, res.Body.Nodes[0].ID)
	}
}
func TestCreateFindDataResponse_found(t *testing.T) {
	localIPstr := "172.19.0.3" + ":" + "10001"
	data := "This_is_the_data"
	hash := HashData(data)

	me := NewKademliaNode(localIPstr)

	network := &Network{}
	network.Node = &me
	network.Node.DataStore[hash] = []byte(data)
	testID := NewKademliaID("4bc578bd59ddbb005fc9ec86e3f44b9d60cf3f70")

	testRes := Response{
		ID: testID,
		Body: Msgbody{
			Hash: hash,
		},
	}

	res := network.CreateFindDataResponse(testRes)

	// check RPC
	if res.RPC != "find_data" {
		t.Errorf("RPC = %s; want: find_data", res.RPC)
	}
	// check destination ID
	if !res.ID.Equals(testRes.ID) {
		t.Errorf("ID = %d; want: %d", res.ID, testRes.ID)
	}
	// check sender's ID
	if !res.SendingContact.ID.Equals(network.Node.Me.ID) {
		t.Errorf("ID = %d; want: %d", res.SendingContact.ID, network.Node.Me.ID)
	}

	// check body of RPC
	if string(res.Body.Data) != data {
		t.Errorf("Expected: %s | Actual: %s", data, string(res.Body.Data))
	}
}

func TestCreateStoreResponse(t *testing.T) {
	localIPstr := "172.19.0.3" + ":" + "10001"
	data := "This_is_the_data"

	me := NewKademliaNode(localIPstr)

	network := &Network{}
	network.Node = &me

	testRes := Response{
		ID: network.Node.Me.ID,
		Body: Msgbody{
			Data: []byte(data),
		},
	}

	res := network.CreateStoreResponse(testRes)
	// check RPC
	if res.RPC != "store_data" {
		t.Errorf("RPC = %s; want: store_data", res.RPC)
	}
	// check destination ID
	if !res.ID.Equals(testRes.ID) {
		t.Errorf("ID = %d; want: %d", res.ID, network.Node.Me.ID)
	}
	// check sender's ID
	if !res.SendingContact.ID.Equals(network.Node.Me.ID) {
		t.Errorf("ID = %d; want: %d", res.SendingContact.ID, network.Node.Me.ID)
	}
}
