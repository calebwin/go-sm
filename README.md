## What it is
go-sm is a library for generating persistent finite-state machines in the Go programming language. go-sm currently supports lifecycle callbacks, state history storage, generated .dot graph visualizations.

## How to use it
A basic finite-state machine with two states and two transitions can be created as follows with a simple 3-step process.
```golang
import "github.com/calebwin/go-sm/fsm"

myFSM := fsm.Generate("locked") // 1) generate a new finite-state machine with an inital state of "locked"

myTransitions := []fsm.Transition{
  fsm.Transition{
  	"coin", 
	[]fsm.State{fsm.State{"locked"}, fsm.State{"un-locked"},}, 
	fsm.State{"un-locked"}
  },
  fsm.Transition{
  	"push", 
	[]fsm.State{fsm.State{"locked"}, fsm.State{"un-locked"},}, 
	fsm.State{"locked"}
  },
}
myFSM = fsm.SetTransitions(myFSM, myTransitions) // 2) define possible transitions within the finite-state machine

myFSM = fsm.Execute(myFSM, "coin") // 3) execute the transition named "coin"

myFSM.state // "un-locked"
```

### Callback Functions
Callback functions can be defined for the following 4 lifecycle events.
- `onBeforeTransition`
- `onAfterTransition`
- `onEnterState`
- `onLeaveState`
They can be defined with go-sm as follows.
```golang
myFSM := fsm.Generate("locked")

myTransitions := []fsm.Transition{
  fsm.Transition{"coin", []fsm.State{fsm.State{"locked"}, fsm.State{"un-locked"},}, fsm.State{"un-locked"}},
  fsm.Transition{"push", []fsm.State{fsm.State{"locked"}, fsm.State{"un-locked"},}, fsm.State{"locked"}},
}
myFSM = fsm.SetTransitions(myFSM, myTransitions) 

coins := 0
myFSM = fsm.SetCallbacks(myFSM,
  func(transition string) {}, // onBeforeTransition
  func(transition string) { // onAfterTransition
    if transition == "coin" {
      money += 1
    }
  },
  func(state string) {}, // onEnterState
  func(state string) {}, // onLeaveState
)

myFSM = fsm.Execute(myFSM, "coin") 
```

### State History
State history can be maintained with go-sm as follows.
```golang
myFSM := fsm.Generate("locked", true) // first flag set to true to indicate history should be maintained

myTransitions := []fsm.Transition{
  fsm.Transition{"coin", []fsm.State{fsm.State{"locked"}, fsm.State{"un-locked"},}, fsm.State{"un-locked"}},
  fsm.Transition{"push", []fsm.State{fsm.State{"locked"}, fsm.State{"un-locked"},}, fsm.State{"locked"}},
}
myFSM = fsm.SetTransitions(myFSM, myTransitions) 

myFSM = fsm.Execute(myFSM, "coin") 

myFSM.history // {"locked", "un-locked",}
```
State history can be used to undo/redo state transitions.
```golang
myFSM := fsm.Generate("locked", true) // first flag set to true to indicate history should be maintained

myTransitions := []fsm.Transition{
  fsm.Transition{"coin", []fsm.State{fsm.State{"locked"}, fsm.State{"un-locked"},}, fsm.State{"un-locked"}},
  fsm.Transition{"push", []fsm.State{fsm.State{"locked"}, fsm.State{"un-locked"},}, fsm.State{"locked"}},
}
myFSM = fsm.SetTransitions(myFSM, myTransitions) 

myFSM = fsm.Execute(myFSM, "coin")
myFSM = fsm.Execute(myFSM, "coin")

myFSM = fsm.HistoryBack(myFSM, 2) // undo 2 state transitions

myFSM.state // "locked"
```
State history can also be cleared.
```golang
myFSM := fsm.Generate("locked", true) // first flag set to true to indicate history should be maintained

myTransitions := []fsm.Transition{
  fsm.Transition{"coin", []fsm.State{fsm.State{"locked"}, fsm.State{"un-locked"},}, fsm.State{"un-locked"}},
  fsm.Transition{"push", []fsm.State{fsm.State{"locked"}, fsm.State{"un-locked"},}, fsm.State{"locked"}},
}
myFSM = fsm.SetTransitions(myFSM, myTransitions) 

myFSM = fsm.Execute(myFSM, "coin")
myFSM = fsm.Execute(myFSM, "coin")

myFSM = fsm.ClearHistory(myFSM)

myFSM.history // {"un-locked",}
```

### Visualizations
Visualizations of finite-state machines can be generated as .dot files as follows.
```golang
myFSM := fsm.Generate("locked", true) // first flag set to true to indicate history should be maintained

myTransitions := []fsm.Transition{
  fsm.Transition{"coin", []fsm.State{fsm.State{"locked"}, fsm.State{"un-locked"},}, fsm.State{"un-locked"}},
  fsm.Transition{"push", []fsm.State{fsm.State{"locked"}, fsm.State{"un-locked"},}, fsm.State{"locked"}},
}
myFSM = fsm.SetTransitions(myFSM, myTransitions) 

fsm.GenerateVisualization(myFSM, "myVisualization.dot")
```
The above call to `GenerateVisualization` will result in the following file created.
```
digraph {
	locked;
	un-locked;

	"locked" -> "un-locked" [ label=" coin " ];
	"un-locked" -> "un-locked" [ label=" coin " ];
	"locked" -> "locked" [ label=" push " ];
	"un-locked" -> "locked" [ label=" push " ];
}
```
