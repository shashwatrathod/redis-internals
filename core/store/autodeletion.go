package store

import "log"

// Auto-deletion mechanism deletes expired keys form the datastore when executed.
type AutoDeletionStrategy interface {
	// Executes the Autodeletion strategy onto the inteface
	Execute(dstore Store)
}

// AutoDeletionStrategy that randomly samples N keys with an expiry to
// evict the expired keys. Repeats the process until certain threshold is met.
type RandomSampleAutoDeletionStrategy struct {
	// number of keys to be sampled
	sampleSize int

	// fraction of keys found to be expired after which the process is rerun again.
	rerunThreshold float32
}

func NewRandomSampleAutoDeletionStrategy(sampleSize int, rerunThreshold float32) *RandomSampleAutoDeletionStrategy {
	if rerunThreshold >= 0.99 || rerunThreshold <= 0 {
		log.Panicf("Rerun threshold must be in the range (0, 0.99)")
	}
	return &RandomSampleAutoDeletionStrategy{
		sampleSize:     sampleSize,
		rerunThreshold: rerunThreshold,
	}
}

func (strategy *RandomSampleAutoDeletionStrategy) expireSample(dstore Store) float32 {
	var nSearched int = 0

	var keysToBeDeleted []string = make([]string, 0)

	// find the keys to be deleted
	dstore.ForEach(func(key string, val *Value) bool {
		if val.Expiry != nil {
			nSearched++

			if val.Expiry.IsExpired() {
				keysToBeDeleted = append(keysToBeDeleted, key)
			}
		}

		if nSearched == strategy.sampleSize {
			return false
		}

		return true
	})

	// delete expired keys
	var nExpired int = 0

	for _, key := range keysToBeDeleted {
		if dstore.Delete(key) {
			nExpired++
		}
	}

	if nSearched > 0 {
		return float32(nExpired) / float32(nSearched)
	} else {
		return 0
	}
}

func (strategy *RandomSampleAutoDeletionStrategy) Execute(dstore Store) {
	for {
		fracExpired := strategy.expireSample(dstore)

		log.Printf("auto-deleted %.2f %% of the sample.\n", fracExpired*100)

		if fracExpired < strategy.rerunThreshold {
			break
		}
	}
}
