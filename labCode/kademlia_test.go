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
