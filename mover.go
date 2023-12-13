package elecrash

import (
	"sync"
	"time"
)

type Mover struct {
	e     *Elevator
	floor int

	lock   *sync.Mutex
	onDone func(bool)

	stop chan struct{}
}

func NewMover(e *Elevator, floor int, lock *sync.Mutex, onDone func(bool)) *Mover {
	return &Mover{
		e:      e,
		floor:  floor,
		lock:   lock,
		onDone: onDone,
		stop:   make(chan struct{}),
	}
}

func (m *Mover) Start() {
	move := 1 // move 1 floor up
	if m.e.Floor() > m.floor {
		move = -1
	}

	go func() {
		reachedFloor := false
		defer m.onDone(reachedFloor)
		for m.e.Floor() != m.floor {
			select {
			case <-m.stop:
				return
			case <-time.Tick(time.Second):
				m.lock.Lock()
				m.e.SetFloor(m.e.Floor() + move)
				m.lock.Unlock()
			}
		}
		reachedFloor = true
	}()
}

func (m *Mover) Stop() {
	close(m.stop)
}
