package test

func LengthOfLastWord(s string) int {
	res := 0
	lenS := len(s) - 1
	for lenS >= 0 && s[lenS] == ' ' {
		lenS--
	}

	for lenS >= 0 && s[lenS] != ' ' {
		res++
		lenS--
	}
	return res
}
