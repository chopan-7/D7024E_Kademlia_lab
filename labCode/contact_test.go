package labCode

import "testing"

func TestContact(t *testing.T) {
	testIP := "127.0.0.1"
	testID := NewRandomKademliaID()
	testContact := NewContact(testID, testIP)

	testContact.CalcDistance(testID)

}
