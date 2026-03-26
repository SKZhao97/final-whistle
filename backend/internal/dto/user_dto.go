package dto

import "time"

type UserProfileSummaryDTO struct {
	User               UserSummaryDTO  `json:"user"`
	CheckInCount       int             `json:"checkInCount"`
	AvgMatchRating     *float64        `json:"avgMatchRating,omitempty"`
	FavoriteTeamID     *uint           `json:"favoriteTeamId,omitempty"`
	FavoriteTeam       *TeamSummaryDTO `json:"favoriteTeam,omitempty"`
	MostUsedTagID      *uint           `json:"mostUsedTagId,omitempty"`
	MostUsedTag        *TagDTO         `json:"mostUsedTag,omitempty"`
	RecentCheckInCount int             `json:"recentCheckInCount"`
}

type UserCheckInHistoryItemDTO struct {
	ID             uint             `json:"id"`
	MatchID        uint             `json:"matchId"`
	WatchedType    string           `json:"watchedType"`
	SupporterSide  string           `json:"supporterSide"`
	MatchRating    int              `json:"matchRating"`
	HomeTeamRating int              `json:"homeTeamRating"`
	AwayTeamRating int              `json:"awayTeamRating"`
	ShortReview    *string          `json:"shortReview,omitempty"`
	WatchedAt      time.Time        `json:"watchedAt"`
	CreatedAt      time.Time        `json:"createdAt"`
	UpdatedAt      time.Time        `json:"updatedAt"`
	Tags           []TagDTO         `json:"tags"`
	Match          MatchListItemDTO `json:"match"`
}

type UserCheckInHistoryResponseDTO struct {
	Items    []UserCheckInHistoryItemDTO `json:"items"`
	Page     int                         `json:"page"`
	PageSize int                         `json:"pageSize"`
	Total    int64                       `json:"total"`
}
