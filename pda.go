package pda

import "errors"

// PushdownAutomata is a pushdown state machine that manipulates an internal stack
// of states.
type PushdownAutomata struct {
	states *stateStack
}

// NewPushdownAutomata creates a new PushdownAutomata with an empty state stack
func NewPushdownAutomata() *PushdownAutomata {
	return &PushdownAutomata{
		states: newStateStack(),
	}
}

// PushState pushes a new state unto the stack of this automata, calling onPause on the
// previous top state (if it exists) and then onEnter on the new state
func (p *PushdownAutomata) PushState(state State) {
	if previous := p.states.Peek(); previous != nil {
		previous.onPause()
	}
	p.states.Push(state)
	state.onEnter()
}

func (p *PushdownAutomata) PopState() (State, error) {
	toPop := p.states.Pop()
	if toPop == nil {
		return nil, errors.New("Trying to pop from an empty PushdownAutomata")
	}
	toPop.onExit()
	if uncovered := p.states.Peek(); uncovered != nil {
		uncovered.onResume()
	}
	return toPop, nil
}

func (p PushdownAutomata) Current() (State, error) {
	current := p.states.Peek()
	if current == nil {
		return nil, errors.New("No current state")
	}
	return current, nil
}
