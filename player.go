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

func (p *Player) GetName() string {
	return p.Name
}

func (p *Player) GetShortDesc() string {
	return p.Name
}

func (p *Player) GetLongDesc() string {
	return p.Name
}

func (p *Player) GetRoom() *Room {
	return p.Room
}

func (p *Player) SetRoom(r *Room) {
	p.Room = r
}

// assert Player implements Char
var _ Char = &Player{}
