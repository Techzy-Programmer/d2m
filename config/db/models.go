package db

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var dbi *gorm.DB // Database instance
var migTables = []interface{}{&Config[any]{}, &Deployment{}, &DeploymentLog{}}

func init() {
	configPath, err := getUserConfigPath("d2m")
	if err != nil {
		log.Fatalf("Failed to get config path: %v", err)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := os.MkdirAll(configPath, 0755); err != nil {
			log.Fatalf("Failed to create config directory: %v", err)
		}
	}

	var dbErr error
	dbFile := filepath.Join(configPath, "d2m.db")
	dbi, dbErr = gorm.Open(sqlite.Open(dbFile), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if dbErr != nil {
		log.Fatalf("Failed to open database: %v", dbErr)
	}

	dbi.AutoMigrate(migTables...)
}

type Config[T any] struct {
	Key      string `gorm:"primaryKey"`
	Value    T      `gorm:"-"`            // Ignore this field for GORM's migrations; will handle manually
	RawValue []byte `gorm:"column:value"` // Store the JSON-encoded value as a byte slice
}

func (Config[T]) TableName() string {
	// It is required to define table name here since we are using a generic type
	return "configs" // Define static table name
}

// Custom serialization and deserialization methods for the `Value` field since it is a generic type and cannot be handled by GORM directly

func (c *Config[T]) BeforeSave(tx *gorm.DB) (err error) {
	c.RawValue, err = json.Marshal(c.Value)
	return
}

func (c *Config[T]) AfterFind(tx *gorm.DB) (err error) {
	return json.Unmarshal(c.RawValue, &c.Value)
}

type Deployment struct {
	ID         uint `gorm:"primaryKey"`
	Branch     string
	StartAt    time.Time
	EndAt      time.Time
	CommitHash string
	CommitMsg  string
	Repo       string
	Status     string
	Logs       []DeploymentLog `gorm:"foreignKey:DeployID"`
}

type DeploymentLog struct {
	ID        uint `gorm:"primaryKey"`
	Level     uint
	Title     string
	Message   string
	Steps     string
	Timestamp int64
	DeployID  uint
}
