package _026

func groupAnagrams(strs []string) [][]string {
	res := make([][]string, 0)
	tmp := make(map[[26]int][]string, 0)
	for _, str := range strs {
		cnt := [26]int{}
		for _, c := range str {
			cnt[c-'a']++
		}
		tmp[cnt] = append(tmp[cnt], str)
	}
	for _, a := range tmp {
		res = append(res, a)
	}
	return res
}

// golang 的map key 可以是任意类型
