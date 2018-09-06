package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/rs/cors"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Snack struct for db
type Snack struct {
	gorm.Model
	Name        string
	Description string
	Price       uint
	Img         string
	Perishable  bool
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123456"
	dbname   = "snacks_dev"
)

var db *gorm.DB
var err error

func main() {
	fmt.Println("Running main func")
	router := mux.NewRouter()

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err = gorm.Open(
		"postgres",
		psqlInfo)

	if err != nil {
		panic("failed to connect database")
	}

	defer db.Close()

	db.AutoMigrate(&Snack{})
	s1 := Snack{Name: "Pork Rinds", Description: "Mauris lacinia sapien quis libero. Nam dui. Proin leo odio, porttitor id, consequat in, consequat ut, nulla. Sed accumsan felis.", Price: 8, Img: "https://az808821.vo.msecnd.net/content/images/thumbs/0000398_salt-pepper-pork-rinds-2-oz_560.jpeg", Perishable: true}
	s2 := Snack{Name: "Soup - Campbells Beef Noodle", Description: "Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Etiam vel augue. Vestibulum rutrum rutrum neque.", Price: 26, Img: "https://images-na.ssl-images-amazon.com/images/I/71MavWF1P9L._SY550_.jpg", Perishable: false}
	s3 := Snack{Name: "Chicken - Chicken Phyllo", Description: "Donec vitae nisi. Nam ultrices, libero non mattis pulvinar, nulla pede ullamcorper augue, a suscipit nulla elit ac nulla.", Price: 5, Img: "https://tmbidigitalassetsazure.blob.core.windows.net/secure/RMS/attachments/37/1200x1200/exps191978_SD163575C10_07_6b.jpg", Perishable: true}
	db.Create(&s1)
	db.Create(&s2)
	db.Create(&s3)

	router.HandleFunc("/snacks", GetSnacks).Methods("GET")
	router.HandleFunc("/snacks/{id}", GetSnack).Methods("GET")
	router.HandleFunc("/snacks", CreateSnack).Methods("POST")
	router.HandleFunc("/snacks/{id}", UpdateSnack).Methods("PUT")
	router.HandleFunc("/snacks/{id}", DeleteSnack).Methods("DELETE")

	handler := cors.Default().Handler(router)

	fmt.Println("Server running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", handler))
}

// GetSnacks function in main
func GetSnacks(w http.ResponseWriter, r *http.Request) {
	var snacks []Snack
	db.Find(&snacks)
	json.NewEncoder(w).Encode(&snacks)
}

// GetSnack function in main
func GetSnack(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var snack Snack
	db.First(&snack, params["id"])
	json.NewEncoder(w).Encode(&snack)
}

// CreateSnack function in main
func CreateSnack(w http.ResponseWriter, r *http.Request) {
	var snack Snack
	fmt.Println("This is r.Body", r.Body)
	json.NewDecoder(r.Body).Decode(&snack)
	db.Create(&snack)
	json.NewEncoder(w).Encode(&snack)
}

// UpdateSnack function in main
func UpdateSnack(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Updating snack yayyyy")
	var snack Snack
	json.NewDecoder(r.Body).Decode(&snack)
	params := mux.Vars(r)
	db.Find(&snack, "id = ?", params["id"])
	// db.Model(&snack).Updates(Snack{Name: snack.Name, Description: snack.Description, Price: snack.Price, Img: snack.Img, Perishable: snack.Perishable})
	db.Model(&snack).Update("Name", snack.Name)
	db.Model(&snack).Update("Description", snack.Description)
	db.Model(&snack).Update("Price", snack.Price)
	db.Model(&snack).Update("Img", snack.Img)
	
	fmt.Println("I am snack in UpdateSnack", snack)
	json.NewEncoder(w).Encode(&snack)
}

// DeleteSnack function in main
func DeleteSnack(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var snack Snack
	deleted := db.First(&snack, params["id"])
	fmt.Println("This is deleted", deleted)
	db.Delete(&snack)

	var snacks []Snack
	db.Find(&snacks)
	json.NewEncoder(w).Encode(&deleted)
}
