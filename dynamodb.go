package eventstrangler

import (
	"github.com/aws/aws-sdk-go/aws"
	awserr2 "github.com/aws/aws-sdk-go/aws/awserr"
	session2 "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DynamoDBStore struct {
	tableName string
	db        *dynamodb.DynamoDB
}

func (s *DynamoDBStore) Exists(hashKey string) (bool, error) {
	_, err := s.Get(hashKey)
	if err != nil {
		if _, ok := err.(StoreEventNotFoundError); ok {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *DynamoDBStore) Get(hashKey string) (*Record, error) {
	result, err := s.db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(s.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"HashKey": {
				S: aws.String(hashKey),
			},
		},
	})
	if err != nil {
		if awsErr, ok := err.(awserr2.Error); ok && awsErr.Code() == dynamodb.ErrCodeResourceNotFoundException {
			return nil, StoreEventNotFoundError{
				HashKey: hashKey,
			}
		}
		return nil, err
	}

	record := Record{}
	if err = dynamodbattribute.UnmarshalMap(result.Item, &record); err != nil {
		return nil, err
	}
	return &record, nil
}

func (s *DynamoDBStore) Put(_ string, record *Record, _ int) error {
	item, err := dynamodbattribute.MarshalMap(record)
	if err != nil {
		return err
	}

	if _, err = s.db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(s.tableName),
		Item:      item,
	}); err != nil {
		return err
	}
	return nil
}

func (s *DynamoDBStore) Delete(hashKey string) error {
	_, err := s.db.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(s.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"HashKey": {
				S: aws.String(hashKey),
			},
		},
	})

	if err != nil {
		if awsErr, ok := err.(awserr2.Error); ok && awsErr.Code() == dynamodb.ErrCodeResourceNotFoundException {
			return StoreEventNotFoundError{
				HashKey: hashKey,
			}
		}
		return err
	}
	return nil
}

func (s *DynamoDBStore) Close() error {
	return nil
}

func NewDynamoDBStore(tableName string, profile string) (*DynamoDBStore, error) {
	session, err := session2.NewSessionWithOptions(session2.Options{
		SharedConfigState: session2.SharedConfigEnable,
		Profile:           profile,
	})
	if err != nil {
		return nil, err
	}
	db := dynamodb.New(session)

	return &DynamoDBStore{
		tableName: tableName,
		db:        db,
	}, nil
}
