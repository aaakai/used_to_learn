package _026

func lengthOfLongestSubstring(s string) int {
	keyMap := make(map[byte]int)
	tag, res := -1, 0
	for i := 0; i < len(s); i++ {
		if i != 0 {
			delete(keyMap, s[i-1])
		}
		for tag+1 < len(s) && keyMap[s[tag+1]] == 0 {
			keyMap[s[tag+1]]++
			tag++
		}
		res = max3(res, tag-i+1)
	}
	return res
}

func max3(a, b int) int {
	if a > b {
		return a
	}
	return b
}
