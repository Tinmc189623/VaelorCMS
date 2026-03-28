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

// ExampleModel 示例模型 - 演示如何使用 ORM
type ExampleModel struct {
	BaseModel
	Name  string `json:"name" db:"name"`
	Email string `json:"email" db:"email"`
}

// TableName 返回表名
func (e *ExampleModel) TableName() string {
	return "examples"
}

// ExampleUser 示例用户模型
type ExampleUser struct {
	BaseModel
	Username string         `json:"username" db:"username"`
	Email    string         `json:"email" db:"email"`
	Posts    []*ExamplePost `json:"posts,omitempty"`
}

// TableName 返回表名
func (u *ExampleUser) TableName() string {
	return "example_users"
}

// ExamplePost 示例文章模型
type ExamplePost struct {
	BaseModel
	Title   string        `json:"title" db:"title"`
	Content string        `json:"content" db:"content"`
	UserID  int64         `json:"user_id" db:"user_id"`
	Author  *ExampleUser  `json:"author,omitempty"`
	Tags    []*ExampleTag `json:"tags,omitempty"`
}

// TableName 返回表名
func (p *ExamplePost) TableName() string {
	return "example_posts"
}

// ExampleTag 示例标签模型
type ExampleTag struct {
	BaseModel
	Name  string         `json:"name" db:"name"`
	Posts []*ExamplePost `json:"posts,omitempty"`
}

// TableName 返回表名
func (t *ExampleTag) TableName() string {
	return "example_tags"
}

// ExamplePostTag 文章-标签中间表
type ExamplePostTag struct {
	PostID int64 `json:"post_id" db:"post_id"`
	TagID  int64 `json:"tag_id" db:"tag_id"`
}

// TableName 返回表名
func (pt *ExamplePostTag) TableName() string {
	return "example_post_tags"
}

// GetID 实现 Model 接口
func (pt *ExamplePostTag) GetID() int64 {
	return 0
}

// SetID 实现 Model 接口
func (pt *ExamplePostTag) SetID(id int64) {
}

// UsageExample 使用示例
func UsageExample(db *sql.DB) {
	// 1. 创建 ORM 实例
	orm := NewORM(db)

	// 2. 基础 CRUD 操作
	{
		// 创建记录
		user := &ExampleUser{
			Username: "john_doe",
			Email:    "john@example.com",
		}
		id, err := orm.Create(user)
		_ = id
		_ = err

		// 根据 ID 查找
		var foundUser ExampleUser
		found, err := orm.Find(&foundUser, 1)
		_ = found
		_ = err

		// 更新记录
		foundUser.Email = "newemail@example.com"
		rowsAffected, err := orm.Update(&foundUser)
		_ = rowsAffected
		_ = err

		// 删除记录
		rowsAffected, err = orm.Delete(&foundUser)
		_ = rowsAffected
		_ = err
	}

	// 3. 使用查询构建器
	{
		// 查询所有用户
		var users []*ExampleUser
		qb := NewQueryBuilder((&ExampleUser{}).TableName())
		err := orm.FindAll(&users, qb)
		_ = err

		// 条件查询
		qb = NewQueryBuilder((&ExampleUser{}).TableName()).
			Where("created_at > ?", time.Now().AddDate(0, -1, 0)).
			OrderByDesc("created_at").
			Limit(10)
		err = orm.FindAll(&users, qb)
		_ = err

		// WHERE IN 查询
		qb = NewQueryBuilder((&ExampleUser{}).TableName()).
			WhereIn("id", []interface{}{1, 2, 3})
		err = orm.FindAll(&users, qb)
		_ = err

		// 统计记录数
		count, err := orm.Count(&ExampleUser{}, qb)
		_ = count
		_ = err
	}

	// 4. 事务操作
	{
		tm := NewTransactionManager(db)

		// 使用 Transaction 方法
		err := tm.Transaction(func(orm *ORM) error {
			user := &ExampleUser{Username: "transaction_user"}
			_, err := orm.Create(user)
			if err != nil {
				return err
			}

			post := &ExamplePost{
				Title:   "Transaction Post",
				Content: "Content",
				UserID:  user.ID,
			}
			_, err = orm.Create(post)
			return err
		})
		_ = err

		// 使用 TransactionWithResult 方法
		result, err := tm.TransactionWithResult(func(orm *ORM) (interface{}, error) {
			user := &ExampleUser{Username: "result_user"}
			_, err := orm.Create(user)
			return user, err
		})
		_ = result
		_ = err

		// 手动控制事务
		tx, err := tm.Begin()
		if err != nil {
			return
		}
		defer tx.Rollback()

		txORM := tx.GetORM()
		user := &ExampleUser{Username: "manual_tx_user"}
		_, err = txORM.Create(user)
		if err != nil {
			return
		}

		err = tx.Commit()
		_ = err
	}

	// 5. 关系加载
	{
		loader := NewRelationLoader(orm)

		// 注册关系
		loader.RegisterRelation("Posts", Relation{
			Type:         HasMany,
			Field:        "Posts",
			RelatedModel: &ExamplePost{},
			ForeignKey:   "user_id",
			LocalKey:     "ID",
		})

		loader.RegisterRelation("Author", Relation{
			Type:         BelongsTo,
			Field:        "Author",
			RelatedModel: &ExampleUser{},
			ForeignKey:   "user_id",
			LocalKey:     "ID",
		})

		loader.RegisterRelation("Tags", Relation{
			Type:            ManyToMany,
			Field:           "Tags",
			RelatedModel:    &ExampleTag{},
			PivotTable:      "example_post_tags",
			PivotLocalKey:   "post_id",
			PivotForeignKey: "tag_id",
		})

		// 加载单个模型的关系
		var user ExampleUser
		found, _ := orm.Find(&user, 1)
		if found {
			loader.LoadRelation(&user, "Posts")
		}

		// 加载模型切片的关系
		var users []*ExampleUser
		orm.FindAll(&users)
		loader.LoadRelationForSlice(&users, "Posts")
	}
}
