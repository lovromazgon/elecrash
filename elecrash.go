package elecrash

import (
	"sync"
	"time"

	ui "github.com/gizak/termui/v3"
)

type Elecrash struct {
	bg        *Background
	elevators []*Elevator
	movers    []*Mover
	selected  int
	sync.Mutex
}

func NewElecrash(elevators, floors int) *Elecrash {
	bg := NewBackground(elevators, floors)
	e := make([]*Elevator, elevators)
	for i := range e {
		e[i] = NewElevator(i, floors)
	}

	return &Elecrash{
		bg:        bg,
		elevators: e,
		movers:    make([]*Mover, elevators),
	}
}

func (e *Elecrash) Run() {
	ui.Render(e.bg)
	e.Render()
	for range time.NewTicker(time.Millisecond * 100).C {
		e.Render()
	}
}

func (e *Elecrash) Render() {
	e.Lock()
	defer e.Unlock()
	for _, e := range e.elevators {
		ui.Render(e)
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
	if e.movers[e.selected] != nil {
		return // can't override mover
	}
	if e.elevators[e.selected].Floor() == floor {
		return // nothing to move
	}

	m := NewMover(e.elevators[e.selected], floor, &e.Mutex, func(_ bool) {
		e.Lock()
		defer e.Unlock()
		e.movers[e.selected] = nil
	})
	m.Start()
	e.movers[e.selected] = m
}
