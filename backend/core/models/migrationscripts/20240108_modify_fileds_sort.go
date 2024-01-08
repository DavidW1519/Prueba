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

package migrationscripts

import (
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
)

type modfiyFieldsSort struct{}

func (u *modfiyFieldsSort) Up(baseRes context.BasicRes) errors.Error {
	db := baseRes.GetDal()
	// issues
	err := db.Exec("alter table issues modify original_type varchar(500) after type;")
	if err != nil {
		return err
	}
	err = db.Exec("alter table issues modify story_point DOUBLE after original_status;")
	if err != nil {
		return err
	}
	// pull_requests
	err = db.Exec("alter table pull_requests modify original_status varchar(100) after status;")
	if err != nil {
		return err
	}
	err = db.Exec("alter table pull_requests modify type varchar(100) after original_status;")
	if err != nil {
		return err
	}
	err = db.Exec("alter table pull_requests modify component varchar(100) after type;")
	if err != nil {
		return err
	}
	// cicd deployment commits
	err = db.Exec("alter table cicd_deployment_commits modify original_status varchar(100) after status;")
	if err != nil {
		return err
	}
	err = db.Exec("alter table cicd_deployment_commits modify original_result varchar(100) after result;")
	if err != nil {
		return err
	}
	err = db.Exec("alter table cicd_deployment_commits modify duration_sec DOUBLE after finished_date;")
	if err != nil {
		return err
	}
	err = db.Exec("alter table cicd_deployment_commits modify queued_date DATETIME after duration_sec;")
	if err != nil {
		return err
	}
	err = db.Exec("alter table cicd_deployment_commits modify queued_duration_sec DOUBLE after queued_date;")
	if err != nil {
		return err
	}
	// cicd deployments
	err = db.Exec("alter table cicd_deployments modify original_status varchar(100) after status;")
	if err != nil {
		return err
	}
	err = db.Exec("alter table cicd_deployments modify original_result varchar(100) after result;")
	if err != nil {
		return err
	}
	// cicd pipelines
	err = db.Exec("alter table cicd_pipelines modify original_status varchar(100) after status;")
	if err != nil {
		return err
	}
	err = db.Exec("alter table cicd_pipelines modify original_result varchar(100) after result;")
	if err != nil {
		return err
	}
	err = db.Exec("alter table cicd_pipelines modify duration_sec DOUBLE after finished_date;")
	if err != nil {
		return err
	}
	err = db.Exec("alter table cicd_pipelines modify started_date DATETIME after duration_sec;")
	if err != nil {
		return err
	}
	err = db.Exec("alter table cicd_pipelines modify queued_date DATETIME after started_date;")
	if err != nil {
		return err
	}
	err = db.Exec("alter table cicd_pipelines modify queued_duration_sec DOUBLE after queued_date;")
	if err != nil {
		return err
	}
	// cicd tasks
	err = db.Exec("alter table cicd_tasks modify original_status varchar(100) after status;")
	if err != nil {
		return err
	}
	err = db.Exec("alter table cicd_tasks modify original_result varchar(100) after result;")
	if err != nil {
		return err
	}
	err = db.Exec("alter table cicd_tasks modify created_date DATETIME after finished_date;")
	if err != nil {
		return err
	}
	err = db.Exec("alter table cicd_tasks modify duration_sec DOUBLE after created_date;")
	if err != nil {
		return err
	}
	err = db.Exec("alter table cicd_tasks modify queued_date DATETIME after duration_sec;")
	if err != nil {
		return err
	}
	err = db.Exec("alter table cicd_tasks modify queued_duration_sec DOUBLE after queued_date;")
	if err != nil {
		return err
	}

	return nil
}

func (*modfiyFieldsSort) Version() uint64 {
	return 20240108000007
}

func (*modfiyFieldsSort) Name() string {
	return "fix some tables fields sort"
}
