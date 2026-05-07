package entity

import "fmt"

const defaultPoolCap = 1024

// Pool manages a fixed-size slab of Entity objects to avoid GC pressure.
// All allocation and release must happen on the game loop goroutine.
type Pool struct {
	slab    []Entity
	free    []int // indices into slab that are available
	nextID  uint32
}

func NewPool(capacity int) *Pool {
	if capacity <= 0 {
		capacity = defaultPoolCap
	}
	slab := make([]Entity, capacity)
	free := make([]int, capacity)
	for i := range free {
		free[i] = i
		slab[i].poolIndex = i
	}
	return &Pool{slab: slab, free: free, nextID: 1}
}

// Acquire returns a zeroed Entity from the pool. Returns an error when full.
func (p *Pool) Acquire(kind Kind) (*Entity, error) {
	if len(p.free) == 0 {
		return nil, fmt.Errorf("entity pool exhausted (cap=%d)", len(p.slab))
	}
	idx := p.free[len(p.free)-1]
	p.free = p.free[:len(p.free)-1]

	e := &p.slab[idx]
	e.reset()
	e.ID = p.nextID
	e.Kind = kind
	e.alive = true
	p.nextID++
	return e, nil
}

// Release returns an entity back to the pool.
func (p *Pool) Release(e *Entity) {
	e.reset()
	p.free = append(p.free, e.poolIndex)
}

// Len returns the number of currently live (acquired) entities.
func (p *Pool) Len() int { return len(p.slab) - len(p.free) }
