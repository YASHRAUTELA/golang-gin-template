package controllers

import (
	"fmt"
	"net/http"
	"server/config"
	"server/models"
	"server/types"

	"github.com/gin-gonic/gin"
)

// @Summary Login
// @Description User login
// @ID login
// @Accept  json
// @Produce  json
// @Param credentials body types.LoginPayload true "User credentials"
// @Success 200 {object} types.AuthResponse
// @Failure 400 {object} map[string]string
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var loginData types.LoginPayload
	if err := c.ShouldBindJSON(&loginData); err != nil {
		validationError := config.ValidationErrors(err, c)
		if len(validationError) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"errors": validationError})
			return
		}
	}

	// Retreive user data
	userData, _ := models.FetchUserByEmail(loginData.Email)
	if userData == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "data": nil, "message": "Invalid email or password"})
		return
	}

	// Validate user credentials
	isValidPassword := models.CheckHashPassword(loginData.Password, userData.Password)
	if !isValidPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "data": nil, "message": "Invalid email or password"})
		return
	}
	user := types.User{
		ID:       userData.ID,
		Name:     userData.Name,
		Email:    loginData.Email,
		Password: loginData.Password,
	}

	// Generate JWT token
	tokenString, tokenError := models.CreateJWTToken(user)
	if tokenError != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": nil, "message": "Failed to register user"})
		return
	}
	response := types.AuthResponse{
		Token: tokenString,
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": response, "message": "User login successful."})
}

// @Summary Register
// @Description User registration
// @ID register
// @Accept  json
// @Produce  json
// @Param user body types.RegisterPayload true "User info"
// @Success 200 {object} types.AuthResponse
// @Failure 400 {object} map[string]string
// @Router /auth/register [post]
func Register(c *gin.Context) {
	var registerData types.RegisterPayload
	if err := c.ShouldBindJSON(&registerData); err != nil {
		validationError := config.ValidationErrors(err, c)
		if len(validationError) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"errors": validationError})
			return
		}
	}

	user := types.User{
		Name:     registerData.Name,
		Email:    registerData.Email,
		Password: registerData.Password,
	}

	// Retreive user data
	userData, _ := models.FetchUserByEmail(registerData.Email)
	if userData != nil {
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "data": nil, "message": "User with same email already exists"})
		return
	}

	// Generate hash password
	hashPassword, hashPasswordError := models.HashPassword(user.Password)
	if hashPasswordError != nil {
		fmt.Println(hashPasswordError)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": nil, "message": "Failed to register user"})
		return
	}

	user.Password = hashPassword

	// Save user data
	saveUserData, err := models.InsertUser(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": nil, "message": "Failed to register user"})
		return
	}

	// Generate JWT token
	user.ID = saveUserData.ID
	tokenString, tokenError := models.CreateJWTToken(user)
	if tokenError != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": nil, "message": "Failed to register user"})
		return
	}

	response := types.AuthResponse{
		Token: tokenString,
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": response, "message": "User registration successful."})
}

// @Summary Social Login
// @Description User registration/login with Social providers
// @ID sociallogin
// @Accept  json
// @Produce  json
// @Param user body types.SocialLoginPayload true "User info"
// @Success 200 {object} types.AuthResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/sociallogin [post]
func SocialLogin(c *gin.Context) {
	var payload types.SocialLoginPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		validationError := config.ValidationErrors(err, c)
		if len(validationError) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"errors": validationError})
			return
		}
	}

	if payload.Provider == "google" {
		authData, authError := models.SocialLoginWithGoogle(payload.Token)
		if authError != nil {
			fmt.Println(authError)
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "data": nil, "message": "Failed to register user"})
			return
		}

		// Generate JWT token
		user := types.User{
			ID:    authData.ID,
			Name:  authData.Name,
			Email: authData.Email,
		}

		tokenString, tokenError := models.CreateJWTToken(user)
		if tokenError != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "data": nil, "message": "Failed to register user"})
			return
		}

		response := types.AuthResponse{
			Token: tokenString,
		}
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": response, "message": "User registration successful."})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": nil, "message": "User registration failed."})
}
