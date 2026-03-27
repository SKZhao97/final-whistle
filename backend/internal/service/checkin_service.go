package service

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"final-whistle/backend/internal/dto"
	"final-whistle/backend/internal/model"
	"final-whistle/backend/internal/repository"
	"gorm.io/gorm"
)

const (
	maxShortReviewLength = 280
	maxPlayerNoteLength  = 80
)

var (
	ErrCheckInAlreadyExists = errors.New("check-in already exists")
	ErrCheckInMissing       = errors.New("check-in not found")
)

type CheckInValidationError struct {
	Message string
	Details map[string]interface{}
}

func (e *CheckInValidationError) Error() string {
	return e.Message
}

type CheckInService interface {
	GetMyCheckIn(matchID, userID uint, locale string) (*dto.CheckInDetailDTO, error)
	CreateCheckIn(matchID, userID uint, req dto.UpsertCheckInRequestDTO, locale string) (*dto.CheckInDetailDTO, error)
	UpdateCheckIn(matchID, userID uint, req dto.UpsertCheckInRequestDTO, locale string) (*dto.CheckInDetailDTO, error)
}

type checkInService struct {
	repo repository.CheckInRepository
}

func NewCheckInService(repo repository.CheckInRepository) CheckInService {
	return &checkInService{repo: repo}
}

