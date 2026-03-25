package service

import (
	"errors"
	"testing"

	"final-whistle/backend/internal/model"
	"final-whistle/backend/internal/repository"
	"gorm.io/gorm"
)

type fakePlayerRepository struct {
	player        *model.Player
	findErr       error
	recentMatches []repository.PlayerRecentMatchRecord
	recentErr     error
	summary       repository.PlayerDetailRatingSummary
	summaryErr    error
}

func (f *fakePlayerRepository) FindByID(id uint) (*model.Player, error) { return f.player, f.findErr }
func (f *fakePlayerRepository) ListRecentRatedMatches(playerID uint, limit int) ([]repository.PlayerRecentMatchRecord, error) {
	return f.recentMatches, f.recentErr
}
func (f *fakePlayerRepository) GetRatingSummary(playerID uint) (repository.PlayerDetailRatingSummary, error) {
	return f.summary, f.summaryErr
}

func TestPlayerServiceNotFound(t *testing.T) {
	svc := NewPlayerService(&fakePlayerRepository{findErr: gorm.ErrRecordNotFound})
	_, err := svc.GetPlayerDetail(1)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestPlayerServiceSuccess(t *testing.T) {
	svc := NewPlayerService(&fakePlayerRepository{
		player: &model.Player{
			ID:   1,
			Name: "Player",
			Slug: "player",
			Team: model.Team{ID: 2, Name: "Team", Slug: "team"},
		},
	})
	result, err := svc.GetPlayerDetail(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ID != 1 || result.Team.ID != 2 {
		t.Fatalf("unexpected result: %#v", result)
	}
}
