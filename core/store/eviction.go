package store

type EvictionStrategy interface {
	Evict(dstore Store)
}
