package fsm

import (
  "testing"
)

func TestBasic(t *testing.T) {
  myFSM := generate("locked")

  myTransitions := map[string][]Transition{
    // define a transition called "coin" with 2 state-to-state transitions
    "coin" : []Transition {
      Transition {
        "locked",
        "un-locked",
      },
      Transition {
        "un-locked",
        "un-locked",
      },
    },
    // define a transition called "push" with 2 state-to-state transitions
    "push" : []Transition {
      Transition {
        "un-locked",
        "locked",
      },
      Transition {
        "locked",
        "locked",
      },
    },
  }
  myFSM = setTransitions(myFSM, myTransitions)

  myFSM = transition(myFSM, "coin")
  myFSM = transition(myFSM, "coin")

  if myFSM.is("locked") == true {
    t.Fail()
  }
}


func TestLifecycle(t *testing.T) {
  // generate a new FSM with an initial state of "locked"
  myFSM := generate("locked")

  // define all possible transitions
  myTransitions := map[string][]Transition{
    // define a transition called "coin" with 2 state-to-state transitions
    "coin" : []Transition {
      Transition {
        "locked",
        "un-locked",
      },
      Transition {
        "un-locked",
        "un-locked",
      },
    },
    // define a transition called "push" with 2 state-to-state transitions
    "push" : []Transition {
      Transition {
        "un-locked",
        "locked",
      },
      Transition {
        "locked",
        "locked",
      },
    },
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

  if myFSM.is("locked") == true {
    t.Fail()
  }

  if money != 2 {
    t.Fail()
  }

}
