package utils

import (
	"math"
	"time"

	"github.com/shashwatrathod/redis-internals/config"
)

// Represents the Expiry Time of a Key in the store.
type ExpiryTime struct {
	expireAtTimestamp time.Time
}

// Returns the ExpiryTime after expiryInSeconds seconds from Now.
func FromExpiryInSeconds(expiryInSeconds int64) *ExpiryTime {
	now := time.Now()
	now = now.Add(time.Duration(expiryInSeconds) * time.Second)
	return &ExpiryTime{expireAtTimestamp: now}
}

// Returns the ExpiryTime after expiryInMs milliseconds from Now.
func FromExpiryInMilliseconds(expiryInMs int64) *ExpiryTime {
	now := time.Now()
	now = now.Add(time.Duration(expiryInMs) * time.Millisecond)
	return &ExpiryTime{expireAtTimestamp: now}
}

func (et ExpiryTime) getTimeRemaining() time.Duration {
	return time.Until(et.expireAtTimestamp)
}

// Returns the time remaining in reaching expiry from Now in seconds.
// Returns -2 if the expiry has already passed.
func (et ExpiryTime) GetTimeRemainingInSeconds() int64 {
	timeRemaining := et.getTimeRemaining()

	if timeRemaining < 0 {
		return -2
	}

	return int64(math.Round(timeRemaining.Seconds()))
}

// Returns true if the expiry time has already passed.
func (et ExpiryTime) IsExpired() bool {
	return et.getTimeRemaining() <= 0
}

// Returns the expiry timestamp.
func (et ExpiryTime) GetExpireAtTimestamp() time.Time {
	return et.expireAtTimestamp
}

// the LRU-Time is represented as a 32-bit integer.
type LRUTime uint32

// only the least-significant 24-bits of the LRU-time contain any data.
// this is inline with redis's implementation of a 24-bit lru clock, making a trade-off between resolution and memory consumption.
// https://github.com/redis/redis/blob/unstable/src/evict.c

// returns the current time as an LRUTime value.
// the time is obtained from the current Unix timestamp and is masked
// with the LRUTimeResolution configuration value to fit the LRU time format.
func GetCurrentLruTime() LRUTime {
	return ToLRUTime(time.Now())
}

// converts the given time instance into LRU time.
func ToLRUTime(t time.Time) LRUTime {
	return LRUTime(uint32(t.Unix()) & config.LRUTimeResolution)
}
