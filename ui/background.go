package ui

import (
	"image"

	ui "github.com/gizak/termui/v3"
)

const (
	upDownWidth   = 3
	floorWidth    = 3
	elevatorWidth = 5
)

type Background struct {
	ui.Block
	elevators, floors int
}

var _ ui.Drawable = (*Background)(nil)

func NewBackground(elevators, floors int) *Background {
	if elevators < 2 || floors < 4 {
		panic("we need at least 2 elevators and 4 floors")
	}

	width := (upDownWidth+2)*2 +
		(floorWidth+1)*2 +
		((elevatorWidth + 1) * elevators) - 1
	height := 4 + 2*floors + 2

	block := ui.NewBlock()
	block.SetRect(0, 0, width, height)

	return &Background{
		Block:     *block,
		elevators: elevators,
		floors:    floors,
	}
}

// Draw implements the Drawable interface.
func (b *Background) Draw(buf *ui.Buffer) {
	b.Block.Draw(buf)

	verticalCell := ui.Cell{ui.VERTICAL_LINE, b.BorderStyle}
	horizontalCell := ui.Cell{ui.HORIZONTAL_LINE, b.BorderStyle}
	verticalRightCell := ui.Cell{ui.VERTICAL_RIGHT, b.BorderStyle}
	verticalLeftCell := ui.Cell{ui.VERTICAL_LEFT, b.BorderStyle}
	horizontalDownCell := ui.Cell{ui.HORIZONTAL_DOWN, b.BorderStyle}
	horizontalUpCell := ui.Cell{ui.HORIZONTAL_UP, b.BorderStyle}
	crossCell := ui.Cell{'┼', b.BorderStyle}

	rect := b.Block.GetRect()

	// header
	buf.SetCell(verticalRightCell, image.Pt(rect.Min.X, rect.Min.Y+3))
	buf.Fill(horizontalCell, image.Rect(rect.Min.X+1, rect.Min.Y+3, rect.Max.X-1, rect.Min.Y+4))
	buf.SetCell(verticalLeftCell, image.Pt(rect.Max.X-1, rect.Min.Y+3))

	// up/down
	for i, c := range []rune("UP↑") {
		buf.SetCell(ui.NewCell(c), image.Pt(rect.Min.X+1+i, rect.Min.Y+1))
	}
	for i, c := range []rune("DN↓") {
		buf.SetCell(ui.NewCell(c), image.Pt(rect.Max.X-4+i, rect.Min.Y+1))
	}

	// floors
	buf.SetCell(ui.NewCell('F'), image.Pt(rect.Min.X+3+upDownWidth, rect.Min.Y+1))
	buf.SetCell(ui.NewCell('F'), image.Pt(rect.Max.X-4-upDownWidth, rect.Min.Y+1))
	for i := 0; i < b.floors; i++ {
		floor := floorToRune(b.floors - 1 - i)
		// left side
		buf.SetCell(verticalRightCell, image.Pt(rect.Min.X, rect.Min.Y+5+(2*i)))
		buf.SetCell(ui.NewCell(floor), image.Pt(rect.Min.X+3+upDownWidth, rect.Min.Y+5+(2*i)))

		// right side
		buf.SetCell(verticalLeftCell, image.Pt(rect.Max.X-1, rect.Min.Y+5+(2*i)))
		buf.SetCell(ui.NewCell(floor), image.Pt(rect.Max.X-4-upDownWidth, rect.Min.Y+5+(2*i)))
	}

	// vertical lines
	lineOffsets := []int{floorWidth + 1, rect.Max.X - floorWidth - 2}
	for i := 0; i <= b.elevators; i++ {
		lineOffsets = append(lineOffsets, rect.Min.X+upDownWidth+floorWidth+2+(6*i))
	}

	for _, offset := range lineOffsets {
		buf.SetCell(horizontalDownCell, image.Pt(offset, rect.Min.Y))
		buf.Fill(verticalCell, image.Rect(offset, rect.Min.Y+1, offset+1, rect.Max.Y-1))
		buf.SetCell(horizontalUpCell, image.Pt(offset, rect.Max.Y-1))
		for j := 0; j <= b.floors; j++ {
			buf.SetCell(crossCell, image.Pt(offset, rect.Min.Y+3+(2*j)))
		}
	}
}

/*
   ┌───┬─────┬─────┬─────┬─────┬─────┬───┐
   │ F │     │     │     │     │     │ F │
   │   │     │     │     │     │     │   │
   ├───┼─────┼─────┼─────┼─────┼─────┼───┤
   │   │┌───┐│     │     │     │     │   │
   ├ 9 ┼│   │┼     ┼     ┼     ┼     ┼ 9 ┤
   │   │└───┘│     │     │     │     │   │
   ├ 8 ┼     ┼     ┼     ┼     ┼     ┼ 8 ┤
   │   │     │     │     │     │     │   │
   ├ 7 ┼     ┼     ┼     ┼     ┼     ┼ 7 ┤
   │   │     │     │     │     │     │   │
   ├ 6 ┼     ┼     ┼     ┼     ┼     ┼ 6 ┤
   │   │     │     │     │     │     │   │
   ├ 5 ┼     ┼     ┼     ┼     ┼     ┼ 5 ┤
   │   │     │     │     │     │     │   │
   ├ 4 ┼     ┼     ┼     ┼     ┼     ┼ 4 ┤
   │   │     │     │     │     │     │   │
   ├ 3 ┼     ┼     ┼     ┼     ┼     ┼ 3 ┤
   │   │     │     │     │     │     │   │
   ├ 2 ┼     ┼     ┼     ┼     ┼     ┼ 2 ┤
   │   │     │     │     │     │     │   │
   ├ 1 ┼     ┼     ┼     ┼     ┼     ┼ 1 ┤
   │   │     │     │     │     │     │   │
   ├ G ┼     ┼     ┼     ┼     ┼     ┼ G ┤
   │   │     │     │     │     │     │   │
   └───┴─────┴─────┴─────┴─────┴─────┴───┘
*/
