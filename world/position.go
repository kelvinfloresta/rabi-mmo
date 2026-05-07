package world

// Position represents a tile coordinate in the world, including floor.
type Position struct {
	X, Y  int16
	Floor int8
}

func (p Position) Equal(other Position) bool {
	return p.X == other.X && p.Y == other.Y && p.Floor == other.Floor
}

func (p Position) Add(dx, dy int16) Position {
	return Position{X: p.X + dx, Y: p.Y + dy, Floor: p.Floor}
}
