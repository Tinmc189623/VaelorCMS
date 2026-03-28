/*
 * VaelorCMS - 现代化内容管理系统
 *
 * Copyright © 2025-2026 Nexsteaduser. All rights reserved.
 *
 * 作者: Tinmc189623
 * 团队: Nexsteaduser
 *
 * 本程序是自由软件: 你可以重新分发和/或修改
 * 它在 GNU Affero 通用公共许可证的条款下发布,
 * 版本 3 或 (根据你的选择) 任何更高版本。
 *
 * 本程序是希望它有用,
 * 但没有任何保证; 甚至没有适销性或
 * 特定用途的默示保证。见
 * GNU Affero 通用公共许可证获取更多细节。
 *
 * 你应该收到 GNU Affero 通用公共许可证的副本
 * 与此程序一起。如果没有, 请见 <https://www.gnu.org/licenses/>.
 */

package vaelorcore

import (
	"database/sql"
	"reflect"
	"strings"
	"time"
)

// DBExecutor 数据库执行器接口
// 支持普通连接和事务
type DBExecutor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// ORM 主结构体
type ORM struct {
	db DBExecutor
}

// NewORM 创建新的 ORM 实例
// 参数: db - 数据库执行器
// 返回: ORM 实例
func NewORM(db DBExecutor) *ORM {
	return &ORM{db: db}
}

// SetDB 设置数据库执行器
// 参数: db - 数据库执行器
func (o *ORM) SetDB(db DBExecutor) {
	o.db = db
}

