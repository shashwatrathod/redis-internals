package store

import (
	"time"

	"github.com/shashwatrathod/redis-internals/config"
)

type DataStore struct {
	data                 map[string]*Value
	autoDeletionStrategy AutoDeletionStrategy
	evictionStrategy     EvictionStrategy
	keyMetadata          map[string]*KeyMetadata
}

func (s *DataStore) Put(key string, value *Value) {

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
		keyMetadata.LastAccessedTimestamp = time.Now()
	}

	s.data[key] = value
	s.keyMetadata[key] = keyMetadata
}

func (s *DataStore) Get(key string) *Value {
	val := s.data[key]

	// Passively delete a key if it is found to be expired.
	if val != nil && val.Expiry != nil && val.Expiry.IsExpired() {
		s.Delete(key)
	}

	if metadata := s.GetKeyMetadata(key); metadata != nil {
		metadata.LastAccessedTimestamp = time.Now()
	}

	return s.data[key]
}

func (s *DataStore) Delete(key string) bool {
	if _, exists := s.data[key]; exists {
		delete(s.data, key)
		delete(s.keyMetadata, key)
		return true
	}
	return false
}

func (s *DataStore) Reset() {
	s.data = make(map[string]*Value)
	s.keyMetadata = make(map[string]*KeyMetadata)
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

func (s *DataStore) Evict() {
	s.evictionStrategy.Execute(s)
}

func (s *DataStore) KeyCount() int {
	return len(s.data)
}
