package fsm

import (
  "testing"
)

func TestBasic(t *testing.T) {
  myFSM := Generate("locked")

  // define all possible transitions
  myTransitions := []Transition{
    Transition{"coin", []State{State{"locked"}, State{"un-locked"},}, State{"un-locked"}},
    Transition{"push", []State{State{"locked"}, State{"un-locked"},}, State{"locked"}},
  }
  myFSM = SetTransitions(myFSM, myTransitions)

  myFSM = Execute(myFSM, "coin")
  myFSM = Execute(myFSM, "coin")

  if myFSM.Is("locked") == true {
    t.Fail()
  }
}

func TestLifecycle(t *testing.T) {
  // generate a new FSM with an initial state of "locked"
  myFSM := Generate("locked")

  // define all possible transitions
  myTransitions := []Transition{
    Transition{"coin", []State{State{"locked"}, State{"un-locked"},}, State{"un-locked"}},
    Transition{"push", []State{State{"locked"}, State{"un-locked"},}, State{"locked"}},
  }

  // add transition rules to FSM
  myFSM = SetTransitions(myFSM, myTransitions)

  var money int = 0

  // add lifecycle Callbacks to FSM
  myFSM = SetCallbacks(myFSM,
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
  myFSM = Execute(myFSM, "coin")
  myFSM = Execute(myFSM, "coin")

  if myFSM.Is("locked") == true {
    t.Fail()
  }

  if money != 2 {
    t.Fail()
  }

}

func TestHistory(t *testing.T) {
  myFSM := Generate("locked", true)

  // define all possible transitions
  myTransitions := []Transition{
    Transition{"coin", []State{State{"locked"}, State{"un-locked"},}, State{"un-locked"}},
    Transition{"push", []State{State{"locked"}, State{"un-locked"},}, State{"locked"}},
  }
  myFSM = SetTransitions(myFSM, myTransitions)

  myFSM = Execute(myFSM, "coin")
  myFSM = Execute(myFSM, "coin")

  if len(myFSM.history) != 3 {
    t.Fail()
  }

  myFSM = ClearHistory(myFSM)

  if len(myFSM.history) != 1 {
    t.Fail()
  }

  myFSM = Execute(myFSM, "push")
  myFSM = Execute(myFSM, "push")
  myFSM = HistoryBack(myFSM, 2)

  if myFSM.state != "un-locked" {
    t.Fail()
  }
}

func TestVisualization(t *testing.T) {
  myFSM := Generate("locked")

  // define all possible transitions
  myTransitions := []Transition{
    Transition{"coin", []State{State{"locked"}, State{"un-locked"},}, State{"un-locked"}},
    Transition{"push", []State{State{"locked"}, State{"un-locked"},}, State{"locked"}},
  }
  myFSM = SetTransitions(myFSM, myTransitions)

  myFSM = Execute(myFSM, "coin")
  myFSM = Execute(myFSM, "coin")

  GenerateVisualization(myFSM, "testVisualization.txt")
}
