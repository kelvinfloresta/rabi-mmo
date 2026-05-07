package game

import (
	"time"

	"github.com/kelvinfloresta/rabi-mmo/entity"
	"github.com/kelvinfloresta/rabi-mmo/world"
)

const tickRate = 100 * time.Millisecond // 10 ticks/s

// Command is an action queued by an external system (network handler, AI, etc.)
// and consumed by the game loop each tick.
type Command struct {
	EntityID uint32
	Type     CommandType
	DX, DY   int16 // for move commands
	TargetID uint32 // for attack commands
}

type CommandType uint8

const (
	CmdMove CommandType = iota
	CmdAttack
	CmdStop
)

// Game is the top-level single-threaded game state.
type Game struct {
	World    *world.World
	Pool     *entity.Pool
	entities map[uint32]*entity.Entity // live entity index

	cmdQueue []Command
	tick     uint64
}

func New(w *world.World, poolCap int) *Game {
	return &Game{
		World:    w,
		Pool:     entity.NewPool(poolCap),
		entities: make(map[uint32]*entity.Entity),
		cmdQueue: make([]Command, 0, 256),
	}
}

// EnqueueCommand is the ONLY entry point for external input.
// Must be called from the game loop goroutine (or before Run starts).
func (g *Game) EnqueueCommand(cmd Command) {
	g.cmdQueue = append(g.cmdQueue, cmd)
}

// SpawnEntity creates and registers a new entity at the given position.
func (g *Game) SpawnEntity(kind entity.Kind, pos world.Position, hp, mp, atk, def int) (*entity.Entity, error) {
	f, err := g.World.Floor(pos.Floor)
	if err != nil {
		return nil, err
	}

	e, err := g.Pool.Acquire(kind)
	if err != nil {
		return nil, err
	}
	e.Pos = pos
	e.MaxHP, e.HP = hp, hp
	e.MaxMP, e.MP = mp, mp
	e.BaseAtk = atk
	e.BaseDef = def

	if !f.PlaceEntity(int(pos.X), int(pos.Y), e.ID) {
		g.Pool.Release(e)
		return nil, &PositionBlockedError{Pos: pos}
	}

	g.entities[e.ID] = e
	return e, nil
}

// DespawnEntity removes an entity from the world and returns it to the pool.
func (g *Game) DespawnEntity(id uint32) {
	e, ok := g.entities[id]
	if !ok {
		return
	}
	f, err := g.World.Floor(e.Pos.Floor)
	if err == nil {
		f.RemoveEntity(int(e.Pos.X), int(e.Pos.Y))
	}
	delete(g.entities, id)
	g.Pool.Release(e)
}

// Run starts the blocking single-threaded game loop.
// Cancel the returned stop channel to exit.
func (g *Game) Run(stop <-chan struct{}) {
	ticker := time.NewTicker(tickRate)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			g.tick++
			g.processTick()
		}
	}
}

func (g *Game) processTick() {
	for _, cmd := range g.cmdQueue {
		g.handleCommand(cmd)
	}
	g.cmdQueue = g.cmdQueue[:0]
}

func (g *Game) handleCommand(cmd Command) {
	e, ok := g.entities[cmd.EntityID]
	if !ok || !e.IsAlive() {
		return
	}

	switch cmd.Type {
	case CmdMove:
		g.handleMove(e, cmd.DX, cmd.DY)
	case CmdAttack:
		g.handleAttack(e, cmd.TargetID)
	}
}

func (g *Game) handleMove(e *entity.Entity, dx, dy int16) {
	newPos, ok := g.World.MoveEntity(e.ID, e.Pos, dx, dy)
	if ok {
		e.Pos = newPos
	}
}

func (g *Game) handleAttack(attacker *entity.Entity, targetID uint32) {
	target, ok := g.entities[targetID]
	if !ok || !target.IsAlive() {
		return
	}

	dmg := target.TakeDamage(attacker.BaseAtk, attacker.AtkElem)
	_ = dmg // hook for event broadcast later

	if !target.IsAlive() {
		g.DespawnEntity(target.ID)
	}
}

// PositionBlockedError is returned when SpawnEntity cannot place on a tile.
type PositionBlockedError struct {
	Pos world.Position
}

func (e *PositionBlockedError) Error() string {
	return "position is blocked or unwalkable"
}
