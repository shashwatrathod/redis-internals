package store

import (
	"time"

	"github.com/shashwatrathod/redis-internals/config"
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

	// executes the function for each key value pair in the datastore.
	// the fn should return false if the iteration is to be terminated early, else true.
	ForEach(func(key string, value *Value) bool)

	// returns the Metadata for the given key. Metadata contains information like the creation and last-access timestamps.
	GetKeyMetadata(key string) *KeyMetadata

	// identifies a best-candidate key to be evicted from the datastore, and evicts it.
	// should be used when the storage hits its maximum limit.
	Evict()
}

// Represents a Value that can be stored in the datastore.
type Value struct {
	Value     interface{}
	ValueType SupportedDatatypes
	Expiry    *utils.ExpiryTime
}

// contains information like last-accessed ts and created ts for a key in the store.
type KeyMetadata struct {
	// when the key was last accessed. gets updated everytime the key gets updated or fetched (via Get)
	LastAccessedTimestamp time.Time
	// when the key was created. it is set when a key gets created. does not get updated if the value is updated.
	CreatedTimestamp time.Time
}

// returns a new instance of the KeyMetadata with all timestamps set to current time.
func newKeyMetadata() *KeyMetadata {
	return &KeyMetadata{
		LastAccessedTimestamp: time.Now(),
		CreatedTimestamp:      time.Now(),
	}
}

var storeInstance *SimpleDataStore

func GetStore() *SimpleDataStore {
	if storeInstance == nil {
		storeInstance = &SimpleDataStore{
			data:                 make(map[string]*Value),
			keyMetadata:          make(map[string]*KeyMetadata),
			autoDeletionStrategy: NewRandomSampleAutoDeletionStrategy(AUTO_EXPIRE_SEARCH_LIMIT, AUTO_EXPIRE_ALLOWABLE_EXPIRE_FRACTION), // TODO Make this configurable through additional config params or constructors.
			evictionStrategy:     NewAllKeysLRUEvictionStrategy(config.LRUEvictionSampleSize),
			nKeys:                0,
		}
	}

	return storeInstance
}
