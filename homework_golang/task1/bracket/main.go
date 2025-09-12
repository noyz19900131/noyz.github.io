/*
给定一个只包括 '('，')'，'{'，'}'，'['，']' 的字符串，判断字符串是否有效
*/
package main

import "fmt"

func main() {
	fmt.Println(isValid("()"))
	fmt.Println(isValid("()[]{}"))
	fmt.Println(isValid("(]"))
	fmt.Println(isValid("([)]"))
	fmt.Println(isValid("{[]}"))
	fmt.Println(isValid("()[]{}"))
	fmt.Println(isValid("([{}])"))
}

func isValid(str string) bool {
	if len(str)%2 != 0 {
		return false
	}
	mp := map[rune]rune{
		')': '(',
		']': '[',
		'}': '{',
	}
	stack := []rune{}
	for _, ch := range str {
		if ch == '(' || ch == '[' || ch == '{' {
			stack = append(stack, ch)
		} else {
			if len(stack) == 0 || mp[ch] != stack[len(stack)-1] {
				return false
			}
			stack = stack[:len(stack)-1]
		}
	}
	return len(stack) == 0
}
