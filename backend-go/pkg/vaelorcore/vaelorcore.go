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

// Package vaelorcore 提供了一个轻量级、功能完整的 ORM 框架
//
// 主要功能包括:
//   - 完整的 CRUD 操作
//   - 强大的查询构建器
//   - 关系映射（一对多、多对多、一对一、属于关系）
//   - 事务支持
//   - 事务管理器
//
// 示例:
//
//	// 创建 ORM 实例
//	orm := vaelorcore.NewORM(db)
//
//	// 查询记录
//	var user models.User
//	found, err := orm.Find(&user, 1)
//
//	// 使用查询构建器
//	qb := vaelorcore.NewQueryBuilder("users").Where("is_active = ?", true)
//	var users []*models.User
//	err = orm.FindAll(&users, qb)
//
//	// 创建记录
//	user := &models.User{Username: "test"}
//	id, err := orm.Create(user)
//
//	// 事务操作
//	tm := vaelorcore.NewTransactionManager(db)
//	err := tm.Transaction(func(orm *vaelorcore.ORM) error {
//	    // 在事务中执行操作
//	    return nil
//	})
//
//	// 关系加载
//	loader := vaelorcore.NewRelationLoader(orm)
//	loader.RegisterRelation("Articles", vaelorcore.Relation{
//	    Type:         vaelorcore.HasMany,
//	    Field:        "Articles",
//	    RelatedModel: &models.Article{},
//	    ForeignKey:   "author_id",
//	})
//	err = loader.LoadRelation(&user, "Articles")
package vaelorcore
