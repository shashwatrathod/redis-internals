package core

import "github.com/shashwatrathod/redis-internals/utils"

// Represents the DataTypes currently supported by the Application.
type SupportedDatatypes int

const (
	String  SupportedDatatypes = SupportedDatatypes(BulkString)
	Integer                    = RespInteger
	Array                      = RespArray
)

// Represents a Value that can be stored in the datastore.
type Value struct {
	value     interface{}
	valueType SupportedDatatypes
	expiry    *utils.ExpiryTime
}

var store map[string]*Value

func init() {
	store = make(map[string]*Value)
}

func Put(key string, value *Value) {
	store[key] = value
}

func Get(key string) *Value {
	return store[key]
}
