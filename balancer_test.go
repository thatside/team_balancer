package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTeamBalancer_DoBalancing(t *testing.T) {
	t.Skip("Unfinished")

	balancer := TeamBalancer{}
	balancer.Init(3, 1)
}

type TestResult struct {
	Inserted,
	Finished bool
}

type TestTuple struct {
	Description    string
	BaseRoster     *MatchRoster
	InsertedPlayer *Player
	Result         TestResult
}

func provideTestData() []TestTuple {
	return []TestTuple{
		{
			Description:    "Empty teams insert one player",
			BaseRoster:     CreateEmptyRoster(1, 0),
			InsertedPlayer: CreatePlayer(1),
			Result: TestResult{
				Inserted: true,
				Finished: false,
			},
		},
		{
			Description: "One full team insert one player with correct rank",
			BaseRoster: func() *MatchRoster {
				roster := CreateEmptyRoster(1, 0)
				roster.TeamA.AddPlayer(CreatePlayer(1))
				return roster
			}(),
			InsertedPlayer: CreatePlayer(1),
			Result: TestResult{
				Inserted: true,
				Finished: true,
			},
		}, {
			Description: "Two full teams insert one player with correct rank",
			BaseRoster: func() *MatchRoster {
				roster := CreateEmptyRoster(1, 0)
				roster.TeamA.AddPlayer(CreatePlayer(1))
				roster.TeamB.AddPlayer(CreatePlayer(1))
				return roster
			}(),
			InsertedPlayer: CreatePlayer(1),
			Result: TestResult{
				Inserted: false,
				Finished: true,
			},
		}, {
			Description: "One full team, one half empty team, insert fitting player",
			BaseRoster: func() *MatchRoster {
				roster := CreateEmptyRoster(2, 1)
				roster.TeamA.AddPlayer(CreatePlayer(1))
				roster.TeamA.AddPlayer(CreatePlayer(2))
				roster.TeamB.AddPlayer(CreatePlayer(1))
				return roster
			}(),
			InsertedPlayer: CreatePlayer(2),
			Result: TestResult{
				Inserted: true,
				Finished: true,
			},
		}, {
			Description: "One full team, one half empty team, insert not fitting player",
			BaseRoster: func() *MatchRoster {
				roster := CreateEmptyRoster(2, 1)
				roster.TeamA.AddPlayer(CreatePlayer(1))
				roster.TeamA.AddPlayer(CreatePlayer(2))
				roster.TeamB.AddPlayer(CreatePlayer(1))
				return roster
			}(),
			InsertedPlayer: CreatePlayer(3),
			Result: TestResult{
				Inserted: false,
				Finished: false,
			},
		},
	}
}

func TestDoBalanceSingleInsert(t *testing.T) {
	for _, tuple := range provideTestData() {
		t.Run(tuple.Description, func(t *testing.T) {
			r := &TestResult{}
			r.Inserted, r.Finished = DoBalance(tuple.BaseRoster, *tuple.InsertedPlayer)
			assert.Equal(t, tuple.Result, *r)
		})
	}
}
