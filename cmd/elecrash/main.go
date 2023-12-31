package main

import (
	"flag"
	"log"
	"strconv"

	ui "github.com/gizak/termui/v3"
	"github.com/lovromazgon/elecrash"
)

func main() {
	var (
		elevators  = flag.Int("e", 4, "number of elevators")
		floors     = flag.Int("f", 8, "number of floors")
		difficulty = flag.Int("d", 1, "difficulty (1-10)")
	)
	flag.Parse()

	if *elevators < 2 || *elevators > 10 {
		log.Fatal("need between 2 and 10 elevators")
	}
	if *floors < 4 || *floors > 10 {
		log.Fatal("need between 4 and 10 floors")
	}
	if *difficulty < 1 || *difficulty > 10 {
		log.Fatal("need between 1 and 10 difficulty")
	}

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	closeLogger := elecrash.InitLogger()
	defer closeLogger()

	game := elecrash.NewElecrash(*elevators, *floors, float64(*difficulty))
	go game.Run()

	for e := range ui.PollEvents() {
		if e.Type != ui.KeyboardEvent {
			continue
		}
		switch e.ID {
		case "<Left>", "h", "H":
			game.Left()
		case "<Right>", "l", "L":
			game.Right()
		case "1", "2", "3", "4", "5", "6", "7", "8", "9":
			floor, err := strconv.Atoi(e.ID)
			if err != nil {
				panic(err)
			}
			game.ToFloor(floor)
		case "G", "g":
			game.ToFloor(0)
		case "q":
			return
		}
	}
}
