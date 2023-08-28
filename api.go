//hello babe you need a api to connect to a database and do some shit on it?

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Blog struct {
	ID       uint   `gorm:"primaryKey"`
	BlogName string `gorm:"column:blogname"`
	Blog     string `gorm:"column:blog"`
}

var db *gorm.DB

func main() {
	dsn := "host=localhost user=postgres password=123456789 dbname=blogs port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("Failed to connect to the database:", err)
		return
	}

	// Connection successful, print "Connected" to the console
	fmt.Println("Connected to the database. table:")

	var blogs []Blog
	if err := db.Find(&blogs).Error; err != nil {
		fmt.Println("Error retrieving blogs:", err)
		return
	}

	// Print the retrieved blogs
	for _, b := range blogs {
		fmt.Println(b.ID, b.BlogName, b.Blog)
	}
	r := mux.NewRouter()

	// Register the route handler for fetching all blogs
	r.HandleFunc("/blogs", getAllBlogs).Methods("GET")
	r.HandleFunc("/blogs", createBlog).Methods("POST")
	r.HandleFunc("/blogs/{id}", deleteBlog).Methods("DELETE")

	// Start the HTTP server using the Gorilla Mux router
	fmt.Println("Server started on :8000")
	log.Panic(http.ListenAndServe(":8000", r))

}

func getAllBlogs(w http.ResponseWriter, r *http.Request) {
	// Fetch all blogs from the database
	log.Println("hello from function")
	var blogs []Blog
	if err := db.Find(&blogs).Error; err != nil {
		http.Error(w, "Error retrieving blogs", http.StatusInternalServerError)
		return
	}

	// Convert the blogs to JSON and send the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(blogs); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func createBlog(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get the data for the new blog
	var newBlog Blog
	err := json.NewDecoder(r.Body).Decode(&newBlog)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	// Insert the new blog into the database
	if err := db.Create(&newBlog).Error; err != nil {
		http.Error(w, "Error creating blog", http.StatusInternalServerError)
		return
	}

	// Return a success response
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Blog created successfully"))
}

func deleteBlog(w http.ResponseWriter, r *http.Request) {
	// Get the blog ID from the URL parameters
	vars := mux.Vars(r)
	blogID := vars["id"]

	// Delete the blog entry from the database based on its ID
	if err := db.Delete(&Blog{}, blogID).Error; err != nil {
		http.Error(w, "Error deleting blog", http.StatusInternalServerError)
		return
	}

	// Return a success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Blog deleted successfully"))
}
