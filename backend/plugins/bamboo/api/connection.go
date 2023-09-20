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
	"context"
	"net/http"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bamboo/models"
	"github.com/apache/incubator-devlake/server/api/shared"
)

type BambooTestConnResponse struct {
	shared.ApiBody
	Connection *models.BambooConn
}

// @Summary test bamboo connection
// @Description Test bamboo Connection
// @Tags plugins/bamboo
// @Param body body models.BambooConn true "json body"
// @Success 200  {object} BambooTestConnResponse "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/bamboo/test [POST]
func TestConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	// decode
	var err errors.Error
	var connection models.BambooConn
	if err = api.Decode(input.Body, &connection, vld); err != nil {
		return nil, err
	}

	// test connection
	_, err = api.NewApiClientFromConnection(
		context.TODO(),
		basicRes,
		&connection,
	)
	if err != nil {
		return nil, err
	}

	body := BambooTestConnResponse{}
	body.Success = true
	body.Message = "success"
	body.Connection = &connection

	return &plugin.ApiResourceOutput{Body: body, Status: http.StatusOK}, nil
}

// @Summary create bamboo connection
// @Description Create bamboo connection
// @Tags plugins/bamboo
// @Param body body models.BambooConnection true "json body"
// @Success 200  {object} models.BambooConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /plugins/bamboo/connections [POST]
func PostConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	// update from request and save to database
	connection := &models.BambooConnection{}
	err := connectionHelper.Create(connection, input)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connection, Status: http.StatusOK}, nil
}

// @Summary patch bamboo connection
// @Description Patch bamboo connection
// @Tags plugins/bamboo
// @Param body body models.BambooConnection true "json body"
// @Param connectionId path int true "connection ID"
// @Success 200  {object} models.BambooConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /plugins/bamboo/connections/{connectionId} [PATCH]
func PatchConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.BambooConnection{}
	err := connectionHelper.Patch(connection, input)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connection}, nil
}

// @Summary delete a bamboo connection
// @Description Delete a bamboo connection
// @Tags plugins/bamboo
// @Param connectionId path int true "connection ID"
// @Success 200  {object} models.BambooConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 409  {object} services.BlueprintProjectPairs "References exist to this connection"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /plugins/bamboo/connections/{connectionId} [DELETE]
func DeleteConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return connectionHelper.Delete(&models.BambooConnection{}, input)
}

// @Summary get all bamboo connections
// @Description Get all bamboo connections
// @Tags plugins/bamboo
// @Success 200  {object} []models.BambooConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /plugins/bamboo/connections [GET]
func ListConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var connections []models.BambooConnection
	err := connectionHelper.List(&connections)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connections, Status: http.StatusOK}, nil
}

// @Summary get bamboo connection detail
// @Description Get bamboo connection detail
// @Tags plugins/bamboo
// @Param connectionId path int true "connection ID"
// @Success 200  {object} models.BambooConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /plugins/bamboo/connections/{connectionId} [GET]
func GetConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.BambooConnection{}
	err := connectionHelper.First(connection, input.Params)
	return &plugin.ApiResourceOutput{Body: connection}, err
}
