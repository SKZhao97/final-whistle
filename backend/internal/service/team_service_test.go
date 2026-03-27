package service

import (
	"errors"
	"testing"
	"time"

	"final-whistle/backend/internal/model"
	"final-whistle/backend/internal/repository"
	"gorm.io/gorm"
)

type fakeTeamRepository struct {
	team       *model.Team
	findErr    error
	matches    []model.Match
	matchesErr error
	rating     repository.TeamDetailRatingSummary
	ratingErr  error
}

func (f *fakeTeamRepository) FindByID(id uint) (*model.Team, error) { return f.team, f.findErr }
func (f *fakeTeamRepository) ListRecentMatches(teamID uint, limit int) ([]model.Match, error) {
	return f.matches, f.matchesErr
}
func (f *fakeTeamRepository) GetRatingSummary(teamID uint) (repository.TeamDetailRatingSummary, error) {
	return f.rating, f.ratingErr
}

func TestTeamServiceNotFound(t *testing.T) {
	svc := NewTeamService(&fakeTeamRepository{findErr: gorm.ErrRecordNotFound}, &fakeMatchRepository{})
	_, err := svc.GetTeamDetail(1, "en")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestTeamServiceSuccess(t *testing.T) {
	arsenalZh := "阿森纳"
	chelseaZh := "切尔西"
	round := "Matchday 1"
	team := &model.Team{ID: 1, Name: "Arsenal", NameZh: &arsenalZh, Slug: "arsenal"}
	match := model.Match{
		ID:          2,
		Competition: "Premier League",
		Season:      "2024-2025",
		Round:       &round,
		Status:      model.MatchStatusFinished,
		KickoffAt:   time.Now(),
		HomeTeam:    *team,
		AwayTeam:    model.Team{ID: 3, Name: "Chelsea", NameZh: &chelseaZh, Slug: "chelsea"},
	}
	svc := NewTeamService(
		&fakeTeamRepository{team: team, matches: []model.Match{match}},
		&fakeMatchRepository{aggregates: map[uint]repository.MatchAggregateRecord{2: {MatchID: 2}}},
	)
	result, err := svc.GetTeamDetail(1, "zh")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ID != 1 || len(result.RecentMatches) != 1 {
		t.Fatalf("unexpected result: %#v", result)
	}
	if result.Name != "阿森纳" {
		t.Fatalf("expected localized team name, got %#v", result.Name)
	}
	if result.RecentMatches[0].Competition != "英超" {
		t.Fatalf("expected localized competition, got %#v", result.RecentMatches[0].Competition)
	}
}
