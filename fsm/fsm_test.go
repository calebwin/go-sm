package fsm

import (
  "testing"
)

func TestBasic(t *testing.T) {
  myFSM := generate("locked")

  // define all possible transitions
  myTransitions := []Transition{
    Transition{"coin", []State{State{"locked"}, State{"un-locked"},}, State{"un-locked"}},
    Transition{"push", []State{State{"locked"}, State{"un-locked"},}, State{"locked"}},
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
  myTransitions := []Transition{
    Transition{"coin", []State{State{"locked"}, State{"un-locked"},}, State{"un-locked"}},
    Transition{"push", []State{State{"locked"}, State{"un-locked"},}, State{"locked"}},
  }

  // add transition rules to FSM
  myFSM = setTransitions(myFSM, myTransitions)

  var money int = 0

  // add lifecycle Callbacks to FSM
  myFSM = setCallbacks(myFSM,
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

func TestHistory(t *testing.T) {
  myFSM := generate("locked", true)

  // define all possible transitions
  myTransitions := []Transition{
    Transition{"coin", []State{State{"locked"}, State{"un-locked"},}, State{"un-locked"}},
    Transition{"push", []State{State{"locked"}, State{"un-locked"},}, State{"locked"}},
  }
  myFSM = setTransitions(myFSM, myTransitions)

  myFSM = transition(myFSM, "coin")
  myFSM = transition(myFSM, "coin")

  if len(myFSM.history) != 3 {
    t.Fail()
  }

  myFSM = clearHistory(myFSM)

  if len(myFSM.history) != 1 {
    t.Fail()
  }
}
