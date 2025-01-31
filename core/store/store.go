package store

import (
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

type Store struct {
	data map[string]*Value
}

// Represents a Value that can be stored in the datastore.
type Value struct {
	Value     interface{}
	ValueType SupportedDatatypes
	Expiry    *utils.ExpiryTime
}

var storeInstance *Store

func GetStore() *Store {
	if storeInstance == nil {
		storeInstance = &Store{
			data: make(map[string]*Value),
		}
	}

	return storeInstance
}

func (s *Store) Put(key string, value *Value) {
	s.data[key] = value
}

func (s *Store) Get(key string) *Value {
	return s.data[key]
}

func (s *Store) Delete(key string) bool {
	if _, exists := s.data[key]; exists {
		delete(s.data, key)
		return true
	}
	return false
}
