package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func InitDB(url string) *dynamodb.DynamoDB {
	if url == "" {
		url = "http://localhost:8000"
	}

	config := &aws.Config{
		Region:   aws.String("us-west-2"),
		Endpoint: aws.String(url),
	}

	return dynamodb.New(session.Must(session.NewSession(config)))
}

func CreateStructure(db *dynamodb.DynamoDB) {
	const table = "Keywords"
	const index = "KeywordsIndex"

	if db == nil {
		panic("DB has not been initialized with InitDB")
	}

	db.DeleteTable(&dynamodb.DeleteTableInput{TableName: aws.String(table)})

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("email"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("keyword"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("email"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("keyword"),
				KeyType:       aws.String("RANGE"),
			},
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String(index),
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("keyword"),
						KeyType:       aws.String("HASH"),
					},
					{
						AttributeName: aws.String("email"),
						KeyType:       aws.String("RANGE"),
					},
				},
				Projection: &dynamodb.Projection{
					NonKeyAttributes: nil,
					ProjectionType:   aws.String(dynamodb.ProjectionTypeKeysOnly),
				},
				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(1),
					WriteCapacityUnits: aws.Int64(1),
				},
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		TableName: aws.String(table),
	}

	if _, err := db.CreateTable(input); err != nil {
		panic(err)
	}
}
