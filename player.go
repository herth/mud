package main

import (
	"fmt"
	"net"
)

// Player object
type Player struct {
	Name       string
	Connection net.Conn
	Handler    InputHandler
	Valid      bool
	Level      int
	Room       *Room
}

func (p *Player) Printf(f string, args ...interface{}) {
	fmt.Fprintf(p.Connection, f, args...)
}
