package labCode

import (
	"testing"
	"time"
)

func TestDataExpired(t *testing.T) {
	ds := NewDataStore()
	testData := []byte("karlsson_på_taket")
	key, _ := ds.addData(testData)

	items, removed := ds.dataExpired()
	if items == 0 && removed == 1 {
		t.Errorf("Failed: Item removed before TTL expiration. Removed: %d", removed)
	}

	// expire stored data object
	ds.setTTL(key, -10000000) // increase the
	items1, removed1 := ds.dataExpired()

	// test if data object has been removed
	if items1 != 0 && removed1 == 1 {
		t.Errorf("Failed: Item was not removed after TTL expiration. Removed: %d", removed1)
	}
}

func TestSetTTL(t *testing.T) {
	// DataStore
	ds := NewDataStore()
	expected_data := []byte("karlsson_på_taket")

	// adding data to the DS.Store
	key, trueExpTime := ds.addData(expected_data)
	expTime, _ := ds.getExpT(key)

	// test if expiration time matches
	if trueExpTime != expTime {
		t.Errorf("The expiration time doesn't match. Want: %s, Got: %s", trueExpTime, expTime)
	}

}

// func Test
func TestHasExpired(t *testing.T) {
	var now time.Time
	now = time.Now()
	notexpiredTime := now.Add(time.Second * 2)
	expiredTime := now.Add(time.Second * -20000)

	if hasExpired(notexpiredTime) {
		t.Errorf("Fail: true when time hasn't expired.")
	}

	if !hasExpired(expiredTime) {
		t.Errorf("Fail: true when time has expired.")
	}
}
