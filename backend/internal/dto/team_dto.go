package dto

type TeamRatingSummaryDTO struct {
	AvgRating   *float64 `json:"avgRating"`
	RatingCount int64    `json:"ratingCount"`
}

type TeamDetailDTO struct {
	ID            uint                 `json:"id"`
	Name          string               `json:"name"`
	ShortName     *string              `json:"shortName,omitempty"`
	Slug          string               `json:"slug"`
	LogoURL       *string              `json:"logoUrl,omitempty"`
	RecentMatches []MatchListItemDTO   `json:"recentMatches"`
	RatingSummary TeamRatingSummaryDTO `json:"ratingSummary"`
}
