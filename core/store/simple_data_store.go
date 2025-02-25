package store

import (
	"log"

	"github.com/shashwatrathod/redis-internals/config"
	"github.com/shashwatrathod/redis-internals/utils"
)

type DataStore struct {
	data                 map[string]*Value
	autoDeletionStrategy AutoDeletionStrategy
	evictionStrategy     EvictionStrategy
	keyMetadata          map[string]*KeyMetadata
	expiries             map[string]*int64
}

func (s *DataStore) Put(key string, value string, expiry *utils.ExpiryTime) {
	if s.KeyCount() >= config.MaxKeys {
		s.Evict()
	}

	// todo: better error handling to account for inefficient eviction.
	if s.KeyCount() >= config.MaxKeys {
		return
	}

	var keyMetadata *KeyMetadata = newKeyMetadata()

	if s.data[key] != nil && s.keyMetadata[key] != nil {
		keyMetadata = s.GetKeyMetadata(key)
		// Update the LastAccessedTs to Now if the key already exists.
		keyMetadata.LastAccessedTimestamp = utils.GetCurrentLruTime()
	}

	s.data[key] = &Value{
		Value:     value,
		ValueType: String, // TODO: Add more types
	}
	s.keyMetadata[key] = keyMetadata

	s.SetExpiry(key, expiry)
}

func (s *DataStore) GetExpiry(key string) *int64 {
	exp, exists := s.expiries[key]

	if !exists || exp == nil {
		return nil
	}

	return exp
}

func (s *DataStore) Get(key string) *Value {
	_, exists := s.data[key]

	// Passively delete a key if it is found to be expired.
	if exists && s.isExpired(key) {
		s.Delete(key)
	}

	if metadata := s.GetKeyMetadata(key); metadata != nil {
		metadata.LastAccessedTimestamp = utils.GetCurrentLruTime()
	}

	return s.data[key]
}

// returns whether the given key has expired. returns false if the key doesn't exist,
// or if there is no expiry set on the key.
func (s *DataStore) isExpired(key string) bool {
	exp, exists := s.expiries[key]

	if !exists || exp == nil {
		return false
	}

	return utils.FromExpiryInUnixTime(*exp).IsExpired()
}

func (s *DataStore) SetExpiry(key string, expiry *utils.ExpiryTime) {

	if exists := s.Get(key); exists == nil {
		return
	}

	if expiry == nil {
		delete(s.expiries, key)
		return
	}

	timestamp := expiry.ToUnixTimestamp()
	s.expiries[key] = &timestamp
}

func (s *DataStore) Delete(key string) bool {
	if _, exists := s.data[key]; exists {
		delete(s.data, key)
		delete(s.keyMetadata, key)
		delete(s.expiries, key)
		return true
	}
	return false
}

func (s *DataStore) Reset() {
	s.data = make(map[string]*Value)
	s.keyMetadata = make(map[string]*KeyMetadata)
	s.expiries = make(map[string]*int64)
}

func (s *DataStore) AutoDeleteExpiredKeys() {
	s.autoDeletionStrategy.Execute(s)
}

func (s *DataStore) ForEach(fn func(key string, value *Value) bool) {
	for k, v := range s.data {
		if !fn(k, v) {
			break
		}
	}
}

func (s *DataStore) GetKeyMetadata(key string) *KeyMetadata {
	return s.keyMetadata[key]
}

func (s *DataStore) Evict() int {
	nKeysEvicted, err := s.evictionStrategy.Execute(s)

	if err != nil {
		log.Printf("Encountered an error while trying to evict keys : %s", err.Error())
	}

	return nKeysEvicted
}

func (s *DataStore) KeyCount() int {
	return len(s.data)
}
