package db

import (
	"github.com/aliatli/reactor/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

func NewDatabase() (*Database, error) {
	db, err := gorm.Open(sqlite.Open("reactor.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&models.State{})
	if err != nil {
		return nil, err
	}

	return &Database{db}, nil
}

func (db *Database) SaveState(state *models.State) error {
	// First try to find the state, including soft deleted ones
	var existingState models.State
	result := db.Unscoped().Where("name = ?", state.Name).First(&existingState)

	if result.Error == nil {
		// If found (even if deleted), hard delete it first
		if err := db.Unscoped().Where("name = ?", state.Name).Delete(&models.State{}).Error; err != nil {
			return err
		}
	}

	// Create new state
	return db.Create(state).Error
}

func (db *Database) GetAllStates() ([]models.State, error) {
	var states []models.State
	err := db.Find(&states).Error
	return states, err
}

func (db *Database) DeleteState(name string) error {
	// Hard delete the state
	return db.Unscoped().Where("name = ?", name).Delete(&models.State{}).Error
}
