package main

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
	p.Printf("Please answer with Y or N.\n")
}
