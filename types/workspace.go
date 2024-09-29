package types

import "time"

type Workspace struct {
	ID        int        `json:"id" gorm:"primary_key"`
	Name      string     `json:"name"`
	Avatar    string     `json:"avatar"`
	Role      string     `json:"role"`
	Status    bool       `json:"status"`
	CreatedAt *time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type WorkspacePayload struct {
	Name string `form:"name" binding:"required"`
	Role string `form:"role" binding:"required"`
}

func (e *Workspace) TableName() string {
	return "workspaces"
}

type WorkspaceUser struct {
	ID          int        `json:"id" gorm:"primary_key"`
	WorkspaceId int        `json:"workspace_id"`
	UserId      int        `json:"user_id"`
	IsOwner     bool       `json:"is_owner"`
	IsPrimary   bool       `json:"is_primary"`
	CreatedAt   *time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (e *WorkspaceUser) TableName() string {
	return "workspace_users"
}
