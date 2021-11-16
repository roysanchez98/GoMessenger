package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

/*
	Checks if get or post request is received and runs respective method
	if method fails prints error
	If no get or post method outputs error
 */
func handleMessages(w http.ResponseWriter, r *http.Request )  {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	// Check for type of post or get request
	// Choose format of request
	if r.Method == http.MethodGet {
		log.Printf("Entered get request: %#v\n", r)
		err := onGetMessage(w,r) // TODO: RETURNS ERROR
		if err != nil {
			log.Printf("Calling get message failed %v\n", err)
		}
	} else if r.Method == http.MethodPost {
		log.Printf("Entered post request: %#v\n", r)
		err := onPostMessage(w, r)
		if err != nil {
			log.Printf("Calling post message failed %v\n", err)
		}
	} else { // TODO: send more substantial error
		log.Fatalf("Failed to handle message %#v", r)
	}
}

/*
	Handle get request for /messages/messages
	Expects form data for amount or timestamp
	Queries database and retrieves result
	Otherwise, outputs error
*/
func onGetMessage(w http.ResponseWriter, r *http.Request) error {
	// request amount and Timestamp from client header
	err := r.ParseForm()
	if err != nil {
		return err
	}
	amount := r.Form["amount"]
	timestamp := r.Form["since"]
	var messagesToSend []message

	// Use requested data on Database get method to retrieve Message, otherwise return error
	if amount != nil {
		intAmount, err := strconv.ParseInt(amount[0], 10,32)
		if err == nil {
			messagesToSend = GetRecentMessages(intAmount)
		} else {
			return err
		}
	} else if timestamp != nil {
		timeLong, err := strconv.ParseInt(timestamp[0], 10, 64)
		if err == nil {
			messagesToSend = GetMessagesAfter(timeLong)
		} else {
			return err
		}
	} else {
		return fmt.Errorf("get request invalid")
	}
	// Convert Message to JSON
	jsonText, err := json.Marshal(messagesToSend)
	if err == nil {
		_, err := w.Write(jsonText)
		return err
	}
	return err
}

/*
	Post request used to send data to database /messages/messages
	Expects form data username and message
	Returns error otherwise
*/
func onPostMessage(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseMultipartForm(1024)
	if err != nil {
		return err
	}
	username := r.Form["Username"] // TODO: check for username
	mess := r.Form["Message"]

	messageId := time.Now().UnixMilli()

	if len(username) == 0 || len(mess) == 0 {
		return errors.New("error retrieving data")
	}

	postMessage := message{
		MessageId: messageId,
		Username: username[0],
		Message: mess[0],
	}

	PutMessage(postMessage)

	// Convert Message to JSON
	jsonText, err := json.Marshal(postMessage)
	if err == nil {
		_, err := w.Write(jsonText)
		println(jsonText)
		return err
	}
	return err
}

/*
	converts a message to JSON
 */
func (messageToConvert *message) toJSON(w http.ResponseWriter) error {
	encoder := json.NewEncoder(w)
	err := encoder.Encode(messageToConvert)
	return err
}




