package models

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"server/config"
	"server/types"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// SocialLoginWithGoogle authenticates a user with Google using an OAuth token.
// It retrieves the user's information from Google's userinfo endpoint and
// either creates a new user or updates an existing one in the local database.
//
// Parameters:
//   - token: A string containing the Google OAuth token.
//
// Returns:
//   - *types.User: A pointer to the User object if authentication is successful.
//   - error: An error if any step in the process fails, nil otherwise.
//
// The function performs the following steps:
// 1. Sends a GET request to Google's userinfo endpoint with the provided token.
// 2. Parses the response to extract user information.
// 3. Checks if the user exists in the local database using their email.
// 4. Updates or creates the user record with the latest information from Google.
//
// Possible errors include network issues, invalid tokens, or database errors.
// If an error occurs, it returns nil for the user and an error describing the failure.
func SocialLoginWithGoogle(token string) (*types.User, error) {
	const userInfoEndpoint string = "https://www.googleapis.com/oauth2/v2/userinfo"
	client := &http.Client{}
	req, err := http.NewRequest("GET", userInfoEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read user info response: %w", err)
	}

	var googleUser types.GoogleUser
	if err := json.Unmarshal(body, &googleUser); err != nil {
		return nil, fmt.Errorf("failed to decode user info response: %w", err)
	}

	var user types.User
	config.DB.Where("email=?", googleUser.Email).First(&user)
	user.Name = googleUser.Name
	user.Email = googleUser.Email
	user.Avatar = googleUser.Picture

	// Add or Update user data
	config.DB.Save(&user)
	return &user, nil
}

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
