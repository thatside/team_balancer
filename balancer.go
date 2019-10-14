package main

import (
	"math"
	"math/rand"
	"strconv"
	"time"
)

type Player struct {
	Name string
	Rank uint8
}

func CreatePlayer(rank uint8) *Player {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	return &Player{
		Name: "Player" + strconv.Itoa(r1.Intn(100)),
		Rank: rank,
	}
}

type Team struct {
	Players *[]Player
}

func (t *Team) getRankRange() (min, max uint8) {
	min, max = ^uint8(0), uint8(0)
	for _, p := range *t.Players {
		if p.Rank < min {
			min = p.Rank
		}
		if p.Rank > max {
			max = p.Rank
		}
	}

	return min, max
}

type MatchRoster struct {
	TeamA,
	TeamB,
	CurrentTeam,
	OppositeTeam *Team
	TeamSize,
	RankRange uint8
}

func (r *MatchRoster) getRankRange() (min, max uint8) {
	min, max = ^uint8(0), uint8(0)
	aMin, aMax := r.TeamA.getRankRange()
	bMin, bMax := r.TeamB.getRankRange()

	min = bMin
	if aMin < bMin {
		min = aMin
	}

	max = bMax
	if aMax > bMax {
		max = aMax
	}

	return min, max
}

func (r *MatchRoster) isFirstInsert() bool {
	min, max := r.getRankRange()

	return min == ^uint8(0) && max == 0
}

func (r *MatchRoster) canAddPlayer(rank uint8) bool {
	if r.isFirstInsert() {
		return true
	}

	min, max := r.getRankRange()

	minDoesntFitRange := math.Abs(float64(min)-float64(rank)) > float64(r.RankRange)
	maxDoesntFitRange := math.Abs(float64(max)-float64(rank)) > float64(r.RankRange)
	if minDoesntFitRange || maxDoesntFitRange {
		return false
	}

	if len(*r.CurrentTeam.Players) == int(r.TeamSize) {
		return false
	}

	return true
}

func CreateEmptyRoster(teamSize, rankRange uint8) *MatchRoster {
	var playersA, playersB []Player
	roster := &MatchRoster{
		TeamA: &Team{Players: &playersA},
		TeamB: &Team{Players: &playersB},
	}

	roster.CurrentTeam = roster.TeamA
	roster.OppositeTeam = roster.TeamB
	roster.TeamSize = teamSize
	roster.RankRange = rankRange

	return roster
}

func (t *Team) AddPlayer(player *Player) {
	*t.Players = append(*t.Players, *player)
}

type BalancesTeams interface {
	Init(teamSize, rankRange uint8)
	AddPlayer(player Player) (addedToTeam bool)
}

type TeamBalancer struct {
	insertChan   chan Player
	insertedChan chan bool
	rosterChan   chan *MatchRoster
	Roster       *MatchRoster
}

func (tb *TeamBalancer) Init(teamSize, rankRange uint8) {
	roster := CreateEmptyRoster(teamSize, rankRange)

	tb.insertChan = make(chan Player)
	tb.insertedChan = make(chan bool)
	tb.rosterChan = make(chan *MatchRoster)

	tb.Roster = roster

	tb.runBalancing()
}

func (tb *TeamBalancer) AddPlayer(player Player) (addedToTeam, finishedBalancing bool) {
	tb.insertChan <- player

	select {
	case inserted := <-tb.insertedChan:
		return inserted, false
	case roster := <-tb.rosterChan:
		tb.Roster = roster
		close(tb.insertedChan)
		close(tb.insertChan)
		close(tb.rosterChan)
		return false, true
	}
}

func (tb *TeamBalancer) runBalancing() {
	go func(roster *MatchRoster, insertChan <-chan Player, insertedChan chan<- bool, rosterChan chan<- *MatchRoster) {
		for {
			player := <-insertChan

			inserted, finished := DoBalance(roster, player)

			if finished == true {
				rosterChan <- roster
				break
			}

			insertedChan <- inserted
			continue
		}

	}(tb.Roster, tb.insertChan, tb.insertedChan, tb.rosterChan)
}

func DoBalance(roster *MatchRoster, player Player) (inserted, finished bool) {
	roster.CurrentTeam, roster.OppositeTeam = roster.TeamA, roster.TeamB
	if len(*roster.TeamA.Players) > len(*roster.TeamB.Players) {
		roster.CurrentTeam, roster.OppositeTeam = roster.OppositeTeam, roster.CurrentTeam
	}

	currentTeam := roster.CurrentTeam
	oppositeTeam := roster.OppositeTeam

	inserted = false

	if roster.canAddPlayer(player.Rank) {
		currentTeam.AddPlayer(&player)
		inserted = true
	}

	finished = len(*currentTeam.Players) == int(roster.TeamSize) && len(*oppositeTeam.Players) == int(roster.TeamSize)

	return inserted, finished
}
