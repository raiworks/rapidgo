package migrations

import "gorm.io/gorm"

func init() {
	Register(Migration{
		Version: "20260308000001_add_soft_deletes",
		Up: func(db *gorm.DB) error {
			type User struct {
				DeletedAt gorm.DeletedAt `gorm:"index"`
			}
			type Post struct {
				DeletedAt gorm.DeletedAt `gorm:"index"`
			}
			if err := db.Migrator().AddColumn(&User{}, "DeletedAt"); err != nil {
				return err
			}
			if err := db.Migrator().CreateIndex(&User{}, "DeletedAt"); err != nil {
				return err
			}
			if err := db.Migrator().AddColumn(&Post{}, "DeletedAt"); err != nil {
				return err
			}
			return db.Migrator().CreateIndex(&Post{}, "DeletedAt")
		},
		Down: func(db *gorm.DB) error {
			type User struct {
				DeletedAt gorm.DeletedAt
			}
			type Post struct {
				DeletedAt gorm.DeletedAt
			}
			if err := db.Migrator().DropColumn(&Post{}, "DeletedAt"); err != nil {
				return err
			}
			return db.Migrator().DropColumn(&User{}, "DeletedAt")
		},
	})
}
