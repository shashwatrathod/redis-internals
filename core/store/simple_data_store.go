package store

import "log"

type SimpleDataStore struct {
	data map[string]*Value
}

func (s *SimpleDataStore) Put(key string, value *Value) {
	s.data[key] = value
}

func (s *SimpleDataStore) Get(key string) *Value {
	val := s.data[key]

	// Passively delete a key if it found to be expired.
	if val != nil && val.Expiry != nil && val.Expiry.IsExpired() {
		s.Delete(key)
	}

	return s.data[key]
}

func (s *SimpleDataStore) Delete(key string) bool {
	if _, exists := s.data[key]; exists {
		delete(s.data, key)
		return true
	}
	return false
}

func (s *SimpleDataStore) Reset() {
	s.data = make(map[string]*Value)
}

func (s *SimpleDataStore) expireSample() float32 {
	var nExpired int = 0
	var nSearched int = 0

	for key, val := range s.data {
		if val.Expiry != nil {
			nSearched++

			if val.Expiry.IsExpired() {
				s.Delete(key)
				nExpired++
			}
		}

		if nSearched == AUTO_EXPIRE_SEARCH_LIMIT {
			break
		}
	}

	return float32(nExpired) / float32(AUTO_EXPIRE_SEARCH_LIMIT)
}

func (s *SimpleDataStore) AutoDeleteExpiredKeys() {
	for {
		fracExpired := s.expireSample()

		if fracExpired < AUTO_EXPIRE_ALLOWABLE_EXPIRE_FRACTION {
			break
		}
	}

	log.Println("auto-deleted the expired but undeleted keys. total keys", len(s.data))
}
