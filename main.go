package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func main() {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	balancer := TeamBalancer{}

	balancer.Init(3, 1)

	for i := 0; i < 6; i++ {
		player := Player{"Player" + strconv.Itoa(i + 1), uint8(r1.Intn(6) + 1)}
		fmt.Println(player)
		added, finished := balancer.AddPlayer(player)

		fmt.Printf("Added %t, finished %t\n", added, finished)
	}

	fmt.Println(balancer.Roster, *balancer.Roster.TeamA.Players, *balancer.Roster.TeamB.Players)

}
