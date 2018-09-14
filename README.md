# go-sm
A finite state machine library for the Go programming language

# Usage
```
// generate a new finite state machine with an initial state of "locked"
myFSM := fsm.generate("locked")

// define all possible transitions between states
myTransitions := map[string][]fsm.Transition{
  // define a transition called "coin" with 2 state-to-state transitions
  "coin" : []fsm.Transition {
    fsm.Transition {
      "locked",
      "un-locked",
    },
    fsm.Transition {
      "un-locked",
      "un-locked",
    },
  },
  // define a transition called "push" with 2 state-to-state transitions
  "push" : []fsm.Transition {
    fsm.Transition {
      "un-locked",
      "locked",
    },
    fsm.Transition {
      "locked",
      "locked",
    },
  },
}

// add transition rules to finite state machine
myFSM = fsm.setTransitions(myFSM, myTransitions)

var money int = 0

// add lifecycle events to FSM
myFSM = fsm.setEvents(myFSM,
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
myFSM = fsm.transition(myFSM, "coin")
myFSM = fsm.transition(myFSM, "coin")

// money will equal 2 because it got incremented each time time the "coin" transition was executed
```
