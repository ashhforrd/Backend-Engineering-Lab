package lifecycle

import "sync/atomic"

type State struct {
	ready atomic.Bool
}

func NewState() *State {
	state := &State{}
	state.ready.Store(true)

	return state
}

func (s *State) IsReady() bool {
	return s.ready.Load()
}

func (s *State) BeginShutdown() {
	s.ready.Store(false)
}
