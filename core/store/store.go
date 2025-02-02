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

type Store interface {
	Put(key string, value *Value)
	Get(key string) *Value
	Delete(key string) bool
	Reset()
}

// Represents a Value that can be stored in the datastore.
type Value struct {
	Value     interface{}
	ValueType SupportedDatatypes
	Expiry    *utils.ExpiryTime
}

var storeInstance *SimpleDataStore

func GetStore() *SimpleDataStore {
	if storeInstance == nil {
		storeInstance = &SimpleDataStore{
			data: make(map[string]*Value),
		}
	}

	return storeInstance
}
