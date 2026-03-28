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

import "errors"

// ORM 错误定义
var (
	ErrMustBePointer    = errors.New("目标必须是指针类型")
	ErrMustBeSlice      = errors.New("目标必须是切片指针")
	ErrMustBeStruct     = errors.New("目标必须是结构体指针")
	ErrInvalidModel     = errors.New("无效的模型类型")
	ErrIDRequired       = errors.New("模型 ID 不能为空")
	ErrTxNotStarted     = errors.New("事务未开始")
	ErrTxAlreadyClosed  = errors.New("事务已关闭")
	ErrRelationNotFound = errors.New("关系未找到")
	ErrInvalidRelation  = errors.New("无效的关系类型")
)
