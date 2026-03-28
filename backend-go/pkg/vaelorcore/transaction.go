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
)

// Transaction 事务封装
type Transaction struct {
	tx     *sql.Tx
	closed bool
	orm    *ORM
}

// NewTransaction 创建新的事务实例
// 参数: tx - 数据库事务
// 返回: 事务实例
func NewTransaction(tx *sql.Tx) *Transaction {
	return &Transaction{
		tx:     tx,
		closed: false,
		orm:    NewORM(tx),
	}
}

// GetORM 获取事务内的 ORM 实例
// 返回: ORM 实例
func (t *Transaction) GetORM() *ORM {
	return t.orm
}

// Commit 提交事务
// 返回: 可能的错误
func (t *Transaction) Commit() error {
	if t.closed {
		return ErrTxAlreadyClosed
	}
	err := t.tx.Commit()
	t.closed = true
	return err
}

// Rollback 回滚事务
// 返回: 可能的错误
func (t *Transaction) Rollback() error {
	if t.closed {
		return ErrTxAlreadyClosed
	}
	err := t.tx.Rollback()
	t.closed = true
	return err
}

// IsClosed 检查事务是否已关闭
// 返回: 是否已关闭
func (t *Transaction) IsClosed() bool {
	return t.closed
}

// TransactionManager 事务管理器
type TransactionManager struct {
	db *sql.DB
}

// NewTransactionManager 创建新的事务管理器
// 参数: db - 数据库连接
// 返回: 事务管理器实例
func NewTransactionManager(db *sql.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

// Begin 开始一个新事务
// 返回: 事务实例和可能的错误
func (tm *TransactionManager) Begin() (*Transaction, error) {
	tx, err := tm.db.Begin()
	if err != nil {
		return nil, err
	}
	return NewTransaction(tx), nil
}

// Transaction 在事务中执行函数
// 参数: fn - 要执行的函数，接收事务内的 ORM 实例
// 返回: 可能的错误
func (tm *TransactionManager) Transaction(fn func(*ORM) error) error {
	tx, err := tm.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(tx.GetORM()); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// TransactionWithResult 在事务中执行函数并返回结果
// 参数: fn - 要执行的函数，接收事务内的 ORM 实例，返回结果和错误
// 返回: 结果和可能的错误
func (tm *TransactionManager) TransactionWithResult(fn func(*ORM) (interface{}, error)) (interface{}, error) {
	tx, err := tm.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	result, err := fn(tx.GetORM())
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

// Savepoint 创建保存点
// 参数: name - 保存点名称
// 返回: 可能的错误
func (t *Transaction) Savepoint(name string) error {
	if t.closed {
		return ErrTxAlreadyClosed
	}
	_, err := t.tx.Exec("SAVEPOINT " + name)
	return err
}

// RollbackTo 回滚到保存点
// 参数: name - 保存点名称
// 返回: 可能的错误
func (t *Transaction) RollbackTo(name string) error {
	if t.closed {
		return ErrTxAlreadyClosed
	}
	_, err := t.tx.Exec("ROLLBACK TO " + name)
	return err
}

// ReleaseSavepoint 释放保存点
// 参数: name - 保存点名称
// 返回: 可能的错误
func (t *Transaction) ReleaseSavepoint(name string) error {
	if t.closed {
		return ErrTxAlreadyClosed
	}
	_, err := t.tx.Exec("RELEASE SAVEPOINT " + name)
	return err
}
