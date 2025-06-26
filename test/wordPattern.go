package test

import "strings"

/*
给定一种规律 pattern 和一个字符串 s ，判断 s 是否遵循相同的规律。

这里的 遵循 指完全匹配，例如， pattern 里的每个字母和字符串 s 中的每个非空单词之间存在着双向连接的对应规律。



示例1:

输入: pattern = "abba", s = "dog cat cat dog"
输出: true
示例 2:

输入:pattern = "abba", s = "dog cat cat fish"
输出: false
示例 3:

输入: pattern = "aaaa", s = "dog cat cat dog"
输出: false
*/

func WordPattern(pattern string, s string) bool {
	mapC := make(map[string]byte)
	mapS := make(map[byte]string)
	words := strings.Split(s, " ")
	if len(pattern) != len(words) {
		return false
	}
	for k, v := range words {
		str := pattern[k]
		if mapC[v] > 0 && mapC[v] != str || mapS[str] != "" && mapS[str] != v {
			return false
		}
		mapC[v] = str
		mapS[str] = v
	}
	return true
}
