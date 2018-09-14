package fsm

import (
  //"fmt"
)

type onBeforeTransitionEvent func(transition string)
type onAfterTransitionEvent func(transition string)
type onEnterStateEvent func(state string)
type onLeaveStateEvent func(state string)

type Transition struct {
  from string
  to string
}

type FSM struct {
  state string
  transitions map[string][]Transition
  onBeforeTransition onBeforeTransitionEvent
  onAfterTransition onAfterTransitionEvent
  onEnterState onEnterStateEvent
  onLeaveState onLeaveStateEvent
}

func (fsm *FSM) is(state string) bool {
  return fsm.state == state
}

func (fsm *FSM) can(state string) bool {
  for _, transitions := range fsm.transitions {
    for _, transition := range transitions {
      if transition.from == fsm.state && transition.to == state {
        return true
      }
    }
  }
  return false
}

func (fsm *FSM) cannot(state string) bool {
  for _, transitions := range fsm.transitions {
    for _, transition := range transitions {
      if transition.from == fsm.state && transition.to == state {
        return false
      }
    }
  }
  return true
}

func (fsm *FSM) validTransitions() map[string][]Transition {
  validTransitions := make(map[string][]Transition)

  for transitionName, transitions := range fsm.transitions {
    for _, transition := range transitions {
      if transition.from == fsm.state {
        validTransitions[transitionName] = append(validTransitions[transitionName], Transition {
          transition.from,
          transition.to,
        })
      }
    }
  }

  return validTransitions
}

func (fsm *FSM) allTransitions() map[string][]Transition {
  allTransitions := make(map[string][]Transition)

  for transitionName, transitions := range fsm.transitions {
    for _, transition := range transitions {
      allTransitions[transitionName] = append(allTransitions[transitionName], Transition {
        transition.from,
        transition.to,
      })
    }
  }

  return allTransitions
}

func (fsm *FSM) allStates() []string {
  allStates := make([]string, 0)

  for _, transitions := range fsm.transitions {
    for _, transition := range transitions {
      var fromStateAdded bool = false
      var toStateAdded bool = true

      for _, state := range allStates {
        if state == transition.from {
          fromStateAdded = true
        }
        if state == transition.to {
          toStateAdded = true
        }
      }

      if !fromStateAdded {
        allStates = append(allStates, transition.from)
      }

      if !toStateAdded {
        allStates = append(allStates, transition.to)
      }
    }
  }

  return allStates
}

func generate(initialState string) FSM {
  return FSM {
    initialState,
    make(map[string][]Transition),
    func (transition string) {},
    func (transition string) {},
    func (state string) {},
    func (state string) {},
  }
}

func setTransitions(fsm FSM, newTransitions map[string][]Transition) FSM {
  return FSM {
    fsm.state,
    newTransitions,
    fsm.onBeforeTransition,
    fsm.onAfterTransition,
    fsm.onEnterState,
    fsm.onLeaveState,
  }
}

func setEvents(fsm FSM, newOnBeforeTransition onBeforeTransitionEvent, newOnAfterTransition onAfterTransitionEvent, newOnEnterState onEnterStateEvent, newOnLeaveState onLeaveStateEvent) FSM {
    return FSM {
      fsm.state,
      fsm.transitions,
      newOnBeforeTransition,
      newOnAfterTransition,
      newOnEnterState,
      newOnLeaveState,
    }
}

func transition(fsm FSM, transition string) FSM {
  validTransitions := fsm.validTransitions()

  if len(validTransitions) > 0 {
    fsm.onLeaveState(fsm.state)
    fsm.onBeforeTransition(transition)
    fsm.onAfterTransition(transition)
    fsm.onEnterState(validTransitions[transition][0].to)

    return FSM {
      validTransitions[transition][0].to,
      fsm.transitions,
      fsm.onBeforeTransition,
      fsm.onAfterTransition,
      fsm.onEnterState,
      fsm.onLeaveState,
    }
  }

  return fsm
}
