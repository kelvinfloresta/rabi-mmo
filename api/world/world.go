package world

import "fmt"

// World holds all floors and acts as the authoritative collision/movement layer.
type World struct {
	floors map[int8]*Floor
}

func New() *World {
	return &World{floors: make(map[int8]*Floor)}
}

func (w *World) AddFloor(floorIndex int8, f *Floor) {
	w.floors[floorIndex] = f
}

func (w *World) Floor(floorIndex int8) (*Floor, error) {
	f, ok := w.floors[floorIndex]
	if !ok {
		return nil, fmt.Errorf("floor %d does not exist", floorIndex)
	}
	return f, nil
}

// CanMove returns true if an entity can move from pos by (dx,dy) on the same floor.
func (w *World) CanMove(pos Position, dx, dy int16) bool {
	f, err := w.Floor(pos.Floor)
	if err != nil {
		return false
	}
	nx, ny := int(pos.X+dx), int(pos.Y+dy)
	return f.CanEnter(nx, ny)
}

// MoveEntity relocates an entity, updating both source and destination tiles.
// Returns the new position and true on success.
func (w *World) MoveEntity(id uint32, from Position, dx, dy int16) (Position, bool) {
	f, err := w.Floor(from.Floor)
	if err != nil {
		return from, false
	}

	nx, ny := int(from.X+dx), int(from.Y+dy)
	if !f.PlaceEntity(nx, ny, id) {
		return from, false
	}
	f.RemoveEntity(int(from.X), int(from.Y))

	return Position{X: from.X + dx, Y: from.Y + dy, Floor: from.Floor}, true
}

// TileAt returns the tile at an absolute world position.
func (w *World) TileAt(pos Position) (*Tile, error) {
	f, err := w.Floor(pos.Floor)
	if err != nil {
		return nil, err
	}
	return f.At(int(pos.X), int(pos.Y))
}
