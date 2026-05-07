package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kelvinfloresta/rabi-mmo/combat"
	"github.com/kelvinfloresta/rabi-mmo/entity"
	"github.com/kelvinfloresta/rabi-mmo/game"
	"github.com/kelvinfloresta/rabi-mmo/mapfile"
	"github.com/kelvinfloresta/rabi-mmo/world"
)

func main() {
	w := world.New()

	for _, spec := range []struct {
		index int8
		path  string
	}{
		{0, "maps/floor0.map"},
		{1, "maps/floor1.map"},
	} {
		floor, err := mapfile.LoadFloor(spec.path)
		if err != nil {
			log.Fatalf("load floor %d: %v", spec.index, err)
		}
		w.AddFloor(spec.index, floor)
	}

	g := game.New(w, 1024)

	// Spawn a player and a monster as a smoke-test.
	player, err := g.SpawnEntity(entity.KindPlayer,
		world.Position{X: 2, Y: 2, Floor: 0},
		100, 50, 10, 5)
	if err != nil {
		log.Fatalf("spawn player: %v", err)
	}
	player.AtkElem = combat.ElementFire
	player.DefElem = combat.ElementPhysical

	monster, err := g.SpawnEntity(entity.KindMonster,
		world.Position{X: 5, Y: 2, Floor: 0},
		60, 0, 8, 3)
	if err != nil {
		log.Fatalf("spawn monster: %v", err)
	}
	monster.DefElem = combat.ElementIce

	// Try to move player into a wall — should be silently blocked.
	g.EnqueueCommand(game.Command{EntityID: player.ID, Type: game.CmdMove, DX: -2, DY: 0})

	// Move player right.
	g.EnqueueCommand(game.Command{EntityID: player.ID, Type: game.CmdMove, DX: 1, DY: 0})

	// Attack the monster (Fire vs Ice → 150% damage).
	g.EnqueueCommand(game.Command{EntityID: player.ID, Type: game.CmdAttack, TargetID: monster.ID})

	stop := make(chan struct{})
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sig
		fmt.Println("\nshutting down")
		close(stop)
	}()

	fmt.Printf("rabi-mmo running — player id=%d  monster id=%d\n", player.ID, monster.ID)
	fmt.Println("press Ctrl+C to stop")

	g.Run(stop)
}
