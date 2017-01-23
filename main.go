package main

// Start with:
// default via 10.0.0.3 dev eth0  proto babel onlink
// 10.0.0.1 via 10.0.0.3 dev eth0  proto babel onlink
// 10.0.0.2 via 10.0.0.2 dev eth0  proto babel onlink
// 10.0.0.3 via 10.0.0.3 dev eth0  proto babel onlink
// 10.0.0.5 via 10.0.0.5 dev eth0  proto babel onlink
// 10.0.0.6 via 10.0.0.2 dev eth0  proto babel onlink
// 10.0.1.0/24 via 10.0.0.3 dev eth0  proto babel onlink
// 10.0.1.2 via 10.0.0.3 dev eth0  proto babel onlink

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type RouteMap struct {
	Interfaces map[string]*Interface
}

type Interface struct {
	Neighbors map[string]*Neighbor
}

type Neighbor struct {
	Routes map[string]*Route
}

type Route struct {
	Protocol    string
	Interface   string
	Neighbor    string
	Destination string
}

func main() {
	raw, err := exec.Command("ip", "route").Output()
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(raw), "\n")

	var routes []*Route

	for _, line := range lines {
		if len(line) > 0 {
			route := buildRoute(line)
			routes = append(routes, &route)
		}
	}

	routeMap := mapRoutes(routes)

	b, err := json.MarshalIndent(routeMap, "\n", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(string(b))
}

func buildRoute(line string) Route {
	fields := strings.Split(line, " ")
	var route Route
	route.Destination = fields[0]
	for i, field := range fields {
		switch field {

		case "proto":
			route.Protocol = fields[i+1]

		case "dev":
			route.Interface = fields[i+1]

		case "via":
			route.Neighbor = fields[i+1]

		}
	}

	return route
}

func mapRoutes(routeArr []*Route) *RouteMap {
	rtm := &RouteMap{
		Interfaces: make(map[string]*Interface),
	}

	for _, route := range routeArr {
		if rtm.Interfaces[route.Interface] == nil {
			rtm.Interfaces[route.Interface] = &Interface{
				Neighbors: make(map[string]*Neighbor),
			}
		}

		if rtm.Interfaces[route.Interface].Neighbors[route.Neighbor] == nil {
			rtm.Interfaces[route.Interface].Neighbors[route.Neighbor] = &Neighbor{
				Routes: make(map[string]*Route),
			}
		}

		rtm.Interfaces[route.Interface].Neighbors[route.Neighbor].Routes[route.Destination] = route
	}
	return rtm
}
