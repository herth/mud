package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func readWord(reader *bufio.Reader) (word string, eof bool) {
	var b bytes.Buffer
	eof = skipWS(reader)
	for {
		rune, size, _ := reader.ReadRune()
		if size == 0 {
			eof = true
			break
		}
		if rune == 32 {
			break
		}
		if rune == 13 {
			break
		}
		if rune == 10 {
			break
		}
		b.WriteRune(rune)
	}
	return b.String(), eof
}

func readLong(reader *bufio.Reader) (word string, eof bool) {
	var b bytes.Buffer
	eof = skipWS(reader)
	for {
		rune, size, _ := reader.ReadRune()
		if size == 0 {
			eof = true
			break
		}
		if rune == '~' {
			break
		}
		b.WriteRune(rune)
	}
	return b.String(), eof
}

func readEOL(reader *bufio.Reader) (word string, eof bool) {
	var b bytes.Buffer
	eof = false
	for {
		rune, size, _ := reader.ReadRune()
		if size == 0 {
			eof = true
			break
		}
		if rune == '\r' {
			continue
		}
		if rune == '\n' {
			break
		}
		b.WriteRune(rune)
	}
	return b.String(), eof
}

func skipTo(reader *bufio.Reader, stopRune rune) (eof bool) {
	eof = false
	for {
		rune, size, _ := reader.ReadRune()
		if size == 0 {
			eof = true
			break
		}
		if rune == stopRune {
			reader.UnreadRune()
			break
		}
	}
	return eof
}

func skipWS(reader *bufio.Reader) (eof bool) {
	eof = false
	for {
		rune, size, _ := reader.ReadRune()
		if size == 0 {
			eof = true
			break
		}
		if rune == 32 {
			continue
		}
		if rune == 13 {
			continue
		}
		if rune == 10 {
			continue
		}
		if rune == 9 {
			continue
		}
		reader.UnreadRune()
		break
	}
	return eof
}

// skipp all whitespace, but stop at line end
func skipWStoEOL(reader *bufio.Reader) (eof bool) {
	eof = false
	for {
		rune, size, _ := reader.ReadRune()
		if size == 0 {
			eof = true
			break
		}
		if rune == 32 {
			continue
		}
		if rune == 13 {
			continue
		}
		if rune == 9 {
			continue
		}
		reader.UnreadRune()
		break
	}
	return eof
}

// World The whole world
type World struct {
	Players [](*Player)
	Rooms   map[string]*Room
}

var world World

// Mob  the mobile
type Mob struct {
	ID          string
	Name        string
	ShortDesc   string
	LongDesc    string
	Description string
	Race        string
}

// Item the item
type Item struct {
	ID        string
	Name      string
	ShortDesc string
	LongDesc  string
}

const (
	NORTH = 0
	EAST  = 1
	SOUTH = 2
	WEST  = 3
	UP    = 4
	DOWN  = 5
)

var DirName []string = []string{"north", "east", "south", "west", "up", "down"}

// Exit room exit
type Exit struct {
	Description string
	Keyword     string
	State       string
	Key         string
	To          string
}

// ExtraDescription x
type ExtraDescription struct {
	Keyword     string
	Description string
}

// Room room
type Room struct {
	ID          string
	Name        string
	Description string
	Exit        [6]Exit
	Extra       []ExtraDescription
	Area        *Area
	X, Y, Z     int
}

func (r *Room) GetID() string {
	id := r.ID
	if id[0] != ':' {
		id = fmt.Sprintf(":%d:%s", r.Area.Nr, r.ID)
	}
	return id
}

// Area the whole area
type Area struct {
	Name     string
	FileName string
	Nr       int
	Flags    string
	Mobs     [](*Mob)
	Items    [](*Item)
	Rooms    [](*Room)
}

var areas [](*Area)

func (mob *Mob) read(r *bufio.Reader) {
	mob.Name, _ = readLong(r)
	mob.ShortDesc, _ = readLong(r)
	mob.LongDesc, _ = readLong(r)
	mob.Description, _ = readLong(r)
	mob.Race, _ = readLong(r)
	skipTo(r, '#')
	return
}

func (area *Area) readMobs(r *bufio.Reader) {
	for {
		vnum, eof := readWord(r)
		if vnum == "#0" {
			return
		}

		if eof {
			fmt.Println("Eof while reading mobiles")
			return
		}
		if len(vnum) < 1 {
			fmt.Println("empty vnum")
			return
		}

		//fmt.Printf("vnum:%v<<<\n", vnum)
		mob := Mob{ID: vnum[1:]}
		mob.read(r)
		//fmt.Println("read mob:", mob.Vnum, mob.Name)
		area.Mobs = append(area.Mobs, &mob)
		//fmt.Println("mobs:", area.Mobs)
	}
}

func (item *Item) read(r *bufio.Reader) {
	item.Name, _ = readLong(r)
	item.ShortDesc, _ = readLong(r)
	item.LongDesc, _ = readLong(r)
	skipTo(r, '#')
	return
}

func (area *Area) readItems(r *bufio.Reader) {
	for {
		vnum, eof := readWord(r)
		if vnum == "#0" {
			return
		}
		if eof {
			fmt.Println("Eof while reading items in area", area.FileName)
			//panic("items eof")
			return
		}
		if len(vnum) < 1 {
			fmt.Println("empty vnum")
			return
		}

		//fmt.Printf("vnum:%v<<<\n", vnum)
		item := Item{ID: vnum[1:]}
		item.read(r)
		//fmt.Println("read item:", item.Vnum, item.Name)
		area.Items = append(area.Items, &item)
		//fmt.Println("items:", area.Items)
	}
}

