package seed

import (
	"final-whistle/backend/internal/model"
	"gorm.io/gorm"
)

// ResetDatabase deletes all data from all tables and re-seeds the database
// This is intended for development use only
func ResetDatabase(db *gorm.DB) error {
	// Disable foreign key checks for PostgreSQL
	if err := db.Exec("SET session_replication_role = 'replica'").Error; err != nil {
		// If PostgreSQL-specific command fails, try to delete in order
		return resetWithOrderedDeletes(db)
	}

	// Delete all data from tables in reverse dependency order
	tables := []interface{}{
		&model.CheckInTag{},
		&model.PlayerRating{},
		&model.CheckIn{},
		&model.Session{},
		&model.MatchPlayer{},
		&model.Match{},
		&model.Player{},
		&model.Tag{},
		&model.Team{},
		&model.User{},
	}

	for _, table := range tables {
		if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(table).Error; err != nil {
			return err
		}
	}

	// Re-enable foreign key checks
	db.Exec("SET session_replication_role = 'origin'")

	// Re-seed the database
	return SeedAll(db)
}

// resetWithOrderedDeletes deletes data in the correct order when foreign key
// constraint disabling is not available
func resetWithOrderedDeletes(db *gorm.DB) error {
	// Delete in order of foreign key dependencies
	if err := db.Where("1=1").Delete(&model.CheckInTag{}).Error; err != nil {
		return err
	}
	if err := db.Where("1=1").Delete(&model.PlayerRating{}).Error; err != nil {
		return err
	}
	if err := db.Where("1=1").Delete(&model.CheckIn{}).Error; err != nil {
		return err
	}
	if err := db.Where("1=1").Delete(&model.Session{}).Error; err != nil {
		return err
	}
	if err := db.Where("1=1").Delete(&model.MatchPlayer{}).Error; err != nil {
		return err
	}
	if err := db.Where("1=1").Delete(&model.Match{}).Error; err != nil {
		return err
	}
	if err := db.Where("1=1").Delete(&model.Player{}).Error; err != nil {
		return err
	}
	if err := db.Where("1=1").Delete(&model.Tag{}).Error; err != nil {
		return err
	}
	if err := db.Where("1=1").Delete(&model.Team{}).Error; err != nil {
		return err
	}
	if err := db.Where("1=1").Delete(&model.User{}).Error; err != nil {
		return err
	}

	// Re-seed the database
	return SeedAll(db)
}
