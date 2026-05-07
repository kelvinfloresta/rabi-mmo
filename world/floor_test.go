package world

import "testing"

func TestFloorCollision(t *testing.T) {
	f := NewFloor(5, 5)

	// Center is floor by default — should be walkable.
	if !f.CanEnter(2, 2) {
		t.Fatal("center tile should be walkable")
	}

	// Place a wall and verify it blocks entry.
	_ = f.SetTerrain(2, 2, TerrainWall)
	if f.CanEnter(2, 2) {
		t.Fatal("wall tile should not be walkable")
	}

	// Out-of-bounds should return false, not panic.
	if f.CanEnter(-1, 0) || f.CanEnter(0, 99) {
		t.Fatal("out-of-bounds should not be enterable")
	}
}

func TestEntityOccupancy(t *testing.T) {
	f := NewFloor(5, 5)

	if !f.PlaceEntity(1, 1, 42) {
		t.Fatal("should place entity on empty floor tile")
	}

	// Same tile must be blocked by the existing entity.
	if f.PlaceEntity(1, 1, 99) {
		t.Fatal("should not place second entity on occupied tile")
	}

	f.RemoveEntity(1, 1)
	if !f.CanEnter(1, 1) {
		t.Fatal("tile should be free after RemoveEntity")
	}
}

func TestDoorPassability(t *testing.T) {
	f := NewFloor(3, 3)
	_ = f.SetTerrain(1, 1, TerrainDoor)

	tile, _ := f.At(1, 1)
	if f.CanEnter(1, 1) {
		t.Fatal("closed door should not be enterable")
	}

	tile.DoorOpen = true
	if !f.CanEnter(1, 1) {
		t.Fatal("open door should be enterable")
	}
}
