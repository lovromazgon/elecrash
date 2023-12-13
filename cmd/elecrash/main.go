package main

import (
	"log"

	ui "github.com/gizak/termui/v3"
	"github.com/lovromazgon/elecrash"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	b := elecrash.NewBackground(5, 4)
	ui.Render(b)

	for e := range ui.PollEvents() {
		if e.Type == ui.KeyboardEvent {
			break
		}
	}
}
