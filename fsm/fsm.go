package fsm

import (
  "os"
  "fmt"
  "strings"
)

// types of callback functions
type onBeforeTransitionCallback func(transition string)
type onAfterTransitionCallback func(transition string)
type onEnterStateCallback func(state string)
type onLeaveStateCallback func(state string)

// definition of a State
type State struct {
  name string
}

// definition of a Transition
type Transition struct {
  name string
  from []State // origin states of transition
  to State // desitation state of transition
}

// definition of an FSM
type FSM struct {
  state string
  transitions map[string]Transition

  // callback methods
  onBeforeTransition onBeforeTransitionCallback
  onAfterTransition onAfterTransitionCallback
  onEnterState onEnterStateCallback
  onLeaveState onLeaveStateCallback

  // history
  history []string
  historyPos int
}

// checks if name of current state of FSM is the given state name
func (fsm *FSM) Is(state string) bool {
  return fsm.state == state
}

// checks if state of given name can be reached from current state of FSM
func (fsm *FSM) Can(state string) bool {
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

// checks if state of given name cannot be reached from current state of FSM
func (fsm *FSM) Cannot(state string) bool {
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

// returns a list of names of valid transitions that can be made from current state of FSM
func (fsm *FSM) ValidTransitions() []string {
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

// returns a list of names of all transitions that can be made from all states of FSM
func (fsm *FSM) AllTransitions() []string {
  allTransitions := []string{}

  for transitionName, _ := range fsm.transitions {
    allTransitions = append(allTransitions, transitionName)
  }

  return allTransitions
}

// returns a list of states in FSM
func (fsm *FSM) AllStates() []string {
  allStates := []string{}
  statesAdded := map[string]bool{}

  for _, transition := range fsm.transitions {
    for _, from := range transition.from {
      if statesAdded[from.name] != true {
        allStates = append(allStates, from.name)
        statesAdded[from.name] = true
      }
    }
    if statesAdded[transition.to.name] != true {
      allStates = append(allStates, transition.to.name)
      statesAdded[transition.to.name] = true
    }
  }

  return allStates
}

// returns a new generated FSM with given name of initial state and given flags
// first flag is a boolean that indicates whether or not state history should be tracked
func Generate(initialState string, flags ...bool) FSM {
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

// returns a new FSM with given list of transitions
func SetTransitions(fsm FSM, newTransitions []Transition) FSM {
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

// returns new FSM with given callback functions
func SetCallbacks(fsm FSM, newOnBeforeTransition onBeforeTransitionCallback, newOnAfterTransition onAfterTransitionCallback, newOnEnterState onEnterStateCallback, newOnLeaveState onLeaveStateCallback) FSM {
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

// returns new FSM with transition of given name executed
func Execute(fsm FSM, transition string) FSM {
  validTransitionNames := fsm.ValidTransitions()

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

// returns new FSM with state history cleared
func ClearHistory(fsm FSM) FSM  {
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

// returns new FSM with state stepped back given number of times
func HistoryBack(fsm FSM, numSteps int) FSM  {
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

// returns new FSM with state stepped forwards given number of times
func HistoryForward(fsm FSM, numSteps int) FSM  {
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

// returns new FSM with state history limited to the last given number of states
func LimitHistory(fsm FSM, limit int) FSM  {
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

// generates a visualization of FSM
func GenerateVisualization(fsm FSM, newFilePath string) {
  file, err := os.Create(newFilePath)

  if err != nil {
    panic(err)
  }
  defer file.Close()

  fmt.Fprintf(file, "digraph {\n")

  fmt.Fprintf(file, "\t" + strings.Join(fsm.AllStates(), ";\n\t"))
  fmt.Fprintf(file, ";\n\n")

  allTransitionNames := fsm.AllTransitions()
  for _, transitionName := range allTransitionNames {
    for _, from := range fsm.transitions[transitionName].from {
      fmt.Fprintf(file, "\t\"" + from.name + "\" -> \"" + fsm.transitions[transitionName].to.name + "\" [ label=\" " + transitionName + " \" ];\n")
    }
  }

  fmt.Fprintf(file, "}")
}
