package elecrash

import (
	"math/rand"
	"slices"
	"sync"
	"time"

	"github.com/gizak/termui/v3"
	"github.com/lovromazgon/elecrash/ui"
)

const renderTick = time.Millisecond * 100

type Elecrash struct {
	bg        *ui.Background
	elevators []*Elevator
	floors    int
	selected  int

	waitingPeople []*Person
	peopleUp      *ui.People
	peopleDown    *ui.People

	r                *rand.Rand
	spawnProbability float64

	totalLoaded         int
	totalRides          int
	averageWaitDuration time.Duration
	averageRideDuration time.Duration

	score *ui.Score

	sync.Mutex
}

func NewElecrash(elevators, floors int, spawnRatePerSecond float64) *Elecrash {
	bg := ui.NewBackground(elevators, floors)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	elecrash := &Elecrash{
		bg:               bg,
		floors:           floors,
		r:                r,
		waitingPeople:    make([]*Person, 0),
		peopleUp:         ui.NewPeople(floors, 2),
		peopleDown:       ui.NewPeople(floors, bg.GetRect().Max.X-3),
		spawnProbability: spawnRatePerSecond * renderTick.Seconds(),
		score:            ui.NewScore(elevators, floors),
	}

	e := make([]*Elevator, elevators)
	for i := range e {
		e[i] = NewElevator(elecrash, i, floors)
	}

	elecrash.elevators = e
	return elecrash
}

func (e *Elecrash) Run() {
	termui.Render(e.bg)
	termui.Render(e.score)
	e.Render()
	for range time.NewTicker(renderTick).C {
		if f := e.r.Float64(); f < e.spawnProbability {
			currentFloor := int((f / e.spawnProbability) * float64(e.floors))
			targetFloor := currentFloor
			for targetFloor == currentFloor {
				// ensure we generate a target floor that is different than current floor
				targetFloor = int(e.r.Float64() * float64(e.floors))
			}
			spawned := e.SpawnPerson(currentFloor, targetFloor)
			if !spawned {
				// game over
				ui.RenderGameOver(e.bg.GetRect())
				return
			}
		}
		e.Render()
	}
}

func (e *Elecrash) SpawnPerson(currentFloor, targetFloor int) bool {
	e.Lock()
	defer e.Unlock()

	spawned := false
	if targetFloor > currentFloor {
		spawned = e.peopleUp.Add(currentFloor, 1)
	} else {
		spawned = e.peopleDown.Add(currentFloor, 1)
	}
	if !spawned {
		logger.Info("failed to spawn person", "currentFloor", currentFloor, "targetFloor", targetFloor)
		return spawned
	}

	e.waitingPeople = append(e.waitingPeople, NewPerson(currentFloor, targetFloor))
	logger.Info("spawned person", "currentFloor", currentFloor, "targetFloor", targetFloor)
	return true
}

func (e *Elecrash) Render() {
	e.Lock()
	defer e.Unlock()
	for _, e := range e.elevators {
		e.Render()
	}
	termui.Render(e.peopleUp)
	termui.Render(e.peopleDown)
}

func (e *Elecrash) Left() {
	e.Lock()
	defer e.Unlock()
	if e.selected == 0 {
		return // can't go more left
	}
	e.elevators[e.selected].Deselect()
	e.selected -= 1
	e.elevators[e.selected].Select()
}

func (e *Elecrash) Right() {
	e.Lock()
	defer e.Unlock()
	if e.selected == len(e.elevators)-1 {
		return // can't go more right
	}
	e.elevators[e.selected].Deselect()
	e.selected += 1
	e.elevators[e.selected].Select()
}

func (e *Elecrash) ToFloor(floor int) {
	e.Lock()
	defer e.Unlock()

	e.elevators[e.selected].Move(floor)
}

func (e *Elecrash) LoadPeople(elv *Elevator) {
	e.Lock()
	defer e.Unlock()

	limit := maxPeople - len(elv.people)
	if limit == 0 {
		return // elevator is full
	}
	people := make([]*Person, 0, limit)
	loaded := make([]int, 0, limit)
	goingUp := 0
	goingDown := 0
	for i, p := range e.waitingPeople {
		if p.currentFloor == elv.targetFloor {
			people = append(people, p)
			loaded = append(loaded, i)
			if p.targetFloor > p.currentFloor {
				goingUp += 1
			} else {
				goingDown += 1
			}
			if len(people) == limit {
				break // can't load more people
			}
		}
	}

	if len(people) == 0 {
		return
	}

	elv.LoadPeople(people)
	loadTime := time.Now()

	slices.Reverse(loaded)
	var waitDurationTotal time.Duration
	for _, i := range loaded {
		e.waitingPeople[i].ridingTime = loadTime
		waitDurationTotal += e.waitingPeople[i].ridingTime.Sub(e.waitingPeople[i].waitingTime)

		e.waitingPeople[i] = e.waitingPeople[len(e.waitingPeople)-1]
		e.waitingPeople = e.waitingPeople[:len(e.waitingPeople)-1]
	}

	e.peopleUp.Add(elv.targetFloor, -goingUp)
	e.peopleDown.Add(elv.targetFloor, -goingDown)

	e.averageWaitDuration = time.Duration(float64((e.averageWaitDuration*time.Duration(e.totalLoaded))+waitDurationTotal) / float64(e.totalLoaded+len(loaded)))
	e.totalLoaded += len(loaded)
	e.updateScore()

	logger.Info("loaded people into elevator", "lane", elv.lane, "goingUp", goingUp, "goingDown", goingDown)
}

func (e *Elecrash) UnloadPeople(elv *Elevator) {
	e.Lock()
	defer e.Unlock()

	unloaded := elv.UnloadPeople()
	if len(unloaded) == 0 {
		return
	}

	unloadTime := time.Now()
	var rideDurationTotal time.Duration
	for _, p := range unloaded {
		p.arrivedTime = unloadTime
		rideDurationTotal += p.arrivedTime.Sub(p.ridingTime)
	}

	e.averageRideDuration = time.Duration(float64((e.averageRideDuration*time.Duration(e.totalRides))+rideDurationTotal) / float64(e.totalRides+len(unloaded)))
	e.totalRides += len(unloaded)
	e.updateScore()

	logger.Info("unloaded people from elevator", "lane", elv.lane, "total", len(unloaded))
}

func (e *Elecrash) updateScore() {
	logger.Info("SCORE", "total rides", e.totalRides, "avg wait", e.averageRideDuration, "avg ride", e.averageRideDuration)
	e.score.SetTotalRides(e.totalRides)
	e.score.SetAverageWaitDuration(e.averageWaitDuration)
	e.score.SetAverageRideDuration(e.averageRideDuration)
	termui.Render(e.score)
}
