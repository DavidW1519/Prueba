/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package api

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/config"
	"github.com/apache/incubator-devlake/models/common"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
)

var vld = validator.New()

func TestConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {

	// decode
	var err error
	var connection models.TestConnectionRequest
	err = mapstructure.Decode(input.Body, &connection)
	if err != nil {
		return nil, err
	}
	// validate
	err = vld.Struct(connection)
	if err != nil {
		return nil, err
	}
	// test connection
	apiClient, err := helper.NewApiClient(
		connection.Endpoint,
		map[string]string{
			"Authorization": fmt.Sprintf("Basic %v", connection.Auth),
		},
		3*time.Second,
		connection.Proxy,
		nil,
	)
	if err != nil {
		return nil, err
	}
	res, err := apiClient.Get("api/2/serverInfo", nil, nil)
	if err != nil {
		return nil, err
	}
	// check if `/rest/` was missing
	if res.StatusCode == http.StatusNotFound && !strings.HasSuffix(connection.Endpoint, "/rest/") {
		endpointUrl, err := url.Parse(connection.Endpoint)
		if err != nil {
			return nil, err
		}
		refUrl, err := url.Parse("/rest/")
		if err != nil {
			return nil, err
		}
		restUrl := endpointUrl.ResolveReference(refUrl)
		return nil, errors.NewNotFound(fmt.Sprintf("Seems like an invalid Endpoint URL, please try %s", restUrl.String()))
	}
	resBody := &models.JiraServerInfo{}
	err = helper.UnmarshalResponse(res, resBody)
	if err != nil {
		return nil, err
	}
	// check version
	if resBody.DeploymentType == models.DeploymentServer {
		// only support 8.x.x or higher
		if versions := resBody.VersionNumbers; len(versions) == 3 && versions[0] < 8 {
			return nil, fmt.Errorf("Support JIRA Server 8+ only")
		}
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	return nil, nil
}

func findConnectionByInputParam(input *core.ApiResourceInput) (*models.JiraConnection, error) {
	jiraConnectionId, err := getJiraConnectionIdByInputParam(input)
	if err != nil {
		return nil, fmt.Errorf("invalid connectionId")
	}
	return getJiraConnectionById(jiraConnectionId)
}

func getJiraConnectionIdByInputParam(input *core.ApiResourceInput) (uint64, error) {
	connectionId := input.Params["connectionId"]
	if connectionId == "" {
		return 0, fmt.Errorf("missing connectionId")
	}
	return strconv.ParseUint(connectionId, 10, 64)
}

func getJiraConnectionById(id uint64) (*models.JiraConnection, error) {
	jiraConnection := &models.JiraConnection{}
	err := db.First(jiraConnection, id).Error
	if err != nil {
		return nil, err
	}

	// decrypt
	v := config.GetConfig()
	encKey := v.GetString(core.EncodeKeyEnvStr)
	jiraConnection.BasicAuthEncoded, err = core.Decrypt(encKey, jiraConnection.BasicAuthEncoded)
	if err != nil {
		return nil, err
	}

	return jiraConnection, nil
}
func mergeFieldsToJiraConnection(jiraConnection *models.JiraConnection, connections ...map[string]interface{}) error {
	// decode
	for _, connection := range connections {
		err := mapstructure.Decode(connection, jiraConnection)
		if err != nil {
			return err
		}
	}

	// validate
	vld := validator.New()
	err := vld.Struct(jiraConnection)
	if err != nil {
		return err
	}

	return nil
}

func refreshAndSaveJiraConnection(jiraConnection *models.JiraConnection, data map[string]interface{}) error {
	var err error
	// update fields from request body
	err = mergeFieldsToJiraConnection(jiraConnection, data)
	if err != nil {
		return err
	}

	encKey, err := getEncKey()
	if err != nil {
		return err
	}
	jiraConnection.BasicAuthEncoded, err = core.Encrypt(encKey, jiraConnection.BasicAuthEncoded)
	if err != nil {
		return err
	}

	// transaction for nested operations
	tx := db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	if jiraConnection.ID > 0 {
		err = tx.Save(jiraConnection).Error
	} else {
		err = tx.Create(jiraConnection).Error
	}
	if err != nil {
		if common.IsDuplicateError(err) {
			return fmt.Errorf("jira connection with name %s already exists", jiraConnection.Name)
		}
		return err
	}
	// perform optional operation
	typeMappings := data["typeMappings"]
	if typeMappings != nil {
		err = saveTypeMappings(tx, jiraConnection.ID, typeMappings)
		if err != nil {
			return err
		}
	}

	jiraConnection.BasicAuthEncoded, err = core.Decrypt(encKey, jiraConnection.BasicAuthEncoded)
	if err != nil {
		return err
	}
	return nil
}

/*
POST /plugins/jira/connections
{
	"name": "jira data connection name",
	"endpoint": "jira api endpoint, i.e. https://merico.atlassian.net/rest",
	"basicAuthEncoded": "generated by `echo -n <jira login email>:<jira token> | base64`",
	"epicKeyField": "name of customfield of epic key",
	"storyPointField": "name of customfield of story point",
	"typeMappings": { // optional, send empty object to delete all typeMappings of the data connection
		"userType": {
			"standardType": "devlake standard type",
			"statusMappings": {  // optional, send empt object to delete all status mapping for the user type
				"userStatus": {
					"standardStatus": "devlake standard status"
				}
			}
		}
	}
}
*/
func PostConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// create a new connection
	jiraConnection := &models.JiraConnection{}

	// update from request and save to database
	err := refreshAndSaveJiraConnection(jiraConnection, input.Body)
	if err != nil {
		return nil, err
	}

	return &core.ApiResourceOutput{Body: jiraConnection, Status: http.StatusOK}, nil
}

