package main

func main() {
	display := new_display()
	state := new_state()
	input := new_input()

	for state.plan() {
		input.update()
		display.update()
	}

	display.shutdown()
}
