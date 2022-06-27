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

package e2ehelper

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/apache/incubator-devlake/config"
	"github.com/apache/incubator-devlake/helpers/pluginhelper"
	"github.com/apache/incubator-devlake/impl/dalgorm"
	"github.com/apache/incubator-devlake/logger"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/runner"
	"github.com/apache/incubator-devlake/utils"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

// DataFlowTester provides a universal data integrity validation facility to help `Plugin` verifying records between
// each step.
//
// How it works:
//
//   1. Flush specified `table` and import data from a `csv` file by `ImportCsv` method
//   2. Execute specified `subtask` by `Subtask` method
//   3. Verify actual data from specified table against expected data from another `csv` file
//   4. Repeat step 2 and 3
//
// Recommended Usage:

//   1. Create a folder under your plugin root folder. i.e. `plugins/gitlab/e2e/ to host all your e2e-tests`
//   2. Create a folder named `tables` to hold all data in `csv` format
//   3. Create e2e test-cases to cover all possible data-flow routes
//
// Example code:
//
//   See [Gitlab Project Data Flow Test](plugins/gitlab/e2e/project_test.go) for detail
//
// DataFlowTester use `N`
type DataFlowTester struct {
	Cfg    *viper.Viper
	Db     *gorm.DB
	Dal    dal.Dal
	T      *testing.T
	Name   string
	Plugin core.PluginMeta
	Log    core.Logger
}

// NewDataFlowTester create a *DataFlowTester to help developer test their subtasks data flow
func NewDataFlowTester(t *testing.T, pluginName string, pluginMeta core.PluginMeta) *DataFlowTester {
	err := core.RegisterPlugin(pluginName, pluginMeta)
	if err != nil {
		panic(err)
	}
	cfg := config.GetConfig()
	e2eDbUrl := cfg.GetString(`E2E_DB_URL`)
	if e2eDbUrl == `` {
		panic(fmt.Errorf(`e2e can only run with E2E_DB_URL, please set it in .env`))
	}
	cfg.Set(`DB_URL`, cfg.GetString(`E2E_DB_URL`))
	db, err := runner.NewGormDb(cfg, logger.Global)
	if err != nil {
		panic(err)
	}
	return &DataFlowTester{
		Cfg:    cfg,
		Db:     db,
		Dal:    dalgorm.NewDalgorm(db),
		T:      t,
		Name:   pluginName,
		Plugin: pluginMeta,
		Log:    logger.Global,
	}
}

// ImportCsvIntoRawTable imports records from specified csv file into target raw table, note that existing data would be deleted first.
func (t *DataFlowTester) ImportCsvIntoRawTable(csvRelPath string, rawTableName string) {
	csvIter := pluginhelper.NewCsvFileIterator(csvRelPath)
	defer csvIter.Close()
	t.FlushRawTable(rawTableName)
	// load rows and insert into target table
	for csvIter.HasNext() {
		toInsertValues := csvIter.Fetch()
		toInsertValues[`data`] = json.RawMessage(toInsertValues[`data`].(string))
		result := t.Db.Table(rawTableName).Create(toInsertValues)
		if result.Error != nil {
			panic(result.Error)
		}
		assert.Equal(t.T, int64(1), result.RowsAffected)
	}
}

// ImportCsvIntoTabler imports records from specified csv file into target tabler, note that existing data would be deleted first.
func (t *DataFlowTester) ImportCsvIntoTabler(csvRelPath string, dst schema.Tabler) {
	csvIter := pluginhelper.NewCsvFileIterator(csvRelPath)
	defer csvIter.Close()
	t.FlushTabler(dst)
	// load rows and insert into target table
	for csvIter.HasNext() {
		toInsertValues := csvIter.Fetch()
		for i := range toInsertValues {
			if toInsertValues[i].(string) == `` {
				toInsertValues[i] = nil
			}
		}
		result := t.Db.Model(dst).Create(toInsertValues)
		if result.Error != nil {
			panic(result.Error)
		}
		assert.Equal(t.T, int64(1), result.RowsAffected)
	}
}

