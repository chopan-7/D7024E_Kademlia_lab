package labCode

import (
	"log"
	"os"
	"time"
)

/*
Datastore holds the data as key/value pair for each node with functionalities that checks for object expiration
and object expiration delay.
*/

const TTL time.Duration = 60 // the time which a key/value pair expires after publication date

type DataStore struct {
	Store  map[string][]byte
	Expire map[string]time.Time
	Log    *log.Logger
}

// Create a new DataStore object
func NewDataStore() *DataStore {
	DS := &DataStore{}
	DS.Store = make(map[string][]byte)
	DS.Expire = make(map[string]time.Time)

	// Data log
	file, err := os.OpenFile("store_log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	DS.Log = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	return DS
}

// SetTTL set the expiration time for the key/value pair
func (DS *DataStore) setTTL(key string) (exp time.Time) {
	t := time.Now()
	exp = t.Add(time.Second * TTL)
	DS.Expire[key] = exp
	return
}

// storeData saves the given data as a key/value pair to the nodes datastore and returns the key
func (DS *DataStore) addData(data []byte) (key string, exp time.Time) {
	key = HashData(string(data))
	exp = DS.setTTL(key)
	DS.Store[key] = data
	DS.Log.Printf("New data object stored with key: %s", key)
	return
}

// dataExpired checks if the data objects in its own DataStore has expired and removes all expired data object.
func (DS *DataStore) dataExpired() (items int, removed int) {
	items = len(DS.Store) // number of data obejcts in DataStore
	removed = 0           // number of data objects removed from the DataStore

	if items > 0 {
		for key, _ := range DS.Store {
			exp := DS.Expire[key] // TTL of the data object
			if hasExpired(exp) {
				// remove the object from DataStore and DataExpire
				delete(DS.Store, key)
				delete(DS.Expire, key)
				DS.Log.Printf("Data object (key: %s) has been removed due to expired TTL.", key)

				// update return values
				items = len(DS.Store)
				removed++
			}
		}
	}
	return
}

// getData(key) returns the value and boolean
func (DS *DataStore) getData(key string) (val []byte, hasVal bool) {
	val, hasVal = DS.Store[key]
	return
}

// hasExpired compares the time.Time value with time.Now()
func hasExpired(expTime time.Time) bool {
	if time.Now().After(expTime) {
		return true
	} else {
		return false
	}
}
