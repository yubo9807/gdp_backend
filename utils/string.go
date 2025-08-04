package utils

import (
	"strconv"
	"strings"
)

// VersionCompare 比较两个版本号
// 返回值:
//
//	-1: v1 < v2
//	 0: v1 == v2
//	 1: v1 > v2
func VersionCompare(v1, v2 string) int {
	// 处理空字符串情况
	if v1 == "" && v2 == "" {
		return 0
	}
	if v1 == "" {
		return -1
	}
	if v2 == "" {
		return 1
	}

	// 按 '.' 分割版本号
	v1Parts := strings.Split(v1, ".")
	v2Parts := strings.Split(v2, ".")

	// 确定需要比较的最大长度
	maxLen := len(v1Parts)
	if len(v2Parts) > maxLen {
		maxLen = len(v2Parts)
	}

	// 逐段比较
	for i := 0; i < maxLen; i++ {
		// 获取当前段的值，如果索引超出切片长度，则视为 0
		var num1, num2 int
		var err error

		if i < len(v1Parts) {
			// 将字符串转换为整数，Go 中没有 ~~ 操作符，使用 strconv.Atoi
			num1, err = strconv.Atoi(v1Parts[i])
			if err != nil {
				// 如果转换失败（例如非数字字符），可以按需处理，这里简单视为 0 或返回错误
				// 为了与原逻辑（~~ 会转为 0）一致，我们设为 0
				num1 = 0
			}
		} else {
			num1 = 0 // 缺失的段视为 0
		}

		if i < len(v2Parts) {
			num2, err = strconv.Atoi(v2Parts[i])
			if err != nil {
				num2 = 0
			}
		} else {
			num2 = 0 // 缺失的段视为 0
		}

		// 比较当前段
		if num1 > num2 {
			return 1
		}
		if num1 < num2 {
			return -1
		}
	}

	// 所有段都相等
	return 0
}

// 返回版本号列表中的最大版本号
func VersionMax(vs []string) string {
	max := vs[0]
	for _, v := range vs {
		if VersionCompare(v, max) > 0 {
			max = v
		}
	}
	return max
}
