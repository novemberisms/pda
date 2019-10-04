package pda

import (
	"testing"
)

type testState struct {
	name             string
	onEnterCallback  func()
	onExitCallback   func()
	onResumeCallback func()
	onPauseCallback  func()
}

func (ts testState) onEnter() {
	ts.onEnterCallback()
}

func (ts testState) onExit() {
	ts.onExitCallback()
}

func (ts testState) onResume() {
	ts.onResumeCallback()
}

func (ts testState) onPause() {
	ts.onPauseCallback()
}

func assert(t *testing.T, condition bool, message string) {
	t.Helper()
	if !condition {
		t.Error(message)
	}
}

func compare(t *testing.T, result []string, expected []string, message string) {
	t.Helper()
	if len(result) != len(expected) {
		t.Error(message)
		return
	}
	for i := 0; i < len(result); i++ {
		if result[i] != expected[i] {
			t.Error(message)
		}
	}
}

func TestCreation(t *testing.T) {
	p := NewPushdownAutomata()
	current, _ := p.Current()
	assert(t, current == nil, "Must not have an initial state")
}

func TestPush(t *testing.T) {
	p := NewPushdownAutomata()
	onEnterCalled, onPauseCalled := false, false
	initialState := testState{
		name: "initial",
		onEnterCallback: func() {
			onEnterCalled = true
		},
		onPauseCallback: func() {
			onPauseCalled = true
		},
	}

	p.PushState(initialState)
	curr, _ := p.Current()
	assert(t, curr.(testState).name == initialState.name, "p.Current() not accurate")
	assert(t, onEnterCalled == true, "onEnter not called")
	assert(t, onPauseCalled == false, "onPause called too early")

	nextState := testState{
		name:            "next",
		onEnterCallback: func() {},
	}

	p.PushState(nextState)
	curr, _ = p.Current()

	assert(t, curr.(testState).name == nextState.name, "p.Current() not accurate")
	assert(t, onPauseCalled == true, "onPause not called")
}

func TestPop(t *testing.T) {
	p := NewPushdownAutomata()

	exitStack := make([]string, 0)
	pauseStack := make([]string, 0)

	addStack := func(stack *[]string, message string) {
		*stack = append(*stack, message)
	}

	doNothing := func() {}

	stateA := testState{
		name:             "A",
		onEnterCallback:  doNothing,
		onResumeCallback: doNothing,
		onExitCallback: func() {
			addStack(&exitStack, "A")
		},
		onPauseCallback: func() {
			addStack(&pauseStack, "A")
		},
	}

	stateB := testState{
		name:             "B",
		onEnterCallback:  doNothing,
		onResumeCallback: doNothing,
		onExitCallback: func() {
			addStack(&exitStack, "B")
		},
		onPauseCallback: func() {
			addStack(&pauseStack, "B")
		},
	}

	stateC := testState{
		name:             "C",
		onEnterCallback:  doNothing,
		onResumeCallback: doNothing,
		onExitCallback: func() {
			addStack(&exitStack, "C")
		},
		onPauseCallback: func() {
			addStack(&pauseStack, "C")
		},
	}

	compare(t, exitStack, []string{}, "exit stack should have nothing")
	compare(t, pauseStack, []string{}, "pause stack should have nothing")

	p.PushState(stateA)
	p.PushState(stateB)

	compare(t, exitStack, []string{}, "exit stack should have nothing")
	compare(t, pauseStack, []string{"A"}, "pause stack should have A")

	p.PushState(stateC)

	compare(t, exitStack, []string{}, "exit stack should have nothing")
	compare(t, pauseStack, []string{"A", "B"}, "pause stack should have A and B")

	popped, err := p.PopState()

	assert(t, err == nil, "error should not have been thrown")
	assert(t, popped.(testState).name == stateC.name, "popped should be state C")
	compare(t, exitStack, []string{"C"}, "exit stack should have C")
	compare(t, pauseStack, []string{"A", "B"}, "pause stack should have A and B")

	popped, err = p.PopState()

	assert(t, err == nil, "error should not have been thrown")
	assert(t, popped.(testState).name == stateB.name, "popped should be state B")
	compare(t, exitStack, []string{"C", "B"}, "exit stack should have C and B")
	compare(t, pauseStack, []string{"A", "B"}, "pause stack should have A and B")

	popped, err = p.PopState()

	assert(t, err == nil, "error should not have been thrown")
	assert(t, popped.(testState).name == stateA.name, "popped should be state A")
	compare(t, exitStack, []string{"C", "B", "A"}, "exit stack should have C, B, and A")
	compare(t, pauseStack, []string{"A", "B"}, "pause stack should have A and B")

	popped, err = p.PopState()
	assert(t, err != nil, "error should have been thrown when popping empty pda")
	assert(t, popped == nil, "no state should have been returned")
}
