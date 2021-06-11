package main

import (
	"context"
	"encoding/json"
	"net/http"
	"simple-apigateway-lambda/pkg/store"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type lambdaHandler struct {
	store *store.Store
}

const tableName = "DynamoSample"

func (h *lambdaHandler) Handle(_ context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var obj store.Object
	err := json.Unmarshal([]byte(req.Body), &obj)
	if err != nil {
		resp := events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       string("error"),
		}

		return resp, nil
	}

	obj.Path = req.Path

	err = h.store.PutObject(obj)
	if err != nil {
		resp := events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       string("error inserting"),
		}

		return resp, nil
	}

	b, _ := json.Marshal(obj)
	resp := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(b),
	}

	return resp, nil
}

func main() {
	sess := session.Must(session.NewSession())

	handler := lambdaHandler{
		store: store.NewStore(tableName, dynamodb.New(sess)),
	}

	lambda.Start(handler.Handle)
}
