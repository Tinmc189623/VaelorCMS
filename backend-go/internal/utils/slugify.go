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

package utils

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/gosimple/slug"
)

// GenerateSlug 生成 URL 友好的别名
// 参数:
//   text: 要转换的文本
// 返回:
//   string: URL 友好的别名
func GenerateSlug(text string) string {
	// 使用 gosimple/slug 库生成 URL 友好的别名
	slug.Lowercase = true
	return slug.Make(text)
}

// GenerateUniqueSlug 生成唯一的 URL 友好别名
// 参数:
//   text: 要转换的文本
//   existsFunc: 检查别名是否已存在的函数
// 返回:
//   string: 唯一的 URL 友好别名
func GenerateUniqueSlug(text string, existsFunc func(string) bool) string {
	baseSlug := GenerateSlug(text)
	if baseSlug == "" {
		baseSlug = "untitled"
	}

	// 如果基础别名不存在，直接返回
	if !existsFunc(baseSlug) {
		return baseSlug
	}

	// 尝试添加数字后缀
	suffix := 1
	for {
		candidate := baseSlug + "-" + strings.TrimLeft(strconv.Itoa(suffix), "0")
		if !existsFunc(candidate) {
			return candidate
		}
		suffix++
	}
}

// IsValidSlug 检查字符串是否是有效的 URL 别名
// 参数:
//   s: 要检查的字符串
// 返回:
//   bool: 是否是有效的 URL 别名
func IsValidSlug(s string) bool {
	if s == "" {
		return false
	}

	// 只允许小写字母、数字、连字符和下划线
	pattern := regexp.MustCompile(`^[a-z0-9-_]+$`)
	return pattern.MatchString(s)
}
