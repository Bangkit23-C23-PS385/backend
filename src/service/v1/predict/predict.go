package predict

import (
	"backend/src/constant"
	httpPredict "backend/src/entity/v1/http/predict"
	predictRepo "backend/src/repository/v1/predict"
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
)

type Service struct {
	repo predictRepo.Repositorier
}

func NewService(repo predictRepo.Repositorier) *Service {
	return &Service{
		repo: repo,
	}
}

type Servicer interface {
	GetSymptoms() (resps httpPredict.SymptomsResponse, err error)
	SubmitData(req httpPredict.PredictSymptoms) (resp httpPredict.DiseaseResponse, err error)
}

func (svc Service) GetSymptoms() (resps httpPredict.SymptomsResponse, err error) {
	entities, err := svc.repo.GetSymptoms()
	if err != nil {
		err = errors.Wrap(err, "repo: get symptoms")
		return
	}

	resps = httpPredict.SymptomsResponse{}
	for _, entity := range entities {
		symptom := httpPredict.Symptom{
			ID:        int(entity.ID),
			SymptomEN: entity.SymptomEN,
			SymptomID: entity.SymptomID,
		}

		resps.Symptoms = append(resps.Symptoms, symptom)
	}

	return
}

func (svc Service) SubmitData(req httpPredict.PredictSymptoms) (resp httpPredict.DiseaseResponse, err error) {
	data, _ := json.Marshal(req)
	predictedDisease := struct {
		Disease string `json:"disease"`
	}{}

	cctx, cancel := context.WithCancel(context.Background())
	client, err := pubsub.NewClient(cctx, constant.GCPProjectID)
	if err != nil {
		err = errors.Wrap(err, "error creating pubsub client")
		return
	}

	topic := client.Topic(constant.GCPTopicSubmitData)
	result := topic.Publish(cctx, &pubsub.Message{
		Data: data,
	})
	id, err := result.Get(cctx)
	if err != nil {
		err = errors.Wrap(err, "pubsub: failed to get result")
		return
	}
	fmt.Printf("Published message with ID: %s\n", id)
	fmt.Println("Waiting response from model...")

	subName := constant.GCPSubscriptionPredict
	subscription := client.Subscription(subName)
	_ = subscription.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		msg.Ack()
		fmt.Printf("Received message: %s\n", string(msg.Data))
		_ = json.Unmarshal(msg.Data, &predictedDisease)
		cancel()
	})

	diseaseEntity, err := svc.repo.GetDiseaseDetail(predictedDisease.Disease)
	if err != nil {
		err = errors.Wrap(err, "repo: get disease detail")
		return
	}
	resp = httpPredict.DiseaseResponse{
		Disease:     diseaseEntity.DiseaseID,
		Description: diseaseEntity.Description,
		Precaution1: diseaseEntity.Precaution1,
		Precaution2: diseaseEntity.Precaution2,
		Precaution3: diseaseEntity.Precaution3,
		Precaution4: diseaseEntity.Precaution4,
	}

	return
}
