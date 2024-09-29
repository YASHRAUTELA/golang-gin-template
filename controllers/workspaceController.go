package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"server/config"
	"server/models"
	"server/types"
	"server/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary Create Workspace
// @Description Create User Workspace
// @ID create-workspace
// @Accept  json
// @Produce  json
// @Param details body types.WorkspacePayload true "Workspace Details"
// @Success 200 {object} types.Workspace
// @Failure 400 {object} map[string]string
// @Router /workspace [post]
func CreateWorkspace(c *gin.Context) {
	var workspacePayload types.WorkspacePayload
	if err := c.ShouldBind(&workspacePayload); err != nil {
		validationError := config.ValidationErrors(err, c)
		if len(validationError) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"errors": validationError})
			return
		}
	}

	// Validate avatar form data
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": nil, "message": "Failed to retrieve avatar field data"})
		return
	}

	// Get user context data (from auth middleware)
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

	// Check for existing workspace
	existingWorkspace, _ := models.GetWorkspace(workspacePayload.Name)
	if existingWorkspace != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": nil, "message": "Workspace with same name already exists"})
		return
	}

	uniqueFileName := fmt.Sprintf("%s%s", uuid.New().String(), filepath.Ext(file.Filename))

	// Save the file locally temporarily
	tempFile := filepath.Join("/tmp", uniqueFileName)
	if err := c.SaveUploadedFile(file, tempFile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": nil, "message": "Failed to save file temporarily"})
		return
	}

	// Upload file to S3
	objectKey := uniqueFileName
	err = utils.UploadFileToS3(objectKey, tempFile)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": nil, "message": "Failed to upload file to S3"})
		return
	}

	// Generate the file URL
	var S3_BUCKET_NAME string = os.Getenv("AWS_S3_BUCKET_NAME")
	avatarUrl := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", S3_BUCKET_NAME, objectKey)

	// Save the avatar URL to the workspace record
	workspace, err := models.SaveWorkspace(workspacePayload.Name, workspacePayload.Role, avatarUrl)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": nil, "message": "Failed to save workspace data"})
		return
	}

	// Get User Workspaces
	userPrimaryWorkspacesCount, err := models.GetUserPrimaryWorkspacesCount(userData.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": nil, "message": "Failed to save workspace user's data"})
		return
	}

	// Save workspace user details
	_, err = models.SaveWorkspaceUser(workspace.ID, userData.ID, true, userPrimaryWorkspacesCount == 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": nil, "message": "Failed to save workspace user's data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": workspace, "message": "Workspace created successfully"})
}
