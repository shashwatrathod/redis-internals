package store

import (
	"time"

	"github.com/shashwatrathod/redis-internals/config"
)

type SimpleDataStore struct {
	data                 map[string]*Value
	autoDeletionStrategy AutoDeletionStrategy
	evictionStrategy     EvictionStrategy
	keyMetadata          map[string]*KeyMetadata
	nKeys                int
}

func (s *SimpleDataStore) Put(key string, value *Value) {

	if s.nKeys >= config.MaxKeys {
		s.Evict()
	}

	// todo: better error handling to account for inefficient eviction.
	if s.nKeys >= config.MaxKeys {
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

func (s *SimpleDataStore) Get(key string) *Value {
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

func (s *SimpleDataStore) Delete(key string) bool {
	if _, exists := s.data[key]; exists {
		delete(s.data, key)
		delete(s.keyMetadata, key)
		return true
	}
	return false
}

func (s *SimpleDataStore) Reset() {
	s.data = make(map[string]*Value)
	s.keyMetadata = make(map[string]*KeyMetadata)
	s.nKeys = 0
}

func (s *SimpleDataStore) AutoDeleteExpiredKeys() {
	s.autoDeletionStrategy.Execute(s)
}

func (s *SimpleDataStore) ForEach(fn func(key string, value *Value) bool) {
	for k, v := range s.data {
		if !fn(k, v) {
			break
		}
	}
}

func (s *SimpleDataStore) GetKeyMetadata(key string) *KeyMetadata {
	return s.keyMetadata[key]
}

func (s *SimpleDataStore) Evict() {
	s.evictionStrategy.Execute(s)
}
