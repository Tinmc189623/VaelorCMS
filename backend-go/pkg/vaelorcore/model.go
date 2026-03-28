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
	"time"
)

// Model 基础模型接口
// 所有数据模型都应该实现此接口
type Model interface {
	TableName() string
	GetID() int64
	SetID(id int64)
}

// BaseModel 基础模型结构体
// 包含常用的时间戳字段和 ID 字段
type BaseModel struct {
	ID        int64     `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// GetID 获取模型 ID
// 返回: 模型的 ID 值
func (bm *BaseModel) GetID() int64 {
	return bm.ID
}

// SetID 设置模型 ID
// 参数: id - 要设置的 ID 值
func (bm *BaseModel) SetID(id int64) {
	bm.ID = id
}

// TimestampsModel 带软删除的基础模型
type TimestampsModel struct {
	BaseModel
	DeletedAt sql.NullTime `json:"deleted_at" db:"deleted_at"`
}

// RelationType 关系类型定义
type RelationType string

const (
	// HasOne 一对一关系
	HasOne RelationType = "has_one"
	// HasMany 一对多关系
	HasMany RelationType = "has_many"
	// BelongsTo 属于关系
	BelongsTo RelationType = "belongs_to"
	// ManyToMany 多对多关系
	ManyToMany RelationType = "many_to_many"
)

// Relation 关系定义结构体
type Relation struct {
	// Type 关系类型
	Type RelationType
	// Field 模型中存储关联数据的字段名
	Field string
	// RelatedModel 关联的模型类型
	RelatedModel Model
	// ForeignKey 外键字段名
	ForeignKey string
	// LocalKey 本地主键字段名
	LocalKey string
	// PivotTable 多对多关系的中间表名
	PivotTable string
	// PivotLocalKey 中间表中本地键字段名
	PivotLocalKey string
	// PivotForeignKey 中间表中外键字段名
	PivotForeignKey string
}
