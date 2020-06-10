package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	
	"github.com/gorilla/mux"

    "fmt"
    "os"
	"strconv"
	"net/http"
	"io/ioutil"
	"log"
	"encoding/json"
)

var sess = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))

var svc = dynamodb.New(sess)

var tableName = "Movies"

type movie struct {
	Year int `json:"Year"`
	Title string `json:"Title"`
}

// AddMovie adds a dynamodb item into the backend db
func AddMovie(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		return
	}

	// Read request body.
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Body read error, %v", err)
		w.WriteHeader(500) // Return 500 Internal Server Error.
		return
	}

	var movie movie

	if err = json.Unmarshal(body, &movie); err != nil {
		log.Printf("Body parse error, %v", err)
		w.WriteHeader(400) // Return 400 Bad Request.
		return
	}

	av, err := dynamodbattribute.MarshalMap(movie)
	if err != nil {
   		log.Println("Got error marshalling new movie item:")
		log.Println(err.Error())
    	os.Exit(1)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
    	log.Println("Got error calling PutItem:")
    	log.Println(err.Error())
    	os.Exit(1)
	}

	year := strconv.Itoa(movie.Year)

	fmt.Println("Successfully added '" + movie.Title + "' (" + year + ") to table " + tableName)
}

func main() {
	r := mux.NewRouter()
    r.HandleFunc("/movie", AddMovie).Methods("POST")
    http.ListenAndServe(":8080", r)
}