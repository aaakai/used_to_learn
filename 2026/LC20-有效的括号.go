package _026

func isValid(s string) bool {
	n := len(s)
	if n%2 == 1 {
		return false
	}
	charMap := map[byte]byte{
		')': '(',
		']': '[',
		'}': '{',
	}
	stack := []byte{}

	for i := 0; i < n; i++ {
		if _, ok := charMap[s[i]]; ok {
			if len(stack) == 0 || stack[len(stack)-1] != charMap[s[i]] {
				return false
			}
			stack = stack[:len(stack)-1]
		} else {
			stack = append(stack, s[i])
		}
	}
	return len(stack) == 0
}
