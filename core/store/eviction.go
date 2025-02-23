package store

import (
	"github.com/shashwatrathod/redis-internals/utils"
)

// eviction strategy selects the best key candidate to be deleted from the datastore
// and deletes the key when triggered.
type EvictionStrategy interface {
	// executes the eviction strategy on the datastore. returns the number of keys evicted
	Execute(dstore Store) (int, error)
}

// AllKeysLRUEvictionStrategy implements Redis's allkeys-lru eviction strategy.
// it uses an approximated LRU algorithm to sample "N" keys and evict the least recently used keys out of the sample from the datastore.
type AllKeysLRUEvictionStrategy struct {
	SampleSize int
}

func NewAllKeysLRUEvictionStrategy(sampleSize int) *AllKeysLRUEvictionStrategy {
	return &AllKeysLRUEvictionStrategy{
		SampleSize: sampleSize,
	}
}

func (strategy *AllKeysLRUEvictionStrategy) findLeastRecentlyUsedKey(dstore Store) *string {
	var leastRecentlyUsedKey *string = nil
	var earliestAccessTime utils.LRUTime = utils.GetCurrentLruTime()

	nKeysScanned := 0

	dstore.ForEach(func(key string, val *Value) bool {
		if nKeysScanned >= strategy.SampleSize {
			return false
		}

		keyMetadata := dstore.GetKeyMetadata(key)

		if keyMetadata != nil && uint32(keyMetadata.LastAccessedTimestamp) < uint32(earliestAccessTime) {
			earliestAccessTime = keyMetadata.LastAccessedTimestamp
			leastRecentlyUsedKey = &key
			nKeysScanned++
		}

		return true
	})

	return leastRecentlyUsedKey
}

func (strategy *AllKeysLRUEvictionStrategy) Execute(dstore Store) (int, error) {
	leastRecentlyUsedKey := strategy.findLeastRecentlyUsedKey(dstore)

	if leastRecentlyUsedKey != nil {
		dstore.Delete(*leastRecentlyUsedKey)
	}

	// TODO: implement
	return 0, nil
}
