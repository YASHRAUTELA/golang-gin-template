package models

import (
	"fmt"
	"os"
	"server/config"
	"server/types"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// VerifyJWTToken validates the JWT token from the jwtToken string param.
//
// Parameters:
//   - jwtToken: The string containing the JWT token.
//
// Returns:
//   - *types.User: An User object if the token is valid.
//   - error: An error object if there is an issue validating the token.
func VerifyJWTToken(jwtToken string) (*types.User, error) {
	secret := []byte(os.Getenv("SECRET"))
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return false, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid or expired token")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user := &types.User{
			ID:    int(claims["id"].(float64)),
			Name:  claims["name"].(string),
			Email: claims["email"].(string),
		}
		return user, nil
	} else {
		return nil, fmt.Errorf("invalid token claims")
	}
}

// CreateJWTToken creates a new JWT token for the given user.
//
// Parameters:
//   - user: An User object representing the user for whom the token is being created.
//
// Returns:
//   - string: The JWT token string.
//   - error: An error object if there is an issue creating the token.
func CreateJWTToken(user types.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"name":     user.Name,
		"email":    user.Email,
		"password": user.Password,
	})
	secret := []byte(os.Getenv("SECRET"))
	tokenString, err := token.SignedString(secret)
	return tokenString, err
}

// HashPassword generates a hashed password from the given password string.
//
// Parameters:
//   - password: The string representing the password to be hashed.
//
// Returns:
//   - string: The hashed password string.
//   - error: An error object if there is an issue hashing the password.
func HashPassword(password string) (string, error) {
	bytes, error := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), error
}

// CheckHashPassword compares the given password string with the provided hash string.
//
// Parameters:
//   - password: The string representing the password to be checked.
//   - hash: The string representing the hash to be compared against.
//
// Returns:
//   - bool: True if the password matches the hash, false otherwise.
//   - error: An error object if there is an issue comparing the password and hash.
func CheckHashPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

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
