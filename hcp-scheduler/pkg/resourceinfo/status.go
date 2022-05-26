package resourceinfo

import (
	"errors"
	"sync"
)

var (
	// ErrNotFound is the not found error message.
	ErrNotFound = errors.New("not found")
)

// StateData is a generic type for arbitrary data stored in CycleState.
type StateData interface {
	// Clone is an interface to make a copy of StateData. For performance reasons,
	// clone should make shallow copies for members (e.g., slices or maps) that are not
	// impacted by PreFilter's optional AddPod/RemovePod methods.
	Clone() StateData
}

// StateKey is the type of keys stored in CycleState.
type StateKey string

// 전체 스케줄링 레이어에 대한 기록
type CycleStatus struct {
	mx             sync.RWMutex
	storage        map[StateKey]StateData
	AvailableNodes int
	TotalNumNodes  int
}

func NewCycleStatus(totalNum int) *CycleStatus {
	return &CycleStatus{
		TotalNumNodes:  totalNum,
		AvailableNodes: totalNum,
	}
}

func (c *CycleStatus) MinusOneAvailableNodes() {
	c.AvailableNodes--
}

func (c *CycleStatus) IsAnyClusters() bool {
	if c.AvailableNodes == 0 {
		return false
	} else {
		return true
	}
}

// Clone creates a copy of CycleState and returns its pointer. Clone returns
// nil if the context being cloned is nil.
func (c *CycleStatus) Clone() *CycleStatus {
	if c == nil {
		return nil
	}
	copy := NewCycleStatus(c.TotalNumNodes)
	for k, v := range c.storage {
		copy.Write(k, v.Clone())
	}
	return copy
}

// Read retrieves data with the given "key" from CycleState. If the key is not
// present an error is returned.
// This function is thread safe by acquiring an internal lock first.
func (c *CycleStatus) Read(key StateKey) (StateData, error) {
	c.mx.RLock()
	defer c.mx.RUnlock()
	if v, ok := c.storage[key]; ok {
		return v, nil
	}
	return nil, ErrNotFound
}

// Write stores the given "val" in CycleState with the given "key".
// This function is thread safe by acquiring an internal lock first.
func (c *CycleStatus) Write(key StateKey, val StateData) {
	c.mx.Lock()
	c.storage[key] = val
	c.mx.Unlock()
}
