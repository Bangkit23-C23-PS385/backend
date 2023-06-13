package predict

type Symptoms struct {
	ID        uint   `gorm:"primarykey"`
	SymptomEN string `gorm:"column:symptom_en"`
	SymptomID string `gorm:"column:symptom_id"`
}

func (Symptoms) TableName() string {
	return `symptoms_mapping`
}
