package elecrash

import ui "github.com/gizak/termui/v3"

type Elevator struct {
	ui.Block
	elevatorBlock ui.Block

	lane     int
	floor    int
	maxFloor int
	selected bool
	dirty    bool
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

	selected := lane == 0

	elevatorBlock := ui.NewBlock()
	elevatorBlock.SetRect(
		5+(lane*6),
		4+(maxFloor-1)*2,
		5+(lane*6)+elevatorWidth,
		4+(maxFloor-1)*2+3,
	)
	if selected {
		elevatorBlock.BorderStyle.Fg = ui.ColorGreen
	}

	return &Elevator{
		Block:         *b,
		lane:          lane,
		floor:         0,
		maxFloor:      maxFloor,
		selected:      selected,
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
	e.selected = true
	e.elevatorBlock.BorderStyle.Fg = ui.ColorGreen
}

func (e *Elevator) Deselect() {
	e.selected = false
	e.elevatorBlock.BorderStyle.Fg = ui.Theme.Block.Border.Fg
}

func (e *Elevator) Draw(buf *ui.Buffer) {
	if e.dirty {
		// redraw lane
		e.Block.Draw(buf)
	}
	e.elevatorBlock.Draw(buf)
	// TODO draw people
}
