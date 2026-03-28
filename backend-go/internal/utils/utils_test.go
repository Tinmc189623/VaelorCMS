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
	"testing"
)

// TestPasswordHash 测试密码哈希和验证功能
func TestPasswordHash(t *testing.T) {
	password := "testpassword123"

	// 测试生成密码哈希
	hashed, err := GetPasswordHash(password)
	if err != nil {
		t.Fatalf("生成密码哈希失败: %v", err)
	}

	// 测试验证正确的密码
	if !VerifyPassword(password, hashed) {
		t.Error("验证正确的密码失败")
	}

	// 测试验证错误的密码
	if VerifyPassword("wrongpassword", hashed) {
		t.Error("验证错误的密码应该失败")
	}
}

// TestSlugify 测试 URL 别名生成功能
func TestSlugify(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"Hello World", "hello-world"},
		{"VaelorCMS 内容管理系统", "vaelorcms-nei-rong-guan-li-xi-tong"},
		{"Test 123!@#", "test-123"},
		{"  Leading and Trailing Spaces  ", "leading-and-trailing-spaces"},
	}

	for _, tc := range testCases {
		result := GenerateSlug(tc.input)
		if result == "" {
			t.Errorf("GenerateSlug(%q) 返回空字符串", tc.input)
		}
	}
}

// TestIsValidSlug 测试别名验证功能
func TestIsValidSlug(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"hello-world", true},
		{"test_123", true},
		{"Hello-World", false}, // 不允许大写
		{"test 123", false},    // 不允许空格
		{"test!@#", false},      // 不允许特殊字符
		{"", false},             // 空字符串
	}

	for _, tc := range testCases {
		result := IsValidSlug(tc.input)
		if result != tc.expected {
			t.Errorf("IsValidSlug(%q) = %v, expected %v", tc.input, result, tc.expected)
		}
	}
}
