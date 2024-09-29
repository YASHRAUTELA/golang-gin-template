package models

import (
	"fmt"
	"server/config"
	"server/types"
	"server/utils"
)

func GetSession(title string) (*types.Session, error) {
	var session types.Session
	result := config.DB.Where("title=?", title).First(&session)
	if result.Error != nil {
		return nil, result.Error
	}
	return &session, nil
}

func FetchUserSessions(userId int) (*[]types.Session, error) {
	var sessions []types.Session

	result := config.DB.Raw(fmt.Sprintf("SELECT * FROM %s WHERE created_by=? OR id=ANY(SELECT session_id FROM %s WHERE user_id=?)", utils.SESSIONS_TABLE, utils.SESSION_COLLABORATORS_TABLE), userId, userId).Scan(&sessions)
	if result.Error != nil {
		return nil, result.Error
	}

	return &sessions, nil
}

func FetchProcessedSessionCollaborators(sessionId int) *[]types.SessionCollaboratorResponse {
	var sessionCollaborators []types.SessionCollaboratorResponse
	config.DB.Table(fmt.Sprintf("%s sc", utils.SESSION_COLLABORATORS_TABLE)).Select("sc.id, sc.session_id, sc.user_id, u.name, u.email, u.avatar, sc.status").Joins(fmt.Sprintf("JOIN %s u ON u.id=sc.user_id", utils.USERS_TABLE)).Where("sc.session_id=?", sessionId).Scan(&sessionCollaborators)
	return &sessionCollaborators
}

func FetchProcessedSessionsCollaborators(sessionIds []int) *[]types.SessionCollaboratorResponse {
	var sessionCollaborators []types.SessionCollaboratorResponse
	config.DB.Table(fmt.Sprintf("%s sc", utils.SESSION_COLLABORATORS_TABLE)).Select("sc.id, sc.session_id, sc.user_id, u.name, u.email, u.avatar, sc.status").Joins(fmt.Sprintf("JOIN %s u ON u.id=sc.user_id", utils.USERS_TABLE)).Where("sc.session_id IN ?", sessionIds).Scan(&sessionCollaborators)
	return &sessionCollaborators
}

func FetchProcessedUserSessions(userId int) (*[]types.ProcessedSessionResponse, error) {
	var sessions []types.ProcessedSession
	var processedSessions []types.ProcessedSessionResponse

	err := config.DB.Table(utils.SESSIONS_TABLE).
		Select("id, title, objective, stage, datetime, duration, status, created_by").
		Where("created_by=?", userId).
		Or(fmt.Sprintf("id=ANY(SELECT session_id FROM %s WHERE user_id=?)", utils.SESSION_COLLABORATORS_TABLE), userId).
		Find(&sessions).Error
	if err != nil {
		return nil, err
	}

	var sessionIds []int
	processedSessionIndexMap := make(map[int]int)
	for index, session := range sessions {
		processedSessionIndexMap[session.ID] = index
		sessionIds = append(sessionIds, session.ID)
		processedSession := types.ProcessedSessionResponse{
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
		processedSessions = append(processedSessions, processedSession)

	}

	sessionCollaborators := FetchProcessedSessionsCollaborators(sessionIds)

	for _, sessionCollaborator := range *sessionCollaborators {
		arrIndex := processedSessionIndexMap[sessionCollaborator.SessionId]
		collaborator := types.CollaboratorResponse{
			ID:     sessionCollaborator.ID,
			UserId: sessionCollaborator.UserId,
			Name:   sessionCollaborator.Name,
			Email:  sessionCollaborator.Email,
			Avatar: sessionCollaborator.Avatar,
			Status: sessionCollaborator.Status,
		}
		processedSessions[arrIndex].Collaborators = append(processedSessions[arrIndex].Collaborators, collaborator)
	}

	return &processedSessions, nil
}

func SaveSession(title string, objective string, stage string, datetime string, duration int, userId int) (*types.Session, error) {
	session := types.Session{
		Title:     title,
		Objective: objective,
		Stage:     stage,
		Datetime:  datetime,
		Duration:  duration,
		CreatedBy: userId,
	}
	result := config.DB.Save(&session)
	if result.Error != nil {
		return nil, result.Error
	}
	return &session,
		nil
}

func SaveSessionCollaborators(sessionId int, userIds []int) (*[]types.SessionCollaborator, error) {
	collaborators := make([]types.SessionCollaborator, len(userIds))

	for index, userId := range userIds {
		collaborators[index].SessionId = sessionId
		collaborators[index].UserId = userId
	}

	result := config.DB.Create(collaborators)
	if result.Error != nil {
		return nil, result.Error
	}

	return &collaborators,
		nil
}

func GetSessionData(sessionId int) (*types.Session, error) {
	var session types.Session
	result := config.DB.First(&session, sessionId)
	if result.Error != nil {
		return nil, result.Error
	}
	return &session,
		nil
}

func GetSessionAttachments(sessionId int, category string) (*[]types.SessionAttachment, error) {
	var sessionAttachments []types.SessionAttachment
	result := config.DB.Where("session_id=?", sessionId).Find(&sessionAttachments)
	if result.Error != nil {
		return nil, result.Error
	}
	return &sessionAttachments,
		nil
}

func SaveSessionAttachments(sessionId int, fileDataArr []types.SessionAttachmentFileData) error {
	sessionAttachments := []*types.SessionAttachment{}

	for _, fileData := range fileDataArr {
		sessionAttachment := types.SessionAttachment{
			SessionId: sessionId,
			Name:      fileData.Name,
			Url:       fileData.Url,
			Category:  utils.SESSION_CATEGORY_ADD_FILES,
			Status:    true,
		}
		sessionAttachments = append(sessionAttachments, &sessionAttachment)
	}

	result := config.DB.Create(sessionAttachments)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
