package main

import (
	"os"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
  	"gorm.io/gorm"
)

type Person struct {
	gorm.Model

	Name  string
	Email string

	Movies []Movie
}

type Movie struct {
	gorm.Model

	Title    string
	Rating   int
	Year     int
	PersonID int
}

var (
	person = &Person{Name: "Jane", Email: "janedoe@gmail.com"}
	movies = []Movie{
		{Title: "The Witcher", Rating: 8, Year: 2019, PersonID: 1},
		{Title: "Knives out", Rating: 8, Year: 2018, PersonID: 1},
		{Title: "Heart break gallery", Rating: 6, Year: 2019, PersonID: 1},
	}
)

var db *gorm.DB
var err error

func main() {
	
	postgresConnection := os.Getenv("POSTGRES_CONNECTION")


	db, err = gorm.Open(postgres.Open(postgresConnection))
	if err != nil {
		log.Fatal(err)
		panic(err)
	} else {
		fmt.Println("Successfully connected to database!")

	}
	
	err = db.DB().Ping()
	if err != nil {
		panic(err)
	}

	defer db.Close()

	db.AutoMigrate(&Person{})
	db.AutoMigrate(&Movie{})

	// API routes
	router := mux.NewRouter()

	router.HandleFunc("/people", getPeople).Methods("GET")
	router.HandleFunc("/person/{id}", getPerson).Methods("GET")
	router.HandleFunc("/movie/{id}", getMovie).Methods("GET")
	router.HandleFunc("/movies", getMovies).Methods("GET")
	router.HandleFunc("/create/person", createPerson).Methods("POST")
	router.HandleFunc("/create/movie", createMovie).Methods("POST")
	router.HandleFunc("/delete/person/{id}", deletePerson).Methods("DELETE")
	router.HandleFunc("/delete/movie/{id}", deleteMovie).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))

}

// Model controllers
func getPeople(w http.ResponseWriter, r *http.Request) {
	var people []Person
	db.Find(&people)

	json.NewEncoder(w).Encode(&people)
}

func getPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var person Person
	var movies []Movie

	db.First(&person, params["id"])
	db.Model(&person).Related(&movies)

	person.Movies = movies

	json.NewEncoder(w).Encode(&person)
}

func createPerson(w http.ResponseWriter, r *http.Request) {
	var person Person
	json.NewDecoder(r.Body).Decode(&person)

	createdPerson := db.Create(&person)
	err = createdPerson.Error
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(&person)
	}
}

func deletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var person Person

	db.First(&person, params["id"])
	db.Delete(&person)

	json.NewEncoder(w).Encode(&person)
}

func getMovies(w http.ResponseWriter, r *http.Request) {
	var movies []Movie

	db.Find(&movies)

	json.NewEncoder(w).Encode(&movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var movie Movie

	db.First(&movie, params["id"])

	json.NewEncoder(w).Encode(&movie)
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	var movie Movie
	json.NewDecoder(r.Body).Decode(&movie)

	createdMovie := db.Create(&movie)
	err = createdMovie.Error
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(&movie)
	}
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var movie Movie

	db.First(&movie, params["id"])
	db.Delete(&movie)

	json.NewEncoder(w).Encode(&movie)
}
