package elecrash

import (
	"sync"
	"time"

	"github.com/gizak/termui/v3"
	"github.com/lovromazgon/elecrash/ui"
)

type ElevatorState string

const (
	Idle    ElevatorState = "idle"
	Opening ElevatorState = "opening"
	Closing ElevatorState = "closing"
	Moving  ElevatorState = "moving"
)

//   /<------------------\
// Idle -> Opening -> Closing
//    \-> Moving |

type Elevator struct {
	lane        int
	targetFloor int
	state       ElevatorState

	uiElevator *ui.Elevator

	sync.Mutex
}

func NewElevator(lane, maxFloor int) *Elevator {
	uiElevator := ui.NewElevator(lane, maxFloor)
	if lane == 0 {
		uiElevator.Select()
	}
	return &Elevator{
		lane:        lane,
		targetFloor: 0,
		state:       Idle,

		uiElevator: uiElevator,
	}
}

func (m *Elevator) Open() {
	m.Lock()
	defer m.Unlock()
	if m.state != Idle {
		logger.Info("can't open if not idle", "state", m.state)
		return
	}
	logger.Info("elevator opening")
	m.state = Opening
	m.uiElevator.ShowOpening()
	go func() {
		<-time.Tick(time.Second)
		// TODO load people
		m.Close()
	}()
}

func (m *Elevator) Close() {
	m.Lock()
	defer m.Unlock()
	if m.state != Opening {
		logger.Info("can't close if not opening", "state", m.state)
		return
	}
	logger.Info("elevator closing")
	m.state = Closing
	m.uiElevator.ShowClosing()
	go func() {
		<-time.Tick(time.Second)
		m.Lock()
		m.state = Idle
		m.uiElevator.ShowIdle()
		m.Unlock()
	}()
}

func (m *Elevator) Move(floor int) {
	m.Lock()
	defer m.Unlock()
	if m.state != Idle {
		logger.Info("can't move if not idle", "state", m.state)
		return
	}
	if m.uiElevator.Floor() == floor {
		logger.Info("nothing to move", "floor", floor)
		return
	}

	logger.Info("elevator moving")
	m.state = Moving
	m.targetFloor = floor
	m.uiElevator.ShowMoving(m.targetFloor - m.uiElevator.Floor())
	go m.tickMove()
}

func (m *Elevator) tickMove() {
	move := 1 // move 1 floor up
	if m.uiElevator.Floor() > m.targetFloor {
		move = -1
	}

	for m.uiElevator.Floor() != m.targetFloor {
		<-time.Tick(time.Second)
		m.Lock()
		m.uiElevator.SetFloor(m.uiElevator.Floor() + move)
		m.Unlock()
	}

	logger.Info("elevator idle")
	m.Lock()
	m.state = Idle
	m.Unlock()
	m.Open()
}

func (m *Elevator) Render() {
	termui.Render(m.uiElevator)
}

func (m *Elevator) Select() {
	m.uiElevator.Select()
}

func (m *Elevator) Deselect() {
	m.uiElevator.Deselect()
}
