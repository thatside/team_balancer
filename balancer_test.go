package main

import (
	"strconv"
	"testing"
)



func TestTeamBalancer_AddPlayer_(t *testing.T) {
	balancer := TeamBalancer{}
	balancer.Init(3, 1)

	added, finished := balancer.AddPlayer(Player{"Player" + strconv.Itoa(i + 1), uint8(r1.Intn(6) + 1)})
}