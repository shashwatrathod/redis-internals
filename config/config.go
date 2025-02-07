package config

var Host string = "localhost"
var Port int = 7379

var LogRequest bool = false

// storage config
var MaxKeys int = 10

// eviction policy config parameters

// number of keys to sample while selecting the best candidate for removal.
var LRUEvictionSampleSize int = 5
