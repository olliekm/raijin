package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"raijin/auth"
	"raijin/user"

	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

func main() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	// Use the database
	db := client.Database("chatapp") // Replace "chatapp" with your database name

	// Example: Insert a user document
	usersCollection := db.Collection("users") // Create a collection named "users"
	user := bson.D{{Key: "username", Value: "testuser"}, {Key: "password", Value: "testpass"}}
	insertResult, err := usersCollection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted user with ID: %v\n", insertResult.InsertedID)

	// Disconnect from MongoDB
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()
	setupAPI()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Adjust to your frontend URL
		AllowCredentials: true,
	})
	handler := c.Handler(http.DefaultServeMux)

	log.Fatal(http.ListenAndServe(":8080", handler))
	log.Println("Server started on port 8080")
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	var newUser user.User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate username (check for length and alpha numeric)
	if !user.IsUsernameValid(newUser.Username) {
		http.Error(w, "Invalid username. Must be alphanumeric and max 15 characters.", http.StatusBadRequest)
		return
	}

	taken, err := user.IsUsernameTaken(db, newUser.Username) // Pass the db pointer
	if err != nil {
		http.Error(w, "Error checking username availability", http.StatusInternalServerError)
		return
	}
	if taken {
		http.Error(w, "Username is already taken.", http.StatusConflict)
		return
	}

	hashedPassword, err := auth.HashPassword(newUser.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	newUser.Password = hashedPassword // Store the hashed password

	// Insert user into MongoDB
	userDoc := bson.D{
		{Key: "username", Value: newUser.Username},
		{Key: "display_name", Value: newUser.DisplayName},
		{Key: "password", Value: newUser.Password}, // Store the hashed password instead
		{Key: "profile_pic", Value: newUser.ProfilePic},
	}
	usersCollection := db.Collection("users")

	insertResult, err := usersCollection.InsertOne(context.TODO(), userDoc)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(insertResult.InsertedID) // Respond with the created user ID
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	var loginUser user.User
	if err := json.NewDecoder(r.Body).Decode(&loginUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve the user from the database
	usersCollection := db.Collection("users")
	var storedUser user.User
	err := usersCollection.FindOne(context.TODO(), bson.D{{Key: "username", Value: loginUser.Username}}).Decode(&storedUser)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Verify the password
	err = auth.VerifyPassword(storedUser.Password, loginUser.Password)
	if err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateJWT(storedUser.Username)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}
	// Successful login
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func getUserProfile(w http.ResponseWriter, r *http.Request) {
	// Extract the token from the request header
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Authorization header is required", http.StatusUnauthorized)
		return
	}

	// Validate the token
	username, err := auth.ValidateJWT(tokenString)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Retrieve the user from the database
	usersCollection := db.Collection("users")
	var userProfile user.User
	err = usersCollection.FindOne(context.TODO(), bson.D{{Key: "username", Value: username}}).Decode(&userProfile)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Respond with the user profile
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userProfile)
}

func setupAPI() {

	manager := NewManager()
	http.HandleFunc("/register", registerUser) // Register user endpoint
	http.HandleFunc("/login", loginUser)       // Log in user endpoint
	http.HandleFunc("/ws", manager.serveWS)
	http.HandleFunc("/profile", getUserProfile) // New profile retrieval endpoint

}
