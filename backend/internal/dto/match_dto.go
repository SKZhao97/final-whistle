package dto

import "time"

type TeamSummaryDTO struct {
	ID        uint    `json:"id"`
	Name      string  `json:"name"`
	ShortName *string `json:"shortName,omitempty"`
	Slug      string  `json:"slug"`
	LogoURL   *string `json:"logoUrl,omitempty"`
}

type PlayerSummaryDTO struct {
	ID        uint           `json:"id"`
	Name      string         `json:"name"`
	Slug      string         `json:"slug"`
	Position  *string        `json:"position,omitempty"`
	AvatarURL *string        `json:"avatarUrl,omitempty"`
	Team      TeamSummaryDTO `json:"team"`
}

type MatchAggregateSummaryDTO struct {
	MatchRatingAvg    *float64 `json:"matchRatingAvg"`
	HomeTeamRatingAvg *float64 `json:"homeTeamRatingAvg"`
	AwayTeamRatingAvg *float64 `json:"awayTeamRatingAvg"`
	CheckInCount      int64    `json:"checkInCount"`
}

type MatchListItemDTO struct {
	ID          uint                     `json:"id"`
	Competition string                   `json:"competition"`
	Season      string                   `json:"season"`
	Round       *string                  `json:"round,omitempty"`
	Status      string                   `json:"status"`
	KickoffAt   time.Time                `json:"kickoffAt"`
	HomeTeam    TeamSummaryDTO           `json:"homeTeam"`
	AwayTeam    TeamSummaryDTO           `json:"awayTeam"`
	HomeScore   *int                     `json:"homeScore,omitempty"`
	AwayScore   *int                     `json:"awayScore,omitempty"`
	Aggregates  MatchAggregateSummaryDTO `json:"aggregates"`
}

type MatchListResponseDTO struct {
	Items    []MatchListItemDTO `json:"items"`
	Page     int                `json:"page"`
	PageSize int                `json:"pageSize"`
	Total    int64              `json:"total"`
}

type MatchPlayerRatingSummaryDTO struct {
	Player      PlayerSummaryDTO `json:"player"`
	AvgRating   *float64         `json:"avgRating"`
	RatingCount int64            `json:"ratingCount"`
}

type MatchRecentReviewDTO struct {
	ID          uint           `json:"id"`
	User        UserSummaryDTO `json:"user"`
	MatchRating int            `json:"matchRating"`
	ShortReview string         `json:"shortReview"`
	Tags        []TagDTO       `json:"tags"`
	CreatedAt   time.Time      `json:"createdAt"`
}

type MatchDetailDTO struct {
	ID            uint                          `json:"id"`
	Competition   string                        `json:"competition"`
	Season        string                        `json:"season"`
	Round         *string                       `json:"round,omitempty"`
	Status        string                        `json:"status"`
	KickoffAt     time.Time                     `json:"kickoffAt"`
	HomeTeam      TeamSummaryDTO                `json:"homeTeam"`
	AwayTeam      TeamSummaryDTO                `json:"awayTeam"`
	HomeScore     *int                          `json:"homeScore,omitempty"`
	AwayScore     *int                          `json:"awayScore,omitempty"`
	Venue         *string                       `json:"venue,omitempty"`
	Aggregates    MatchAggregateSummaryDTO      `json:"aggregates"`
	PlayerRatings []MatchPlayerRatingSummaryDTO `json:"playerRatings"`
	RecentReviews []MatchRecentReviewDTO        `json:"recentReviews"`
}
