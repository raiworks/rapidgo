package migrations

import "gorm.io/gorm"

func init() {
	Register(Migration{
		Version: "20260307000001_create_jobs_tables",
		Up: func(db *gorm.DB) error {
			// Create jobs table.
			if err := db.Exec(`CREATE TABLE IF NOT EXISTS jobs (
				id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
				queue       VARCHAR(255) NOT NULL,
				type        VARCHAR(255) NOT NULL,
				payload     TEXT NOT NULL,
				attempts    INT UNSIGNED NOT NULL DEFAULT 0,
				max_attempts INT UNSIGNED NOT NULL DEFAULT 3,
				available_at DATETIME NOT NULL,
				reserved_at  DATETIME NULL,
				created_at   DATETIME NOT NULL,
				INDEX idx_jobs_queue (queue),
				INDEX idx_jobs_available_at (available_at),
				INDEX idx_jobs_reserved_at (reserved_at)
			)`).Error; err != nil {
				return err
			}

			// Create failed_jobs table.
			return db.Exec(`CREATE TABLE IF NOT EXISTS failed_jobs (
				id        BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
				queue     VARCHAR(255) NOT NULL,
				type      VARCHAR(255) NOT NULL,
				payload   TEXT NOT NULL,
				error     TEXT NOT NULL,
				failed_at DATETIME NOT NULL
			)`).Error
		},
		Down: func(db *gorm.DB) error {
			if err := db.Exec(`DROP TABLE IF EXISTS failed_jobs`).Error; err != nil {
				return err
			}
			return db.Exec(`DROP TABLE IF EXISTS jobs`).Error
		},
	})
}
