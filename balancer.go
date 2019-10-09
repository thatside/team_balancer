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
	BasicRank uint8
}

type MatchRoster struct {
	TeamA,
	TeamB,
	CurrentTeam,
	OppositeTeam *Team
	TeamSize,
	RankRange uint8
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
	insertChan chan Player
	insertedChan chan bool
	rosterChan chan *MatchRoster
	Roster *MatchRoster
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
	case inserted := <- tb.insertedChan:
		return inserted, false
	case roster := <- tb.rosterChan:
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
	currentTeam := roster.CurrentTeam
	oppositeTeam := roster.OppositeTeam

	if len(*currentTeam.Players) == 0 && len(*oppositeTeam.Players) == 0 {
		currentTeam.BasicRank = player.Rank
		oppositeTeam.BasicRank = player.Rank
	}

	if len(*currentTeam.Players) == 0 && len(*oppositeTeam.Players) != 0 {
		if math.Abs(float64(oppositeTeam.BasicRank-player.Rank)) > float64(roster.RankRange) {
			return false, false
		}
	}

	if math.Abs(float64(oppositeTeam.BasicRank-player.Rank)) > float64(roster.RankRange) || math.Abs(float64(currentTeam.BasicRank-player.Rank)) > float64(roster.RankRange) {
		return false, false
	}

	currentTeam.AddPlayer(&player)

	if len(*currentTeam.Players) == int(roster.TeamSize) && len(*oppositeTeam.Players) == int(roster.TeamSize) {
		return false, true
	}

	roster.CurrentTeam, roster.OppositeTeam = roster.OppositeTeam, roster.CurrentTeam

	return true, false
}
