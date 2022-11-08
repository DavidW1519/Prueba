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

package tasks

import (
	"fmt"
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models/common"
	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/pagerduty/models"
	"reflect"
	"time"
)

var ConvertIncidentsMeta = core.SubTaskMeta{
	Name:             "convertIncidents",
	EntryPoint:       ConvertIncidents,
	EnabledByDefault: true,
	Description:      "Convert incidents into domain layer table issues",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}

type (
	// IncidentWithUser struct that represents the joined query result
	IncidentWithUser struct {
		common.NoPKModel
		*models.Incident
		*models.User
		AssignedAt time.Time
	}
)

func ConvertIncidents(taskCtx core.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*PagerDutyTaskData)
	cursor, err := db.Cursor(
		dal.Select("pi.*, pu.*, pa.assigned_at"),
		dal.From("_tool_pagerduty_incident AS pi"),
		dal.Join(`LEFT JOIN _tool_pagerduty_assignment AS pa ON pa.incident_number = pi.number`),
		dal.Join(`LEFT JOIN _tool_pagerduty_user AS pu ON pa.user_id = pu.id`),
		dal.Where("pi.connection_id = ?", data.Options.ConnectionId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()
	seenIncidents := map[int]*IncidentWithUser{}
	idGen := didgen.NewDomainIdGenerator(&models.Incident{})
	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: models.PagerDutyParams{
				ConnectionId: data.Options.ConnectionId,
				Stream:       models.IncidentStream,
			},
			Table: RAW_INCIDENTS_TABLE,
		},
		InputRowType: reflect.TypeOf(IncidentWithUser{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			combined := inputRow.(*IncidentWithUser)
			incident := combined.Incident
			user := combined.User
			if seen, ok := seenIncidents[incident.Number]; ok {
				if combined.AssignedAt.Before(seen.AssignedAt) {
					// skip this one (it's an older assignee)
					return nil, nil
				}
			}
			status := getStatus(incident)
			leadTime, resolutionDate := getTimes(incident)
			domainIssue := &ticket.Issue{
				DomainEntity: domainlayer.DomainEntity{
					Id: idGen.Generate(data.Options.ConnectionId, incident.Number),
				},
				Url:             incident.Url,
				IssueKey:        fmt.Sprintf("%d", incident.Number),
				Description:     incident.Summary,
				Type:            ticket.INCIDENT,
				Status:          status,
				OriginalStatus:  string(incident.Status),
				ResolutionDate:  resolutionDate,
				CreatedDate:     &incident.CreatedDate,
				UpdatedDate:     &incident.UpdatedDate,
				LeadTimeMinutes: leadTime,
				Priority:        string(incident.Urgency),
				AssigneeId:      user.Id,
				AssigneeName:    user.Name,
			}
			seenIncidents[incident.Number] = combined
			return []interface{}{
				domainIssue,
			}, nil
		},
	})
	if err != nil {
		return err
	}
	return converter.Execute()
}

func getStatus(incident *models.Incident) string {
	if incident.Status == models.IncidentStatusTriggered {
		return ticket.TODO
	}
	if incident.Status == models.IncidentStatusAcknowledged {
		return ticket.IN_PROGRESS
	}
	if incident.Status == models.IncidentStatusResolved {
		return ticket.DONE
	}
	panic("unknown incident status encountered")
}

func getTimes(incident *models.Incident) (int64, *time.Time) {
	var leadTime int64
	var resolutionDate *time.Time
	if incident.Status == models.IncidentStatusResolved {
		resolutionDate = &incident.UpdatedDate
		leadTime = int64(resolutionDate.Sub(incident.CreatedDate).Minutes())
	}
	return leadTime, resolutionDate
}