func (s *checkInService) GetMyCheckIn(matchID, userID uint, locale string) (*dto.CheckInDetailDTO, error) {
	if _, err := s.loadMatch(matchID); err != nil {
		return nil, err
	}

	checkIn, err := s.repo.FindCheckInByUserAndMatch(userID, matchID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return mapCheckInDetail(checkIn, locale), nil
}

func (s *checkInService) CreateCheckIn(matchID, userID uint, req dto.UpsertCheckInRequestDTO, locale string) (*dto.CheckInDetailDTO, error) {
	req = normalizeUpsertRequest(req)

	match, err := s.loadMatch(matchID)
	if err != nil {
		return nil, err
	}
	if err := validateMatchIsFinished(match); err != nil {
		return nil, err
	}
	if err := s.validateUpsertPayload(matchID, req); err != nil {
		return nil, err
	}

	if _, err := s.repo.FindCheckInByUserAndMatch(userID, matchID); err == nil {
		return nil, ErrCheckInAlreadyExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	checkIn := buildCheckInModel(userID, matchID, req)
	tagIDs := slices.Clone(req.Tags)

	var created *model.CheckIn
	if err := s.repo.WithTransaction(func(txRepo repository.CheckInRepository) error {
		if err := txRepo.CreateCheckIn(checkIn); err != nil {
			return err
		}

		if err := txRepo.ReplacePlayerRatings(checkIn.ID, buildPlayerRatings(checkIn.ID, req.PlayerRatings)); err != nil {
			return err
		}
		if err := txRepo.ReplaceCheckInTags(checkIn.ID, tagIDs); err != nil {
			return err
		}

		loaded, err := txRepo.FindCheckInByUserAndMatch(userID, matchID)
		if err != nil {
			return err
		}
		created = loaded
		return nil
	}); err != nil {
		return nil, err
	}

	return mapCheckInDetail(created, locale), nil
}

func (s *checkInService) UpdateCheckIn(matchID, userID uint, req dto.UpsertCheckInRequestDTO, locale string) (*dto.CheckInDetailDTO, error) {
	req = normalizeUpsertRequest(req)

	match, err := s.loadMatch(matchID)
	if err != nil {
		return nil, err
	}
	if err := validateMatchIsFinished(match); err != nil {
		return nil, err
	}
	if err := s.validateUpsertPayload(matchID, req); err != nil {
		return nil, err
	}

	checkIn, err := s.repo.FindCheckInByUserAndMatch(userID, matchID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCheckInMissing
		}
		return nil, err
	}

	checkIn.WatchedType = model.WatchedType(strings.TrimSpace(req.WatchedType))
	checkIn.SupporterSide = model.SupporterSide(strings.TrimSpace(req.SupporterSide))
	checkIn.MatchRating = req.MatchRating
	checkIn.HomeTeamRating = req.HomeTeamRating
	checkIn.AwayTeamRating = req.AwayTeamRating
	checkIn.ShortReview = normalizeOptionalString(req.ShortReview)
	checkIn.WatchedAt = req.WatchedAt

	var updated *model.CheckIn
	if err := s.repo.WithTransaction(func(txRepo repository.CheckInRepository) error {
		if err := txRepo.UpdateCheckIn(checkIn); err != nil {
			return err
		}
		if err := txRepo.ReplacePlayerRatings(checkIn.ID, buildPlayerRatings(checkIn.ID, req.PlayerRatings)); err != nil {
			return err
		}
		if err := txRepo.ReplaceCheckInTags(checkIn.ID, slices.Clone(req.Tags)); err != nil {
			return err
		}

		loaded, err := txRepo.FindCheckInByUserAndMatch(userID, matchID)
		if err != nil {
			return err
		}
		updated = loaded
		return nil
	}); err != nil {
		return nil, err
	}

	return mapCheckInDetail(updated, locale), nil
}

func (s *checkInService) loadMatch(matchID uint) (*model.Match, error) {
	match, err := s.repo.FindMatchByID(matchID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return match, nil
}

func validateMatchIsFinished(match *model.Match) error {
	if match.Status != model.MatchStatusFinished {
		return &CheckInValidationError{Message: "check-ins are only allowed for finished matches"}
	}
	return nil
}

func (s *checkInService) validateUpsertPayload(matchID uint, req dto.UpsertCheckInRequestDTO) error {
	if req.WatchedType != string(model.WatchedTypeFull) &&
		req.WatchedType != string(model.WatchedTypePartial) &&
		req.WatchedType != string(model.WatchedTypeHighlights) {
		return &CheckInValidationError{Message: "invalid watchedType"}
	}
	if req.SupporterSide != string(model.SupporterSideHome) &&
		req.SupporterSide != string(model.SupporterSideAway) &&
		req.SupporterSide != string(model.SupporterSideNeutral) {
		return &CheckInValidationError{Message: "invalid supporterSide"}
	}
	if req.WatchedAt.IsZero() {
		return &CheckInValidationError{Message: "watchedAt is required"}
	}
	if err := validateRatingRange("matchRating", req.MatchRating); err != nil {
		return err
	}
	if err := validateRatingRange("homeTeamRating", req.HomeTeamRating); err != nil {
		return err
	}
	if err := validateRatingRange("awayTeamRating", req.AwayTeamRating); err != nil {
		return err
	}
	if req.ShortReview != nil && len(*req.ShortReview) > maxShortReviewLength {
		return &CheckInValidationError{Message: "shortReview must be 280 characters or fewer"}
	}
	playerIDs := make([]uint, 0, len(req.PlayerRatings))
	seenPlayers := make(map[uint]struct{}, len(req.PlayerRatings))
	for _, item := range req.PlayerRatings {
		if item.PlayerID == 0 {
			return &CheckInValidationError{Message: "playerRatings[].playerId is required"}
		}
		if _, exists := seenPlayers[item.PlayerID]; exists {
			return &CheckInValidationError{Message: "duplicate player ratings are not allowed"}
		}
		seenPlayers[item.PlayerID] = struct{}{}
		playerIDs = append(playerIDs, item.PlayerID)
		if err := validateRatingRange("playerRatings[].rating", item.Rating); err != nil {
			return err
		}
		if item.Note != nil && len(strings.TrimSpace(*item.Note)) > maxPlayerNoteLength {
			return &CheckInValidationError{Message: "playerRatings[].note must be 80 characters or fewer"}
		}
	}

	eligiblePlayerIDs, err := s.repo.GetEligiblePlayerIDs(matchID, playerIDs)
	if err != nil {
		return err
	}
	for _, playerID := range playerIDs {
		if _, ok := eligiblePlayerIDs[playerID]; !ok {
			return &CheckInValidationError{Message: "playerRatings include players not in this match"}
		}
	}

	if hasDuplicates(req.Tags) {
		return &CheckInValidationError{Message: "duplicate tags are not allowed"}
	}
	activeTags, err := s.repo.GetActiveTagsByIDs(req.Tags)
	if err != nil {
		return err
	}
	if len(activeTags) != len(req.Tags) {
		return &CheckInValidationError{Message: "tags include invalid or inactive tag ids"}
	}

	return nil
}

func normalizeUpsertRequest(req dto.UpsertCheckInRequestDTO) dto.UpsertCheckInRequestDTO {
	req.WatchedType = strings.TrimSpace(req.WatchedType)
	req.SupporterSide = strings.TrimSpace(req.SupporterSide)
	req.ShortReview = normalizeOptionalString(req.ShortReview)

	for i := range req.PlayerRatings {
		req.PlayerRatings[i].Note = normalizeOptionalString(req.PlayerRatings[i].Note)
	}

	return req
}

func validateRatingRange(field string, rating int) error {
	if rating < 1 || rating > 10 {
		return &CheckInValidationError{Message: fmt.Sprintf("%s must be between 1 and 10", field)}
	}
	return nil
}

func buildCheckInModel(userID, matchID uint, req dto.UpsertCheckInRequestDTO) *model.CheckIn {
	return &model.CheckIn{
		UserID:         userID,
		MatchID:        matchID,
		WatchedType:    model.WatchedType(strings.TrimSpace(req.WatchedType)),
		SupporterSide:  model.SupporterSide(strings.TrimSpace(req.SupporterSide)),
		MatchRating:    req.MatchRating,
		HomeTeamRating: req.HomeTeamRating,
		AwayTeamRating: req.AwayTeamRating,
		ShortReview:    normalizeOptionalString(req.ShortReview),
		WatchedAt:      req.WatchedAt,
	}
}

func buildPlayerRatings(checkInID uint, ratings []dto.PlayerRatingInputDTO) []model.PlayerRating {
	result := make([]model.PlayerRating, 0, len(ratings))
	for _, rating := range ratings {
		result = append(result, model.PlayerRating{
			CheckInID: checkInID,
			PlayerID:  rating.PlayerID,
			Rating:    rating.Rating,
			Note:      normalizeOptionalString(rating.Note),
		})
	}
	return result
}

func mapCheckInDetail(checkIn *model.CheckIn, locale string) *dto.CheckInDetailDTO {
	if checkIn == nil {
		return nil
	}

	tags := make([]dto.TagDTO, 0, len(checkIn.Tags))
	for _, tag := range checkIn.Tags {
		tags = append(tags, toTagDTO(tag, locale))
	}

	playerRatings := make([]dto.CheckInPlayerRatingDTO, 0, len(checkIn.PlayerRatings))
	for _, rating := range checkIn.PlayerRatings {
		playerRatings = append(playerRatings, dto.CheckInPlayerRatingDTO{
			ID: rating.ID,
			Player: dto.PlayerSummaryDTO{
				ID:        rating.Player.ID,
				Name:      rating.Player.Name,
				Slug:      rating.Player.Slug,
				Position:  rating.Player.Position,
				AvatarURL: rating.Player.AvatarURL,
				Team:      toTeamSummaryDTO(rating.Player.Team, locale),
			},
			Rating: rating.Rating,
			Note:   rating.Note,
		})
	}

	return &dto.CheckInDetailDTO{
		ID:             checkIn.ID,
		MatchID:        checkIn.MatchID,
		WatchedType:    string(checkIn.WatchedType),
		SupporterSide:  string(checkIn.SupporterSide),
		MatchRating:    checkIn.MatchRating,
		HomeTeamRating: checkIn.HomeTeamRating,
		AwayTeamRating: checkIn.AwayTeamRating,
		ShortReview:    checkIn.ShortReview,
		WatchedAt:      checkIn.WatchedAt,
		Tags:           tags,
		PlayerRatings:  playerRatings,
		CreatedAt:      checkIn.CreatedAt,
		UpdatedAt:      checkIn.UpdatedAt,
	}
}

func normalizeOptionalString(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func hasDuplicates(values []uint) bool {
	seen := make(map[uint]struct{}, len(values))
	for _, value := range values {
		if _, exists := seen[value]; exists {
			return true
		}
		seen[value] = struct{}{}
	}
	return false
}
