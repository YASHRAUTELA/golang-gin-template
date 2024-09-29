package models

import (
	"server/config"
	"server/types"
)

func GetWorkspace(name string) (*types.Workspace, error) {
	var workspace types.Workspace
	result := config.DB.Where("name=?", name).First(&workspace)
	if result.Error != nil {
		return nil, result.Error
	}
	return &workspace, nil
}

func GetUserWorkspaces(userId int) (*types.WorkspaceUser, error) {
	var userWorkspaces types.WorkspaceUser
	result := config.DB.Where("user_id=?", userId).Find(&userWorkspaces)
	if result.Error != nil {
		return nil, result.Error
	}
	return &userWorkspaces, nil
}

func GetUserPrimaryWorkspacesCount(userId int) (int64, error) {
	var count int64
	var userWorkspaces types.WorkspaceUser
	result := config.DB.Where("user_id=?", userId).Where("is_primary", true).Find(&userWorkspaces).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

func SaveWorkspace(name string, role string, avatar string) (*types.Workspace, error) {
	workspace := types.Workspace{
		Name:   name,
		Role:   role,
		Avatar: avatar,
	}

	result := config.DB.Save(&workspace)
	if result.Error != nil {
		return nil, result.Error
	}
	return &workspace, nil
}

func SaveWorkspaceUser(workspaceId int, userId int, isOwner bool, isPrimary bool) (*types.WorkspaceUser, error) {
	workspaceUser := types.WorkspaceUser{
		WorkspaceId: workspaceId,
		UserId:      userId,
		IsOwner:     isOwner,
		IsPrimary:   isPrimary,
	}
	result := config.DB.Save(&workspaceUser)
	if result.Error != nil {
		return nil, result.Error
	}
	return &workspaceUser, nil
}
