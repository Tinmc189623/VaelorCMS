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
	"log"
	"vaelorcms/internal/config"
)

// CreateTables 创建所有数据库表
// 该函数会创建应用程序所需的所有表，如果表已存在则不会重建
func CreateTables() error {
	log.Println("开始创建数据库表...")

	if err := createUsersTable(); err != nil {
		return err
	}

	if err := createCategoriesTable(); err != nil {
		return err
	}

	if err := createTagsTable(); err != nil {
		return err
	}

	if err := createContentsTable(); err != nil {
		return err
	}

	if err := createContentCategoriesTable(); err != nil {
		return err
	}

	if err := createContentTagsTable(); err != nil {
		return err
	}

	if err := createArticlesTable(); err != nil {
		return err
	}

	if err := createPagesTable(); err != nil {
		return err
	}

	if err := createMediaTable(); err != nil {
		return err
	}

	if err := createSettingsTable(); err != nil {
		return err
	}

	log.Println("所有数据库表创建完成")
	return nil
}

// createUsersTable 创建用户表
func createUsersTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		full_name TEXT,
		avatar TEXT,
		bio TEXT,
		role TEXT NOT NULL DEFAULT 'user',
		is_active BOOLEAN NOT NULL DEFAULT 1,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := Exec(query)
	if err != nil {
		return err
	}

	log.Println("用户表创建成功")
	return nil
}

// createCategoriesTable 创建分类表
func createCategoriesTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		slug TEXT NOT NULL UNIQUE,
		description TEXT,
		parent_id INTEGER,
		sort_order INTEGER DEFAULT 0,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (parent_id) REFERENCES categories(id) ON DELETE SET NULL
	)`

	_, err := Exec(query)
	if err != nil {
		return err
	}

	log.Println("分类表创建成功")
	return nil
}

// createTagsTable 创建标签表
func createTagsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS tags (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		slug TEXT NOT NULL UNIQUE,
		description TEXT,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := Exec(query)
	if err != nil {
		return err
	}

	log.Println("标签表创建成功")
	return nil
}

// createContentsTable 创建内容表
func createContentsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS contents (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		slug TEXT NOT NULL UNIQUE,
		content TEXT,
		excerpt TEXT,
		status TEXT NOT NULL DEFAULT 'draft',
		type TEXT NOT NULL DEFAULT 'post',
		author_id INTEGER NOT NULL,
		featured_image TEXT,
		views INTEGER DEFAULT 0,
		comment_count INTEGER DEFAULT 0,
		published_at DATETIME,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
	)`

	_, err := Exec(query)
	if err != nil {
		return err
	}

	log.Println("内容表创建成功")
	return nil
}

// createContentCategoriesTable 创建内容-分类关联表
func createContentCategoriesTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS content_categories (
		content_id INTEGER NOT NULL,
		category_id INTEGER NOT NULL,
		PRIMARY KEY (content_id, category_id),
		FOREIGN KEY (content_id) REFERENCES contents(id) ON DELETE CASCADE,
		FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
	)`

	_, err := Exec(query)
	if err != nil {
		return err
	}

	log.Println("内容-分类关联表创建成功")
	return nil
}

// createContentTagsTable 创建内容-标签关联表
func createContentTagsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS content_tags (
		content_id INTEGER NOT NULL,
		tag_id INTEGER NOT NULL,
		PRIMARY KEY (content_id, tag_id),
		FOREIGN KEY (content_id) REFERENCES contents(id) ON DELETE CASCADE,
		FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
	)`

	_, err := Exec(query)
	if err != nil {
		return err
	}

	log.Println("内容-标签关联表创建成功")
	return nil
}

// createArticlesTable 创建文章表
func createArticlesTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS articles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		content_id INTEGER NOT NULL UNIQUE,
		is_featured BOOLEAN DEFAULT 0,
		is_sticky BOOLEAN DEFAULT 0,
		read_time INTEGER,
		FOREIGN KEY (content_id) REFERENCES contents(id) ON DELETE CASCADE
	)`

	_, err := Exec(query)
	if err != nil {
		return err
	}

	log.Println("文章表创建成功")
	return nil
}

// createPagesTable 创建页面表
func createPagesTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS pages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		content_id INTEGER NOT NULL UNIQUE,
		template TEXT,
		is_homepage BOOLEAN DEFAULT 0,
		sort_order INTEGER DEFAULT 0,
		FOREIGN KEY (content_id) REFERENCES contents(id) ON DELETE CASCADE
	)`

	_, err := Exec(query)
	if err != nil {
		return err
	}

	log.Println("页面表创建成功")
	return nil
}

// createMediaTable 创建媒体表
func createMediaTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS media (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		description TEXT,
		file_name TEXT NOT NULL,
		file_path TEXT NOT NULL,
		file_type TEXT NOT NULL,
		file_size INTEGER NOT NULL,
		mime_type TEXT,
		width INTEGER,
		height INTEGER,
		author_id INTEGER,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE SET NULL
	)`

	_, err := Exec(query)
	if err != nil {
		return err
	}

	log.Println("媒体表创建成功")
	return nil
}

// createSettingsTable 创建设置表
func createSettingsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS settings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		key TEXT NOT NULL UNIQUE,
		value TEXT,
		group_name TEXT,
		description TEXT,
		type TEXT DEFAULT 'string',
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := Exec(query)
	if err != nil {
		return err
	}

	log.Println("设置表创建成功")
	return nil
}

// DropAllTables 删除所有数据库表
// 警告: 此操作不可逆，请谨慎使用
func DropAllTables() error {
	log.Println("警告: 开始删除所有数据库表...")

	tables := []string{
		"settings",
		"media",
		"pages",
		"articles",
		"content_tags",
		"content_categories",
		"contents",
		"tags",
		"categories",
		"users",
	}

	for _, table := range tables {
		query := "DROP TABLE IF EXISTS " + table
		_, err := Exec(query)
		if err != nil {
			return err
		}
		log.Printf("表 %s 已删除", table)
	}

	log.Println("所有数据库表已删除")
	return nil
}

// InitializeDatabase 完整初始化数据库
// 包括创建连接池和创建表
func InitializeDatabase(cfg *config.Config) error {
	if err := InitDB(cfg); err != nil {
		return err
	}

	if err := CreateTables(); err != nil {
		return err
	}

	log.Println("数据库初始化完成")
	return nil
}
