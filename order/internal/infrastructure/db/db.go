package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func createDSN(c *Config) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host,
		c.Port,
		c.Username,
		c.Password,
		c.Database,
	)
}

func NewDB(config *Config) (*gorm.DB, error) {
	dsn := createDSN(config)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
