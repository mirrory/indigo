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
	"os"
	"strconv"
	"context"
	"strings"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
)

// Models

type Command struct {
	CommandText string `json:"command"`
	Flags string `json:"flags"`
}

type Dialogue struct {
	Introduction string `json:"introduction"`
	SettlementNew string `json:"settlement-new"`
	BusinessNew string `json:"business-new"`
	ChemicalNew string `json:"chemical-new"`
	EducationLearn string `json:"education-learn"`
	BudgetSet string `json:"budget-set"`
	PersonNew string `json:"person-new"`
	PersonNewName string `json:"person-new-name"`
	PersonTalk string `json:"person-talk"`
	PersonEat string `json:"person-eat"`
	PersonCook string `json:"person-cook"`
	PersonAddToParty string `json:"person-add-to-party"`
	PersonExercise string `json:"person-exercise"`
	PersonSetLanguage string `json:"person-set-language"`
	LanguageMode string `json:"language-mode"`
	PersonEstablishRelationship string `json:"person-establish-relationship"`
	PersonRelationshipType string `json:"person-relationship-type"`
	PersonSetBirthday string `json:"person-set-birthday"`
	PersonSetBio string `json:"person-set-bio"`
	PersonSetAlignment string `json:"person-set-alignment"`
	LawNew string `json:"law-new"`
	TimeTravel string `json:"time-travel"`
	ReligionNew string `json:"religion-new"`
	AdventureNew string `json:"adventure-new"`
	SettingsConfirm string `json:"settings-confirm"`
	SpeciesNew string `json:"species-new"`
	M string `json:"m"`
	Help string `json:"help"`
	StatsRandom string `json:"stats-random"`
	ParticleNew string `json:"particle-new"`
	MusicControl string `json:"music-control"`
	PaintingNew string `json:"painting-new"`
	PlanetNew string `json:"planet-new"`
	ThemeChange string `json:"theme-change"`
	DiseaseNew string `json:"disease-new"`
	TechnologyNew string `json:"technology-new"`
	LandNew string `json:"land-new"`
	SaveFile string `json:"save-file"`
	PsychologyInspect string `json:"psychology-inspect"`
	LoadFile string `json:"load-file"`
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

	// TODO: Refactor to manage opening/closing this file A LOT in parallel
	jsonFile, err := os.Open("dialogue.json")
	if err != nil {
    	fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var dialogue Dialogue
	json.Unmarshal(byteValue, &dialogue)

	// Here is where we can decide the response based on the request.
	var re Response
	var ret string

	commandFlags := strings.Split(cmd.Flags, " ")
	flagMap := map[string]string{}

	if commandFlags != nil {
		for i := 0; i < len(commandFlags); i++ {
			if (i % 2 != 0) {
				flagName := commandFlags[i-1][1:len(commandFlags[i-1])]
				flagMap[flagName] = commandFlags[i]
			} 
		}
	}

	// Decide what to do based on command
	if cmd.CommandText == "a" {
		ret = dialogue.SettlementNew
	} else if cmd.CommandText == "b" {
		ret = dialogue.BusinessNew
	} else if cmd.CommandText == "c" {
		ret = dialogue.ChemicalNew + "Added: " + WriteToDB()
	} else if cmd.CommandText == "d" {
		ret = dialogue.EducationLearn
	} else if cmd.CommandText == "e" {
		ret = dialogue.BudgetSet
	} else if cmd.CommandText == "f" {
		// Person/people management command.
		if personName, ok := flagMap["n"]; ok {
			ret = fmt.Sprintf(dialogue.PersonNewName, personName)
		} else if personName, ok := flagMap["t"]; ok {
			ret = fmt.Sprintf(dialogue.PersonTalk, personName)
		} else if meal, ok := flagMap["f"]; ok {
			ret = fmt.Sprintf(dialogue.PersonEat, meal)
		} else if personName, ok := flagMap["o"]; ok {
			ret = fmt.Sprintf(dialogue.PersonCook, personName)
		} else if personName, ok := flagMap["a"]; ok {
			ret = fmt.Sprintf(dialogue.PersonAddToParty, personName)
		} else if personName, ok := flagMap["e"]; ok {
			ret = fmt.Sprintf(dialogue.PersonExercise, personName)
		} else if personName, ok := flagMap["l"]; ok {
			ret = fmt.Sprintf(dialogue.PersonSetLanguage, personName)
		} else if p, ok := flagMap["g"]; ok {
			ret = fmt.Sprintf(dialogue.LanguageMode, p)
		} else if personName, ok := flagMap["i"]; ok {
			ret = fmt.Sprintf(dialogue.PersonEstablishRelationship, personName)
		} else if personName, ok := flagMap["y"]; ok {
			ret = fmt.Sprintf(dialogue.PersonRelationshipType, personName)
		} else if personName, ok := flagMap["h"]; ok {
			ret = fmt.Sprintf(dialogue.PersonSetBirthday, personName)
		} else if personName, ok := flagMap["b"]; ok {
			ret = fmt.Sprintf(dialogue.PersonSetBio, personName)
		} else if personName, ok := flagMap["m"]; ok {
			ret = fmt.Sprintf(dialogue.PersonSetAlignment, personName)
		} else {
			ret = dialogue.PersonNew
		}
	} else if cmd.CommandText == "g" {
		ret = dialogue.LawNew
	} else if cmd.CommandText == "h" {
		ret = dialogue.TimeTravel
	} else if cmd.CommandText == "i" {
		ret = dialogue.ReligionNew
	} else if cmd.CommandText == "j" {
		ret = dialogue.AdventureNew
	} else if cmd.CommandText == "k" {
		ret = dialogue.SettingsConfirm
	} else if cmd.CommandText == "l" {
		ret = dialogue.SpeciesNew
	} else if cmd.CommandText == "m" {
		ret = dialogue.M
	} else if cmd.CommandText == "n" {
		ret = dialogue.Help
	} else if cmd.CommandText == "o" {
		ret = dialogue.StatsRandom + "There Are: " + strconv.FormatInt(ReadFromDB(), 10)
	} else if cmd.CommandText == "p" {
		ret = dialogue.ParticleNew
	} else if cmd.CommandText == "q" {
		ret = dialogue.MusicControl
	} else if cmd.CommandText == "r" {
		ret = dialogue.PaintingNew
	} else if cmd.CommandText == "s" {
		ret = dialogue.PlanetNew
	} else if cmd.CommandText == "t" {
		ret = dialogue.ThemeChange
	} else if cmd.CommandText == "u" {
		ret = dialogue.DiseaseNew
	} else if cmd.CommandText == "v" {
		ret = dialogue.TechnologyNew
	} else if cmd.CommandText == "w" {
		ret = dialogue.LandNew
	} else if cmd.CommandText == "x" {
		ret = dialogue.SaveFile
	} else if cmd.CommandText == "y" {
		ret = dialogue.PsychologyInspect
	} else if cmd.CommandText == "z" {
		ret = dialogue.LoadFile
	} else if cmd.CommandText == "welcome" {
		ret = dialogue.Introduction
	} else {
		ret = "Not a command: " + cmd.CommandText
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