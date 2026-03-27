package testutil

import (
	"context"
	"time"

	syncprovider "final-whistle/backend/internal/sync/provider"
)

type FakeProvider struct{}

func NewFakeProvider() *FakeProvider {
	return &FakeProvider{}
}

func (p *FakeProvider) ListCompetitionTeams(ctx context.Context, competitionCode string) ([]syncprovider.Team, error) {
	shortArsenal := "ARS"
	shortLiverpool := "LIV"
	return []syncprovider.Team{
		{ExternalID: "1", Name: "Arsenal", ShortName: &shortArsenal},
		{ExternalID: "2", Name: "Liverpool", ShortName: &shortLiverpool},
	}, nil
}

func (p *FakeProvider) ListTeamPlayers(ctx context.Context, teamExternalID string) ([]syncprovider.Player, error) {
	positionRW := "Right Winger"
	positionAM := "Attacking Midfield"
	positionFW := "Forward"
	positionCB := "Centre-Back"
	switch teamExternalID {
	case "1":
		return []syncprovider.Player{
			{ExternalID: "101", TeamExternalID: "1", Name: "Bukayo Saka", Position: &positionRW},
			{ExternalID: "102", TeamExternalID: "1", Name: "Martin Odegaard", Position: &positionAM},
		}, nil
	case "2":
		return []syncprovider.Player{
			{ExternalID: "201", TeamExternalID: "2", Name: "Mohamed Salah", Position: &positionFW},
			{ExternalID: "202", TeamExternalID: "2", Name: "Virgil van Dijk", Position: &positionCB},
		}, nil
	default:
		return []syncprovider.Player{}, nil
	}
}

func (p *FakeProvider) ListCompetitionMatches(ctx context.Context, competitionCode string, dateFrom, dateTo string) ([]syncprovider.Match, error) {
	shortArsenal := "ARS"
	shortLiverpool := "LIV"
	homeScore := 2
	awayScore := 1
	round := "Matchday 30"
	venue := "Emirates Stadium"
	kickoff := time.Date(2026, 3, 27, 20, 0, 0, 0, time.UTC)
	return []syncprovider.Match{
		{
			ExternalID:  "5001",
			Competition: "Premier League",
			Season:      "2026",
			Round:       &round,
			Status:      "FINISHED",
			KickoffAt:   kickoff,
			HomeTeam: syncprovider.MatchTeam{
				ExternalID: "1",
				Name:       "Arsenal",
				ShortName:  &shortArsenal,
			},
			AwayTeam: syncprovider.MatchTeam{
				ExternalID: "2",
				Name:       "Liverpool",
				ShortName:  &shortLiverpool,
			},
			HomeScore: &homeScore,
			AwayScore: &awayScore,
			Venue:     &venue,
		},
	}, nil
}

func (p *FakeProvider) GetMatchRoster(ctx context.Context, externalMatchID string) ([]syncprovider.MatchRosterPlayer, error) {
	return []syncprovider.MatchRosterPlayer{
		{PlayerExternalID: "101", TeamExternalID: "1"},
		{PlayerExternalID: "102", TeamExternalID: "1"},
		{PlayerExternalID: "201", TeamExternalID: "2"},
		{PlayerExternalID: "202", TeamExternalID: "2"},
	}, nil
}
