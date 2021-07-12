// main.go
// Entry point for the server

// Reference:
// https://thenewstack.io/make-a-restful-json-api-go/

package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"time"
	"encoding/json"
	"io"
	"io/ioutil"
	"strconv"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
)

// Models

type Command struct {
	CommandText string `json:"command"`
	Flags string `json:"flags"`
}

type Response struct {
	ResponseText string `json:"response"`
	ImageFileName string `json:"imagefile"`
}

type Job struct {
	JobID int
	XWorkerID int
	DateAssigned time.Time
	DateUpdated time.Time
	DateScheduled time.Time
	IsRecurring bool
	NextScheduledRun time.Time
	DateCompleted time.Time
	IsPast bool
	IsInProgress bool
	PercentDone int
}

type Jobs []Job

// Routes (Views)

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf("%s\t%s\t%s\t%s", r.Method, r.RequestURI, name, time.Since(start))
	})
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Root)

	router.HandleFunc("/command", ProcessCommand)

	// var handler http.Handler
	// handler = ProcessCommandsWrapper(http.HandlerFunc(ProcessCommands))
	// handler = Logger(handler, "ProcessCommands")

	router.Methods("POST", "OPTIONS").Path("/commands").Name("ProcessCommands").Handler(http.HandlerFunc(ProcessCommands))
	return router
}

func main(){
	router := NewRouter()
	fmt.Println("Indigo API piapi started on 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// Controllers

func Root(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "Project Indigo API Route - Success")
}

func ProcessCommand(w http.ResponseWriter, r *http.Request){
	// 	vars := mux.Vars(r)
	//  command := vars["command"]
	resp := Response{ResponseText:"Hi from API"}
	json.NewEncoder(w).Encode(resp)
}

func ProcessCommandsWrapper(h http.Handler) {}

func ProcessCommands (w http.ResponseWriter, r *http.Request){
	var cmd Command 
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {
		panic(err)
	}

	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &cmd); err != nil {
		// Note - OPTIONS route can't handle this error case
		/* w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)

		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		} */
	}

	// Here is where we can decide the response based on the request.
	var re Response
	var ret string

	// Decide what to do based on command
	if cmd.CommandText == "c" {
		ret = "Added: " + WriteToDB()
	} else if cmd.CommandText == "r" {
		ret = "There Are: " + strconv.FormatInt(ReadFromDB(), 10)
	} else {
		ret = "Cmd echo: " + cmd.CommandText
	}

	re = Response{ResponseText:ret,ImageFileName:"2.png"} 

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
    w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	// Set to status success
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(re); err != nil {
		panic(err)
	}
}

// Handles writing to database
// Todo DEFER a call to disconnect the client https://pkg.go.dev/go.mongodb.org/mongo-driver#readme-installation
func WriteToDB() string {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	// Todo auth root:root@
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	// Todo figure out how to convert err to string
	if (err != nil) { return "" }
	collection := client.Database("indigo").Collection("main")
	res, err := collection.InsertOne(context.Background(), bson.M{"entity": "me"})
	if (err != nil) { return "" }
	id := res.InsertedID
	if (id != nil) { 
		return "Success" 
	} else { 
		return "Failure" 
	}
}

// Handles reading from database
// Todo DEFER a call to disconnect the client https://pkg.go.dev/go.mongodb.org/mongo-driver#readme-installation
func ReadFromDB() int64 {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	// Todo auth root:root@
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if (err != nil) { return -1 }
	collection := client.Database("indigo").Collection("main")
	opts := options.Count().SetMaxTime(2 * time.Second)
	count, err := collection.CountDocuments(context.Background(), bson.D{{"entity", "me"}}, opts)
	if err != nil {
		log.Fatal(err)
		return -1
	}
	return count
	/* cur, err := collection.Find(context.Background(), bson.D{})
	if err != nil { log.Fatal(err) }
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		raw := cur.Current 
		return raw
	}
	if err := cur.Err(); err != nil {
		return err
	} */
}