package elecrash

import "time"

const (
	maxPeople = 6
)

type PersonState string

const (
	Waiting PersonState = "waiting"
	Riding  PersonState = "riding"
	Arrived PersonState = "arrived"
)

type Person struct {
	currentFloor int
	targetFloor  int
	elevator     *Elevator
	waitingTime  time.Time
	ridingTime   time.Time
	arrivedTime  time.Time

	state PersonState
}

func NewPerson(currentFloor, targetFloor int) *Person {
	return &Person{
		currentFloor: currentFloor,
		targetFloor:  targetFloor,
		waitingTime:  time.Now(),
		state:        Waiting,
	}
}
