package config

var Host string = "localhost"
var Port int = 7379

// logs the raw request body if set to true.
var LogRequest bool = false

// storage config

// maximum number of keys that can store in the store before eviction kicks in
var MaxKeys int = 100

// eviction policy config parameters

// ratio of keys to be evicted everytime the an eviction is performed.
var EvictionRatio float32 = 0.20

// number of keys to sample while selecting the best candidate for removal.
var LRUEvictionSampleSize int = 5