// MigrateRawTableAndFlush migrate table and deletes all records from specified table
func (t *DataFlowTester) FlushRawTable(rawTableName string) {
	// flush target table
	err := t.Db.Migrator().DropTable(rawTableName)
	if err != nil {
		panic(err)
	}
	err = t.Db.Table(rawTableName).AutoMigrate(&helper.RawData{})
	if err != nil {
		panic(err)
	}
}

// FlushTabler migrate table and deletes all records from specified table
func (t *DataFlowTester) FlushTabler(dst schema.Tabler) {
	// flush target table
	err := t.Db.Migrator().DropTable(dst)
	if err != nil {
		panic(err)
	}
	err = t.Db.AutoMigrate(dst)
	if err != nil {
		panic(err)
	}
}

// Subtask executes specified subtasks
func (t *DataFlowTester) Subtask(subtaskMeta core.SubTaskMeta, taskData interface{}) {
	subtaskCtx := helper.NewStandaloneSubTaskContext(t.Cfg, t.Log, t.Db, context.Background(), t.Name, taskData)
	err := subtaskMeta.EntryPoint(subtaskCtx)
	if err != nil {
		panic(err)
	}
}

func (t *DataFlowTester) getPkFields(dst schema.Tabler) []string {
	columnTypes, err := t.Db.Migrator().ColumnTypes(dst)
	var pkFields []string
	if err != nil {
		panic(err)
	}
	for _, columnType := range columnTypes {
		if isPrimaryKey, _ := columnType.PrimaryKey(); isPrimaryKey {
			pkFields = append(pkFields, columnType.Name())
		}
	}
	return pkFields
}

// CreateSnapshot reads rows from database and write them into .csv file.
func (t *DataFlowTester) CreateSnapshot(dst schema.Tabler, csvRelPath string, targetfields []string) {
	location, _ := time.LoadLocation(`UTC`)
	pkFields := t.getPkFields(dst)
	allFields := append(pkFields, targetfields...)
	allFields = utils.StringsUniq(allFields)
	dbCursor, err := t.Dal.Cursor(
		dal.Select(strings.Join(allFields, `,`)),
		dal.From(dst.TableName()),
		dal.Orderby(strings.Join(pkFields, `,`)),
	)
	if err != nil {
		panic(err)
	}

	columns, err := dbCursor.Columns()
	if err != nil {
		panic(err)
	}
	csvWriter := pluginhelper.NewCsvFileWriter(csvRelPath, columns)
	defer csvWriter.Close()

	// define how to scan value
	columnTypes, _ := dbCursor.ColumnTypes()
	forScanValues := make([]interface{}, len(allFields))
	for i, columnType := range columnTypes {
		if columnType.ScanType().Name() == `Time` || columnType.ScanType().Name() == `NullTime` {
			forScanValues[i] = new(sql.NullTime)
		} else if columnType.ScanType().Name() == `bool` {
			forScanValues[i] = new(bool)
		} else {
			forScanValues[i] = new(string)
		}
	}

	for dbCursor.Next() {
		err = dbCursor.Scan(forScanValues...)
		if err != nil {
			panic(err)
		}
		values := make([]string, len(allFields))
		for i := range forScanValues {
			switch forScanValues[i].(type) {
			case *sql.NullTime:
				value := forScanValues[i].(*sql.NullTime)
				if value.Valid {
					values[i] = value.Time.In(location).Format("2006-01-02T15:04:05.000-07:00")
				} else {
					values[i] = ``
				}
			case *bool:
				if *forScanValues[i].(*bool) {
					values[i] = `1`
				} else {
					values[i] = `0`
				}
			case *string:
				values[i] = fmt.Sprint(*forScanValues[i].(*string))
			}
		}
		csvWriter.Write(values)
	}
}

