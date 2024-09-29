package models

import (
	"server/config"
	"server/types"
)

// InsertUser inserts a new user into the database.
//
// Parameters:
//   - user: A pointer to the User object to be inserted.
//
// Returns:
//   - *User: A pointer to the inserted User object.
//   - error: An error object if there is an issue inserting the user.
func InsertUser(user *types.User) (*types.User, error) {
	result := config.DB.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

// Function to fetch user detail by id
//
// Parameters:
//   - id: An integer representing the user's unique identifier.
//
// Returns:
//   - *types.User: A pointer to the User object if found.
//   - error: An error object if there is an issue retrieving the user.
func FetchUser(id int) (*types.User, error) {
	var user types.User
	result := config.DB.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// Function to fetch user detail by email
//
// Parameters:
//   - email: An email string representing the user's unique email ID.
//
// Returns:
//   - *types.User: A pointer to the User object if found.
//   - error: An error object if there is an issue retrieving the user.
func FetchUserByEmail(email string) (*types.User, error) {
	var user types.User
	result := config.DB.Where("email=?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
