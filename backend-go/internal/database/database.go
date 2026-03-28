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
	"log"
	"strings"
	"sync"
	"time"

	_ "modernc.org/sqlite"
	"vaelorcms/internal/config"
)

// 全局数据库连接池
var (
	dbPool *sql.DB
	once   sync.Once
)

// InitDB 初始化数据库连接池
// 该函数会创建数据库连接池并启用外键约束
func InitDB(cfg *config.Config) error {
	var initErr error
	once.Do(func() {
		initErr = initDatabasePool(cfg)
	})
	return initErr
}

// initDatabasePool 内部函数，初始化数据库连接池
func initDatabasePool(cfg *config.Config) error {
	dbPath := parseDatabaseURL(cfg.DatabaseURL)
	
	var err error
	dbPool, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("打开数据库失败: %w", err)
	}

	// 配置连接池
	dbPool.SetMaxOpenConns(cfg.DBPoolSize)
	dbPool.SetMaxIdleConns(cfg.DBPoolSize)
	dbPool.SetConnMaxLifetime(time.Duration(cfg.DBPoolRecycle) * time.Second)
	dbPool.SetConnMaxIdleTime(30 * time.Minute)

	// 测试连接
	if err = dbPool.Ping(); err != nil {
		return fmt.Errorf("数据库连接失败: %w", err)
	}

	// 启用外键约束
	if err = enableForeignKeys(); err != nil {
		return fmt.Errorf("启用外键约束失败: %w", err)
	}

	log.Println("数据库连接池初始化成功")
	return nil
}

// parseDatabaseURL 解析数据库 URL，提取 SQLite 文件路径
func parseDatabaseURL(databaseURL string) string {
	if strings.HasPrefix(databaseURL, "sqlite:///") {
		return strings.TrimPrefix(databaseURL, "sqlite:///")
	}
	if strings.HasPrefix(databaseURL, "sqlite://") {
		return strings.TrimPrefix(databaseURL, "sqlite://")
	}
	return databaseURL
}

// enableForeignKeys 启用 SQLite 外键约束
func enableForeignKeys() error {
	_, err := dbPool.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return err
	}
	log.Println("SQLite 外键约束已启用")
	return nil
}

// GetDB 获取数据库连接
// 返回一个可用于数据库操作的 *sql.DB 实例
func GetDB() *sql.DB {
	if dbPool == nil {
		log.Fatal("数据库连接未初始化，请先调用 InitDB()")
	}
	return dbPool
}

// CloseDB 关闭数据库连接池
// 在应用程序退出时调用
func CloseDB() error {
	if dbPool != nil {
		log.Println("正在关闭数据库连接池...")
		return dbPool.Close()
	}
	return nil
}

// GetConnection 获取一个数据库连接
// 使用完毕后必须调用 conn.Close() 释放连接
func GetConnection() (*sql.Conn, error) {
	if dbPool == nil {
		return nil, fmt.Errorf("数据库连接未初始化")
	}
	return dbPool.Conn(nil)
}

// ExecContext 执行不返回行的 SQL 语句
// 参数: query - SQL 语句, args - 参数
// 返回: 影响的行数和可能的错误
func Exec(query string, args ...interface{}) (sql.Result, error) {
	if dbPool == nil {
		return nil, fmt.Errorf("数据库连接未初始化")
	}
	return dbPool.Exec(query, args...)
}

// QueryContext 执行查询并返回多行结果
// 参数: query - SQL 语句, args - 参数
// 返回: 查询结果和可能的错误
func Query(query string, args ...interface{}) (*sql.Rows, error) {
	if dbPool == nil {
		return nil, fmt.Errorf("数据库连接未初始化")
	}
	return dbPool.Query(query, args...)
}

// QueryRowContext 执行查询并返回单行结果
// 参数: query - SQL 语句, args - 参数
// 返回: 单行结果
func QueryRow(query string, args ...interface{}) *sql.Row {
	if dbPool == nil {
		log.Fatal("数据库连接未初始化")
	}
	return dbPool.QueryRow(query, args...)
}

// BeginTx 开始一个事务
// 返回: 事务对象和可能的错误
func BeginTx() (*sql.Tx, error) {
	if dbPool == nil {
		return nil, fmt.Errorf("数据库连接未初始化")
	}
	return dbPool.Begin()
}
