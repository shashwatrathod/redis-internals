package store

type SimpleDataStore struct {
	data map[string]*Value
}

func (s *SimpleDataStore) Put(key string, value *Value) {
	s.data[key] = value
}

func (s *SimpleDataStore) Get(key string) *Value {
	return s.data[key]
}

func (s *SimpleDataStore) Delete(key string) bool {
	if _, exists := s.data[key]; exists {
		delete(s.data, key)
		return true
	}
	return false
}
