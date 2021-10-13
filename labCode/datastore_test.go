package labCode

import (
	"testing"
)

func TestSetTTL(t *testing.T) {
	//New Kademlia object
	kad := NewKademliaNode("127.0.0.1:8001")
	data := []byte("karlsson_på_taket")

	// adding data to the DS.Store
	key, trueExpTime := kad.DS.addData(data)
	expTime, _ := kad.DS.getExpT(key)

	if trueExpTime != expTime {
		t.Errorf("The expiration time doesn't match. Want: %s, Got: %s", trueExpTime, expTime)
	}
}

// func Test
