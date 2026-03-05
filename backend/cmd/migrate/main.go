package main

import (
	"partitionlab/internal/app/ds"
	"partitionlab/internal/app/dsn"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(
		&ds.User{},
		&ds.Partition{},
		&ds.Calculation{},
		&ds.CalculationItem{},
	)
	if err != nil {
		panic("cant migrate db")
	}

	// Remove duplicates by title (keep the lowest id)
	db.Exec(`WITH d AS (
		SELECT id, title,
			   ROW_NUMBER() OVER (PARTITION BY title ORDER BY id) AS rn
		FROM partitions
	)
	DELETE FROM partitions p
	USING d
	WHERE p.id = d.id AND d.rn > 1`)

	// Ensure unique titles for partitions to avoid duplicates on reseed
	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS ux_partitions_title ON partitions (title)")
}
