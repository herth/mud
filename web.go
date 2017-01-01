package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func handler(w http.ResponseWriter, r *http.Request) {
}

func topHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>MUD</h1>\n\n<a href=\"areas\">Areas</a>")
}

func areasHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1><a href=\"/areas\">MUD</a></h1>\n\n")
	for _, area := range areas {
		name := area.Name
		if len(name) < 1 {
			name = area.FileName
		}
		fmt.Fprintf(w, "<a href=\"area/?nr=%d\">%s</a><br>\n", area.Nr, name)
	}
}

func areaHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1><a href=\"/areas\">MUD</a></h1>\n\n")
	nr := r.URL.Query().Get("nr")
	n, ok := strconv.Atoi(nr)
	if ok != nil {
		fmt.Fprintf(w, "No such area.\n")
		return
	}
	var area *Area
	for _, a := range areas {
		if a.Nr == n {
			area = a
			break
		}
	}

	if area == nil {
		fmt.Fprintf(w, "No such area.\n")
		return
	}

	if area != nil {
		fmt.Fprintf(w, "<h1>%s</h1>\n", area.Name)
		fmt.Fprintf(w, "<h2>Rooms</h2>\n")
		for _, r := range area.Rooms {
			fmt.Fprintf(w, "<div id=\"room\">\n")
			fmt.Fprintf(w, "<div id=\"id\"><a name=\"%s\"></a>%s</div><br>\n", r.GetID(), r.ID)
			fmt.Fprintf(w, "<div id=\"name\">%s</div><br>\n", r.Name)
			fmt.Fprintf(w, "<div id=\"desc\">%s</div><br>\n", r.Description)

			fmt.Fprintf(w, "Exits: ")
			for i, exit := range r.Exit {
				to := exit.To
				if len(to) > 0 {
					if to[0] != ':' {
						to = fmt.Sprintf(":%d:%s", r.Area.Nr, to)
					}
					r := world.Rooms[to]
					if r != nil {
						//fmt.Fprintf(w, "<a href=\"/area/?nr=%d#%s\">%s</a> ", r.Area.Nr, to, DirName[i])
						fmt.Fprintf(w, "<a href=\"/room/?nr=%s\">%s</a> ", to, DirName[i])
					}
				}
			}
			fmt.Fprintf(w, "</div>")
		}
	}

}

func roomHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1><a href=\"/areas\">MUD</a></h1>\n\n")
	nr := r.URL.Query().Get("nr")

	room := world.Rooms[nr]

	if room == nil {
		fmt.Fprintf(w, "No such room.\n")
		return
	}

	fmt.Fprintf(w, "<div id=\"room\">\n")
	fmt.Fprintf(w, "<div id=\"area\"><a href=/area/?nr=%d>%s</a></div><br>\n", room.Area.Nr, room.Area.Name)
	fmt.Fprintf(w, "<div id=\"id\"><a name=\"%s\"></a>%s</div><br>\n", room.GetID(), room.ID)
	fmt.Fprintf(w, "<div id=\"name\">%s</div><br>\n", room.Name)
	fmt.Fprintf(w, "<div id=\"desc\">%s</div><br>\n", room.Description)

	fmt.Fprintf(w, "Exits: ")
	for i, exit := range room.Exit {
		to := exit.To
		if len(to) > 0 {
			if to[0] != ':' {
				to = fmt.Sprintf(":%d:%s", room.Area.Nr, to)
			}
			r := world.Rooms[to]
			if r != nil {
				//fmt.Fprintf(w, "<a href=\"/area/?nr=%d#%s\">%s</a> ", r.Area.Nr, to, DirName[i])
				fmt.Fprintf(w, "<a href=\"/room/?nr=%s\">%s</a> ", to, DirName[i])
			}
		}
	}
	fmt.Fprintf(w, "</div>")
}

func startWebServer() {
	http.HandleFunc("/", topHandler)
	http.HandleFunc("/areas", areasHandler)
	http.HandleFunc("/area/", areaHandler)
	http.HandleFunc("/room/", roomHandler)
	http.ListenAndServe(":8080", nil)
}
