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
	// "go.mongodb.org/mongo-driver/mongo"
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
	ret = "Cmd echo: " + cmd.CommandText
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
/* func WriteToDB(){
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:root@localhost:27017"))
	if (err != nil) { return err }
	collection := client.Database("indigo").Collection("main")
	res, err := collection.InsertOne(context.Background(), bson.M{"entity": "me"})
	if err != nil { return err }
	id := res.InsertedID
} */