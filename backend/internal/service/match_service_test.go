package service

import (
	"errors"
	"testing"
	"time"

	"final-whistle/backend/internal/model"
	"final-whistle/backend/internal/repository"
	"gorm.io/gorm"
)

type fakeMatchRepository struct {
	listMatches      []model.Match
	listTotal        int64
	listErr          error
	aggregates       map[uint]repository.MatchAggregateRecord
	aggregatesErr    error
	matchDetail      *model.Match
	matchDetailErr   error
	playerRatings    []repository.MatchPlayerRatingRecord
	playerRatingsErr error
	reviews          []repository.MatchRecentReviewRecord
	reviewsErr       error
}

func (f *fakeMatchRepository) ListMatches(params repository.MatchListParams) ([]model.Match, int64, error) {
	return f.listMatches, f.listTotal, f.listErr
}

func (f *fakeMatchRepository) GetMatchAggregates(matchIDs []uint) (map[uint]repository.MatchAggregateRecord, error) {
	return f.aggregates, f.aggregatesErr
}

func (f *fakeMatchRepository) FindMatchByID(id uint) (*model.Match, error) {
	return f.matchDetail, f.matchDetailErr
}

func (f *fakeMatchRepository) GetPlayerRatingSummary(matchID uint, limit int) ([]repository.MatchPlayerRatingRecord, error) {
	return f.playerRatings, f.playerRatingsErr
}

func (f *fakeMatchRepository) GetRecentReviews(matchID uint, limit int) ([]repository.MatchRecentReviewRecord, error) {
	return f.reviews, f.reviewsErr
}

func TestMatchServiceListMatchesEmpty(t *testing.T) {
	svc := NewMatchService(&fakeMatchRepository{
		listMatches: []model.Match{},
		listTotal:   0,
		aggregates:  map[uint]repository.MatchAggregateRecord{},
	})

	result, err := svc.ListMatches(repository.MatchListParams{Page: 1, PageSize: 20})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Items) != 0 || result.Total != 0 {
		t.Fatalf("expected empty result, got %#v", result)
	}
}

func TestMatchServiceGetMatchDetailNotFound(t *testing.T) {
	svc := NewMatchService(&fakeMatchRepository{matchDetailErr: gorm.ErrRecordNotFound})
	_, err := svc.GetMatchDetail(1)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestMatchServiceGetMatchDetailSuccess(t *testing.T) {
	match := &model.Match{
		ID:          1,
		Competition: "Premier League",
		Season:      "2024-2025",
		Status:      model.MatchStatusFinished,
		KickoffAt:   time.Now(),
		HomeTeam:    model.Team{ID: 1, Name: "Home", Slug: "home"},
		AwayTeam:    model.Team{ID: 2, Name: "Away", Slug: "away"},
	}
	svc := NewMatchService(&fakeMatchRepository{
		matchDetail: match,
		aggregates:  map[uint]repository.MatchAggregateRecord{1: {MatchID: 1, CheckInCount: 0}},
	})

	result, err := svc.GetMatchDetail(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ID != 1 || result.Aggregates.CheckInCount != 0 {
		t.Fatalf("unexpected result: %#v", result)
	}
}
