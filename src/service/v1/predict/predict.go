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
	SubmitData(req httpPredict.PredictSymptoms) (err error)
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

func (svc Service) SubmitData(req httpPredict.PredictSymptoms) (err error) {
	data, _ := json.Marshal(req)

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, constant.GCPProjectID)
	if err != nil {
		err = errors.Wrap(err, "error creating pubsub client")
		return
	}

	topic := client.Topic(constant.GCPTopicSubmitData)
	result := topic.Publish(ctx, &pubsub.Message{
		Data: data,
	})
	id, err := result.Get(ctx)
	if err != nil {
		err = errors.Wrap(err, "pubsub: failed to get result")
		return
	}
	fmt.Printf("Published message with ID: %s\n", id)

	return
}
