package seed

import (
	"gorm.io/gorm"
)

// SeedAll runs all seed functions in the correct order
func SeedAll(db *gorm.DB) error {
	seedFuncs := []func(*gorm.DB) error{
		SeedTeams,
		SeedPlayers,
		SeedMatches,
		SeedMatchPlayers,
		SeedTags,
		SeedDevUsers,
	}

	for _, seedFunc := range seedFuncs {
		if err := seedFunc(db); err != nil {
			return err
		}
	}

	return nil
}

// SeedBaseData runs only the foundational seed data (teams, players, matches, match-players, tags)
func SeedBaseData(db *gorm.DB) error {
	baseSeedFuncs := []func(*gorm.DB) error{
		SeedTeams,
		SeedPlayers,
		SeedMatches,
		SeedMatchPlayers,
		SeedTags,
	}

	for _, seedFunc := range baseSeedFuncs {
		if err := seedFunc(db); err != nil {
			return err
		}
	}

	return nil
}