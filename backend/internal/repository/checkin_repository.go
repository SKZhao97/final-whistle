package repository

import (
	"final-whistle/backend/internal/model"
	"gorm.io/gorm"
)

type CheckInRepository interface {
	FindMatchByID(id uint) (*model.Match, error)
	FindCheckInByUserAndMatch(userID, matchID uint) (*model.CheckIn, error)
	GetEligiblePlayerIDs(matchID uint, playerIDs []uint) (map[uint]struct{}, error)
	GetActiveTagsByIDs(tagIDs []uint) ([]model.Tag, error)
	WithTransaction(fn func(repo CheckInRepository) error) error
	CreateCheckIn(checkIn *model.CheckIn) error
	UpdateCheckIn(checkIn *model.CheckIn) error
	ReplacePlayerRatings(checkInID uint, ratings []model.PlayerRating) error
	ReplaceCheckInTags(checkInID uint, tagIDs []uint) error
}

type GormCheckInRepository struct {
	*BaseRepository
}

func NewCheckInRepository(db *gorm.DB) *GormCheckInRepository {
	return &GormCheckInRepository{BaseRepository: NewBaseRepository(db)}
}

func (r *GormCheckInRepository) FindMatchByID(id uint) (*model.Match, error) {
	var match model.Match
	if err := r.DB.First(&match, id).Error; err != nil {
		return nil, err
	}
	return &match, nil
}

func (r *GormCheckInRepository) FindCheckInByUserAndMatch(userID, matchID uint) (*model.CheckIn, error) {
	var checkIn model.CheckIn
	if err := r.DB.
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Order("tags.sort_order ASC, tags.id ASC")
		}).
		Preload("PlayerRatings", func(db *gorm.DB) *gorm.DB {
			return db.Order("player_ratings.id ASC")
		}).
		Preload("PlayerRatings.Player.Team").
		Where("user_id = ? AND match_id = ?", userID, matchID).
		First(&checkIn).Error; err != nil {
		return nil, err
	}
	return &checkIn, nil
}

func (r *GormCheckInRepository) GetEligiblePlayerIDs(matchID uint, playerIDs []uint) (map[uint]struct{}, error) {
	result := make(map[uint]struct{}, len(playerIDs))
	if len(playerIDs) == 0 {
		return result, nil
	}

	var rows []struct {
		PlayerID uint
	}
	if err := r.DB.
		Table("match_players").
		Select("player_id").
		Where("match_id = ? AND player_id IN ?", matchID, playerIDs).
		Group("player_id").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	for _, row := range rows {
		result[row.PlayerID] = struct{}{}
	}

	return result, nil
}

func (r *GormCheckInRepository) GetActiveTagsByIDs(tagIDs []uint) ([]model.Tag, error) {
	if len(tagIDs) == 0 {
		return []model.Tag{}, nil
	}

	var tags []model.Tag
	if err := r.DB.
		Where("id IN ? AND is_active = ?", tagIDs, true).
		Order("sort_order ASC, id ASC").
		Find(&tags).Error; err != nil {
		return nil, err
	}

	return tags, nil
}

func (r *GormCheckInRepository) WithTransaction(fn func(repo CheckInRepository) error) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		return fn(NewCheckInRepository(tx))
	})
}

func (r *GormCheckInRepository) CreateCheckIn(checkIn *model.CheckIn) error {
	return r.DB.Create(checkIn).Error
}

func (r *GormCheckInRepository) UpdateCheckIn(checkIn *model.CheckIn) error {
	return r.DB.Save(checkIn).Error
}

func (r *GormCheckInRepository) ReplacePlayerRatings(checkInID uint, ratings []model.PlayerRating) error {
	if err := r.DB.Where("check_in_id = ?", checkInID).Delete(&model.PlayerRating{}).Error; err != nil {
		return err
	}
	if len(ratings) == 0 {
		return nil
	}
	return r.DB.Create(&ratings).Error
}

func (r *GormCheckInRepository) ReplaceCheckInTags(checkInID uint, tagIDs []uint) error {
	if err := r.DB.Where("check_in_id = ?", checkInID).Delete(&model.CheckInTag{}).Error; err != nil {
		return err
	}
	if len(tagIDs) == 0 {
		return nil
	}

	joinRows := make([]model.CheckInTag, 0, len(tagIDs))
	for _, tagID := range tagIDs {
		joinRows = append(joinRows, model.CheckInTag{
			CheckInID: checkInID,
			TagID:     tagID,
		})
	}
	return r.DB.Create(&joinRows).Error
}
