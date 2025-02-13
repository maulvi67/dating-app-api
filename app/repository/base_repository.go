package repository

import (
	db "dating-apps/helper/database"

	"gorm.io/gorm"
)

type baseRepository struct {
	db *gorm.DB
}

type BaseRepository interface {
	GetDB() *gorm.DB
}

func NewBaseRepository(db *db.Database) BaseRepository {
	g, ok := (*db).Client().(*gorm.DB)
	if !ok {
		return &baseRepository{}
	}
	return &baseRepository{
		db: g,
	}
}

func (br *baseRepository) GetDB() *gorm.DB {
	return br.db
}
