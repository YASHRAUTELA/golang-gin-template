package models

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"server/config"
	"server/types"
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