// Find 根据 ID 查找记录
// 参数: model - 模型实例（指针）, id - 记录 ID
// 返回: 是否找到和可能的错误
func (o *ORM) Find(model Model, id int64) (bool, error) {
	qb := NewQueryBuilder(model.TableName()).Where("id = ?", id)
	query, args := qb.Build()

	row := o.db.QueryRow(query, args...)
	err := ScanRow(row, model)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// FindFirst 查找第一条记录
// 参数: model - 模型实例（指针）, qb - 查询构建器（可选）
// 返回: 是否找到和可能的错误
func (o *ORM) FindFirst(model Model, qb ...*QueryBuilder) (bool, error) {
	var builder *QueryBuilder
	if len(qb) > 0 && qb[0] != nil {
		builder = qb[0]
	} else {
		builder = NewQueryBuilder(model.TableName())
	}
	builder.Limit(1)

	query, args := builder.Build()

	row := o.db.QueryRow(query, args...)
	err := ScanRow(row, model)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// FindAll 查询所有记录
// 参数: models - 模型切片指针, qb - 查询构建器（可选）
// 返回: 可能的错误
func (o *ORM) FindAll(models interface{}, qb ...*QueryBuilder) error {
	sliceValue := reflect.ValueOf(models)
	if sliceValue.Kind() != reflect.Ptr {
		return ErrMustBePointer
	}
	sliceElem := sliceValue.Elem()
	if sliceElem.Kind() != reflect.Slice {
		return ErrMustBeSlice
	}

	elemType := sliceElem.Type().Elem()
	isPtr := elemType.Kind() == reflect.Ptr
	if isPtr {
		elemType = elemType.Elem()
	}

	var tableName string
	if elemType.Kind() == reflect.Struct {
		tempModel := reflect.New(elemType).Interface()
		if model, ok := tempModel.(Model); ok {
			tableName = model.TableName()
		}
	}

	var builder *QueryBuilder
	if len(qb) > 0 && qb[0] != nil {
		builder = qb[0]
	} else {
		if tableName == "" {
			return ErrInvalidModel
		}
		builder = NewQueryBuilder(tableName)
	}

	query, args := builder.Build()

	rows, err := o.db.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	return ScanRows(rows, models)
}

// Create 创建新记录
// 参数: model - 模型实例（指针）
// 返回: 插入的 ID 和可能的错误
func (o *ORM) Create(model Model) (int64, error) {
	columns, values, err := o.getModelColumnsAndValues(model, true)
	if err != nil {
		return 0, err
	}

	ib := NewInsertBuilder(model.TableName()).Columns(columns...).Values(values...)
	query, args := ib.Build()

	result, err := o.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err == nil {
		model.SetID(id)
	}
	return id, err
}

// Update 更新记录
// 参数: model - 模型实例（指针）
// 返回: 影响的行数和可能的错误
func (o *ORM) Update(model Model) (int64, error) {
	if model.GetID() == 0 {
		return 0, ErrIDRequired
	}

	columns, values, err := o.getModelColumnsAndValues(model, false)
	if err != nil {
		return 0, err
	}

	ub := NewUpdateBuilder(model.TableName())
	for i, col := range columns {
		ub.Set(col, values[i])
	}
	ub.Where("id = ?", model.GetID())

	query, args := ub.Build()
	result, err := o.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// Save 保存记录（根据是否有 ID 决定是创建还是更新）
// 参数: model - 模型实例（指针）
// 返回: 可能的错误
func (o *ORM) Save(model Model) error {
	if model.GetID() == 0 {
		_, err := o.Create(model)
		return err
	}
	_, err := o.Update(model)
	return err
}

// Delete 删除记录
// 参数: model - 模型实例（指针）
// 返回: 影响的行数和可能的错误
func (o *ORM) Delete(model Model) (int64, error) {
	if model.GetID() == 0 {
		return 0, ErrIDRequired
	}

	qb := NewQueryBuilder(model.TableName()).Where("id = ?", model.GetID())
	query, args := qb.BuildDelete()

	result, err := o.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// DeleteByID 根据 ID 删除记录
// 参数: model - 模型实例（仅用于获取表名）, id - 记录 ID
// 返回: 影响的行数和可能的错误
func (o *ORM) DeleteByID(model Model, id int64) (int64, error) {
	qb := NewQueryBuilder(model.TableName()).Where("id = ?", id)
	query, args := qb.BuildDelete()

	result, err := o.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// Count 统计记录数
// 参数: model - 模型实例, qb - 查询构建器（可选）
// 返回: 记录数和可能的错误
func (o *ORM) Count(model Model, qb ...*QueryBuilder) (int64, error) {
	var builder *QueryBuilder
	if len(qb) > 0 && qb[0] != nil {
		builder = qb[0]
	} else {
		builder = NewQueryBuilder(model.TableName())
	}

	query, args := builder.BuildCount()

	var count int64
	err := o.db.QueryRow(query, args...).Scan(&count)
	return count, err
}

// Exists 检查记录是否存在
// 参数: model - 模型实例, qb - 查询构建器（可选）
// 返回: 是否存在和可能的错误
func (o *ORM) Exists(model Model, qb ...*QueryBuilder) (bool, error) {
	count, err := o.Count(model, qb...)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Pluck 获取单个列的值列表
// 参数: model - 模型实例, column - 列名, dest - 目标切片指针, qb - 查询构建器（可选）
// 返回: 可能的错误
func (o *ORM) Pluck(model Model, column string, dest interface{}, qb ...*QueryBuilder) error {
	var builder *QueryBuilder
	if len(qb) > 0 && qb[0] != nil {
		builder = qb[0]
	} else {
		builder = NewQueryBuilder(model.TableName())
	}
	builder.Select(column)

	query, args := builder.Build()

	rows, err := o.db.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Ptr {
		return ErrMustBePointer
	}
	destSlice := destValue.Elem()
	if destSlice.Kind() != reflect.Slice {
		return ErrMustBeSlice
	}

	for rows.Next() {
		elem := reflect.New(destSlice.Type().Elem()).Elem()
		if err := rows.Scan(elem.Addr().Interface()); err != nil {
			return err
		}
		destSlice.Set(reflect.Append(destSlice, elem))
	}

	return rows.Err()
}

// getModelColumnsAndValues 获取模型的列和值
func (o *ORM) getModelColumnsAndValues(model Model, includeID bool) ([]string, []interface{}, error) {
	modelValue := reflect.ValueOf(model)
	if modelValue.Kind() != reflect.Ptr {
		return nil, nil, ErrMustBePointer
	}
	modelElem := modelValue.Elem()
	if modelElem.Kind() != reflect.Struct {
		return nil, nil, ErrMustBeStruct
	}

	var columns []string
	var values []interface{}
	now := time.Now()

	for i := 0; i < modelElem.NumField(); i++ {
		field := modelElem.Type().Field(i)
		fieldValue := modelElem.Field(i)

		if !field.IsExported() {
			continue
		}

		dbTag := field.Tag.Get("db")
		if dbTag == "-" {
			continue
		}

		columnName := dbTag
		if columnName == "" {
			columnName = toSnakeCase(field.Name)
		}

		if columnName == "id" && !includeID {
			continue
		}

		if columnName == "created_at" && includeID {
			if fieldValue.IsZero() {
				fieldValue.Set(reflect.ValueOf(now))
			}
		}
		if columnName == "updated_at" {
			fieldValue.Set(reflect.ValueOf(now))
		}

		columns = append(columns, columnName)
		values = append(values, fieldValue.Interface())
	}

	return columns, values, nil
}

// toSnakeCase 将驼峰命名转换为下划线命名
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, c := range s {
		if i > 0 && 'A' <= c && c <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune('a' + c - 'A')
	}
	return result.String()
}
