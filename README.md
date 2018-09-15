# go-sm
A finite state machine library for the Go programming language

# Usage
```
// generate a new FSM with an initial state of "locked"
myFSM := generate("locked")

// define all possible transitions
myTransitions := []Transition{
  Transition{"coin", []State{State{"locked"}, State{"un-locked"},}, State{"un-locked"}},
  Transition{"push", []State{State{"locked"}, State{"un-locked"},}, State{"locked"}},
}

// add transition rules to FSM
myFSM = setTransitions(myFSM, myTransitions)

var money int = 0

// add lifecycle events to FSM
myFSM = setEvents(myFSM,
  func(transition string) {}, // onBeforeTransition
  func(transition string) { // onAfterTransition
    if transition == "coin" {
      money += 1
    }
  },
  func(state string) {}, // onEnterState
  func(state string) {}, // onLeaveState
)

// execute the "coin" transition twice
myFSM = transition(myFSM, "coin")
myFSM = transition(myFSM, "coin")

// money == 2 after "coin" transition executed twice

```
