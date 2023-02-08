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

package service

import (
	"fmt"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/helpers/pluginhelper"
	"github.com/apache/incubator-devlake/plugins/customize/models"
	"io"
	"strings"
)

// Service wraps database operations
type Service struct {
	dal dal.Dal
}

func NewService(dal dal.Dal) *Service {
	return &Service{dal: dal}
}

// GetFields returns all the customized fields for the table
func (s *Service) GetFields(table string) ([]models.CustomizedField, errors.Error) {
	// the customized fields created before v0.16.0 were not recorded in the table `_tool_customized_field`, we should take care of them
	columns, err := s.dal.GetColumns(&models.Table{Name: table}, func(columnMeta dal.ColumnMeta) bool {
		return strings.HasPrefix(columnMeta.Name(), "x_")
	})
	if err != nil {
		return nil, errors.Default.Wrap(err, "GetColumns error")
	}
	ff, err := s.getCustomizedFields(table)
	if err != nil {
		return nil, err
	}
	fieldMap := make(map[string]models.CustomizedField)
	for _, f := range ff {
		fieldMap[f.ColumnName] = f
	}
	var result []models.CustomizedField
	for _, col := range columns {
		if field, ok := fieldMap[col.Name()]; ok {
			result = append(result, field)
		} else {
			result = append(result, models.CustomizedField{
				ColumnName: col.Name(),
				DataType:   dal.Varchar,
			})
		}
	}
	return result, nil
}
func (s *Service) checkField(table, field string) (bool, errors.Error) {
	if table == "" {
		return false, errors.Default.New("empty table name")
	}
	if !strings.HasPrefix(field, "x_") {
		return false, errors.Default.New("column name should start with `x_`")
	}
	fields, err := s.GetFields(table)
	if err != nil {
		return false, err
	}
	for _, fld := range fields {
		if fld.ColumnName == field {
			return true, nil
		}
	}
	return false, nil
}

// CreateField creates a new column for the table cf.TbName and creates a new record in the table `_tool_customized_fields`
func (s *Service) CreateField(cf *models.CustomizedField) errors.Error {
	exists, err := s.checkField(cf.TbName, cf.ColumnName)
	if err != nil {
		return err
	}
	if exists {
		return errors.BadInput.New(fmt.Sprintf("the column %s already exists", cf.ColumnName))
	}
	err = s.dal.AddColumn(cf.TbName, cf.ColumnName, cf.DataType)
	if err != nil {
		return errors.Default.Wrap(err, "AddColumn error")
	}
	err = s.dal.Create(cf)
	if err != nil {
		return errors.Default.Wrap(err, "create customizedField")
	}
	return nil
}

// DeleteField deletes the `field` form the `table`
func (s *Service) DeleteField(table, field string) errors.Error {
	exists, err := s.checkField(table, field)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}
	err = s.dal.DropColumns(table, field)
	if err != nil {
		return errors.Default.Wrap(err, "DropColumn error")
	}
	return s.dal.Delete(&models.CustomizedField{}, dal.Where("tb_name = ? AND column_name = ?", table, field))
}

func (s *Service) getCustomizedFields(table string) ([]models.CustomizedField, errors.Error) {
	var result []models.CustomizedField
	err := s.dal.All(&result, dal.Where("tb_name = ?", table))
	return result, err
}

func (s *Service) ImportCSV(table, rawDataParams string, file io.ReadCloser) errors.Error {
	var err errors.Error
	switch table {
	case "issues":
		err = s.dal.Delete(&ticket.Issue{}, dal.Where("_raw_data_params = ?", rawDataParams))
		if err != nil {
			return err
		}
		return s.importCSV(file, rawDataParams, s.issueHandler)
	case "issue_commits":
		err = s.dal.Delete(&crossdomain.IssueCommit{}, dal.Where("_raw_data_params = ?", rawDataParams))
		if err != nil {
			return err
		}
		return s.importCSV(file, rawDataParams, s.issueCommitHandler)
	}
	return errors.Default.New(fmt.Sprintf("can not import to the table %s", table))
}

func (s *Service) importCSV(file io.ReadCloser, rawDataParams string, recordHandler func(map[string]interface{}) errors.Error) errors.Error {
	iterator, err := pluginhelper.NewCsvFileIteratorFromFile(file)
	if err != nil {
		return err
	}
	var hasNext bool
	for {
		if hasNext, err = iterator.HasNextWithError(); !hasNext {
			return err
		} else {
			record := iterator.Fetch()
			record["_raw_data_params"] = rawDataParams
			for k, v := range record {
				if v.(string) == "NULL" {
					record[k] = nil
				}
			}
			err = recordHandler(record)
			if err != nil {
				return err
			}
		}
	}
}

func (s *Service) issueHandler(record map[string]interface{}) errors.Error {
	var err errors.Error
	var id string
	if record["id"] == nil {
		return errors.Default.New("record without id")
	}
	id, _ = record["id"].(string)
	if id == "" {
		return errors.Default.New("empty id")
	}
	if record["labels"] != nil {
		labels, ok := record["labels"].(string)
		if !ok {
			return errors.Default.New("labels is not string")
		}
		var issueLabels []*ticket.IssueLabel
		labelSet := make(map[string]struct{}) // for deduplicate label
		for _, label := range strings.Split(labels, ",") {
			label = strings.TrimSpace(label)
			if label == "" {
				continue
			}
			if _, exist := labelSet[label]; !exist {
				issueLabel := &ticket.IssueLabel{
					IssueId:   id,
					LabelName: label,
					NoPKModel: common.NoPKModel{
						RawDataOrigin: common.RawDataOrigin{
							RawDataParams: record["_raw_data_params"].(string),
						},
					},
				}
				issueLabels = append(issueLabels, issueLabel)
				labelSet[label] = struct{}{}
			}
		}
		if len(issueLabels) > 0 {
			err = s.dal.CreateOrUpdate(issueLabels)
			if err != nil {
				return err
			}
		}
	}
	delete(record, "labels")
	return s.dal.CreateWithMap(&ticket.Issue{}, record)
}

func (s *Service) issueCommitHandler(record map[string]interface{}) errors.Error {
	return s.dal.CreateWithMap(&crossdomain.IssueCommit{}, record)
}
