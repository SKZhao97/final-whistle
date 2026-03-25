package dto

type PlayerRecentMatchDTO struct {
	Match       MatchListItemDTO `json:"match"`
	AvgRating   *float64         `json:"avgRating"`
	RatingCount int64            `json:"ratingCount"`
}

type PlayerRatingSummaryDTO struct {
	AvgRating   *float64 `json:"avgRating"`
	RatingCount int64    `json:"ratingCount"`
}

type PlayerDetailDTO struct {
	ID            uint                   `json:"id"`
	Name          string                 `json:"name"`
	Slug          string                 `json:"slug"`
	Position      *string                `json:"position,omitempty"`
	AvatarURL     *string                `json:"avatarUrl,omitempty"`
	Team          TeamSummaryDTO         `json:"team"`
	RecentMatches []PlayerRecentMatchDTO `json:"recentMatches"`
	RatingSummary PlayerRatingSummaryDTO `json:"ratingSummary"`
}
