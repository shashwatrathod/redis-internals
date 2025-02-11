package store_test

import (
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/shashwatrathod/redis-internals/config"
	"github.com/shashwatrathod/redis-internals/core/store"
	"github.com/shashwatrathod/redis-internals/utils"
)

func TestSimpleDataStore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Store Suite")
}

var _ = Describe("SimpleDataStore", func() {
	var (
		dataStore *store.DataStore
	)

	BeforeEach(func() {
		dataStore = store.GetStore()
	})

	AfterEach(func() {
		dataStore.Reset()
	})

	Describe("Put", func() {
		It("should store a value", func() {
			value := &store.Value{
				Value:     "value",
				ValueType: store.String,
				Expiry:    nil,
			}
			dataStore.Put("key", value)
			Expect(dataStore.Get("key")).To(Equal(value))
		})
		It("should evict keys if the datastore is at max capacity", func() {
			// fill the store to its max capacity
			for i := 1; i <= config.MaxKeys; i++ {
				dataStore.Put(fmt.Sprintf("key%d", i+1), &store.Value{
					Value:     fmt.Sprintf("value%d", i+1),
					ValueType: store.String,
					Expiry:    nil,
				})
			}

			key := fmt.Sprintf("key%d", config.MaxKeys+1)
			val := &store.Value{
				Value:     fmt.Sprintf("value%d", config.MaxKeys+1),
				ValueType: store.String,
				Expiry:    nil,
			}
			dataStore.Put(key, val)
			Expect(dataStore.Get(key)).To(Equal(val))
		})
	})

	Describe("Get", func() {
		It("should retrieve a stored value", func() {
			value := &store.Value{
				Value:     "value",
				ValueType: store.String,
				Expiry:    nil,
			}
			dataStore.Put("key", value)
			Expect(dataStore.Get("key")).To(Equal(value))
		})

		It("should return nil for a non-existent key", func() {
			Expect(dataStore.Get("nonexistent")).To(BeNil())
		})

		It("should delete and return nil for an expired key", func() {
			expiryTime := utils.FromExpiryInMilliseconds(-1000)
			value := &store.Value{
				Value:     "value",
				ValueType: store.String,
				Expiry:    expiryTime,
			}
			dataStore.Put("key", value)
			Expect(dataStore.Get("key")).To(BeNil())
		})
	})

	Describe("Delete", func() {
		It("should delete a stored value", func() {
			value := &store.Value{
				Value:     "value",
				ValueType: store.String,
				Expiry:    nil,
			}
			dataStore.Put("key", value)
			Expect(dataStore.Delete("key")).To(BeTrue())
			Expect(dataStore.Get("key")).To(BeNil())
		})

		It("should return false for a non-existent key", func() {
			Expect(dataStore.Delete("nonexistent")).To(BeFalse())
		})
	})
})
