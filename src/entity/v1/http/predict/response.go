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
	Disease string `json:"disease"`
}
