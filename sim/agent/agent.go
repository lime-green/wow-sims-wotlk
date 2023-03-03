package main

import (
	"time"
)

type Agent interface {
	Cast(spell string, session *Session) Response
	Wait(duration int, session *Session) Response
	GetState(session *Session) Response
	Init(session *Session)
}

type BaseAgent struct{}

func (base *BaseAgent) Wait(duration int, session *Session) Response {
	newTime := session.sim.CurrentTime + time.Duration(duration)*time.Millisecond
	session.sim.RunPendingActions(newTime)
	session.sim.Advance(newTime - session.sim.CurrentTime)

	if session.sim.CurrentTime >= session.sim.Duration {
		session.sim.Stop()
	}

	return Response{Success: true}
}
