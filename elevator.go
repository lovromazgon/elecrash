package elecrash

import ui "github.com/gizak/termui/v3"

type Elevator struct {
	ui.Block
	lane     int
	floor    int
	maxFloor int
	selected bool
}

var _ ui.Drawable = (*Elevator)(nil)

func NewElevator(lane int, maxFloor int) *Elevator {
	b := ui.NewBlock()

	// init position on ground floor
	b.SetRect(
		5+(lane*6),
		4+(maxFloor-1)*2,
		5+(lane*6)+elevatorWidth,
		4+(maxFloor-1)*2+3,
	)
	selected := lane == 0
	if selected {
		b.BorderStyle.Fg = ui.ColorGreen
	}

	return &Elevator{
		Block:    *b,
		lane:     lane,
		floor:    0,
		maxFloor: maxFloor,
		selected: selected,
	}
}

func (e *Elevator) SetFloor(floor int) {
	e.floor = floor
	rect := e.GetRect()
	e.SetRect(
		rect.Min.X,
		4+(e.maxFloor-1-floor)*2,
		rect.Max.X,
		4+(e.maxFloor-1-floor)*2+3,
	)
}

func (e *Elevator) Select() {
	e.selected = true
	e.BorderStyle.Fg = ui.ColorGreen
}

func (e *Elevator) Deselect() {
	e.selected = false
	e.BorderStyle.Fg = ui.Theme.Block.Border.Fg
}

func (e *Elevator) Draw(buf *ui.Buffer) {
	e.Block.Draw(buf)
	// TODO draw people
}
