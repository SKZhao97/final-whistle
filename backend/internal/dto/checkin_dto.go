package dto

import "time"

type PlayerRatingInputDTO struct {
	PlayerID uint    `json:"playerId"`
	Rating   int     `json:"rating"`
	Note     *string `json:"note,omitempty"`
}

type CheckInPlayerRatingDTO struct {
	ID     uint             `json:"id"`
	Player PlayerSummaryDTO `json:"player"`
	Rating int              `json:"rating"`
	Note   *string          `json:"note,omitempty"`
}

type CheckInDetailDTO struct {
	ID             uint                     `json:"id"`
	MatchID        uint                     `json:"matchId"`
	WatchedType    string                   `json:"watchedType"`
	SupporterSide  string                   `json:"supporterSide"`
	MatchRating    int                      `json:"matchRating"`
	HomeTeamRating int                      `json:"homeTeamRating"`
	AwayTeamRating int                      `json:"awayTeamRating"`
	ShortReview    *string                  `json:"shortReview,omitempty"`
	WatchedAt      time.Time                `json:"watchedAt"`
	Tags           []TagDTO                 `json:"tags"`
	PlayerRatings  []CheckInPlayerRatingDTO `json:"playerRatings"`
	CreatedAt      time.Time                `json:"createdAt"`
	UpdatedAt      time.Time                `json:"updatedAt"`
}

type UpsertCheckInRequestDTO struct {
	WatchedType    string                 `json:"watchedType"`
	SupporterSide  string                 `json:"supporterSide"`
	MatchRating    int                    `json:"matchRating"`
	HomeTeamRating int                    `json:"homeTeamRating"`
	AwayTeamRating int                    `json:"awayTeamRating"`
	ShortReview    *string                `json:"shortReview,omitempty"`
	WatchedAt      time.Time              `json:"watchedAt"`
	Tags           []uint                 `json:"tags"`
	PlayerRatings  []PlayerRatingInputDTO `json:"playerRatings"`
}
