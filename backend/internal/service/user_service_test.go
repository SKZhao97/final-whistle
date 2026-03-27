package service

import (
	"errors"
	"testing"
	"time"

	"final-whistle/backend/internal/model"
	"final-whistle/backend/internal/repository"
	"gorm.io/gorm"
)

type fakeUserRepository struct {
	profileRecord     *repository.UserProfileSummaryRecord
	profileErr        error
	historyCheckIns   []model.CheckIn
	historyTotal      int64
	historyErr        error
	lastHistoryParams repository.UserCheckInHistoryParams
}

func (f *fakeUserRepository) FindUserByID(id uint) (*model.User, error) {
	if f.profileRecord != nil {
		return &f.profileRecord.User, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (f *fakeUserRepository) GetUserProfileSummary(userID uint, recentSince time.Time) (*repository.UserProfileSummaryRecord, error) {
	return f.profileRecord, f.profileErr
}

func (f *fakeUserRepository) GetUserCheckInHistory(userID uint, params repository.UserCheckInHistoryParams) ([]model.CheckIn, int64, error) {
	f.lastHistoryParams = params
	return f.historyCheckIns, f.historyTotal, f.historyErr
}

func TestUserServiceGetProfileSummary(t *testing.T) {
	favoriteTeam := &model.Team{ID: 10, Name: "Arsenal", Slug: "arsenal"}
	tag := &model.Tag{ID: 7, Name: "Electric", NameEn: "Electric", NameZh: "电光火石", Slug: "electric"}
	repo := &fakeUserRepository{
		profileRecord: &repository.UserProfileSummaryRecord{
			User:               model.User{ID: 1, Name: "Demo User"},
			CheckInCount:       3,
			AvgMatchRating:     ptrFloat64(8.5),
			FavoriteTeamID:     ptrUint(10),
			FavoriteTeam:       favoriteTeam,
			MostUsedTagID:      ptrUint(7),
			MostUsedTag:        tag,
			RecentCheckInCount: 2,
		},
	}

	service := NewUserService(repo)
	result, err := service.GetProfileSummary(1, "zh")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.CheckInCount != 3 || result.RecentCheckInCount != 2 {
		t.Fatalf("unexpected profile counts: %+v", result)
	}
	if result.FavoriteTeam == nil || result.FavoriteTeam.Name != "Arsenal" {
		t.Fatalf("expected favorite team to be mapped")
	}
	if result.MostUsedTag == nil || result.MostUsedTag.Name != "电光火石" {
		t.Fatalf("expected most used tag to be mapped")
	}
}

func TestUserServiceGetProfileSummaryNotFound(t *testing.T) {
	service := NewUserService(&fakeUserRepository{profileErr: gorm.ErrRecordNotFound})
	_, err := service.GetProfileSummary(1, "en")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestUserServiceGetCheckInHistory(t *testing.T) {
	now := time.Now().UTC()
	repo := &fakeUserRepository{
		historyCheckIns: []model.CheckIn{
			{
				ID:             5,
				MatchID:        2,
				WatchedType:    model.WatchedTypeFull,
				SupporterSide:  model.SupporterSideHome,
				MatchRating:    8,
				HomeTeamRating: 9,
				AwayTeamRating: 6,
				WatchedAt:      now,
				CreatedAt:      now,
				UpdatedAt:      now,
				Tags:           []model.Tag{{ID: 1, Name: "Electric", NameEn: "Electric", NameZh: "电光火石", Slug: "electric"}},
				Match: model.Match{
					ID:          2,
					Competition: "Premier League",
					Season:      "2025/26",
					Status:      model.MatchStatusFinished,
					KickoffAt:   now,
					HomeTeam:    model.Team{ID: 1, Name: "Liverpool", Slug: "liverpool"},
					AwayTeam:    model.Team{ID: 2, Name: "Arsenal", Slug: "arsenal"},
				},
			},
		},
		historyTotal: 1,
	}

	service := NewUserService(repo)
	result, err := service.GetCheckInHistory(1, 0, 100, "zh")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if repo.lastHistoryParams.Page != 1 || repo.lastHistoryParams.PageSize != 50 {
		t.Fatalf("expected pagination bounds to be applied, got %+v", repo.lastHistoryParams)
	}
	if len(result.Items) != 1 {
		t.Fatalf("expected one history item")
	}
	if result.Items[0].Match.HomeTeam.Name != "Liverpool" {
		t.Fatalf("expected match context to be mapped")
	}
	if len(result.Items[0].Tags) != 1 {
		t.Fatalf("expected tags to be mapped")
	}
	if result.Items[0].Tags[0].Name != "电光火石" {
		t.Fatalf("expected localized tag label, got %#v", result.Items[0].Tags[0])
	}
}

func ptrUint(value uint) *uint {
	return &value
}

func ptrFloat64(value float64) *float64 {
	return &value
}
