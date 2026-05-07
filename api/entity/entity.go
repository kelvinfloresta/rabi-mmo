package entity

import (
	"github.com/kelvinfloresta/rabi-mmo/combat"
	"github.com/kelvinfloresta/rabi-mmo/world"
)

type Kind uint8

const (
	KindPlayer Kind = iota
	KindMonster
	KindNPC
)

// Entity is the shared base for all world actors.
type Entity struct {
	ID      uint32
	Kind    Kind
	Pos     world.Position
	alive   bool

	// Combat stats
	HP        int
	MaxHP     int
	MP        int
	MaxMP     int
	BaseAtk   int
	BaseDef   int
	AtkElem   combat.Element
	DefElem   combat.Element

	// Pool bookkeeping
	poolIndex int
}

func (e *Entity) IsAlive() bool { return e.alive }

func (e *Entity) TakeDamage(base int, attackElem combat.Element) int {
	dmg := combat.CalcDamage(base, attackElem, e.DefElem)
	e.HP -= dmg
	if e.HP <= 0 {
		e.HP = 0
		e.alive = false
	}
	return dmg
}

func (e *Entity) Heal(amount int) {
	e.HP += amount
	if e.HP > e.MaxHP {
		e.HP = e.MaxHP
	}
}

func (e *Entity) reset() {
	e.ID = 0
	e.alive = false
	e.HP = 0
	e.MaxHP = 0
	e.MP = 0
	e.MaxMP = 0
	e.BaseAtk = 0
	e.BaseDef = 0
	e.AtkElem = combat.ElementPhysical
	e.DefElem = combat.ElementPhysical
	e.Pos = world.Position{}
}
