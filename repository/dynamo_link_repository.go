package repository

import (
	"shorters/domain"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

type dynamoLinkRepository struct {
	tableName string
	client    *dynamodb.DynamoDB
}

func newDynamoClient() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	return dynamodb.New(sess)
}

func NewDynamoLinkRepository() LinkRepository {
	return &dynamoLinkRepository{tableName: "Shorters", client: newDynamoClient()}
}

func (r *dynamoLinkRepository) Find(key string) (*domain.Link, error) {
	p := expression.NamesList(expression.Name("Key"), expression.Name("URL"))
	expr, _ := expression.NewBuilder().WithProjection(p).Build()
	res, err := r.client.GetItem(&dynamodb.GetItemInput{
		ExpressionAttributeNames: expr.Names(),
		Key:                      map[string]*dynamodb.AttributeValue{"Key": {S: &key}},
		ProjectionExpression:     expr.Projection(),
		TableName:                aws.String(r.tableName),
	})
	if err != nil {
		return nil, err
	}
	if len(res.Item) == 0 {
		return nil, domain.LinkNotFoundError{key}
	}
	var link domain.Link
	_ = dynamodbattribute.UnmarshalMap(res.Item, &link)
	return &link, nil
}

func (r *dynamoLinkRepository) FindByUser(email string) ([]*domain.Link, error) {
	f := expression.Name("Creator").Equal(expression.Value(email))
	p := expression.NamesList(expression.Name("Key"), expression.Name("URL"), expression.Name("Visits"))
	expr, _ := expression.NewBuilder().WithFilter(f).WithProjection(p).Build()
	res, err := r.client.Scan(&dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(r.tableName),
	})
	if err != nil {
		return nil, err
	}
	links := make([]*domain.Link, len(res.Items))
	for i, item := range res.Items {
		var link domain.Link
		_ = dynamodbattribute.UnmarshalMap(item, &link)
		links[i] = &link
	}
	return links, nil
}

func (r *dynamoLinkRepository) Store(link *domain.Link) error {
	item, _ := dynamodbattribute.MarshalMap(link)
	_, err := r.client.PutItem(&dynamodb.PutItemInput{
		Item:         item,
		ReturnValues: aws.String(dynamodb.ReturnValueNone),
		TableName:    aws.String(r.tableName),
	})
	return err
}

func (r *dynamoLinkRepository) AddVisits(key string) error {
	u := expression.Add(expression.Name("Visits"), expression.Value(1))
	expr, _ := expression.NewBuilder().WithUpdate(u).Build()
	_, err := r.client.UpdateItem(&dynamodb.UpdateItemInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Key:                       map[string]*dynamodb.AttributeValue{"Key": {S: &key}},
		ReturnValues:              aws.String(dynamodb.ReturnValueNone),
		TableName:                 aws.String(r.tableName),
		UpdateExpression:          expr.Update(),
	})
	return err
}
