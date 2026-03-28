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

package database

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Model 基础模型接口
// 所有数据模型都应该实现此接口
type Model interface {
	TableName() string
}

// QueryBuilder 轻量级查询构建器
// 用于构建 SQL 查询语句
type QueryBuilder struct {
	tableName string
	selects   []string
	wheres    []string
	args      []interface{}
	orderBy   string
	limit     int
	offset    int
}

// NewQueryBuilder 创建新的查询构建器
// 参数: tableName - 表名
// 返回: 查询构建器实例
func NewQueryBuilder(tableName string) *QueryBuilder {
	return &QueryBuilder{
		tableName: tableName,
		selects:   []string{"*"},
		wheres:    []string{},
		args:      []interface{}{},
		limit:     0,
		offset:    0,
	}
}

// Select 设置查询字段
// 参数: fields - 字段列表
// 返回: 查询构建器实例（支持链式调用）
func (qb *QueryBuilder) Select(fields ...string) *QueryBuilder {
	qb.selects = fields
	return qb
}

// Where 添加 WHERE 条件（使用 AND 连接）
// 参数: condition - 条件表达式, args - 参数值
// 返回: 查询构建器实例（支持链式调用）
func (qb *QueryBuilder) Where(condition string, args ...interface{}) *QueryBuilder {
	qb.wheres = append(qb.wheres, condition)
	qb.args = append(qb.args, args...)
	return qb
}

// OrderBy 设置排序
// 参数: order - 排序表达式
// 返回: 查询构建器实例（支持链式调用）
func (qb *QueryBuilder) OrderBy(order string) *QueryBuilder {
	qb.orderBy = order
	return qb
}

// Limit 设置限制行数
// 参数: limit - 限制行数
// 返回: 查询构建器实例（支持链式调用）
func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.limit = limit
	return qb
}

// Offset 设置偏移量
// 参数: offset - 偏移量
// 返回: 查询构建器实例（支持链式调用）
func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.offset = offset
	return qb
}

// Build 构建 SQL 查询语句
// 返回: SQL 语句和参数列表
func (qb *QueryBuilder) Build() (string, []interface{}) {
	var query strings.Builder

	query.WriteString("SELECT ")
	query.WriteString(strings.Join(qb.selects, ", "))
	query.WriteString(" FROM ")
	query.WriteString(qb.tableName)

	if len(qb.wheres) > 0 {
		query.WriteString(" WHERE ")
		query.WriteString(strings.Join(qb.wheres, " AND "))
	}

	if qb.orderBy != "" {
		query.WriteString(" ORDER BY ")
		query.WriteString(qb.orderBy)
	}

	if qb.limit > 0 {
		query.WriteString(" LIMIT ")
		query.WriteString(strconv.Itoa(qb.limit))
	}

	if qb.offset > 0 {
		query.WriteString(" OFFSET ")
		query.WriteString(strconv.Itoa(qb.offset))
	}

	return query.String(), qb.args
}

// BuildCount 构建 COUNT 查询语句
// 返回: SQL 语句和参数列表
func (qb *QueryBuilder) BuildCount() (string, []interface{}) {
	var query strings.Builder

	query.WriteString("SELECT COUNT(*) FROM ")
	query.WriteString(qb.tableName)

	if len(qb.wheres) > 0 {
		query.WriteString(" WHERE ")
		query.WriteString(strings.Join(qb.wheres, " AND "))
	}

	return query.String(), qb.args
}

// InsertBuilder INSERT 语句构建器
type InsertBuilder struct {
	tableName string
	columns   []string
	values    [][]interface{}
}

// NewInsertBuilder 创建新的 INSERT 构建器
// 参数: tableName - 表名
// 返回: INSERT 构建器实例
func NewInsertBuilder(tableName string) *InsertBuilder {
	return &InsertBuilder{
		tableName: tableName,
		columns:   []string{},
		values:    [][]interface{}{},
	}
}

// Columns 设置插入列
// 参数: columns - 列名列表
// 返回: INSERT 构建器实例（支持链式调用）
func (ib *InsertBuilder) Columns(columns ...string) *InsertBuilder {
	ib.columns = columns
	return ib
}

// Values 添加一行数据
// 参数: values - 值列表
// 返回: INSERT 构建器实例（支持链式调用）
func (ib *InsertBuilder) Values(values ...interface{}) *InsertBuilder {
	ib.values = append(ib.values, values)
	return ib
}

