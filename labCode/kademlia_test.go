package labCode

import(
	"fmt"
	"testing"
	"encoding/hex"
)


func TestNewNode(t *testing.T) {
	testAddr := "localhost:8080"
	nn := NewNode(testAddr)			// new node object
	contact := nn.routingtable.me	// contact info
	nodeID, _ := hex.DecodeString(testAddr)	// nodeID in hex
	
	if contact.ID.String() != testAddr {
		t.Errorf("ID = %d; want: %d", contact.ID, testAddr)
	}
}