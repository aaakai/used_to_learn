package _026

func longestConsecutive(nums []int) int {
	tagMap := make(map[int]bool)
	for _, num := range nums {
		tagMap[num] = true
	}
	res := 0
	for num := range tagMap {
		tmp := 1
		if !tagMap[num-1] {
			cur := num
			for tagMap[cur+1] {
				cur++
				tmp++
			}
		}
		if tmp > res {
			res = tmp
		}
	}
	return res
}

// 也可以先排序 循环一次 i到j 升序，单次循环结束后 i=j
