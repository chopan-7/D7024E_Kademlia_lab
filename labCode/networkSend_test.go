package labCode

import (
	"log"
	"net"
	"testing"
)

func TestSendPingMessage(t *testing.T) {
	// Create network and kademlia object
	testIP := GetOutboundIP()
	port := "10002"
	testIPStr := testIP.String() + ":" + port
	node := NewKademliaNode(testIPStr)
	net := &Network{}
	net.Node = &node
	go net.Listen()

	err := net.SendPingMessage(&net.Node.Me)
	if err != nil {
		t.Errorf("Failed: %s", err)
	}
}

func TestSendFincContactMessage(t *testing.T) {
	// Create network and kademlia object
	testIP := GetOutboundIP()
	port := "10003"
	testIPStr := testIP.String() + ":" + port
	node := NewKademliaNode(testIPStr)
	net := &Network{}
	net.Node = &node
	go net.Listen()

	clist, err := net.SendFindContactMessage(&net.Node.Me, net.Node.Me.ID)
	if err != nil {
		t.Errorf("Failed: %s", err)
	}

	contact_found := len(clist)
	if contact_found < 1 {
		t.Errorf("Failed: Did not find any contact. Wanted: > 0, Found: %d", contact_found)
	}
}

func TestSendStoreAndFindMessage(t *testing.T) {
	// Create network and kademlia object
	testIP := GetOutboundIP()
	port := "10004"
	testIPStr := testIP.String() + ":" + port
	node := NewKademliaNode(testIPStr)
	net := &Network{}
	net.Node = &node
	go net.Listen()

	// Store data
	testData := []byte("karlsson_på_taket")
	store_err := net.SendStoreMessage(&net.Node.Me, testData)

	if store_err != nil {
		t.Errorf("Failed: Could not store data. Err: %s", store_err)
	}

	// Find data
	foundData, _, _, _ := net.SendFindDataMessage(&net.Node.Me, HashData(string(testData)))
	if string(foundData) != string(testData) {
		t.Errorf("Failed: Could not find data '%s'.", testData)
	}

}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}
