package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"server/config"
	"server/models"
	"server/types"
	"server/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary Get User Sessions
// @Description Get User Critique Sessions List API
// @ID get-user-sessions
// @Accept  json
// @Produce  json
// @Success 200 {object} []types.ProcessedSessionResponse
// @Failure 400 {object} map[string]string
// @Router /session [get]
// @Security BearerAuth
func GetUserSessions(c *gin.Context) {
	// Get user context data (from auth middleware)
	ctxUserData, ctxUserDataExists := c.Get("user")
	if !ctxUserDataExists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "User not found"})
		return
	}
	data := ctxUserData.(*types.User)
	if data == nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Failed to retrieve user data"})
		return
	}
	userData, _ := models.FetchUser(data.ID)
	if userData == nil {
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "data": nil, "message": "User data not found"})
		return
	}

	// Check for existing session
	sessions, err := models.FetchProcessedUserSessions(data.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": nil, "message": "Failed to retreive user sessions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": sessions, "message": "Session fetched successfully"})
}

// @Summary Get Session Details
// @Description Get Critique Session Details API
// @ID get-session-details
// @Accept  json
// @Produce  json
// @Param sessionId path int true "Session ID"
// @Success 200 {object} types.ProcessedSessionResponse
// @Failure 400 {object} map[string]string
// @Router /session/{sessionId} [get]
// @Security BearerAuth
func GetSession(c *gin.Context) {
	sessionId := c.Param("sessionId")
	intSessionId, err := strconv.Atoi(sessionId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": nil, "message": "Failed to parse sessionId"})
		return
	}

	// Retrieve Session data
	session, err := models.GetSessionData(intSessionId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": nil, "message": "Failed to retrieve session data"})
		return
	}

	processedSessionCollaborators := models.FetchProcessedSessionCollaborators(intSessionId)

	var collaborators []types.CollaboratorResponse

	for _, sessionCollaborator := range *processedSessionCollaborators {
		collaborator := types.CollaboratorResponse{
			ID:     sessionCollaborator.ID,
			UserId: sessionCollaborator.UserId,
			Name:   sessionCollaborator.Name,
			Email:  sessionCollaborator.Email,
			Avatar: sessionCollaborator.Avatar,
			Status: sessionCollaborator.Status,
		}
		collaborators = append(collaborators, collaborator)
	}

	response := types.ProcessedSessionResponse{
		ID:            session.ID,
		Title:         session.Title,
		Objective:     session.Objective,
		Stage:         session.Stage,
		Datetime:      session.Datetime,
		Duration:      session.Duration,
		Status:        session.Status,
		CreatedBy:     session.CreatedBy,
		Collaborators: collaborators,
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": response, "message": "Session data fetched successfully"})
}

// @Summary Create Session
// @Description Create Critique Session API
// @ID create-session
// @Accept  json
// @Produce  json
// @Param details body types.SessionPayload true "Session Details"
// @Success 200 {object} types.ProcessedSessionResponse
// @Failure 400 {object} map[string]string
// @Router /session [post]
// @Security BearerAuth
func CreateSession(c *gin.Context) {
	var sessionPayload types.SessionPayload
	if err := c.ShouldBindJSON(&sessionPayload); err != nil {
		validationError := config.ValidationErrors(err, c)
		if len(validationError) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"errors": validationError})
			return
		}
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

	// Check for existing session
	existingSession, _ := models.GetSession(sessionPayload.Title)
	if existingSession != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": nil, "message": "Session with same title already exists"})
		return
	}

	// Save Session Data
	session, err := models.SaveSession(sessionPayload.Title, sessionPayload.Objective, sessionPayload.Stage, sessionPayload.Datetime, int(sessionPayload.Duration), userData.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": nil, "message": "Failed to save session info"})
		return
	}

	response := types.ProcessedSessionResponse{
		ID:            session.ID,
		Title:         session.Title,
		Objective:     session.Objective,
		Stage:         session.Stage,
		Datetime:      session.Datetime,
		Duration:      session.Duration,
		Status:        session.Status,
		CreatedBy:     session.CreatedBy,
		Collaborators: []types.CollaboratorResponse{},
	}

	// Check existing users with email
	if sessionPayload.Collaborators != "" {
		emails := strings.Split(sessionPayload.Collaborators, ",")
		usersWithEmail, err := models.FetchUserByEmails(emails)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": nil, "message": "Failed to fetch collaborators data"})
			return
		}

		emailUserIdMap := make(map[string]int)
		for _, value := range *usersWithEmail {
			emailUserIdMap[value.Email] = value.ID
		}

		// Non-existing emails
		var nonExistingEmails []string
		for _, email := range emails {
			if emailUserIdMap[email] == 0 {
				nonExistingEmails = append(nonExistingEmails, email)
			}
		}

		// Create new users for non-existing emails
		if len(nonExistingEmails) >= 1 {
			newUsers := make([]types.User, len(nonExistingEmails))
			for i, email := range nonExistingEmails {
				newUsers[i].Email = email
			}
			users, err := models.InsertUsers(&newUsers)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": nil, "message": "Failed to add users"})
				return
			}

			for _, value := range *users {
				emailUserIdMap[value.Email] = value.ID
			}

		}
		var collaboratorUserIds []int
		for email := range emailUserIdMap {
			collaboratorUserIds = append(collaboratorUserIds, emailUserIdMap[email])
		}

		_, err = models.SaveSessionCollaborators(session.ID, collaboratorUserIds)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": nil, "message": "Failed to save session collaborators data"})
			return
		}

		sessionCollaborators := models.FetchProcessedSessionCollaborators(session.ID)
		collaborators := []types.CollaboratorResponse{}
		for _, sessionCollaborator := range *sessionCollaborators {
			collaborator := types.CollaboratorResponse{
				ID:     sessionCollaborator.ID,
				UserId: sessionCollaborator.UserId,
				Name:   sessionCollaborator.Name,
				Email:  sessionCollaborator.Email,
				Avatar: sessionCollaborator.Avatar,
				Status: sessionCollaborator.Status,
			}
			collaborators = append(collaborators, collaborator)
		}
		response.Collaborators = collaborators
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": response, "message": "Session created successfully"})
}

