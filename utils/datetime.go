package utils

import (
	"time"
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

	return int64(timeRemaining.Seconds())
}

// Returns true if the expiry time has already passed.
func (et ExpiryTime) IsExpired() bool {
	return et.getTimeRemaining() <= 0
}

// Returns the expiry timestamp.
func (et ExpiryTime) GetExpireAtTimestamp() time.Time {
	return et.expireAtTimestamp
}
