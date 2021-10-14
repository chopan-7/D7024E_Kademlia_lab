package labCode

import (
	"testing"
	"time"
)

func TestDataExpired(t *testing.T) {
	ds := NewDataStore()
	testData := []byte("karlsson_på_taket")
	key, _ := ds.addData(testData)

	// expire stored data object
	ds.setTTL(key, -10000000000)
	items, removed := ds.dataExpired()

	// test if data object has been removed
	if items != 0 && removed == 1 {
		t.Errorf("Failed: Items removed = %d", removed)
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
	expiredTime := now.Add(time.Second * -2)

	if hasExpired(notexpiredTime) {
		t.Errorf("Fail: true when time hasn't expired.")
	}

	if !hasExpired(expiredTime) {
		t.Errorf("Fail: true when time has expired.")
	}
}
