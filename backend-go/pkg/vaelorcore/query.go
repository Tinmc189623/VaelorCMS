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
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// QueryBuilder 查询构建器
// 用于构建 SQL 查询语句
type QueryBuilder struct {
	tableName string
	selects   []string
	wheres    []whereClause
	joins     []joinClause
	groupBy   []string
	having    []string
	orderBy   string
	limit     int
	offset    int
	args      []interface{}
	argIndex  int
}

// whereClause WHERE 子句
type whereClause struct {
	condition string
	args      []interface{}
	logic     string
}

// joinClause JOIN 子句
type joinClause struct {
	joinType string
	table    string
	on       string
}

// NewQueryBuilder 创建新的查询构建器
// 参数: tableName - 表名
// 返回: 查询构建器实例
func NewQueryBuilder(tableName string) *QueryBuilder {
	return &QueryBuilder{
		tableName: tableName,
		selects:   []string{"*"},
		wheres:    []whereClause{},
		joins:     []joinClause{},
		groupBy:   []string{},
		having:    []string{},
		limit:     0,
		offset:    0,
		args:      []interface{}{},
		argIndex:  0,
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
	qb.wheres = append(qb.wheres, whereClause{
		condition: qb.replacePlaceholders(condition, len(args)),
		args:      args,
		logic:     "AND",
	})
	return qb
}

// OrWhere 添加 OR WHERE 条件
// 参数: condition - 条件表达式, args - 参数值
// 返回: 查询构建器实例（支持链式调用）
func (qb *QueryBuilder) OrWhere(condition string, args ...interface{}) *QueryBuilder {
	qb.wheres = append(qb.wheres, whereClause{
		condition: qb.replacePlaceholders(condition, len(args)),
		args:      args,
		logic:     "OR",
	})
	return qb
}

// WhereIn 添加 WHERE IN 条件
// 参数: column - 列名, values - 值列表
// 返回: 查询构建器实例（支持链式调用）
func (qb *QueryBuilder) WhereIn(column string, values []interface{}) *QueryBuilder {
	if len(values) == 0 {
		return qb.Where("1 = 0")
	}
	placeholders := make([]string, len(values))
	for i := range placeholders {
		placeholders[i] = "?"
	}
	condition := fmt.Sprintf("%s IN (%s)", column, strings.Join(placeholders, ", "))
	return qb.Where(condition, values...)
}

// WhereNotIn 添加 WHERE NOT IN 条件
// 参数: column - 列名, values - 值列表
// 返回: 查询构建器实例（支持链式调用）
func (qb *QueryBuilder) WhereNotIn(column string, values []interface{}) *QueryBuilder {
	if len(values) == 0 {
		return qb
	}
	placeholders := make([]string, len(values))
	for i := range placeholders {
		placeholders[i] = "?"
	}
	condition := fmt.Sprintf("%s NOT IN (%s)", column, strings.Join(placeholders, ", "))
	return qb.Where(condition, values...)
}

// WhereNull 添加 WHERE IS NULL 条件
// 参数: column - 列名
// 返回: 查询构建器实例（支持链式调用）
func (qb *QueryBuilder) WhereNull(column string) *QueryBuilder {
	return qb.Where(fmt.Sprintf("%s IS NULL", column))
}

// WhereNotNull 添加 WHERE IS NOT NULL 条件
// 参数: column - 列名
// 返回: 查询构建器实例（支持链式调用）
func (qb *QueryBuilder) WhereNotNull(column string) *QueryBuilder {
	return qb.Where(fmt.Sprintf("%s IS NOT NULL", column))
}

// WhereBetween 添加 WHERE BETWEEN 条件
// 参数: column - 列名, min - 最小值, max - 最大值
// 返回: 查询构建器实例（支持链式调用）
func (qb *QueryBuilder) WhereBetween(column string, min, max interface{}) *QueryBuilder {
	return qb.Where(fmt.Sprintf("%s BETWEEN ? AND ?", column), min, max)
}

// WhereLike 添加 WHERE LIKE 条件
// 参数: column - 列名, pattern - 匹配模式
// 返回: 查询构建器实例（支持链式调用）
func (qb *QueryBuilder) WhereLike(column string, pattern string) *QueryBuilder {
	return qb.Where(fmt.Sprintf("%s LIKE ?", column), pattern)
}

// Join 添加 JOIN 子句
// 参数: joinType - 连接类型, table - 表名, on - 连接条件
// 返回: 查询构建器实例（支持链式调用）
func (qb *QueryBuilder) Join(joinType, table, on string) *QueryBuilder {
	qb.joins = append(qb.joins, joinClause{
		joinType: joinType,
		table:    table,
		on:       on,
	})
	return qb
}

// InnerJoin 添加 INNER JOIN 子句
// 参数: table - 表名, on - 连接条件
// 返回: 查询构建器实例（支持链式调用）
func (qb *QueryBuilder) InnerJoin(table, on string) *QueryBuilder {
	return qb.Join("INNER", table, on)
}

// LeftJoin 添加 LEFT JOIN 子句
// 参数: table - 表名, on - 连接条件
// 返回: 查询构建器实例（支持链式调用）
func (qb *QueryBuilder) LeftJoin(table, on string) *QueryBuilder {
	return qb.Join("LEFT", table, on)
}

// RightJoin 添加 RIGHT JOIN 子句
// 参数: table - 表名, on - 连接条件
// 返回: 查询构建器实例（支持链式调用）
func (qb *QueryBuilder) RightJoin(table, on string) *QueryBuilder {
	return qb.Join("RIGHT", table, on)
}

// GroupBy 设置 GROUP BY 子句
// 参数: columns - 分组字段列表
// 返回: 查询构建器实例（支持链式调用）
func (qb *QueryBuilder) GroupBy(columns ...string) *QueryBuilder {
	qb.groupBy = columns
	return qb
}

// Having 添加 HAVING 条件
// 参数: condition - 条件表达式, args - 参数值
// 返回: 查询构建器实例（支持链式调用）
func (qb *QueryBuilder) Having(condition string, args ...interface{}) *QueryBuilder {
	qb.having = append(qb.having, qb.replacePlaceholders(condition, len(args)))
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

// OrderByDesc 设置降序排序
// 参数: column - 列名
// 返回: 查询构建器实例（支持链式调用）
func (qb *QueryBuilder) OrderByDesc(column string) *QueryBuilder {
	return qb.OrderBy(fmt.Sprintf("%s DESC", column))
}

// OrderByAsc 设置升序排序
// 参数: column - 列名
// 返回: 查询构建器实例（支持链式调用）
func (qb *QueryBuilder) OrderByAsc(column string) *QueryBuilder {
	return qb.OrderBy(fmt.Sprintf("%s ASC", column))
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

// ForPage 设置分页
// 参数: page - 页码, perPage - 每页条数
// 返回: 查询构建器实例（支持链式调用）
func (qb *QueryBuilder) ForPage(page, perPage int) *QueryBuilder {
	return qb.Limit(perPage).Offset((page - 1) * perPage)
}

// Build 构建 SQL 查询语句
// 返回: SQL 语句和参数列表
func (qb *QueryBuilder) Build() (string, []interface{}) {
	var query strings.Builder
	var args []interface{}

	query.WriteString("SELECT ")
	query.WriteString(strings.Join(qb.selects, ", "))
	query.WriteString(" FROM ")
	query.WriteString(qb.tableName)

	for _, join := range qb.joins {
		query.WriteString(fmt.Sprintf(" %s JOIN %s ON %s", join.joinType, join.table, join.on))
	}

	if len(qb.wheres) > 0 {
		query.WriteString(" WHERE ")
		for i, where := range qb.wheres {
			if i > 0 {
				query.WriteString(fmt.Sprintf(" %s ", where.logic))
			}
			query.WriteString(where.condition)
			args = append(args, where.args...)
		}
	}

	if len(qb.groupBy) > 0 {
		query.WriteString(" GROUP BY ")
		query.WriteString(strings.Join(qb.groupBy, ", "))
	}

	if len(qb.having) > 0 {
		query.WriteString(" HAVING ")
		query.WriteString(strings.Join(qb.having, " AND "))
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

	return query.String(), args
}

// BuildCount 构建 COUNT 查询语句
// 返回: SQL 语句和参数列表
func (qb *QueryBuilder) BuildCount() (string, []interface{}) {
	var query strings.Builder
	var args []interface{}

	query.WriteString("SELECT COUNT(*) FROM ")
	query.WriteString(qb.tableName)

	for _, join := range qb.joins {
		query.WriteString(fmt.Sprintf(" %s JOIN %s ON %s", join.joinType, join.table, join.on))
	}

	if len(qb.wheres) > 0 {
		query.WriteString(" WHERE ")
		for i, where := range qb.wheres {
			if i > 0 {
				query.WriteString(fmt.Sprintf(" %s ", where.logic))
			}
			query.WriteString(where.condition)
			args = append(args, where.args...)
		}
	}

	if len(qb.groupBy) > 0 {
		query.WriteString(" GROUP BY ")
		query.WriteString(strings.Join(qb.groupBy, ", "))
	}

	if len(qb.having) > 0 {
		query.WriteString(" HAVING ")
		query.WriteString(strings.Join(qb.having, " AND "))
	}

	return query.String(), args
}

// BuildDelete 构建 DELETE 查询语句
// 返回: SQL 语句和参数列表
func (qb *QueryBuilder) BuildDelete() (string, []interface{}) {
	var query strings.Builder
	var args []interface{}

	query.WriteString("DELETE FROM ")
	query.WriteString(qb.tableName)

	if len(qb.wheres) > 0 {
		query.WriteString(" WHERE ")
		for i, where := range qb.wheres {
			if i > 0 {
				query.WriteString(fmt.Sprintf(" %s ", where.logic))
			}
			query.WriteString(where.condition)
			args = append(args, where.args...)
		}
	}

	return query.String(), args
}

// replacePlaceholders 替换占位符
func (qb *QueryBuilder) replacePlaceholders(condition string, count int) string {
	result := condition
	for i := 0; i < count; i++ {
		result = strings.Replace(result, "?", "?", 1)
	}
	return result
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
	wheres    []whereClause
	args      []interface{}
}

// NewUpdateBuilder 创建新的 UPDATE 构建器
// 参数: tableName - 表名
// 返回: UPDATE 构建器实例
func NewUpdateBuilder(tableName string) *UpdateBuilder {
	return &UpdateBuilder{
		tableName: tableName,
		sets:      []string{},
		wheres:    []whereClause{},
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
	ub.wheres = append(ub.wheres, whereClause{
		condition: condition,
		args:      args,
		logic:     "AND",
	})
	return ub
}

// Build 构建 UPDATE SQL 语句
// 返回: SQL 语句和参数列表
func (ub *UpdateBuilder) Build() (string, []interface{}) {
	if len(ub.sets) == 0 {
		return "", []interface{}{}
	}

	var query strings.Builder
	var args []interface{}

	query.WriteString("UPDATE ")
	query.WriteString(ub.tableName)
	query.WriteString(" SET ")
	query.WriteString(strings.Join(ub.sets, ", "))
	args = append(args, ub.args...)

	if len(ub.wheres) > 0 {
		query.WriteString(" WHERE ")
		for i, where := range ub.wheres {
			if i > 0 {
				query.WriteString(fmt.Sprintf(" %s ", where.logic))
			}
			query.WriteString(where.condition)
			args = append(args, where.args...)
		}
	}

	return query.String(), args
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
