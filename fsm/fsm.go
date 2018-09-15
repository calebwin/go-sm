package fsm

import (
  //"fmt"
)

type onBeforeTransitionCallback func(transition string)
type onAfterTransitionCallback func(transition string)
type onEnterStateCallback func(state string)
type onLeaveStateCallback func(state string)

// type Transition struct {
//   from string
//   to string
// }

type State struct {
  name string
}

type Transition struct {
  name string
  from []State
  to State
}

type FSM struct {
  state string
  transitions map[string]Transition
  onBeforeTransition onBeforeTransitionCallback
  onAfterTransition onAfterTransitionCallback
  onEnterState onEnterStateCallback
  onLeaveState onLeaveStateCallback
  history []string
  historyPos int
}

func (fsm *FSM) is(state string) bool {
  return fsm.state == state
}

func (fsm *FSM) can(state string) bool {
  for _, transition := range fsm.transitions {
    if transition.to.name == state {
      for _, from := range transition.from {
        if from.name == fsm.state {
          return true
        }
      }
    }
  }
  return false
}

func (fsm *FSM) cannot(state string) bool {
  for _, transition := range fsm.transitions {
    if transition.to.name == state {
      for _, from := range transition.from {
        if from.name == fsm.state {
          return false
        }
      }
    }
  }
  return true
}

func (fsm *FSM) validTransitions() []string {
  validTransitions := []string{}

  for transitionName, transition := range fsm.transitions {
    for _, from := range transition.from {
      if from.name == fsm.state {
        validTransitions = append(validTransitions, transitionName)
      }
    }
  }

  return validTransitions
}

func (fsm *FSM) allTransitions() []string {
  allTransitions := []string{}

  for transitionName, _ := range fsm.transitions {
    allTransitions = append(allTransitions, transitionName)
  }

  return allTransitions
}

func (fsm *FSM) allStates() []string {
  allStates := []string{}
  statesAdded := map[string]bool{}

  for _, transition := range fsm.transitions {
    for _, from := range transition.from {
      if statesAdded[from.name] != true {
        allStates = append(allStates, from.name)
      }
    }
    if statesAdded[transition.to.name] != true {
      allStates = append(allStates, transition.to.name)
    }
  }

  return allStates
}

func generate(initialState string, flags ...bool) FSM {
  if len(flags) > 0 && flags[0] == true { // enable history
    return FSM {
      initialState,
      make(map[string]Transition),
      func (transition string) {},
      func (transition string) {},
      func (state string) {},
      func (state string) {},
      []string {initialState,},
      0,
    }
  }

  return FSM {
    initialState,
    make(map[string]Transition),
    func (transition string) {},
    func (transition string) {},
    func (state string) {},
    func (state string) {},
    []string {},
    0,
  }
}

func setTransitions(fsm FSM, newTransitions []Transition) FSM {
  var newTransitionMap map[string]Transition = make(map[string]Transition)

  for _, transition := range newTransitions {
    newTransitionMap[transition.name] = transition
  }

  newHistory := fsm.history
  if len(fsm.history) > 0 {
    newHistory = newHistory[0 : fsm.historyPos + 1]
  }

  return FSM {
    fsm.state,
    newTransitionMap,
    fsm.onBeforeTransition,
    fsm.onAfterTransition,
    fsm.onEnterState,
    fsm.onLeaveState,
    newHistory,
    fsm.historyPos,
  }
}

func setCallbacks(fsm FSM, newOnBeforeTransition onBeforeTransitionCallback, newOnAfterTransition onAfterTransitionCallback, newOnEnterState onEnterStateCallback, newOnLeaveState onLeaveStateCallback) FSM {
  newHistory := fsm.history
  if len(fsm.history) > 0 {
    newHistory = newHistory[0 : fsm.historyPos + 1]
  }

  return FSM {
    fsm.state,
    fsm.transitions,
    newOnBeforeTransition,
    newOnAfterTransition,
    newOnEnterState,
    newOnLeaveState,
    newHistory,
    fsm.historyPos,
  }
}

func transition(fsm FSM, transition string) FSM {
  validTransitionNames := fsm.validTransitions()

  if len(validTransitionNames) > 0 {
    fsm.onLeaveState(fsm.state)
    fsm.onBeforeTransition(transition)
    fsm.onAfterTransition(transition)
    fsm.onEnterState(fsm.transitions[validTransitionNames[0]].to.name)

    newHistory := fsm.history
    if len(fsm.history) > 0 {
      newHistory = append(newHistory, fsm.transitions[validTransitionNames[0]].to.name)
    }

    if len(fsm.history) > 0 {
      newHistory = newHistory[0 : fsm.historyPos + 1 + 1]
    }

    return FSM {
      fsm.transitions[validTransitionNames[0]].to.name,
      fsm.transitions,
      fsm.onBeforeTransition,
      fsm.onAfterTransition,
      fsm.onEnterState,
      fsm.onLeaveState,
      newHistory,
      fsm.historyPos + 1,
    }
  }

  return fsm
}



func clearHistory(fsm FSM) FSM  {
  return FSM {
    fsm.state,
    fsm.transitions,
    fsm.onBeforeTransition,
    fsm.onAfterTransition,
    fsm.onEnterState,
    fsm.onLeaveState,
    []string {fsm.state},
    0,
  }
}

func historyBack(fsm FSM, numSteps int) FSM  {
  if fsm.historyPos - numSteps >= 0 {
    return FSM {
      fsm.history[fsm.historyPos - numSteps],
      fsm.transitions,
      fsm.onBeforeTransition,
      fsm.onAfterTransition,
      fsm.onEnterState,
      fsm.onLeaveState,
      fsm.history,
      fsm.historyPos - numSteps,
    }
  }

  return fsm
}

func historyForward(fsm FSM, numSteps int) FSM  {
  if fsm.historyPos + numSteps <= len(fsm.history) - 1 {
    return FSM {
      fsm.history[fsm.historyPos + numSteps],
      fsm.transitions,
      fsm.onBeforeTransition,
      fsm.onAfterTransition,
      fsm.onEnterState,
      fsm.onLeaveState,
      fsm.history,
      fsm.historyPos + numSteps,
    }
  }

  return fsm
}

func limitHistory(fsm FSM, limit int) FSM  {
  return FSM {
    fsm.state,
    fsm.transitions,
    fsm.onBeforeTransition,
    fsm.onAfterTransition,
    fsm.onEnterState,
    fsm.onLeaveState,
    fsm.history[len(fsm.history) - limit : len(fsm.history)],
    fsm.historyPos,
  }
}
