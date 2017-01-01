package main

import (
	"fmt"
	"strings"
)

type InputHandler func(p *Player, tokens []string)

func readName(p *Player, tokens []string) {
	name := tokens[0]
	if len(name) > 2 {
		p.Name = name
		p.Printf("You want to be known as '%s'?\n", name)
		p.Handler = confirmName
	} else {
		p.Printf("Illegal name. How do you want to be called?\n")
	}

}

func confirmName(p *Player, tokens []string) {
	tok := tokens[0]
	if tok == "y" || tok == "Y" {
		p.Printf("Ok, you should be known as '%s'.\n", p.Name)
		p.Handler = commandHandler
		p.Valid = true
		p.Room = world.Rooms[":3:3014"]
		showRoom(p, p.Room)
	} else if tok == "n" || tok == "N" {
		p.Printf("What should it be then?\n")
		p.Handler = readName
	} else {
		p.Printf("Please answer with Y or N.\n")
	}
}

func commandHandler(p *Player, tokens []string) {
	cmd := tokens[0]
	for _, command := range Commands {
		if len(cmd) <= len(command.Name) &&
			cmd == command.Name[:len(cmd)] &&
			p.Level >= command.MinLevel {
			command.Fun(p, tokens)
			return
		}
	}

	if len(cmd) > 0 {
		p.Printf("Sorry, I don't know how to '%v'\n", strings.Join(tokens, " "))
	}
}

type Command struct {
	Name     string
	MinLevel int
	Fun      func(*Player, []string)
	Help     string
}

var Commands []Command

func initCommands() {
	Commands = []Command{
		{"help", 0, cmdHelp, ""},
		{"north", 1, goDir("north"), ""},
		{"east", 1, goDir("east"), ""},
		{"south", 1, goDir("south"), ""},
		{"west", 1, goDir("west"), ""},
		{"up", 1, goDir("up"), ""},
		{"down", 1, goDir("down"), ""},
		{"goto", 1, cmdGoto, ""},
		{"who", 1, cmdWho, ""},
		{"look", 1, cmdLook, ""}}
}

func cmdWho(p *Player, tokens []string) {
	p.Printf("\nPlayers online are:\n")
	for i, player := range world.Players {
		if player.Valid {
			p.Printf("%d\t%s\n", i+1, player.Name)
		}
	}
	p.Printf("\n")
}

func showRoom(p *Player, r *Room) {
	p.Printf("\n%s\n%s\n\n%s\n\nExits: ", r.GetID(), r.Name, r.Description)
	nexit := 0
	for i, exit := range r.Exit {
		if len(exit.To) > 0 {
			if nexit > 0 {
				p.Printf(", ")
			}
			p.Printf(DirName[i])
			nexit++
		}
	}
	if nexit == 0 {
		p.Printf("none")
	}
	p.Printf(".\n\n")

}

func cmdLook(p *Player, tokens []string) {
	if p.Room != nil {
		showRoom(p, p.Room)
	} else {
		p.Printf("\nIt is very dark.\n\n")
	}
}

func cmdHelp(p *Player, tokens []string) {
	p.Printf("\nAvailable commands are:\n")
	for _, command := range Commands {
		p.Printf("%s\t%s\n", command.Name, command.Help)
	}
	p.Printf("\n")
}

func cmdGoto(p *Player, tokens []string) {
	if len(tokens) != 2 {
		p.Printf("Go to where?\n\n")
	} else {
		r := world.Rooms[tokens[1]]
		if r != nil {
			p.Room = r
			showRoom(p, r)
			return
		} else {
			p.Printf("No such room!\n\n")
		}
	}
}

type CmdFun func(*Player, []string)

func goDir(dirName string) CmdFun {
	return func(p *Player, tokens []string) {
		for d, exit := range p.Room.Exit {
			if dirName == DirName[d] {
				to := exit.To
				if len(to) > 0 {
					if to[0] != ':' {
						to = fmt.Sprintf(":%d:%s", p.Room.Area.Nr, to)
					}
					r := world.Rooms[to]
					if r != nil {
						p.Room = r
						showRoom(p, r)
						return
					}
				}
			}
		}
		p.Printf("You cannot go %s.\n\n", dirName)
	}
}
