/*
查找字符串数组中的最长公共前缀
*/
package main

import "fmt"

func main() {
	fmt.Println(longestCommonPrefix([]string{"flower", "flow", "flows"}))
}

func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	prefix := strs[0]
	for i := 1; i < len(strs); i++ {
		for k, _ := range strs[i] {
			if prefix[k] != strs[i][k] {
				prefix = prefix[:k]
				break
			}
		}
	}
	return prefix
}
