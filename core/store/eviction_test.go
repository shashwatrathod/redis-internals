package store_test

import (
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/shashwatrathod/redis-internals/core/store"
	mocks "github.com/shashwatrathod/redis-internals/mocks/github.com/shashwatrathod/redis-internals/core/store"
	"github.com/stretchr/testify/mock"
)

func TestEvictionStrategy(t *testing.T) {
	RegisterFailHandler(Fail)
}

var _ = Describe("AllKeysLRUEvictionStrategy", func() {
	var (
		mockStore  *mocks.Store
		strategy   *store.AllKeysLRUEvictionStrategy
		sampleSize int
	)

	BeforeEach(func() {
		mockStore = mocks.NewStore(GinkgoT())
		sampleSize = 3
		strategy = store.NewAllKeysLRUEvictionStrategy(sampleSize)
	})

	It("should evict the least recently used key", func() {
		// Set up mock store with key metadata
		mockStore.On("ForEach", mock.Anything).Run(func(args mock.Arguments) {
			fn := args.Get(0).(func(string, *store.Value) bool)
			fn("key1", &store.Value{})
			fn("key2", &store.Value{})
			fn("key3", &store.Value{})
		}).Return()

		mockStore.On("GetKeyMetadata", "key1").Return(&store.KeyMetadata{
			LastAccessedTimestamp: time.Now().Add(-10 * time.Minute),
		})
		mockStore.On("GetKeyMetadata", "key2").Return(&store.KeyMetadata{
			LastAccessedTimestamp: time.Now().Add(-5 * time.Minute),
		})
		mockStore.On("GetKeyMetadata", "key3").Return(&store.KeyMetadata{
			LastAccessedTimestamp: time.Now().Add(-1 * time.Minute),
		})

		mockStore.On("Delete", "key1").Return(true)

		// Execute eviction strategy
		strategy.Execute(mockStore)

		// Verify that the least recently used key was deleted
		mockStore.AssertCalled(GinkgoT(), "Delete", "key1")

		// Verify that no other key was deleted
		mockStore.AssertNotCalled(GinkgoT(), "Delete", "key2")
		mockStore.AssertNotCalled(GinkgoT(), "Delete", "key3")
	})

	It("should not evict any key if the store is empty", func() {
		// Set up mock store with no keys
		mockStore.On("ForEach", mock.Anything).Return()

		// Execute eviction strategy
		strategy.Execute(mockStore)

		// Verify that no key was deleted
		mockStore.AssertNotCalled(GinkgoT(), "Delete", mock.Anything)
	})

	It("should handle the case where all keys have the same access time", func() {
		// Set up mock store with key metadata
		mockStore.On("ForEach", mock.Anything).Run(func(args mock.Arguments) {
			fn := args.Get(0).(func(string, *store.Value) bool)
			fn("key1", &store.Value{})
			fn("key2", &store.Value{})
			fn("key3", &store.Value{})
		}).Return()

		sameTime := time.Now().Add(-10 * time.Minute)
		mockStore.On("GetKeyMetadata", "key1").Return(&store.KeyMetadata{
			LastAccessedTimestamp: sameTime,
		})
		mockStore.On("GetKeyMetadata", "key2").Return(&store.KeyMetadata{
			LastAccessedTimestamp: sameTime,
		})
		mockStore.On("GetKeyMetadata", "key3").Return(&store.KeyMetadata{
			LastAccessedTimestamp: sameTime,
		})

		mockStore.On("Delete", "key1").Return(true)

		// Execute eviction strategy
		strategy.Execute(mockStore)

		// Verify that one of the keys was deleted
		mockStore.AssertCalled(GinkgoT(), "Delete", "key1")
	})
	It("should only sample 'SampleSize' keys", func() {
		nkeys := 10
		// Set up mock store with more keys than the sample size
		mockStore.On("ForEach", mock.Anything).Run(func(args mock.Arguments) {
			fn := args.Get(0).(func(string, *store.Value) bool)
			for i := 1; i <= nkeys; i++ {
				cont := fn(fmt.Sprintf("key-%d", i), &store.Value{})
				if !cont {
					break
				}
			}
		}).Return()

		for i := 1; i <= nkeys; i++ {
			mockStore.On("GetKeyMetadata", fmt.Sprintf("key-%d", i)).Return(&store.KeyMetadata{
				LastAccessedTimestamp: time.Now().Add(-time.Duration(i) * time.Minute),
			}).Maybe()
		}

		mockStore.On("Delete", mock.Anything).Return(true)

		// Execute eviction strategy
		strategy.Execute(mockStore)

		// Verify that only 'SampleSize' keys were sampled
		mockStore.AssertNumberOfCalls(GinkgoT(), "GetKeyMetadata", sampleSize)
	})
})