// ExportRawTable reads rows from raw table and write them into .csv file.
func (t *DataFlowTester) ExportRawTable(rawTableName string, csvRelPath string) {
	location, _ := time.LoadLocation(`UTC`)
	allFields := []string{`id`, `params`, `data`, `url`, `input`, `created_at`}
	rawRows := &[]helper.RawData{}
	err := t.Dal.All(
		rawRows,
		dal.Select(`id, params, data, url, input, created_at`),
		dal.From(rawTableName),
		dal.Orderby(`id`),
	)
	if err != nil {
		panic(err)
	}

	csvWriter := pluginhelper.NewCsvFileWriter(csvRelPath, allFields)
	defer csvWriter.Close()

	for _, rawRow := range *rawRows {
		csvWriter.Write([]string{
			strconv.FormatUint(rawRow.ID, 10),
			rawRow.Params,
			string(rawRow.Data),
			rawRow.Url,
			string(rawRow.Input),
			rawRow.CreatedAt.In(location).Format("2006-01-02T15:04:05.000-07:00"),
		})
	}
}

func formatDbValue(value interface{}) string {
	location, _ := time.LoadLocation(`UTC`)
	switch value := value.(type) {
	case time.Time:
		return value.In(location).Format("2006-01-02T15:04:05.000-07:00")
	case bool:
		if value {
			return `1`
		} else {
			return `0`
		}
	default:
		if value != nil {
			return fmt.Sprint(value)
		}
	}
	return ``
}

// VerifyTable reads rows from csv file and compare with records from database one by one. You must specified the
// Primary Key Fields with `pkFields` so DataFlowTester could select the exact record from database, as well as which
// fields to compare with by specifying `targetFields` parameter.
func (t *DataFlowTester) VerifyTable(dst schema.Tabler, csvRelPath string, targetFields []string) {
	_, err := os.Stat(csvRelPath)
	if os.IsNotExist(err) {
		t.CreateSnapshot(dst, csvRelPath, targetFields)
		return
	}
	pkFields := t.getPkFields(dst)

	csvIter := pluginhelper.NewCsvFileIterator(csvRelPath)
	defer csvIter.Close()

	expectedTotal := int64(0)
	csvMap := map[string]map[string]interface{}{}
	for csvIter.HasNext() {
		expected := csvIter.Fetch()
		pkValues := make([]string, 0, len(pkFields))
		for _, pkf := range pkFields {
			pkValues = append(pkValues, expected[pkf].(string))
		}
		pkValueStr := strings.Join(pkValues, `-`)
		_, ok := csvMap[pkValueStr]
		assert.False(t.T, ok, fmt.Sprintf(`%s duplicated in csv (with params from csv %s)`, dst.TableName(), pkValues))
		csvMap[pkValueStr] = expected
	}

	dbRows := &[]map[string]interface{}{}
	err = t.Db.Table(dst.TableName()).Find(dbRows).Error
	if err != nil {
		panic(err)
	}
	for _, actual := range *dbRows {
		pkValues := make([]string, 0, len(pkFields))
		for _, pkf := range pkFields {
			pkValues = append(pkValues, formatDbValue(actual[pkf]))
		}
		expected, ok := csvMap[strings.Join(pkValues, `-`)]
		assert.True(t.T, ok, fmt.Sprintf(`%s not found (with params from csv %s)`, dst.TableName(), pkValues))
		if !ok {
			continue
		}
		for _, field := range targetFields {
			assert.Equal(t.T, expected[field], formatDbValue(actual[field]), fmt.Sprintf(`%s.%s not match (with params from csv %s)`, dst.TableName(), field, pkValues))
		}
		expectedTotal++
	}

	var actualTotal int64
	err = t.Db.Table(dst.TableName()).Count(&actualTotal).Error
	if err != nil {
		panic(err)
	}
	assert.Equal(t.T, expectedTotal, actualTotal, fmt.Sprintf(`%s count not match`, dst.TableName()))
}
