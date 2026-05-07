package mapfile

import (
	"bufio"
	"fmt"
	"os"

	"github.com/kelvinfloresta/rabi-mmo/world"
)

// charTerrain maps ASCII characters in map files to terrain types.
var charTerrain = map[byte]world.TerrainType{
	'.': world.TerrainFloor,
	'#': world.TerrainWall,
	'~': world.TerrainWater,
	'^': world.TerrainMountain,
	'L': world.TerrainLava,
	'+': world.TerrainDoor,
	'<': world.TerrainStairUp,
	'>': world.TerrainStairDown,
}

// LoadFloor reads a plain-text map file and returns a populated Floor.
//
// File format (example):
//
//	########
//	#......#
//	#..~~..#
//	#......#
//	########
//
// Each row must have the same length. Unknown characters default to TerrainFloor.
func LoadFloor(path string) (*world.Floor, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("mapfile: open %q: %w", path, err)
	}
	defer f.Close()

	var rows []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		rows = append(rows, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("mapfile: scan %q: %w", path, err)
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("mapfile: %q is empty", path)
	}

	width := len(rows[0])
	for i, row := range rows {
		if len(row) != width {
			return nil, fmt.Errorf("mapfile: row %d has width %d, expected %d", i, len(row), width)
		}
	}
	height := len(rows)

	floor := world.NewFloor(width, height)
	for y, row := range rows {
		for x := range row {
			t, ok := charTerrain[row[x]]
			if !ok {
				t = world.TerrainFloor
			}
			if setErr := floor.SetTerrain(x, y, t); setErr != nil {
				return nil, setErr
			}
		}
	}

	return floor, nil
}
