package main

type State struct {
	mb *Messagebox
}

func new_state() State {
	var this State
	this.mb = new_messagebox()
	this.mb.listen(QuitMsg)
	return this
}

func (this State) plan() bool {
	for _, msg := range (*this.mb.get()) {
		switch msg.token {
		case QuitMsg:
			return false
		}
	}
	return true
}

func (this State) shutdown() {
}