// Build 构建 INSERT SQL 语句
// 返回: SQL 语句和参数列表
func (ib *InsertBuilder) Build() (string, []interface{}) {
	if len(ib.columns) == 0 || len(ib.values) == 0 {
		return "", []interface{}{}
	}

	var query strings.Builder
	var args []interface{}

	query.WriteString("INSERT INTO ")
	query.WriteString(ib.tableName)
	query.WriteString(" (")
	query.WriteString(strings.Join(ib.columns, ", "))
	query.WriteString(") VALUES ")

	placeholders := make([]string, len(ib.columns))
	for i := range placeholders {
		placeholders[i] = "?"
	}

	valueStr := "(" + strings.Join(placeholders, ", ") + ")"
	valueParts := make([]string, len(ib.values))
	for i, vals := range ib.values {
		valueParts[i] = valueStr
		args = append(args, vals...)
	}

	query.WriteString(strings.Join(valueParts, ", "))

	return query.String(), args
}

// UpdateBuilder UPDATE 语句构建器
type UpdateBuilder struct {
	tableName string
	sets      []string
	wheres    []string
	args      []interface{}
}

// NewUpdateBuilder 创建新的 UPDATE 构建器
// 参数: tableName - 表名
// 返回: UPDATE 构建器实例
func NewUpdateBuilder(tableName string) *UpdateBuilder {
	return &UpdateBuilder{
		tableName: tableName,
		sets:      []string{},
		wheres:    []string{},
		args:      []interface{}{},
	}
}

// Set 设置更新字段和值
// 参数: column - 列名, value - 值
// 返回: UPDATE 构建器实例（支持链式调用）
func (ub *UpdateBuilder) Set(column string, value interface{}) *UpdateBuilder {
	ub.sets = append(ub.sets, column+" = ?")
	ub.args = append(ub.args, value)
	return ub
}

// Where 添加 WHERE 条件
// 参数: condition - 条件表达式, args - 参数值
// 返回: UPDATE 构建器实例（支持链式调用）
func (ub *UpdateBuilder) Where(condition string, args ...interface{}) *UpdateBuilder {
	ub.wheres = append(ub.wheres, condition)
	ub.args = append(ub.args, args...)
	return ub
}

// Build 构建 UPDATE SQL 语句
// 返回: SQL 语句和参数列表
func (ub *UpdateBuilder) Build() (string, []interface{}) {
	if len(ub.sets) == 0 {
		return "", []interface{}{}
	}

	var query strings.Builder

	query.WriteString("UPDATE ")
	query.WriteString(ub.tableName)
	query.WriteString(" SET ")
	query.WriteString(strings.Join(ub.sets, ", "))

	if len(ub.wheres) > 0 {
		query.WriteString(" WHERE ")
		query.WriteString(strings.Join(ub.wheres, " AND "))
	}

	return query.String(), ub.args
}

// DeleteBuilder DELETE 语句构建器
type DeleteBuilder struct {
	tableName string
	wheres    []string
	args      []interface{}
}

// NewDeleteBuilder 创建新的 DELETE 构建器
// 参数: tableName - 表名
// 返回: DELETE 构建器实例
func NewDeleteBuilder(tableName string) *DeleteBuilder {
	return &DeleteBuilder{
		tableName: tableName,
		wheres:    []string{},
		args:      []interface{}{},
	}
}

// Where 添加 WHERE 条件
// 参数: condition - 条件表达式, args - 参数值
// 返回: DELETE 构建器实例（支持链式调用）
func (db *DeleteBuilder) Where(condition string, args ...interface{}) *DeleteBuilder {
	db.wheres = append(db.wheres, condition)
	db.args = append(db.args, args...)
	return db
}

// Build 构建 DELETE SQL 语句
// 返回: SQL 语句和参数列表
func (db *DeleteBuilder) Build() (string, []interface{}) {
	var query strings.Builder

	query.WriteString("DELETE FROM ")
	query.WriteString(db.tableName)

	if len(db.wheres) > 0 {
		query.WriteString(" WHERE ")
		query.WriteString(strings.Join(db.wheres, " AND "))
	}

	return query.String(), db.args
}

