package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/rs/cors"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type (
	// Snack struct for db
	Snack struct {
		gorm.Model
		Name        string
		Description string
		Price       uint
		Img         string
		Perishable  bool
		Reviews     []Review
	}

	// Review struct for db
	Review struct {
		gorm.Model
		Title   string
		Text    string
		Rating  int
		SnackID uint64 `gorm:"index"` // Foreign key (belongs to), tag `index` will create index for this column
	}
)

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
	db.DropTableIfExists(&Snack{}, &Review{})

	db.AutoMigrate(&Snack{})
	s1 := Snack{Name: "Pork Rinds", Description: "Mauris lacinia sapien quis libero. Nam dui. Proin leo odio, porttitor id, consequat in, consequat ut, nulla. Sed accumsan felis.", Price: 8, Img: "https://az808821.vo.msecnd.net/content/images/thumbs/0000398_salt-pepper-pork-rinds-2-oz_560.jpeg", Perishable: true}
	s2 := Snack{Name: "Soup - Campbells Beef Noodle", Description: "Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Etiam vel augue. Vestibulum rutrum rutrum neque.", Price: 26, Img: "https://images-na.ssl-images-amazon.com/images/I/71MavWF1P9L._SY550_.jpg", Perishable: false}
	s3 := Snack{Name: "Chicken - Chicken Phyllo", Description: "Donec vitae nisi. Nam ultrices, libero non mattis pulvinar, nulla pede ullamcorper augue, a suscipit nulla elit ac nulla.", Price: 5, Img: "https://tmbidigitalassetsazure.blob.core.windows.net/secure/RMS/attachments/37/1200x1200/exps191978_SD163575C10_07_6b.jpg", Perishable: true}
	db.Create(&s1)
	db.Create(&s2)
	db.Create(&s3)

	db.AutoMigrate(&Review{})
	re1 := Review{Title: "Incredible!", Text: "If it were a person I'd say to it: Is your name Dan Druff? You get into people's hair. I mean like, I'd say that you're funny but looks aren't everything.", Rating: 1, SnackID: 1}
	re2 := Review{Title: "Tasty!", Text: "If it were a person I'd say to this snack: I appreciate all of your opinions. I mean like, You have ten of the best fingers I have ever seen!", Rating: 3, SnackID: 1}
	re3 := Review{Title: "Tasty!", Text: "If it were a person I'd say to it: Learn from your parents' mistakes - use birth control! I mean like, I thought of you all day today. I was at the zoo.", Rating: 2, SnackID: 2}
	re4 := Review{Title: "Refined!", Text: "If it were a person I'd say to this snack: I would share my dessert with you. I mean like, You are a champ!", Rating: 5, SnackID: 3}
	re5 := Review{Title: "Handmade!", Text: "If it were a person I'd say to this snack: Treat yourself to another compliment! I mean like, You smell nice.", Rating: 5, SnackID: 3}

	db.Create(&re1)
	db.Create(&re2)
	db.Create(&re3)
	db.Create(&re4)
	db.Create(&re5)

	// Snacks Routes
	router.HandleFunc("/snacks", GetSnacks).Methods("GET")
	router.HandleFunc("/snacks/{id}", GetSnack).Methods("GET")
	router.HandleFunc("/snacks", CreateSnack).Methods("POST")
	router.HandleFunc("/snacks/{id}", UpdateSnack).Methods("PUT")
	router.HandleFunc("/snacks/{id}", DeleteSnack).Methods("DELETE")

	// Reviews Routes
	router.HandleFunc("/api/snacks/{id}/reviews", GetReviews).Methods("GET")
	router.HandleFunc("/api/snacks/{id}/reviews/{revId}", GetReview).Methods("GET")
	router.HandleFunc("/api/snacks/{id}/reviews", CreateReview).Methods("POST")
	router.HandleFunc("/api/snacks/{id}/reviews/{revId}", UpdateReview).Methods("PUT")
	router.HandleFunc("/api/snacks/{id}/reviews/{revId}", DeleteReview).Methods("DELETE")

	handler := cors.Default().Handler(router)

	fmt.Println("Server running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", handler))
}

//\/\/\ ERRORS /\/\/\

// // ErrorBadRequest for bad request error
// var ErrorBadRequest = errors.New("Bad request")
// var ErrorNotAllowed = errors.New("Not allowed")
// var ErrorInternalServer = errors.New("Internal Server Error")
// var ErrorNotFound = errors.New("Not Found")

