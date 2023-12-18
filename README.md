# Elecrash

A CLI game that requires you to manage a fleet of elevators.

## Quick start

Build the game using `go build -o elecrash cmd/elecrash/main.go` and run the
binary using `./elecrash`.

Quit the game by pressing `q`.

You can configure some aspects of the game:

```
  -d int
    	difficulty (1-10) (default 1)
  -e int
    	number of elevators (default 4)
  -f int
    	number of floors (default 8)
```

## How to play

Your goal is to manage the elevators and move as many people to their destination
as you can. Once there are more than 6 people waiting on a floor the game is over.

Select the elevator using the arrow keys `←` and `→` and press the floor number
to move the elevator to a floor (not that the ground floor is reached by pressing
`g`). Once the elevator reaches the selected floor it will open and close its
doors which will unload people riding the elevator and load new people that are
waiting on that floor. The elevator will only accept a new command once it is done
moving, so choose the floor wisely.

People are represented by dots spawning on random floors. If it spawns in the left
column marked `UP↑`, the person wants to go to a higher floor. Similarly, if it
spawns in the right column marked `DN↓` they want to go to a lower floor. Once the
person enters the elevator, a diamond `◆` will appear on the floor where the person
wants to get off.

https://github.com/lovromazgon/elecrash/assets/8320753/c7cba684-8ef2-4e5d-ac75-08a30feb5dfc

Good luck!
