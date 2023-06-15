package predict

type Symptoms struct {
	ID        uint   `gorm:"primarykey"`
	SymptomEN string `gorm:"column:symptom_en"`
	SymptomID string `gorm:"column:symptom_id"`
}

func (Symptoms) TableName() string {
	return `symptoms_mapping`
}

type Disease struct {
	ID          uint   `gorm:"primarykey"`
	DiseaseEN   string `gorm:"column:disease_en"`
	DiseaseID   string `gorm:"column:disease_id"`
	Description string `gorm:"column:description"`
	Precaution1 string `gorm:"column:precaution_1"`
	Precaution2 string `gorm:"column:precaution_2"`
	Precaution3 string `gorm:"column:precaution_3"`
	Precaution4 string `gorm:"column:precaution_4"`
}

func (Disease) TableName() string {
	return `diseases_mapping`
}
