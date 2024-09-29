package controllers

import (
	"net/http"
	"server/config"
	"server/models"
	"server/types"

	"github.com/gin-gonic/gin"
)

// @Summary Get user
// @Description Get user by ID
// @ID get-user
// @Accept  json
// @Produce  json
// @Success 200 {object} types.UserResponse
// @Failure 400 {object} map[string]string
// @Router /user [get]
// @Security BearerAuth
func GetUser(c *gin.Context) {
	var user types.User
	if err := c.ShouldBindJSON(&user); err != nil {
		validationError := config.ValidationErrors(err, c)
		if len(validationError) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"errors": validationError})
			return
		}
	}
	ctxUserData, ctxUserDataExists := c.Get("user")
	if !ctxUserDataExists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "User not found"})
		return
	}

	data := ctxUserData.(*types.User)
	if data == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve user data"})
		return
	}

	userData, _ := models.FetchUser(data.ID)
	if userData == nil {
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "data": nil, "message": "User data not found"})
		return
	}

	userResponse := types.UserResponse{
		ID:        userData.ID,
		Name:      userData.Name,
		Email:     userData.Email,
		Avatar:    userData.Avatar,
		CreatedAt: userData.CreatedAt,
		UpdatedAt: userData.UpdatedAt,
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": userResponse, "message": "User data fetched successfully"})
}
