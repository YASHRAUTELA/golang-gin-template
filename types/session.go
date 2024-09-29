package types

import (
	"server/utils"
	"time"
)

type Session struct {
	ID        int        `json:"id" gorm:"primary_key"`
	Title     string     `json:"title"`
	Objective string     `json:"objective"`
	Stage     string     `json:"stage"`
	Datetime  string     `json:"datetime"`
	Duration  int        `json:"duration"`
	Status    bool       `json:"status"`
	CreatedBy int        `json:"created_by"`
	CreatedAt *time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type SessionPayload struct {
	Title         string `json:"title" binding:"required"`
	Objective     string `json:"objective" binding:"required"`
	Stage         string `json:"stage" binding:"required"`
	Datetime      string `json:"datetime" binding:"required"`
	Duration      int    `json:"duration" binding:"required"`
	Collaborators string `json:"collaborators"`
}

func (e *Session) TableName() string {
	return utils.SESSIONS_TABLE
}

type SessionAttachment struct {
	ID        int        `json:"id" gorm:"primary_key"`
	SessionId int        `json:"session_id"`
	Url       string     `json:"url"`
	Name      string     `json:"name"`
	Category  string     `json:"category"`
	Status    bool       `json:"status"`
	CreatedAt *time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type SessionAttachmentsPayload struct {
	SessionId int `form:"session_id" binding:"required"`
}

type SessionAttachmentFileData struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func (e *SessionAttachment) TableName() string {
	return utils.SESSION_ATTACHMENTS_TABLE
}

type SessionCollaborator struct {
	ID        int        `json:"id" gorm:"primary_key"`
	SessionId int        `json:"session_id"`
	UserId    int        `json:"user_id"`
	Status    bool       `json:"status"`
	CreatedAt *time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type ProcessedSessionCollaborator struct {
	ID     int    `json:"id" gorm:"primary_key"`
	UserId int    `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"string"`
	Avatar string `json:"avatar"`
	Status bool   `json:"status"`
}

type SessionCollaboratorPayload struct {
	SessionId string `json:"session_id" binding:"required"`
	Title     string `json:"title" binding:"required"`
}

type SessionCollaboratorResponse struct {
	ID        int    `json:"id"`
	SessionId int    `json:"session_id"`
	UserId    int    `json:"user_id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Avatar    string `json:"avatar"`
	Status    bool   `json:"status"`
}

func (e *SessionCollaborator) TableName() string {
	return utils.SESSION_COLLABORATORS_TABLE
}

type SessionResponse struct {
	Session
	Collaborators []SessionCollaborator `json:"collaborators"`
}

type CollaboratorResponse struct {
	ID     int    `json:"id"`
	UserId int    `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
	Status bool   `json:"status"`
}

type ProcessedSession struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Objective string `json:"objective"`
	Stage     string `json:"stage"`
	Datetime  string `json:"datetime"`
	Duration  int    `json:"duration"`
	Status    bool   `json:"status"`
	CreatedBy int    `json:"created_by"`
}

type ProcessedSessionResponse struct {
	ID            int                    `json:"id"`
	Title         string                 `json:"title"`
	Objective     string                 `json:"objective"`
	Stage         string                 `json:"stage"`
	Datetime      string                 `json:"datetime"`
	Duration      int                    `json:"duration"`
	Status        bool                   `json:"status"`
	CreatedBy     int                    `json:"created_by"`
	Collaborators []CollaboratorResponse `json:"collaborators"`
}
