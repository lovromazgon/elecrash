package main

import (
	"flag"
	"log"

	ui "github.com/gizak/termui/v3"
	"github.com/lovromazgon/elecrash"
)

func main() {
	var (
		elevators = flag.Int("e", 4, "number of elevators")
		floors    = flag.Int("f", 8, "number of floors")
	)
	flag.Parse()

	if *elevators < 2 || *elevators > 10 {
		log.Fatal("need between 2 and 10 elevators")
	}
	if *floors < 4 || *floors > 10 {
		log.Fatal("need between 4 and 10 floors")
	}

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	game := elecrash.NewElecrash(*elevators, *floors)
	go game.Run()

	for e := range ui.PollEvents() {
		if e.Type != ui.KeyboardEvent {
			continue
		}
		switch e.ID {
		case "<Left>":
			game.Left()
		case "<Right>":
			game.Right()
		case "q":
			return
		}
	}
}
