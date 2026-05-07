package world

// TerrainType classifies each tile for collision and gameplay logic.
type TerrainType uint8

const (
	TerrainFloor    TerrainType = iota // walkable, empty floor
	TerrainWall                        // impassable solid wall
	TerrainWater                       // passable only by water-walking entities
	TerrainLava                        // passable, deals fire damage per tick
	TerrainMountain                    // impassable natural obstacle
	TerrainDoor                        // passable when open, impassable when closed
	TerrainStairUp                     // transition to floor+1
	TerrainStairDown                   // transition to floor-1
)

// Tile is one cell in the map grid.
type Tile struct {
	Terrain    TerrainType
	DoorOpen   bool   // only meaningful when Terrain == TerrainDoor
	EntityID   uint32 // 0 = empty; ID of entity currently occupying this tile
}

func (t *Tile) Walkable() bool {
	switch t.Terrain {
	case TerrainFloor, TerrainLava, TerrainStairUp, TerrainStairDown:
		return true
	case TerrainDoor:
		return t.DoorOpen
	default:
		return false
	}
}
