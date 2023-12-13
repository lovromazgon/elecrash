package ui

import (
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

	action []rune
}

var _ ui.Drawable = (*Elevator)(nil)

func NewElevator(lane int, maxFloor int) *Elevator {
	b := ui.NewBlock()

	// init position on ground floor
	b.SetRect(
		5+(lane*6),
		4,
		5+(lane*6)+elevatorWidth,
		4+(maxFloor-1)*2+3,
	)
	b.Border = false

	elevatorBlock := ui.NewBlock()
	elevatorBlock.SetRect(
		5+(lane*6),
		4+(maxFloor-1)*2,
		5+(lane*6)+elevatorWidth,
		4+(maxFloor-1)*2+3,
	)

	return &Elevator{
		Block:         *b,
		lane:          lane,
		floor:         0,
		maxFloor:      maxFloor,
		elevatorBlock: *elevatorBlock,
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
		e.Block.Draw(buf)
	}
	e.elevatorBlock.Draw(buf)

	rect := e.elevatorBlock.GetRect()

	// draw action
	if len(e.action) > 1 {
		buf.SetCell(ui.NewCell(e.action[0]), image.Pt(rect.Min.X+1, rect.Min.Y+1))
		buf.SetCell(ui.NewCell(e.action[1]), image.Pt(rect.Max.X-2, rect.Max.Y-2))
	}

	// TODO draw people

}

func (e *Elevator) ShowOpening() {
	e.Lock()
	defer e.Unlock()
	e.action = []rune{'<', '>'}
}

func (e *Elevator) ShowClosing() {
	e.Lock()
	defer e.Unlock()
	e.action = []rune{'>', '<'}
}

func (e *Elevator) ShowIdle() {
	e.Lock()
	defer e.Unlock()
	e.action = nil
}

func (e *Elevator) ShowMoving(move int) {
	e.Lock()
	defer e.Unlock()
	if move > 0 {
		e.action = []rune{'↑', '↑'}
	} else {
		e.action = []rune{'↓', '↓'}
	}
}
