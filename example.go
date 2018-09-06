package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"

// 	"github.com/gorilla/mux"
// 	"github.com/jinzhu/gorm"
// 	"github.com/rs/cors"

// 	_ "github.com/jinzhu/gorm/dialects/postgres"
// )

// // Resource struct for db
// type Resource struct {
// 	gorm.Model

// 	Link        string
// 	Name        string
// 	Author      string
// 	Description string
// }

// const (
// 	host     = "localhost"
// 	port     = 5432
// 	user     = "postgres"
// 	password = "123456"
// 	dbname   = "resources_dev"
// )

// var db *gorm.DB
// var err error

// func main() {
// 	fmt.Println("Running main func")
// 	router := mux.NewRouter()

// 	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
// 		"password=%s dbname=%s sslmode=disable",
// 		host, port, user, password, dbname)

// 	db, err = gorm.Open(
// 		"postgres",
// 		psqlInfo)

// 	if err != nil {
// 		panic("failed to connect database")
// 	}

// 	defer db.Close()

// 	db.AutoMigrate(&Resource{})
// 	db.Create(&Resource{Link: "google.com", Name: "search on google", Author: "Diep", Description: "Google search fun"})
// 	db.Create(&Resource{Link: "facebook.com", Name: "search on facebook", Author: "Michael", Description: "Facebook is fun"})



// 	router.HandleFunc("/resources", GetResources).Methods("GET")
// 	router.HandleFunc("/resources/{id}", GetResource).Methods("GET")
// 	router.HandleFunc("/resources", CreateResource).Methods("POST")
// 	router.HandleFunc("/resources/{id}", DeleteResource).Methods("DELETE")

// 	handler := cors.Default().Handler(router)

// 	fmt.Println("Server running on port 8000")
// 	log.Fatal(http.ListenAndServe(":8000", handler))
// }

// // GetResources function in main
// func GetResources(w http.ResponseWriter, r *http.Request) {
// 	var resources []Resource
// 	db.Find(&resources)
// 	json.NewEncoder(w).Encode(&resources)
// }

// // GetResource function in main
// func GetResource(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
// 	var resource Resource
// 	db.First(&resource, params["id"])
// 	json.NewEncoder(w).Encode(&resource)
// }

// // CreateResource function in main
// func CreateResource(w http.ResponseWriter, r *http.Request) {
// 	var resource Resource
// 	json.NewDecoder(r.Body).Decode(&resource)
// 	db.Create(&resource)
// 	json.NewEncoder(w).Encode(&resource)
// }

// // DeleteResource function in main
// func DeleteResource(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
// 	var resource Resource
// 	db.First(&resource, params["id"])
// 	db.Delete(&resource)

// 	var resources []Resource
// 	db.Find(&resources)
// 	json.NewEncoder(w).Encode(&resources)
// }
