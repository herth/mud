// Copyright Peter Herth 2016

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"time"
)

func readQuoted(reader *bufio.Reader, delim rune) (result string, eof bool) {
	var b bytes.Buffer
	esc := false

	for {
		r, s, _ := reader.ReadRune()
		if s == 0 {
			eof = true
			break
		}
		if esc {
			switch r {
			case 'n':
				b.WriteRune(10)
			case 't':
				b.WriteRune(8)
			default:
				b.WriteRune(r)
			}
			esc = false
		} else {
			switch r {
			case delim:
				return b.String(), eof
			case '\\':
				esc = true
			default:
				b.WriteRune(r)
			}
		}

	}
	return b.String(), eof
}

func readExpr(reader *bufio.Reader) (tokens []string, eof bool) {
	var b bytes.Buffer
	//eof = skipWS(reader)
loop:
	for {
		r, s, _ := reader.ReadRune()
		if s == 0 {
			eof = true
			break
		}
		switch r {
		case 10:
			break loop
		case 13:
			continue
		case 32, 9:
			if b.Len() > 0 {
				tokens = append(tokens, b.String())
				b.Reset()
			}
			skipWStoEOL(reader)
			continue
		case '"', 39:
			res, _ := readQuoted(reader, r)
			tokens = append(tokens, res)
			b.Reset()
			continue
		case '{':
			res, _ := readQuoted(reader, '}')
			tokens = append(tokens, res)
			b.Reset()
			continue
		}
		b.WriteRune(r)
	}
	if b.Len() > 0 {
		tokens = append(tokens, b.String())
	}
	return
}

func handle(c net.Conn) {
	//buffer := make([]byte, 1024)
	player := Player{Name: "JohnDoe", Connection: c, Handler: readName, Level: 1}
	world.Players = append(world.Players, &player)

	r := bufio.NewReader(c)
	player.Printf("How would you like to be called?\n> ")
	for {
		tokens, eof := readExpr(r)
		//fmt.Println("read")
		// for i, tok := range tokens {
		// 	fmt.Println(i, tok, len(tok))
		// }
		if len(tokens) == 1 && tokens[0] == "quit" {
			c.Close()
			break
		}
		if len(tokens) > 0 {
			player.Handler(&player, tokens)
		}
		if eof {
			break
		} else {
			player.Printf("> ")
		}

	}

	log.Printf("socket of %s closed %v\n", player.Name, c.RemoteAddr())
}

func serve() {
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	fmt.Println("listening on port 8000")
	for {
		// Wait for a connection.
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handle(conn)
	}
}

func main() {
	initCommands()
	world.Rooms = make(map[string]*Room, 40000)
	fmt.Println("MUD starting...")
	t1 := time.Now()
	loadAreas()
	fmt.Fprintf(os.Stderr, "loading:  %v\n", time.Since(t1))
	fmt.Println(runtime.NumCPU(), "cpus")
	go startWebServer()
	serve()
}
