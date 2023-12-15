package ui

import (
	"fmt"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type Score struct {
	ui.Grid

	totalRides          *widgets.Paragraph
	averageWaitDuration *widgets.Paragraph
	averageRideDuration *widgets.Paragraph
}

func NewScore(elevators, floors int) *Score {
	width := (upDownWidth+2)*2 +
		(floorWidth+1)*2 +
		((elevatorWidth + 1) * elevators) - 1
	topPadding := 4 + 2*floors + 2

	totalRides := widgets.NewParagraph()
	averageWaitDuration := widgets.NewParagraph()
	averageRideDuration := widgets.NewParagraph()

	grid := ui.NewGrid()
	grid.Set(
		ui.NewCol(1.0/3, totalRides),
		ui.NewCol(1.0/3, averageWaitDuration),
		ui.NewCol(1.0/3, averageRideDuration),
	)
	grid.SetRect(1, topPadding+1, width, topPadding+5)
	s := &Score{
		Grid:                *grid,
		totalRides:          totalRides,
		averageWaitDuration: averageWaitDuration,
		averageRideDuration: averageRideDuration,
	}

	s.SetTotalRides(0)
	s.SetAverageWaitDuration(0)
	s.SetAverageRideDuration(0)

	return s
}

func (s *Score) SetTotalRides(i int) {
	s.totalRides.Text = fmt.Sprintf("Total Rides: %d", i)
}

func (s *Score) SetAverageWaitDuration(t time.Duration) {
	s.averageWaitDuration.Text = fmt.Sprintf("Avg Wait: %v", roundDuration(t))
}

func (s *Score) SetAverageRideDuration(t time.Duration) {
	s.averageRideDuration.Text = fmt.Sprintf("Avg Ride: %v", roundDuration(t))
}

func roundDuration(d time.Duration) time.Duration {
	precision := time.Duration(10)
	switch {
	case d > time.Second:
		d = d.Round(time.Second / precision)
	case d > time.Millisecond:
		d = d.Round(time.Millisecond / precision)
	case d > time.Microsecond:
		d = d.Round(time.Microsecond / precision)
	}
	return d
}
