package predict

import (
	db "backend/src/database"
	dbPredict "backend/src/entity/v1/db/predict"

	"gorm.io/gorm"
)

type Repository struct {
	master *gorm.DB
}

func NewRepository(db db.DB) *Repository {
	return &Repository{
		master: db.Master,
	}
}

type Repositorier interface {
	GetSymptoms() (entities []dbPredict.Symptoms, err error)
}

func (repo Repository) GetSymptoms() (entities []dbPredict.Symptoms, err error) {
	query := repo.master.Find(&entities)
	err = query.Error

	return
}