// var ErrorExistingUser = errors.New("User already exists in the sytem")

//\/\/\ CRUD SNACKS //\/\/\

// GetSnacks function in main
func GetSnacks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var snacks []Snack
	db.Find(&snacks)
	if err := db.Find(&snacks).Error; err != nil {
		http.Error(w, err.Error(), 404)
		return
	}

	json.NewEncoder(w).Encode(&snacks)
}

// GetSnack function in main
func GetSnack(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	var snack Snack
	db.First(&snack, params["id"])
	if err := db.First(&snack).Error; err != nil {
		http.Error(w, "Snack Not Found", 404)
		return
	}

	json.NewEncoder(w).Encode(&snack)
}

// CreateSnack function in main
func CreateSnack(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var snack Snack
	fmt.Println("This is r.Body", r.Body)
	json.NewDecoder(r.Body).Decode(&snack)
	db.Create(&snack)
	json.NewEncoder(w).Encode(&snack)
}

// UpdateSnack function in main
func UpdateSnack(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fmt.Println("Updating snack yayyyy")
	var snack Snack
	var snack2 Snack
	json.NewDecoder(r.Body).Decode(&snack2)
	fmt.Println("I am decoded snack", snack2)
	params := mux.Vars(r)
	db.Find(&snack, "id = ?", params["id"])
	db.Model(&snack).Updates(Snack{Name: snack2.Name, Description: snack2.Description, Price: snack2.Price, Img: snack2.Img, Perishable: snack2.Perishable})
	// db.Model(&snack).Update("Name", snack2.Name)
	// db.Model(&snack).Update("Description", snack2.Description)
	// db.Model(&snack).Update("Price", snack2.Price)
	// db.Model(&snack).Update("Img", snack2.Img)

	fmt.Println("I am updated snack", snack)
	json.NewEncoder(w).Encode(&snack)
}

// DeleteSnack function in main
func DeleteSnack(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	var snack Snack
	deleted := db.First(&snack, params["id"])
	fmt.Println("This is deleted", deleted)
	db.Delete(&snack)

	var snacks []Snack
	db.Find(&snacks)
	json.NewEncoder(w).Encode(&deleted)
}

//\/\/\ CRUD REVIEWS //\/\/\

// GetReviews func gets all reviews per snack
func GetReviews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var reviews []Review
	params := mux.Vars(r)

	db.Find(&reviews, "snack_id = ?", params["id"])
	if len(reviews) <= 0 {
		http.Error(w, "Review not found for this snack", 404)
		return
	}

	json.NewEncoder(w).Encode(&reviews)
}

// GetReview func gets a single review per snack
func GetReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// var reviews []Review
	var review Review
	params := mux.Vars(r)

	fmt.Println("this is db in getReview", db)
	db.Where("snack_id = ? AND id = ?", params["id"], params["revId"]).First(&review)

	if review.ID <= 0 {
		fmt.Println("Error thrown in GetReview by checking ID = 0")
		http.Error(w, "Review is not found for this snack", 404)
		return
	}

	json.NewEncoder(w).Encode(&review)
}

// CreateReview func creates a snack's review
func CreateReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var review Review
	fmt.Printf("%T", r.Body)
	params := mux.Vars(r)

	// review.SnackID = params["id"]
	json.NewDecoder(r.Body).Decode(&review)
	id := params["id"]
	id1, _ := strconv.ParseUint(id, 10, 64)
	review.SnackID = id1
	db.Create(&review)
	json.NewEncoder(w).Encode(&review)
}

// UpdateReview func updates a snack's review
func UpdateReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fmt.Println("Updating snack yayyyy")
	var review Review
	var review2 Review
	json.NewDecoder(r.Body).Decode(&review2)
	fmt.Println("I am decoded snack", review2)
	params := mux.Vars(r)
	db.Find(&review, "id = ?", params["id"])
	id := params["id"]
	id2, _ := strconv.ParseUint(id, 10, 64)
	db.Model(&review).Updates(Review{Title: review2.Title, Text: review2.Text, Rating: review2.Rating, SnackID: id2})

	fmt.Println("I am updated review", review)
	json.NewEncoder(w).Encode(&review)
}

// DeleteReview func deletes a snack's review
func DeleteReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	var review Review
	deleted := db.First(&review, params["id"])
	fmt.Println("This is deleted", deleted)
	db.Delete(&review)

	json.NewEncoder(w).Encode(&deleted)
}
