package combat

type Element uint8

const (
	ElementPhysical Element = iota
	ElementFire
	ElementWater
	ElementEarth
	ElementEnergy
	ElementDeath
	ElementIce
	ElementAir
	ElementVoid
	elementCount
)

// resistanceTable[attacker][defender] = multiplier in percent (100 = normal, 50 = half, 200 = double)
var resistanceTable [elementCount][elementCount]int16

func init() {
	for a := Element(0); a < elementCount; a++ {
		for d := Element(0); d < elementCount; d++ {
			resistanceTable[a][d] = 100
		}
	}

	// Fire melts Ice, Water resists Fire
	resistanceTable[ElementFire][ElementIce] = 150
	resistanceTable[ElementFire][ElementWater] = 50
	resistanceTable[ElementFire][ElementEarth] = 75

	// Water extinguishes Fire, weak vs Energy
	resistanceTable[ElementWater][ElementFire] = 150
	resistanceTable[ElementWater][ElementEnergy] = 50

	// Ice amplified by Fire (already covered), weak vs Earth
	resistanceTable[ElementIce][ElementEarth] = 50
	resistanceTable[ElementIce][ElementAir] = 125

	// Energy weak vs Earth, strong vs Water
	resistanceTable[ElementEnergy][ElementEarth] = 50
	resistanceTable[ElementEnergy][ElementWater] = 150

	// Death strong vs Physical, weak vs Void
	resistanceTable[ElementDeath][ElementPhysical] = 125
	resistanceTable[ElementDeath][ElementVoid] = 50

	// Void amplifies everything on Death-aligned targets
	resistanceTable[ElementVoid][ElementDeath] = 150

	// Air weak vs Earth, strong vs Fire (fans the flames — thematic only)
	resistanceTable[ElementAir][ElementEarth] = 50
	resistanceTable[ElementAir][ElementFire] = 75
}

// CalcDamage returns the final damage after applying elemental resistance.
func CalcDamage(base int, attackElement, defenseElement Element) int {
	mult := int(resistanceTable[attackElement][defenseElement])
	return base * mult / 100
}
