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
	localIPstr := "172.19.0.3" + ":" + "10001" // currentNode IP
	me := NewKademliaNode(localIPstr)

	network := &Network{}
	network.Node = &me
	testID := NewKademliaID("4bc578bd59ddbb005fc9ec86e3f44b9d60cf3f70")

	testRes := Response{
		Body: Msgbody{
			KadID: testID,
		},
	}

	res := network.CreateFindNodeResponse(testRes)

}
