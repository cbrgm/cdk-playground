package store

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/rs/zerolog/log"
)

type Store struct {
	tableName    string
	dynamoClient dynamodbiface.DynamoDBAPI
}

type dynamoObject struct {
	Pk      string `dynamodbav:"pk"`
	Message string `dynamodbav:"message"`
}

func newDynamoObject(obj Object) dynamoObject {
	return dynamoObject{
		Pk:      obj.Path,
		Message: obj.Message,
	}
}

type Object struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

func NewStore(tablename string, client dynamodbiface.DynamoDBAPI) *Store {
	return &Store{
		tableName:    tablename,
		dynamoClient: client,
	}
}

func (s *Store) PutObject(msg Object) error {
	keys, err := dynamodbattribute.MarshalMap(newDynamoObject(msg))
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: &s.tableName,
		Item:      keys,
	}

	_, err = s.dynamoClient.PutItem(input)
	if err != nil {
		log.Fatal().Err(err).Msg("error while updating item")
	}

	log.Info().Msgf("updated: %v", keys)

	return nil
}

func (s *Store) GetObject(path string) (Object, error) {
	log.Info().Str("host", path).Msg("looking up target")

	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(s.tableName),
		KeyConditions: map[string]*dynamodb.Condition{
			"pk": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(path),
					},
				},
			},
		},
	}

	result, err := s.dynamoClient.Query(queryInput)
	if err != nil {
		log.Error().Err(err).Msg("failed to query")
		return Object{}, err
	}

	res := marshalObject(result.Items)
	if len(res) == 0 {
		return Object{}, err
	}

	return res[0], err
}

func marshalObject(items []map[string]*dynamodb.AttributeValue) []Object {
	objects := []dynamoObject{}

	err := dynamodbattribute.UnmarshalListOfMaps(items, &objects)
	if err != nil {
		log.Fatal().Err(err)
		return []Object{}
	}

	transformed := make([]Object, len(objects))
	for i, dynamoObj := range objects {
		transformed[i] = Object{
			Path:    dynamoObj.Pk,
			Message: dynamoObj.Message,
		}
	}

	return transformed
}
