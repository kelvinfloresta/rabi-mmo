package world

import "fmt"

// Floor is one layer of the world (a 2-D grid of tiles).
type Floor struct {
	Width, Height int
	tiles         []Tile // row-major: tiles[y*Width + x]
}

func NewFloor(width, height int) *Floor {
	return &Floor{
		Width:  width,
		Height: height,
		tiles:  make([]Tile, width*height),
	}
}

func (f *Floor) InBounds(x, y int) bool {
	return x >= 0 && y >= 0 && x < f.Width && y < f.Height
}

func (f *Floor) At(x, y int) (*Tile, error) {
	if !f.InBounds(x, y) {
		return nil, fmt.Errorf("position (%d,%d) out of bounds", x, y)
	}
	return &f.tiles[y*f.Width+x], nil
}

// SetTerrain sets the terrain for a tile, used during map loading.
func (f *Floor) SetTerrain(x, y int, t TerrainType) error {
	tile, err := f.At(x, y)
	if err != nil {
		return err
	}
	tile.Terrain = t
	return nil
}

// CanEnter returns whether the tile at (x,y) can be entered by an entity.
// It checks walkability and whether another entity already occupies the tile.
func (f *Floor) CanEnter(x, y int) bool {
	if !f.InBounds(x, y) {
		return false
	}
	tile := &f.tiles[y*f.Width+x]
	return tile.Walkable() && tile.EntityID == 0
}

// PlaceEntity occupies a tile with the given entity ID. Returns false if blocked.
func (f *Floor) PlaceEntity(x, y int, id uint32) bool {
	if !f.CanEnter(x, y) {
		return false
	}
	f.tiles[y*f.Width+x].EntityID = id
	return true
}

// RemoveEntity clears the entity occupying a tile.
func (f *Floor) RemoveEntity(x, y int) {
	if f.InBounds(x, y) {
		f.tiles[y*f.Width+x].EntityID = 0
	}
}
