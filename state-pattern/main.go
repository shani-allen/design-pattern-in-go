package main

import "fmt"

//WHAT
//WHY
//HOW

//WHAT
// comes under behavioural design
//this is when we have to change the behaviour of object based on its internal state.
//E.g. mobile notification system; ring, vibrate, silent
// state pattern have 3 main participants:

//WHY
//to make our code loosely coupled,reliable, maintainable and reduce complexity i.e using this pattern
//we make the separate classes /objects for each state and behaviour and rely on the context to provide
//the behaviour implementation for the same.

type tvState interface {
	state()
}

// concrete implementation
type on struct {
}

// implementing the behavior of ON state
func (o *on) state() {
	fmt.Println("tv is on!")
}

type off struct {
}

// implementing the behavior of off state
func (o *off) state() {
	fmt.Println("tv is off!")
}

//context

type stateContext struct {
	currentTvState tvState
}

func getContext() *stateContext {
	return &stateContext{
		currentTvState: &off{},
	}
}

func (sc *stateContext) setState(state tvState) {
	sc.currentTvState = state
}

func (sc *stateContext) getState() {
	sc.currentTvState.state()
}

func main() {
	tvContext := getContext() //default is off
	tvContext.getState()      //get the state
	//changing the the state
	tvContext.setState(&on{})
	tvContext.getState()
}
