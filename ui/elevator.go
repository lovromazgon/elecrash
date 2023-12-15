package ui

import (
	"fmt"
	"image"

	ui "github.com/gizak/termui/v3"
)

type Elevator struct {
	ui.Block
	elevatorBlock ui.Block

	lane     int
	floor    int
	maxFloor int
	dirty    bool
	people   int

	state  string
	action []rune
}

var _ ui.Drawable = (*Elevator)(nil)

func NewElevator(lane int, maxFloor int) *Elevator {
	b := ui.NewBlock()
	b.SetRect(
		upDownWidth+floorWidth+3+(lane*6),
		1,
		upDownWidth+floorWidth+3+(lane*6)+elevatorWidth,
		4+(maxFloor-1)*2+3,
	)
	b.Border = false

	// init position on ground floor
	elevatorBlock := ui.NewBlock()
	elevatorBlock.SetRect(
		upDownWidth+floorWidth+3+(lane*6),
		4+(maxFloor-1)*2,
		upDownWidth+floorWidth+3+(lane*6)+elevatorWidth,
		4+(maxFloor-1)*2+3,
	)

	return &Elevator{
		Block:         *b,
		lane:          lane,
		floor:         0,
		maxFloor:      maxFloor,
		elevatorBlock: *elevatorBlock,
		dirty:         true,
	}
}

func (e *Elevator) Floor() int {
	return e.floor
}

func (e *Elevator) SetFloor(floor int) {
	e.floor = floor
	rect := e.elevatorBlock.GetRect()
	e.elevatorBlock.SetRect(
		rect.Min.X,
		4+(e.maxFloor-1-floor)*2,
		rect.Max.X,
		4+(e.maxFloor-1-floor)*2+3,
	)
	e.dirty = true
}

func (e *Elevator) Select() {
	e.elevatorBlock.BorderStyle.Fg = ui.ColorGreen
}

func (e *Elevator) Deselect() {
	e.elevatorBlock.BorderStyle.Fg = ui.Theme.Block.Border.Fg
}

func (e *Elevator) Draw(buf *ui.Buffer) {
	if e.dirty {
		// redraw lane
		rect := e.Block.GetRect()
		e.Block.Draw(buf)
		buf.Fill(
			ui.Cell{ui.HORIZONTAL_LINE, e.Block.BorderStyle},
			image.Rect(rect.Min.X, rect.Min.Y+2, rect.Max.X+1, rect.Min.Y+3),
		)
	}
	e.elevatorBlock.Draw(buf)

	rect := e.elevatorBlock.GetRect()

	// draw action
	if len(e.action) > 1 {
		buf.SetCell(ui.NewCell(e.action[0]), image.Pt(rect.Min.X+1, rect.Min.Y+1))
		buf.SetCell(ui.NewCell(e.action[1]), image.Pt(rect.Max.X-2, rect.Min.Y+1))
	}

	// draw state
	for i, c := range []rune(e.state) {
		buf.SetCell(ui.NewCell(c), image.Pt(rect.Min.X+i, 1))
	}

	// draw people
	buf.SetCell(ui.NewCell(Braille(e.people)), image.Pt(rect.Min.X+(elevatorWidth/2), rect.Min.Y+1))
}

func (e *Elevator) ShowOpening() {
	e.Lock()
	defer e.Unlock()
	e.action = []rune{'<', '>'}
	e.state = " OPN "
}

func (e *Elevator) ShowClosing() {
	e.Lock()
	defer e.Unlock()
	e.action = []rune{'>', '<'}
	e.state = " CLS "
}

func (e *Elevator) ShowIdle() {
	e.Lock()
	defer e.Unlock()
	e.action = nil
	e.state = "     "
}

func (e *Elevator) ShowMoving(targetFloor int) {
	e.Lock()
	defer e.Unlock()
	if e.floor < targetFloor {
		e.action = []rune{'↑', '↑'}
		e.state = fmt.Sprintf(" ↑%c  ", floorToRune(targetFloor))
	} else {
		e.action = []rune{'↓', '↓'}
		e.state = fmt.Sprintf(" ↓%c  ", floorToRune(targetFloor))
	}
}

func (e *Elevator) SetPeople(people int) {
	e.Lock()
	defer e.Unlock()
	e.people = people
}
