package user

import (
	"context"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// IsUsernameValid checks if the username is valid
func IsUsernameValid(username string) bool {
	// Check length
	if len(username) > 15 || len(username) == 0 {
		return false
	}

	// Check if alphanumeric
	isAlphanumeric := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return isAlphanumeric.MatchString(username)
}

// IsUsernameTaken checks if the username is already taken in the database
func IsUsernameTaken(db *mongo.Database, username string) (bool, error) {
	usersCollection := db.Collection("users")
	var result struct{} // Use an empty struct for existence check
	err := usersCollection.FindOne(context.TODO(), bson.D{{Key: "username", Value: username}}).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil // Username is not taken
		}
		return false, err // Error occurred while checking
	}

	return true, nil // Username is taken
}
