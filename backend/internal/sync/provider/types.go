package provider

import (
	"context"
	"time"
)

type Team struct {
	ExternalID string
	Name       string
	ShortName  *string
	CrestURL   *string
}

type Player struct {
	ExternalID     string
	TeamExternalID string
	Name           string
	Position       *string
}

type MatchTeam struct {
	ExternalID string
	Name       string
	ShortName  *string
	CrestURL   *string
}

type Match struct {
	ExternalID  string
	Competition string
	Season      string
	Round       *string
	Status      string
	KickoffAt   time.Time
	HomeTeam    MatchTeam
	AwayTeam    MatchTeam
	HomeScore   *int
	AwayScore   *int
	Venue       *string
}

type MatchRosterPlayer struct {
	PlayerExternalID string
	TeamExternalID   string
}

type Client interface {
	ListCompetitionTeams(ctx context.Context, competitionCode string) ([]Team, error)
	ListTeamPlayers(ctx context.Context, teamExternalID string) ([]Player, error)
	ListCompetitionMatches(ctx context.Context, competitionCode string, dateFrom, dateTo string) ([]Match, error)
	GetMatchRoster(ctx context.Context, externalMatchID string) ([]MatchRosterPlayer, error)
}
