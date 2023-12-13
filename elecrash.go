package elecrash

import (
	"sync"
	"time"

	"github.com/gizak/termui/v3"
	"github.com/lovromazgon/elecrash/ui"
)

type Elecrash struct {
	bg        *ui.Background
	elevators []*Elevator
	selected  int
	sync.Mutex
}

func NewElecrash(elevators, floors int) *Elecrash {
	bg := ui.NewBackground(elevators, floors)
	e := make([]*Elevator, elevators)
	for i := range e {
		e[i] = NewElevator(i, floors)
	}

	return &Elecrash{
		bg:        bg,
		elevators: e,
	}
}

func (e *Elecrash) Run() {
	termui.Render(e.bg)
	e.Render()
	for range time.NewTicker(time.Millisecond * 100).C {
		e.Render()
	}
}

func (e *Elecrash) Render() {
	e.Lock()
	defer e.Unlock()
	for _, e := range e.elevators {
		e.Render()
	}
}

func (e *Elecrash) Left() {
	e.Lock()
	defer e.Unlock()
	if e.selected == 0 {
		return // can't go more left
	}
	e.elevators[e.selected].Deselect()
	e.selected -= 1
	e.elevators[e.selected].Select()
}

func (e *Elecrash) Right() {
	e.Lock()
	defer e.Unlock()
	if e.selected == len(e.elevators)-1 {
		return // can't go more right
	}
	e.elevators[e.selected].Deselect()
	e.selected += 1
	e.elevators[e.selected].Select()
}

func (e *Elecrash) ToFloor(floor int) {
	e.Lock()
	defer e.Unlock()

	e.elevators[e.selected].Move(floor)
}