/*
PATCH /plugins/jira/connections/:connectionId
{
	"name": "jira data connection name",
	"endpoint": "jira api endpoint, i.e. https://merico.atlassian.net/rest",
	"basicAuthEncoded": "generated by `echo -n <jira login email>:<jira token> | base64`",
	"epicKeyField": "name of customfield of epic key",
	"storyPointField": "name of customfield of story point",
	"typeMappings": { // optional, send empty object to delete all typeMappings of the data connection
		"userType": {
			"standardType": "devlake standard type",
			"statusMappings": {  // optional, send empt object to delete all status mapping for the user type
				"userStatus": {
					"standardStatus": "devlake standard status"
				}
			}
		}
	}
}
*/
func PatchConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// load from db
	jiraConnection, err := findConnectionByInputParam(input)
	if err != nil {
		return nil, err
	}

	// update from request and save to database
	err = refreshAndSaveJiraConnection(jiraConnection, input.Body)
	if err != nil {
		return nil, err
	}

	return &core.ApiResourceOutput{Body: jiraConnection}, nil
}

/*
DELETE /plugins/jira/connections/:connectionId
*/
func DeleteConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// load from db
	jiraConnectionID, err := getJiraConnectionIdByInputParam(input)
	if err != nil {
		return nil, err
	}
	// cascading delete
	err = db.Where("id = ?", jiraConnectionID).Delete(&models.JiraConnection{}).Error
	if err != nil {
		return nil, err
	}
	err = db.Where("connection_id = ?", jiraConnectionID).Delete(&models.JiraIssueTypeMapping{}).Error
	if err != nil {
		return nil, err
	}
	err = db.Where("connection_id = ?", jiraConnectionID).Delete(&models.JiraIssueStatusMapping{}).Error
	if err != nil {
		return nil, err
	}

	return &core.ApiResourceOutput{Body: jiraConnectionID}, nil
}

/*
GET /plugins/jira/connections
*/
func ListConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	jiraConnections := make([]models.JiraConnection, 0)
	err := db.Find(&jiraConnections).Error
	if err != nil {
		return nil, err
	}
	encKey, err := getEncKey()
	if err != nil {
		return nil, err
	}
	for i := range jiraConnections {
		jiraConnections[i].BasicAuthEncoded, err = core.Decrypt(encKey, jiraConnections[i].BasicAuthEncoded)
		if err != nil {
			return nil, err
		}
	}
	return &core.ApiResourceOutput{Body: jiraConnections}, nil
}

