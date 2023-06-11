package predict

import (
	"backend/src/constant"
	httpPredict "backend/src/entity/v1/http/predict"
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

type Servicer interface {
	SubmitData(req httpPredict.PredictSymptoms) (err error)
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
