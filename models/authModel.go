package models

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"server/config"
	"server/types"
)

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
