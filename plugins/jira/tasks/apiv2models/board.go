package apiv2models

import (
	"github.com/apache/incubator-devlake/plugins/jira/models"
)

type Board struct {
	ID       uint64 `json:"id"`
	Self     string `json:"self"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Location *struct {
		ProjectId      uint   `json:"projectId"`
		DisplayName    string `json:"displayName"`
		ProjectName    string `json:"projectName"`
		ProjectKey     string `json:"projectKey"`
		ProjectTypeKey string `json:"projectTypeKey"`
		AvatarURI      string `json:"avatarURI"`
		Name           string `json:"name"`
	} `json:"location"`
}

func (b Board) ToToolLayer(connectionId uint64) *models.JiraBoard {
	result := &models.JiraBoard{
		ConnectionId: connectionId,
		BoardId:      b.ID,
		Name:         b.Name,
		Self:         b.Self,
		Type:         b.Type,
	}
	if b.Location != nil {
		result.ProjectId = b.Location.ProjectId
	}
	return result
}