// FindByID 根据 ID 查找记录
// 参数: tableName - 表名, id - 记录 ID, dest - 目标对象指针
// 返回: 是否找到和可能的错误
func FindByID(tableName string, id int64, dest interface{}) (bool, error) {
	qb := NewQueryBuilder(tableName).Where("id = ?", id)
	query, args := qb.Build()

	row := QueryRow(query, args...)
	err := ScanRow(row, dest)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// FindAll 查询所有记录
// 参数: tableName - 表名, dest - 目标切片指针
// 返回: 可能的错误
func FindAll(tableName string, dest interface{}) error {
	qb := NewQueryBuilder(tableName)
	query, args := qb.Build()

	rows, err := Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	return ScanRows(rows, dest)
}

// Count 统计记录数
// 参数: tableName - 表名
// 返回: 记录数和可能的错误
func Count(tableName string) (int64, error) {
	qb := NewQueryBuilder(tableName)
	query, args := qb.BuildCount()

	var count int64
	err := QueryRow(query, args...).Scan(&count)
	return count, err
}

// Create 创建新记录
// 参数: tableName - 表名, columns - 列名, values - 值
// 返回: 插入的 ID 和可能的错误
func Create(tableName string, columns []string, values ...interface{}) (int64, error) {
	ib := NewInsertBuilder(tableName).Columns(columns...).Values(values...)
	query, args := ib.Build()

	result, err := Exec(query, args...)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// Update 更新记录
// 参数: tableName - 表名, id - 记录 ID, updates - 更新字段映射
// 返回: 影响的行数和可能的错误
func Update(tableName string, id int64, updates map[string]interface{}) (int64, error) {
	ub := NewUpdateBuilder(tableName)
	for col, val := range updates {
		ub.Set(col, val)
	}
	ub.Where("id = ?", id)

	query, args := ub.Build()
	result, err := Exec(query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// Delete 删除记录
// 参数: tableName - 表名, id - 记录 ID
// 返回: 影响的行数和可能的错误
func Delete(tableName string, id int64) (int64, error) {
	db := NewDeleteBuilder(tableName).Where("id = ?", id)
	query, args := db.Build()

	result, err := Exec(query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// ScanRow 扫描单行结果到结构体
// 参数: row - 行结果, dest - 目标对象指针
// 返回: 可能的错误
func ScanRow(row *sql.Row, dest interface{}) error {
	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Ptr {
		return fmt.Errorf("目标必须是指针类型")
	}

	destElem := destValue.Elem()
	if destElem.Kind() != reflect.Struct {
		return fmt.Errorf("目标必须是结构体指针")
	}

	fields := getStructFields(destElem.Type())
	ptrs := make([]interface{}, len(fields))
	for i, field := range fields {
		ptrs[i] = destElem.FieldByName(field.Name).Addr().Interface()
	}

	return row.Scan(ptrs...)
}

// ScanRows 扫描多行结果到切片
// 参数: rows - 行结果集, dest - 目标切片指针
// 返回: 可能的错误
func ScanRows(rows *sql.Rows, dest interface{}) error {
	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Ptr {
		return fmt.Errorf("目标必须是指针类型")
	}

	destSlice := destValue.Elem()
	if destSlice.Kind() != reflect.Slice {
		return fmt.Errorf("目标必须是切片指针")
	}

	elemType := destSlice.Type().Elem()
	isPtr := elemType.Kind() == reflect.Ptr
	if isPtr {
		elemType = elemType.Elem()
	}

	if elemType.Kind() != reflect.Struct {
		return fmt.Errorf("切片元素必须是结构体或结构体指针")
	}

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	for rows.Next() {
		elem := reflect.New(elemType).Elem()
		ptrs := make([]interface{}, len(columns))
		for i, col := range columns {
			field := elem.FieldByName(toCamelCase(col))
			if field.IsValid() && field.CanSet() {
				ptrs[i] = field.Addr().Interface()
			} else {
				var discard interface{}
				ptrs[i] = &discard
			}
		}

		if err := rows.Scan(ptrs...); err != nil {
			return err
		}

		if isPtr {
			destSlice.Set(reflect.Append(destSlice, elem.Addr()))
		} else {
			destSlice.Set(reflect.Append(destSlice, elem))
		}
	}

	return rows.Err()
}

// getStructFields 获取结构体字段信息
func getStructFields(t reflect.Type) []reflect.StructField {
	var fields []reflect.StructField
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.IsExported() {
			fields = append(fields, field)
		}
	}
	return fields
}

// toCamelCase 将下划线命名转换为驼峰命名
func toCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i := range parts {
		if i == 0 {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		} else {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, "")
}
