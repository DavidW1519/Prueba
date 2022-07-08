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

package dalgorm

import (
	"database/sql"
	"reflect"
	"strings"

	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// Dalgorm FIXME ...
type Dalgorm struct {
	db *gorm.DB
}

func buildTx(tx *gorm.DB, clauses []dal.Clause) *gorm.DB {
	for _, c := range clauses {
		t := c.Type
		d := c.Data
		switch t {
		case dal.JoinClause:
			tx = tx.Joins(d.(dal.DalClause).Expr, d.(dal.DalClause).Params...)
		case dal.WhereClause:
			tx = tx.Where(d.(dal.DalClause).Expr, d.(dal.DalClause).Params...)
		case dal.OrderbyClause:
			tx = tx.Order(d.(string))
		case dal.LimitClause:
			tx = tx.Limit(d.(int))
		case dal.OffsetClause:
			tx = tx.Offset(d.(int))
		case dal.FromClause:
			if str, ok := d.(string); ok {
				tx = tx.Table(str)
			} else {
				tx = tx.Model(d)
			}
		case dal.SelectClause:
			tx = tx.Select(d.(string))
		case dal.GroupbyClause:
			tx = tx.Group(d.(string))
		case dal.HavingClause:
			tx = tx.Having(d.(dal.DalClause).Expr, d.(dal.DalClause).Params...)
		}
	}
	return tx
}

var _ dal.Dal = (*Dalgorm)(nil)

// RawCursor executes raw sql query and returns a database cursor
func (d *Dalgorm) RawCursor(query string, params ...interface{}) (*sql.Rows, error) {
	return d.db.Raw(query, params...).Rows()
}

// Exec executes raw sql query
func (d *Dalgorm) Exec(query string, params ...interface{}) error {
	return d.db.Exec(query, params...).Error
}

// AutoMigrate runs auto migration for given models
func (d *Dalgorm) AutoMigrate(entity interface{}, clauses ...dal.Clause) error {
	return buildTx(d.db, clauses).AutoMigrate(entity)
}

// Cursor returns a database cursor, cursor is especially useful when handling big amount of rows of data
func (d *Dalgorm) Cursor(clauses ...dal.Clause) (*sql.Rows, error) {
	return buildTx(d.db, clauses).Rows()
}

// CursorTx FIXME ...
func (d *Dalgorm) CursorTx(clauses ...dal.Clause) *gorm.DB {
	return buildTx(d.db, clauses)
}

// Fetch loads row data from `cursor` into `dst`
func (d *Dalgorm) Fetch(cursor *sql.Rows, dst interface{}) error {
	return d.db.ScanRows(cursor, dst)
}

// All loads matched rows from database to `dst`, USE IT WITH COUTIOUS!!
func (d *Dalgorm) All(dst interface{}, clauses ...dal.Clause) error {
	return buildTx(d.db, clauses).Find(dst).Error
}

// First loads first matched row from database to `dst`, error will be returned if no records were found
func (d *Dalgorm) First(dst interface{}, clauses ...dal.Clause) error {
	err := buildTx(d.db, clauses).First(dst).Error
	return err
}

// Count total records
func (d *Dalgorm) Count(clauses ...dal.Clause) (int64, error) {
	var count int64
	err := buildTx(d.db, clauses).Count(&count).Error
	return count, err
}

// Pluck used to query single column
func (d *Dalgorm) Pluck(column string, dest interface{}, clauses ...dal.Clause) error {
	return buildTx(d.db, clauses).Pluck(column, dest).Error
}

// Create insert record to database
func (d *Dalgorm) Create(entity interface{}, clauses ...dal.Clause) error {
	return buildTx(d.db, clauses).Create(entity).Error
}

// Update updates record
func (d *Dalgorm) Update(entity interface{}, clauses ...dal.Clause) error {
	return buildTx(d.db, clauses).Save(entity).Error
}

// CreateOrUpdate tries to create the record, or fallback to update all if failed
func (d *Dalgorm) CreateOrUpdate(entity interface{}, clauses ...dal.Clause) error {
	return buildTx(d.db, clauses).Clauses(clause.OnConflict{UpdateAll: true}).Create(entity).Error
}

// CreateIfNotExist tries to create the record if not exist
func (d *Dalgorm) CreateIfNotExist(entity interface{}, clauses ...dal.Clause) error {
	return buildTx(d.db, clauses).Clauses(clause.OnConflict{DoNothing: true}).Create(entity).Error
}

// Delete records from database
func (d *Dalgorm) Delete(entity interface{}, clauses ...dal.Clause) error {
	return buildTx(d.db, clauses).Delete(entity).Error
}

// GetColumns FIXME ...
func (d *Dalgorm) GetColumns(dst schema.Tabler, filter func(columnMeta dal.ColumnMeta) bool) (cms []dal.ColumnMeta, err error) {
	columnTypes, err := d.db.Migrator().ColumnTypes(dst.TableName())
	if err != nil {
		return
	}
	for _, columnType := range columnTypes {
		if filter == nil {
			cms = append(cms, columnType)
		} else if filter(columnType) {
			cms = append(cms, columnType)
		}
	}
	return
}

// GetPrimaryKeyFields get the PrimaryKey from `gorm` tag
func (d *Dalgorm) GetPrimaryKeyFields(t reflect.Type) []reflect.StructField {
	return utils.WalkFields(t, func(field *reflect.StructField) bool {
		return strings.Contains(strings.ToLower(field.Tag.Get("gorm")), "primarykey")
	})
}

// AllTables returns all tables in the database
func (d *Dalgorm) AllTables() ([]string, error) {
	var tableSql string
	if d.db.Dialector.Name() == "mysql" {
		tableSql = "show tables"
	} else {
		tableSql = "select table_name from information_schema.tables where table_schema = 'public' and table_name not like '_devlake%'"
	}
	var tables []string
	err := d.db.Raw(tableSql).Scan(&tables).Error
	if err != nil {
		return nil, err
	}
	var filteredTables []string
	for _, table := range tables {
		if !strings.HasPrefix(table, "_devlake") {
			filteredTables = append(filteredTables, table)
		}
	}
	return filteredTables, nil
}

// NewDalgorm FIXME ...
func NewDalgorm(db *gorm.DB) *Dalgorm {
	return &Dalgorm{db}
}
