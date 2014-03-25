package main

type State struct {

}

func new_state() State {
	var this State
	return this
}

func (this State) plan() bool {
	return true
}

func (this State) shutdown() {
}
