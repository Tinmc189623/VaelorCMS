/*
 * VaelorCMS - 用户服务测试程序
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
	"fmt"
	"vaelorcms/internal/services"
)

func main() {
	fmt.Println("用户服务测试")
	
	// 创建用户服务实例
	userService := services.NewUserService()
	
	// 打印服务实例信息
	fmt.Printf("用户服务实例: %v\n", userService)
	fmt.Println("用户服务加载成功！")
}
