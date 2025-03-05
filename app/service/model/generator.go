package model

import (
	"errors"
	"sync"
	"time"
)

var (
	generators = make(map[string]*Snowflake)
	lock       sync.Mutex
)

func GetOrCreateGenerator(tableName string, nodeID int64) *Snowflake {
	lock.Lock()
	defer lock.Unlock()

	epoch := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	if gen, exists := generators[tableName]; exists {
		return gen
	}

	gen, err := NewSnowflake(epoch, nodeID)
	if err != nil {
		panic(err)
	}

	generators[tableName] = gen
	return gen
}

type Snowflake struct {
	mutex         sync.Mutex
	epoch         int64
	machineID     int64
	sequence      int64
	lastTimestamp int64
}

const (
	timestampBits = 41
	sequenceBits  = 8
	machineIDBits = 2

	sequenceMax    = (1 << sequenceBits) - 1
	machineIDShift = sequenceBits
	timestampShift = sequenceBits + machineIDBits
	timestampMask  = (1 << timestampBits) - 1
)

func NewSnowflake(epoch time.Time, machineID int64) (*Snowflake, error) {
	if machineID < 0 || machineID >= (1<<machineIDBits) {
		return nil, errors.New("machineID out of range")
	}
	return &Snowflake{
		epoch:     epoch.UnixMilli(),
		machineID: machineID,
	}, nil
}

func (s *Snowflake) GenerateID() (int64, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	now := time.Now().UnixMilli()
	if now < s.lastTimestamp {
		return 0, errors.New("clock moved backwards")
	}

	if now == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & sequenceMax
		if s.sequence == 0 {
			for now <= s.lastTimestamp {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		s.sequence = 0
	}

	s.lastTimestamp = now
	id := ((now - s.epoch) << timestampShift) | (s.machineID << machineIDShift) | s.sequence

	return id, nil
}

func (s *Snowflake) ExtractTimeByID(id int64) time.Time {
	timestamp := (id >> timestampShift) & timestampMask
	return time.UnixMilli(s.epoch + timestamp)
}
