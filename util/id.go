package util

import (
	"fmt"
	"sync"
	"time"
)

// Constants defining the Snowflake ID structure
const (
	epoch         = int64(1577836800000) // January 1, 2020, 00:00:00 UTC
	timestampBits = 41                   // 41 bits for timestamp (69 years)
	machineIDBits = 10                   // 10 bits for machine ID (1024 machines)
	sequenceBits  = 12                   // 12 bits for sequence (4096 IDs per ms)
	maxMachineID  = 1<<machineIDBits - 1 // 1023
	maxSequence   = 1<<sequenceBits - 1  // 4095
)

// Snowflake represents a Snowflake ID generator
type Snowflake struct {
	mu            sync.Mutex
	machineID     int64
	lastTimestamp int64
	sequence      int64
}

// NewSnowflake creates a new Snowflake instance with the given machine ID
func NewSnowflake(machineID int64) (*Snowflake, error) {
	if machineID < 0 || machineID > maxMachineID {
		return nil, fmt.Errorf("machineID must be between 0 and %d", maxMachineID)
	}
	return &Snowflake{
		machineID: machineID,
	}, nil
}

// GenerateID generates a unique 64-bit Snowflake ID
func (s *Snowflake) GenerateID() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Current timestamp in milliseconds since the epoch
	now := time.Now().UnixMilli() - epoch

	// Prevent generation if the clock is before the epoch
	if now < 0 {
		panic("clock is before the epoch")
	}

	// Handle sequence increment within the same millisecond
	if now == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & maxSequence
		if s.sequence == 0 {
			// Sequence overflow: wait for the next millisecond
			for now <= s.lastTimestamp {
				now = time.Now().UnixMilli() - epoch
			}
		}
	} else {
		// New millisecond: reset sequence
		s.sequence = 0
	}

	// Update the last timestamp
	s.lastTimestamp = now

	// Construct the 64-bit ID:
	// - Timestamp: left-shifted by 22 bits (machineIDBits + sequenceBits)
	// - Machine ID: left-shifted by 12 bits (sequenceBits)
	// - Sequence: lowest 12 bits
	id := (now << (machineIDBits + sequenceBits)) |
		(s.machineID << sequenceBits) |
		s.sequence

	return id
}
