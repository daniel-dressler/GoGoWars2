package main

type Subsystem interface {
	update()
	shutdown()
}
