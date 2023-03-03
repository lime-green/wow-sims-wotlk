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
	session.sim.Advance(time.Duration(duration) * time.Millisecond)
	session.sim.RunPendingActions(session.sim.CurrentTime)
	return Response{Success: true}
}
