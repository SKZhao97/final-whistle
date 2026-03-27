// Package model 定义数据模型和数据库表结构。
// 本文件包含签到相关模型定义。
package model

import (
	"time"
)

// WatchedType 定义签到时的观看类型。
type WatchedType string

const (
	// WatchedTypeFull 表示观看了整场比赛。
	WatchedTypeFull WatchedType = "FULL"
	// WatchedTypePartial 表示观看了部分比赛。
	WatchedTypePartial WatchedType = "PARTIAL"
	// WatchedTypeHighlights 表示观看了比赛集锦。
	WatchedTypeHighlights WatchedType = "HIGHLIGHTS"
)

// SupporterSide 定义签到时的支持方。
type SupporterSide string

const (
	// SupporterSideHome 表示支持主队。
	SupporterSideHome SupporterSide = "HOME"
	// SupporterSideAway 表示支持客队。
	SupporterSideAway SupporterSide = "AWAY"
	// SupporterSideNeutral 表示中立。
	SupporterSideNeutral SupporterSide = "NEUTRAL"
)

// CheckIn 表示用户对一场比赛的签到记录。
type CheckIn struct {
	ID              uint          `gorm:"primaryKey"`                    // 主键
	UserID          uint          `gorm:"not null;index"`                // 用户ID，外键，索引
	MatchID         uint          `gorm:"not null;index"`                // 比赛ID，外键，索引
	WatchedType     WatchedType   `gorm:"size:20;not null"`              // 观看类型
	SupporterSide   SupporterSide `gorm:"size:20;not null"`              // 支持方
	MatchRating     int           `gorm:"not null;check:match_rating BETWEEN 1 AND 10"`      // 比赛评分，1-10
	HomeTeamRating  int           `gorm:"not null;check:home_team_rating BETWEEN 1 AND 10"`  // 主队评分，1-10
	AwayTeamRating  int           `gorm:"not null;check:away_team_rating BETWEEN 1 AND 10"`  // 客队评分，1-10
	ShortReview     *string       `gorm:"size:280"`                      // 简短评价，可选，最大280字符
	WatchedAt       time.Time     `gorm:"not null"`                      // 观看时间
	CreatedAt       time.Time     `gorm:"autoCreateTime;index"`          // 创建时间，自动生成，索引
	UpdatedAt       time.Time     `gorm:"autoUpdateTime"`                // 更新时间，自动更新

	// 关联关系
	User          User            `gorm:"foreignKey:UserID"`             // 所属用户
	Match         Match           `gorm:"foreignKey:MatchID"`            // 所属比赛
	PlayerRatings []PlayerRating  `gorm:"foreignKey:CheckInID"`          // 球员评分列表
	Tags          []Tag           `gorm:"many2many:checkin_tags;foreignKey:ID;joinForeignKey:CheckInID;joinReferences:TagID"` // 标签列表
}

// TableName 指定GORM使用的数据库表名。
func (CheckIn) TableName() string {
	return "check_ins"
}