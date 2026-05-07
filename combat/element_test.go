package combat

import "testing"

func TestFireVsIce(t *testing.T) {
	dmg := CalcDamage(100, ElementFire, ElementIce)
	if dmg != 150 {
		t.Fatalf("Fire vs Ice: want 150, got %d", dmg)
	}
}

func TestFireVsWater(t *testing.T) {
	dmg := CalcDamage(100, ElementFire, ElementWater)
	if dmg != 50 {
		t.Fatalf("Fire vs Water: want 50, got %d", dmg)
	}
}

func TestPhysicalNeutral(t *testing.T) {
	dmg := CalcDamage(100, ElementPhysical, ElementPhysical)
	if dmg != 100 {
		t.Fatalf("Physical vs Physical: want 100, got %d", dmg)
	}
}
