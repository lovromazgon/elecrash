package ui

import (
	"image"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func RenderGameOver(rect image.Rectangle) {
	b := ui.NewBlock()
	b.SetRect(rect.Min.X, rect.Min.Y, rect.Max.X, rect.Max.Y)

	p := widgets.NewParagraph()
	p.Text = "GAME OVER"
	p.TextStyle.Modifier = ui.ModifierBold
	p.Border = false

	p.SetRect(
		(rect.Min.X+rect.Max.X)/2-len(p.Text)/2-1,
		(rect.Min.Y+rect.Max.Y)/2,
		(rect.Min.X+rect.Max.X)/2+len(p.Text)/2+2,
		(rect.Min.Y+rect.Max.Y)/2+1,
	)

	ui.Render(b)
	ui.Render(p)
}
