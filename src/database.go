package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
	"strconv"
)

type message struct {
	Username string
	MessageId int64
	Message string
}

var sess *session.Session
var service *dynamodb.DynamoDB

/*
	Creates singleton Session from config file
	Creates the session once and returns it
 */
func GetSession() *session.Session {
	if sess != nil {
		return sess
	}
	var err error
	sess, err = session.NewSessionWithOptions(session.Options{Profile: "personal", SharedConfigState: session.SharedConfigEnable})
	if err != nil {
		log.Fatalf("Creating session failed: %v", err)
	}
	return sess
}

/*
	Creates singleton client for DynamoDB from session
	Creates the client once and returns it
 */
func GetClient() *dynamodb.DynamoDB {
	if service != nil {
		return service
	}
	service = dynamodb.New(GetSession())
	return service
}

/*
	Called by the http servers onGetMessage
	Given a timestamp in milliseconds, check for messages that come after given timestamp.
	Returns all results or empty if no results
 */
func GetMessagesAfter(timestamp int64) []message {
	clientDB := GetClient()

	input := dynamodb.ScanInput{
		TableName: aws.String("message"),
		Select: aws.String("ALL_ATTRIBUTES"),
		FilterExpression: aws.String("MessageId >= :ts"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue {
			":ts": {N: aws.String(strconv.FormatInt(timestamp, 10))},
		},
	}

	output, err := clientDB.Scan(&input)
	if err != nil {
		log.Fatalf("Failed to scan data: %v", err)
	}

	var mess message
	result := []message{}
	for _, item := range output.Items {
		err = dynamodbattribute.UnmarshalMap(item, &mess)
		if err != nil {
			log.Fatalf("Unmarshalling failed %v", err)
		}

		result = append(result, mess)
	}
	return result
}

/*
	Called by the http onGetMessage
	Returns up to amount of recent messages
 */
func GetRecentMessages(amount int64) []message {
	recentMessages := []message{}
	clientDB := GetClient()
	var mess message

	input := dynamodb.ScanInput{
		TableName: aws.String("message"),
		Select: aws.String("ALL_ATTRIBUTES"),
		Limit: aws.Int64(amount),
	}

	output, err := clientDB.Scan(&input)
	if err != nil {
		log.Fatalf("Failed to scan data: %v", err)
	}

	for _, item := range output.Items {
		err = dynamodbattribute.UnmarshalMap(item, &mess)
		if err != nil {
			log.Fatalf("Failed to unmarshal: %v", err)
		}
		recentMessages = append(recentMessages, mess)
	}
	return recentMessages
}

/*
	Called by the http onPostMessage
	Given a message inserts message into the database
 */
func PutMessage(mess message) {
	clientDB := GetClient()

	item, err := dynamodbattribute.MarshalMap(mess)
	if err != nil {
		log.Fatalf("Marshal failed %v", err)
	}

	input := dynamodb.PutItemInput{
		Item: item,
		TableName: aws.String("message"),
	}
	_, err = clientDB.PutItem(&input)
	if err != nil {
		log.Fatalf("putitem failed %v", err)
	}
}
