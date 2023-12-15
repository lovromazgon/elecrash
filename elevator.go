package elecrash

import (
	"slices"
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
	people      []*Person

	uiElevator *ui.Elevator

	game *Elecrash

	sync.Mutex
}

func NewElevator(game *Elecrash, lane, maxFloor int) *Elevator {
	uiElevator := ui.NewElevator(lane, maxFloor)
	if lane == 0 {
		uiElevator.Select()
	}
	return &Elevator{
		lane:        lane,
		targetFloor: 0,
		state:       Idle,
		people:      make([]*Person, 0, maxPeople),

		uiElevator: uiElevator,

		game: game,
	}
}

func (e *Elevator) Open() {
	e.Lock()
	defer e.Unlock()
	if e.state != Idle {
		logger.Info("can't open if not idle", "state", e.state)
		return
	}
	logger.Info("elevator opening")
	e.state = Opening
	e.uiElevator.ShowOpening()
	go func() {
		<-time.Tick(time.Second / 2)
		e.game.UnloadPeople(e)
		e.game.LoadPeople(e)
		<-time.Tick(time.Second / 2)
		e.Close()
	}()
}

func (e *Elevator) Close() {
	e.Lock()
	defer e.Unlock()
	if e.state != Opening {
		logger.Info("can't close if not opening", "state", e.state)
		return
	}
	logger.Info("elevator closing")
	e.state = Closing
	e.uiElevator.ShowClosing()
	go func() {
		<-time.Tick(time.Second)
		e.Lock()
		e.state = Idle
		e.uiElevator.ShowIdle()
		e.Unlock()
	}()
}

func (e *Elevator) Move(floor int) {
	e.Lock()
	defer e.Unlock()
	if e.state != Idle {
		logger.Info("can't move if not idle", "state", e.state)
		return
	}

	e.state = Moving
	e.targetFloor = floor

	if e.uiElevator.Floor() != floor {
		logger.Info("elevator moving")
		e.uiElevator.ShowMoving(e.targetFloor)
	}

	go e.tickMove()
}

func (e *Elevator) tickMove() {
	move := 1 // move 1 floor up
	if e.uiElevator.Floor() > e.targetFloor {
		move = -1
	}

	for e.uiElevator.Floor() != e.targetFloor {
		<-time.Tick(time.Second)
		e.Lock()
		e.uiElevator.SetFloor(e.uiElevator.Floor() + move)
		e.Unlock()
	}

	logger.Info("elevator idle")
	e.Lock()
	e.state = Idle
	e.Unlock()
	e.Open()
}

func (e *Elevator) Render() {
	termui.Render(e.uiElevator)
}

func (e *Elevator) Select() {
	e.uiElevator.Select()
}

func (e *Elevator) Deselect() {
	e.uiElevator.Deselect()
}

func (e *Elevator) UnloadPeople() []*Person {
	e.Lock()
	defer e.Unlock()
	var unloaded []int
	var people []*Person
	for i, p := range e.people {
		if p.targetFloor == e.targetFloor {
			unloaded = append(unloaded, i)
			people = append(people, p)
		}
	}
	slices.Reverse(unloaded)
	for _, i := range unloaded {
		e.people[i] = e.people[len(e.people)-1]
		e.people = e.people[:len(e.people)-1]
	}
	e.uiElevator.SetPeople(len(e.people))
	return people
}

func (e *Elevator) LoadPeople(people []*Person) {
	e.Lock()
	defer e.Unlock()

	if len(e.people)+len(people) > maxPeople {
		panic("tried to load more people than allowed into elevator")
	}
	e.people = append(e.people, people...)
	e.uiElevator.SetPeople(len(e.people))
}
