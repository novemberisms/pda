package pda

import "github.com/novemberisms/stack"

// State describes a valid pushdown automata state.
//
// It requires four lifetime callbacks that each receive a reference
// to the pushdown automata that it belongs to.
type State interface {
	onEnter()
	onExit()
	onPause()
	onResume()
}

type stateStack struct {
	*stack.Stack
}

func newStateStack() *stateStack {
	return &stateStack{stack.NewStack(0)}
}

func (s *stateStack) Push(state State) {
	s.Stack.Push(state)
}

func (s stateStack) Peek() State {
	val := s.Stack.Peek()
	if val == nil {
		return nil
	}
	return val.(State)
}

func (s *stateStack) Pop() State {
	val := s.Stack.Pop()
	if val == nil {
		return nil
	}
	return val.(State)
}