func (room *Room) read(r *bufio.Reader) {
	room.Name, _ = readLong(r)
	room.Description, _ = readLong(r)
	readWord(r)
	readWord(r)
	readWord(r)
	for {
		entry, _ := readWord(r)
		//fmt.Println("entry", entry)
		e := entry[0]
		if e == 'D' { // exit
			dir := entry[1] - '0'
			//fmt.Println("dir=", dir)
			desc, _ := readLong(r)
			//fmt.Println("desc=", desc)
			keyword, _ := readLong(r)
			//fmt.Println("keyword=", keyword)
			state, _ := readWord(r)
			//fmt.Println("state=", state)
			key, _ := readWord(r)
			//fmt.Println("key=", key)
			toRoom, _ := readWord(r)
			//fmt.Println("toRooom=", toRoom)
			ex := Exit{Description: desc, Keyword: keyword, State: state, Key: key, To: toRoom}
			//fmt.Println("exx=", ex)

			//fmt.Println()

			room.Exit[dir] = ex
			//fmt.Println("exits", room.Exit)
			continue
		}
		if entry == "S" {
			skipTo(r, '#')
			return
		}
		if entry == "E" {
			k, _ := readLong(r)
			d, _ := readLong(r)
			room.Extra = append(room.Extra, ExtraDescription{Keyword: k, Description: d})
			continue
		}
		if entry == "H" {
			readWord(r)
			continue
		}
		if entry == "M" {
			readWord(r)
			continue
		}
		if entry == "K" {
			readWord(r)
			x, _ := readWord(r)
			room.X, _ = strconv.Atoi(x)
			y, _ := readWord(r)
			room.Y, _ = strconv.Atoi(y)
			z, _ := readWord(r)
			room.Z, _ = strconv.Atoi(z)
			continue
		}

		break
		// case "H":
		// 	break
		// case "M":
		// 	break
		// case "K":
		// 	break
		// case "amap":
		// 	break
		// case "WLH":
		// 	break
		// case "X":
		// 	break
		// case "watches":
		// 	break
		// case "clan":
		// 	break
		// case "O":
		// 	break
		// case "A":
		// 	break
		// case "damage":
		// 	break
		// case "encounters":
		// 	break
		// case "encountr":
		// 	break
		// case "populate":
		// 	break
		// case "religion":
		// 	break

	}
	skipTo(r, '#')
	return
}

func (area *Area) readRooms(r *bufio.Reader) {
	for {
		vnum, eof := readWord(r)
		if vnum == "#0" {
			return
		}

		if eof {
			fmt.Println("Eof while reading rooms")
			return
		}
		if len(vnum) < 1 {
			fmt.Println("empty vnum")
			return
		}

		// //fmt.Printf("vnum:%v<<<\n", vnum)
		room := Room{ID: vnum[1:], Area: area}
		room.read(r)
		//fmt.Println("read room:", room.Vnum, room.Name, room)
		area.Rooms = append(area.Rooms, &room)
	}
}

func (area *Area) readResets(r *bufio.Reader) {
	for {
		code, eof := readWord(r)
		if code == "S" {
			return
		}
		if eof {
			fmt.Println("Eof while reading resets")
			return
		}
		_, _ = readEOL(r)
	}
}

func (area *Area) load() {
	file, err := os.Open("./area/" + area.FileName)
	if err != nil {
		return
	}
	defer file.Close()

	//	fmt.Println("reading", area.FileName)

	r := bufio.NewReader(file)
	for {
		word, eof := readWord(r)
		//fmt.Println("word", word)
		switch word {
		case "#AREA":
			_, _ = readLong(r)
			area.Name, _ = readLong(r)
			_, _ = readLong(r)
			continue
		case "#MOBILES":
			area.readMobs(r)
		case "#OBJECTS":
			area.readItems(r)
		case "#ROOMS":
			area.readRooms(r)
		case "#RESETS":
			area.readResets(r)
		default:
			//fmt.Println("skipping word", word)
		}

		if eof {
			return
		}
	}

}

func loadAreas() {
	file, err := os.Open("./area/area.list")
	n := 0
	if err != nil {
		return
	}
	var c chan (int)
	parallel := true
	if parallel {
		c = make(chan int)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		l := scanner.Text()
		//fmt.Println("area:", l)
		if l[0] != '#' {
			words := bufio.NewScanner(strings.NewReader(l))
			words.Split(bufio.ScanWords)
			words.Scan()
			nr, _ := strconv.Atoi(words.Text())
			words.Scan()
			flags := words.Text()
			words.Scan()
			fileName := words.Text()
			//parts := strings.Split(l, " ")
			//fmt.Println("parts:", len(parts), parts)
			area := Area{Nr: nr, FileName: fileName, Flags: flags}

			areas = append(areas, &area)
			//			area.load()
			//fmt.Println("loaded", area)
			if parallel {
				go func() {
					area.load()
					c <- 1 // mark area loaded
				}()
			} else {
				area.load()
			}
			n++
		}
	}
	if parallel {
		for i := 0; i < n; i++ {
			<-c // wait till all areas sent load signal
		}
	}
	nroom := 0
	for _, a := range areas {
		num := a.Nr
		for _, r := range a.Rooms {
			id := fmt.Sprintf(":%d:%s", num, r.ID)
			world.Rooms[id] = r
			nroom++
		}
	}
	fmt.Println(len(areas), "areas loaded", nroom, "rooms.")
	//fmt.Println(areas)
}
