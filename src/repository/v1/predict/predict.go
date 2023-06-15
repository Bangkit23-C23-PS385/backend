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
	GetDiseaseDetail(diseaseEN string) (entity dbPredict.Disease, err error)
}

func (repo Repository) GetSymptoms() (entities []dbPredict.Symptoms, err error) {
	query := repo.master.Find(&entities)
	err = query.Error

	return
}

func (repo Repository) GetDiseaseDetail(diseaseEN string) (entity dbPredict.Disease, err error) {
	query := repo.master.Where("disease_en = ?", diseaseEN).Take(&entity)
	err = query.Error

	return
}
