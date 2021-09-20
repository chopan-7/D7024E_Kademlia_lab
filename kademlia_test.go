package main

import (
	"kademlia/labCode"
	"testing"
)

// tests if given a contact with an address string can convert the string to the correct UDPAddr
func TestGetUdpAddr(t *testing.T) {
	con := labCode.Contact{
		Address: "192.168.0.1:5000",
	}

	udpAddr := labCode.GetUDPAddrFromContact(&con)

	if udpAddr.IP.String() != "192.168.0.1" {
		t.Errorf("UDP ip address incorrect, got: %s, want: %s", udpAddr.IP.String(), "192.168.0.1")
	}
	if udpAddr.Port != 5000 {
		t.Errorf("UDP Port incorrect, got: %d, want: %d", &udpAddr.Port, 5000)
	}
}

// tests that validation function returns correct value if validated between different RPC response bodies.
func TestValidateRPCID(t *testing.T) {
	sameRPCID := labCode.NewRandomKademliaID()
	wrongRPCID := labCode.NewRandomKademliaID()

	res1 := labCode.Response{
		RPC: "ping",
		ID:  sameRPCID,
	}

	res2 := labCode.Response{
		RPC: "ping",
		ID:  sameRPCID,
	}

	res3 := labCode.Response{
		RPC: "find_node",
		ID:  sameRPCID,
	}

	res4 := labCode.Response{
		RPC: "ping",
		ID:  wrongRPCID,
	}

	res5 := labCode.Response{
		RPC: "find_node",
		ID:  wrongRPCID,
	}

	// These RPCs should validate
	if !labCode.Validate(res1, res2) {
		t.Errorf("Validation failed when it should have succeeded for responses: %v and %v", res1, res2)
	}
	if labCode.Validate(res1, res3) {
		t.Errorf("Validation succeeded when it should have failed for wrong RPC string between the responses. RPC string 1: %v and RPC string 2: %v", res1.RPC, res3.RPC)
	}
	if labCode.Validate(res1, res4) {
		t.Errorf("Validation succeeded when it should have failed for wrong RPC ID between the responses. RPC ID 1: %v and RPC ID 2: %v", res1.ID, res4.ID)
	}
	if labCode.Validate(res1, res4) {
		t.Errorf("Validation succeeded when it should have failed for wrong RPC string and ID between the responses. RPC 1: %v and RPC 2: %v", res1, res5)
	}

}
