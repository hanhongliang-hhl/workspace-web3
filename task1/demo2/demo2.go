package main

import (
	"fmt"
	"strings"
)

/*
	编写一个函数来查找字符串数组中的最长公共前缀。

如果不存在公共前缀，返回空字符串 ""。
*/
func longestCommonPrefix(strs []string) string {
	// 定义前缀,以第一个字符串为基准不断缩减匹配
	prefix := strs[0]
	for i := 1; i < len(strs); i++ {
		// 不断缩短prefix直到匹配当前字符串的前缀
		for !strings.HasPrefix(strs[i], prefix) {
			prefix = prefix[:len(prefix)-1]
			// 如果prefix为空，说明没有公共前缀
			if prefix == "" {
				panic("没有公共前缀")
			}
		}
	}
	return prefix
}
func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	strs := []string{"flower", "flow", "flight", "abc"}
	// strs := []string{"flower", "flow", "flight"}
	newVar := longestCommonPrefix(strs)
	fmt.Println(newVar)
}
