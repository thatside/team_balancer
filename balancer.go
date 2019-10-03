package main

import "math"

type Player struct {
	Name string
	Rank uint8
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
	var playersA, playersB []Player
	roster := &MatchRoster{
		TeamA: &Team{Players: &playersA},
		TeamB: &Team{Players: &playersB},
	}

	roster.CurrentTeam = roster.TeamA
	roster.OppositeTeam = roster.TeamB
	roster.TeamSize = teamSize
	roster.RankRange = rankRange

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
		return false, true
	}
}

func (tb *TeamBalancer) runBalancing() {
	go func(roster *MatchRoster, insertChan <-chan Player, insertedChan chan<- bool, rosterChan chan<- *MatchRoster) {
		for {
			player := <-insertChan

			currentTeam := roster.CurrentTeam
			oppositeTeam := roster.OppositeTeam

			if len(*currentTeam.Players) == 0 && len(*oppositeTeam.Players) == 0 {
				currentTeam.BasicRank = player.Rank
				oppositeTeam.BasicRank = player.Rank
			}

			if len(*currentTeam.Players) == 0 && len(*oppositeTeam.Players) != 0 {
				if math.Abs(float64(oppositeTeam.BasicRank-player.Rank)) > float64(roster.RankRange) {
					insertedChan <- false
					continue
				}
			}

			if math.Abs(float64(oppositeTeam.BasicRank-player.Rank)) > float64(roster.RankRange) || math.Abs(float64(currentTeam.BasicRank-player.Rank)) > float64(roster.RankRange) {
				insertedChan <- false
				continue
			}

			newPlayers := append(*currentTeam.Players, player)
			currentTeam.Players = &newPlayers

			if len(*currentTeam.Players) == int(roster.TeamSize) && len(*oppositeTeam.Players) == int(roster.TeamSize) {
				rosterChan <- roster
				close(insertedChan)
				close(rosterChan)
				break
			}

			roster.CurrentTeam, roster.OppositeTeam = roster.OppositeTeam, roster.CurrentTeam

			insertedChan <- true
			continue
		}

	}(tb.Roster, tb.insertChan, tb.insertedChan, tb.rosterChan)
}
