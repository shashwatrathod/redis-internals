package store

import (
	"github.com/shashwatrathod/redis-internals/core/resp"
	"github.com/shashwatrathod/redis-internals/utils"
)

// Represents the DataTypes currently supported by the Application.
type SupportedDatatypes int

const (
	String  SupportedDatatypes = SupportedDatatypes(resp.BulkString)
	Integer                    = resp.RespInteger
	Array                      = resp.RespArray
)

const (
	AUTO_EXPIRE_SEARCH_LIMIT                      = 20
	AUTO_EXPIRE_ALLOWABLE_EXPIRE_FRACTION float32 = 0.25
)

type Store interface {
	// sets the values of the given key in the store. Overrites the value if the key already exists.
	Put(key string, value *Value)

	// returns the value of the given key if it exists in the store, else returns nil.
	Get(key string) *Value

	// deletes the given key from the store.
	// returns true if the key was present in the store, else false.
	Delete(key string) bool

	// searches a random sample of AUTO_EXPIRE_SEARCH_LIMIT keys with expiry, and
	// purges the expired keys. If the % of expired keys is more than AUTO_EXPIRE_ALLOWABLE_EXPIRE_FRACTION,
	// then repeats the process.
	AutoDeleteExpiredKeys()

	// resets the data in the store. DELETES all the keys. WARNING : Irreversable operation!
	Reset()
}

// Represents a Value that can be stored in the datastore.
type Value struct {
	Value     interface{}
	ValueType SupportedDatatypes
	Expiry    *utils.ExpiryTime
}

var storeInstance *SimpleDataStore

func GetStore() *SimpleDataStore {
	if storeInstance == nil {
		storeInstance = &SimpleDataStore{
			data: make(map[string]*Value),
		}
	}

	return storeInstance
}