// @Summary Get Session Files List
// @Description Get Critique Session Files List API
// @ID get-session-files
// @Accept  json
// @Produce  json
// @Param sessionId path int true "Session ID"
// @Success 200 {object} []types.SessionAttachment
// @Failure 400 {object} map[string]string
// @Router /session/{sessionId}/files [get]
// @Security BearerAuth
func GetSessionFiles(c *gin.Context) {
	sessionId := c.Param("sessionId")
	intSessionId, err := strconv.Atoi(sessionId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": nil, "message": "Failed to parse sessionId"})
		return
	}

	// Retrieve Session File Attachments
	attachments, err := models.GetSessionAttachments(intSessionId, utils.SESSION_CATEGORY_ADD_FILES)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": nil, "message": "Failed to retrieve session file attachements"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": attachments, "message": "Session files fetched successfully"})
}

// @Summary Add Session Files
// @Description Add Files for Critique Session API
// @ID add-session-files
// @Accept  json
// @Produce  json
// @Param details body types.SessionAttachmentsPayload true "Session Files Details"
// @Success 200 {object} []types.SessionAttachment
// @Failure 400 {object} map[string]string
// @Router /session/files [post]
// @Security BearerAuth
func AddSessionFiles(c *gin.Context) {
	var payload types.SessionAttachmentsPayload
	if err := c.ShouldBind(&payload); err != nil {
		validationError := config.ValidationErrors(err, c)
		if len(validationError) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"errors": validationError})
			return
		}
	}
	form, _ := c.MultipartForm()
	files := form.File["images"]
	filePaths := []types.SessionAttachmentFileData{}
	for _, file := range files {
		fileExt := filepath.Ext(file.Filename)
		originalFileName := strings.TrimSuffix(filepath.Base(file.Filename), fileExt)
		now := time.Now()
		fileName := strings.ReplaceAll(strings.ToLower(originalFileName), " ", "-") + fmt.Sprintf("%v", now.Unix()) + fileExt
		filePath := utils.GetS3Url(fileName)

		// Save the file locally temporarily
		tempFile := filepath.Join("/tmp", fileName)
		if err := c.SaveUploadedFile(file, tempFile); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": nil, "message": "Failed to save file temporarily"})
			return
		}

		// Upload file to S3
		err := utils.UploadFileToS3(fileName, tempFile)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": nil, "message": "Failed to upload file to S3"})
			return
		}

		// Append all file paths
		filePaths = append(filePaths, types.SessionAttachmentFileData{
			Name: fileName,
			Url:  filePath,
		})
	}

	// Save session file attachements
	if len(filePaths) > 0 {
		err := models.SaveSessionAttachments(payload.SessionId, filePaths)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": nil, "message": "Failed to save session file attachements"})
			return
		}
	}

	// Retrieve Session Attachments
	attachments, err := models.GetSessionAttachments(payload.SessionId, utils.SESSION_CATEGORY_ADD_FILES)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": nil, "message": "Failed to retrieve session file attachements"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": attachments, "message": "Session files added successfully"})
}
