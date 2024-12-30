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

	// Drop existing tables and recreate them
	if err := db.Migrator().DropTable(&models.State{}); err != nil {
		return nil, err
	}

	// Auto migrate the schema
	if err := db.AutoMigrate(&models.State{}); err != nil {
		return nil, err
	}

	return &Database{db}, nil
}

func (db *Database) SaveState(state *models.State) error {
	// Try to find existing state
	var existingState models.State
	result := db.Where("name = ?", state.Name).First(&existingState)

	if result.Error == nil {
		// State exists, update it including edges
		state.ID = existingState.ID // Keep the same ID
		return db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(state).Error
	} else if result.Error == gorm.ErrRecordNotFound {
		// State doesn't exist, create new
		return db.Create(state).Error
	}

	return result.Error
}

func (db *Database) GetAllStates() ([]models.State, error) {
	var states []models.State
	err := db.Find(&states).Error
	return states, err
}

func (db *Database) DeleteState(name string) error {
	return db.Unscoped().Where("name = ?", name).Delete(&models.State{}).Error
}
