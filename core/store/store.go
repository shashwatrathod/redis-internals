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
	Put(key string, value string, expiry *utils.ExpiryTime)

	// returns the value of the given key if it exists in the store, else returns nil.
	Get(key string) *Value

	// deletes the given key from the store.
	// returns true if the key was present in the store, else false.
	Delete(key string) bool

	// returns the expiry timestamp of the given key.
	// returns null if they key doesn't exist or if there is no expiry set on the key.
	GetExpiry(key string) *int64

	// sets the expiry of the key to the given value if the key exists.
	// removes the expiry on the key if expiry==nil.
	SetExpiry(key string, expiry *utils.ExpiryTime)

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

	// evicts keys from the store based on different eviction strategies.
	// you can tweak the proportion of keys evicted in each run by modifying the config.EvictionRatio parameter.
	// returns the number of keys evicted.
	Evict() int

	// returns the number of keys present in the datastore at the moment.
	KeyCount() int
}

// Represents a Value that can be stored in the datastore.
type Value struct {
	Value     interface{}
	ValueType SupportedDatatypes
}

// contains information like last-accessed ts and created ts for a key in the store.
type KeyMetadata struct {
	// when the key was last accessed. gets updated everytime the key gets updated or fetched (via Get)
	LastAccessedTimestamp utils.LRUTime
	// when the key was created. it is set when a key gets created. does not get updated if the value is updated.
	CreatedTimestamp time.Time
}

// returns a new instance of the KeyMetadata with all timestamps set to current time.
func newKeyMetadata() *KeyMetadata {
	return &KeyMetadata{
		LastAccessedTimestamp: utils.GetCurrentLruTime(),
		CreatedTimestamp:      time.Now(),
	}
}

var storeInstance *DataStore

func GetStore() *DataStore {
	if storeInstance == nil {
		storeInstance = &DataStore{
			data:                 make(map[string]*Value),
			keyMetadata:          make(map[string]*KeyMetadata),
			autoDeletionStrategy: NewRandomSampleAutoDeletionStrategy(AUTO_EXPIRE_SEARCH_LIMIT, AUTO_EXPIRE_ALLOWABLE_EXPIRE_FRACTION), // TODO Make this configurable through additional config params or constructors.
			evictionStrategy:     NewAllKeysLRUEvictionStrategy(config.LRUEvictionSampleSize),
		}
	}

	return storeInstance
}
