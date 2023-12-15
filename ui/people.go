package ui

import (
	"image"

	ui "github.com/gizak/termui/v3"
)

type PeopleDirection string

const (
	maxPeople = 6
)

// People renders a group of people on one floor that want to ride either up or
// down.
type People struct {
	ui.Block

	count []int
}

var _ ui.Drawable = (*People)(nil)

func NewPeople(maxFloor int, xOffset int) *People {
	b := ui.NewBlock()
	b.SetRect(xOffset, 5, xOffset+1, 4+(2*maxFloor))

	return &People{
		Block: *b,
		count: make([]int, maxFloor),
	}
}

func (p *People) Count(floor int) int {
	p.Lock()
	defer p.Unlock()
	return p.count[floor]
}

func (p *People) Add(floor int, i int) bool {
	p.Lock()
	defer p.Unlock()
	count := p.count[floor]
	if count+i > maxPeople || count+i < 0 {
		return false
	}
	p.count[floor] = count + i
	return true
}

func (p *People) Draw(buf *ui.Buffer) {
	rect := p.GetRect()

	// draw people
	for i, c := range p.count {
		buf.SetCell(ui.NewCell(Braille(c)), image.Pt(rect.Min.X, rect.Min.Y+((len(p.count)-1)*2)-(i*2)))
	}
}
