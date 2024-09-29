package models

import (
	"fmt"
	"net/http"
	"os"
	"server/config"
	"server/types"
	"strings"

	"github.com/gin-gonic/gin"
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

func HashPassword(password string) (string, error) {
	bytes, error := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), error
}

func CheckHashPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func InsertUser(user *types.User) (*types.User, error) {
	result := config.DB.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func InsertUsers(users *[]types.User) (*[]types.User, error) {
	result := config.DB.Create(users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func FetchUsers(ctx *gin.Context) {
	var users []types.User
	config.DB.Find(&users)
	ctx.JSON(http.StatusOK, users)
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

// Function to fetch users list with emails
//
// Parameters:
//   - email: An email string representing the user's unique email ID.
//
// Returns:
//   - *[]types.User: A pointer to the User object if found.
//   - error: An error object if there is an issue retrieving the user.
func FetchUserByEmails(emails []string) (*[]types.User, error) {
	var users []types.User
	// Trim spaces in each email
	for i := range emails {
		emails[i] = strings.TrimSpace(emails[i])
	}
	result := config.DB.Where("email IN ?", emails).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return &users, nil
}

func UpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")
	var user types.User
	if err := config.DB.First(&user, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Save(&user)
	ctx.JSON(http.StatusOK, user)
}

func DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	var user types.User
	if err := config.DB.First(&user, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	}
	config.DB.Delete(&user)
	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
