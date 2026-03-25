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
	_, err := svc.GetTeamDetail(1)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestTeamServiceSuccess(t *testing.T) {
	team := &model.Team{ID: 1, Name: "Arsenal", Slug: "arsenal"}
	match := model.Match{
		ID:          2,
		Competition: "Premier League",
		Season:      "2024-2025",
		Status:      model.MatchStatusFinished,
		KickoffAt:   time.Now(),
		HomeTeam:    *team,
		AwayTeam:    model.Team{ID: 3, Name: "Chelsea", Slug: "chelsea"},
	}
	svc := NewTeamService(
		&fakeTeamRepository{team: team, matches: []model.Match{match}},
		&fakeMatchRepository{aggregates: map[uint]repository.MatchAggregateRecord{2: {MatchID: 2}}},
	)
	result, err := svc.GetTeamDetail(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ID != 1 || len(result.RecentMatches) != 1 {
		t.Fatalf("unexpected result: %#v", result)
	}
}
