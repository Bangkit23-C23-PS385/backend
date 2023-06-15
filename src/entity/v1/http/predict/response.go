package predict

type SymptomsResponse struct {
	Symptoms []Symptom `json:"symptoms"`
}

type Symptom struct {
	ID        int    `json:"id"`
	SymptomEN string `json:"symptom_en"`
	SymptomID string `json:"symptom_id"`
}

type DiseaseResponse struct {
	Disease     string `json:"disease"`
	Description string `json:"description"`
	Precaution1 string `json:"precaution_1"`
	Precaution2 string `json:"precaution_2"`
	Precaution3 string `json:"precaution_3"`
	Precaution4 string `json:"precaution_4"`
}
