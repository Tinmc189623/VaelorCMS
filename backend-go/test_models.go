/*
 * VaelorCMS - 模型测试文件
 *
 * Copyright © 2025-2026 Nexsteaduser. All rights reserved.
 *
 * 作者: Tinmc189623
 * 团队: Nexsteaduser
 *
 * 本程序是自由软件: 你可以重新分发和/或修改
 * 它在 GNU Affero 通用公共许可证的条款下发布,
 * 版本 3 或 (根据你的选择) 任何更高版本。
 */

package main

import (
	"vaelorcms/internal/models"
)

func main() {
	// 简单的测试，只是为了验证模型是否可以编译
	_ = &models.User{}
	_ = &models.Article{}
	_ = &models.Category{}
	_ = &models.Tag{}
	_ = &models.Content{}
	_ = &models.Media{}
	_ = &models.Page{}
	_ = &models.Setting{}
	
	println("Models compiled successfully!")
}