/*
GET /plugins/jira/connections/:connectionId


{
	"name": "jira data connection name",
	"endpoint": "jira api endpoint, i.e. https://merico.atlassian.net/rest",
	"basicAuthEncoded": "generated by `echo -n <jira login email>:<jira token> | base64`",
	"epicKeyField": "name of customfield of epic key",
	"storyPointField": "name of customfield of story point",
	"typeMappings": { // optional, send empty object to delete all typeMappings of the data connection
		"userType": {
			"standardType": "devlake standard type",
			"statusMappings": {  // optional, send empt object to delete all status mapping for the user type
				"userStatus": {
					"standardStatus": "devlake standard status"
				}
			}
		}
	}
}
*/
func GetConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	jiraConnection, err := findConnectionByInputParam(input)
	if err != nil {
		return nil, err
	}

	detail := &models.JiraConnectionDetail{
		JiraConnection: *jiraConnection,
		TypeMappings:   make(map[string]map[string]interface{}),
	}

	typeMappings, err := findIssueTypeMappingByConnectionId(jiraConnection.ID)
	if err != nil {
		return nil, err
	}
	for _, jiraTypeMapping := range typeMappings {
		// type mapping
		typeMappingDict := map[string]interface{}{
			"standardType": jiraTypeMapping.StandardType,
		}
		detail.TypeMappings[jiraTypeMapping.UserType] = typeMappingDict

		// status mapping
		statusMappings, err := findIssueStatusMappingByConnectionIdAndUserType(
			jiraConnection.ID,
			jiraTypeMapping.UserType,
		)
		if err != nil {
			return nil, err
		}
		if len(statusMappings) == 0 {
			continue
		}
		statusMappingsDict := make(map[string]interface{})
		for _, jiraStatusMapping := range statusMappings {
			statusMappingsDict[jiraStatusMapping.UserStatus] = map[string]interface{}{
				"standardStatus": jiraStatusMapping.StandardStatus,
			}
		}
		typeMappingDict["statusMappings"] = statusMappingsDict
	}

	return &core.ApiResourceOutput{Body: detail}, nil
}

// GET /plugins/jira/connections/:connectionId/epics
func GetEpicsByConnectionId(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	jiraConnection, err := findConnectionByInputParam(input)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: [1]models.EpicResponse{{
		Id:    1,
		Title: jiraConnection.EpicKeyField,
		Value: jiraConnection.EpicKeyField,
	}}}, nil
}

// GET /plugins/jira/connections/:connectionId/granularities
type GranularitiesResponse struct {
	Id    int
	Title string
	Value string
}

func GetGranularitiesByConnectionId(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	jiraConnection, err := findConnectionByInputParam(input)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: [1]GranularitiesResponse{
		{
			Id:    1,
			Title: "Story Point Field",
			Value: jiraConnection.StoryPointField,
		},
	}}, nil
}

// GET /plugins/jira/connections/:connectionId/boards
func GetBoardsByConnectionId(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connectionId := input.Params["connectionId"]
	if connectionId == "" {
		return nil, fmt.Errorf("missing connectionid")
	}
	jiraConnectionId, err := strconv.ParseUint(connectionId, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid connectionId")
	}
	var jiraBoards []models.JiraBoard
	err = db.Where("connection_Id = ?", jiraConnectionId).Find(&jiraBoards).Error
	if err != nil {
		return nil, err
	}
	var boardResponses []models.BoardResponse
	for _, board := range jiraBoards {
		boardResponses = append(boardResponses, models.BoardResponse{
			Id:    int(board.BoardId),
			Title: board.Name,
			Value: fmt.Sprintf("%v", board.BoardId),
		})
	}
	return &core.ApiResourceOutput{Body: boardResponses}, nil
}

func getEncKey() (string, error) {
	// encrypt
	v := config.GetConfig()
	encKey := v.GetString(core.EncodeKeyEnvStr)
	if encKey == "" {
		// Randomly generate a bunch of encryption keys and set them to config
		encKey = core.RandomEncKey()
		v.Set(core.EncodeKeyEnvStr, encKey)
		err := config.WriteConfig(v)
		if err != nil {
			return encKey, err
		}
	}
	return encKey, nil
}
