package _026

func subarraySum(nums []int, k int) int {
	res, tag := 0, 0
	tmpMap := make(map[int]int)
	tmpMap[0] = 1
	for _, v := range nums {
		tag += v
		if _, ok := tmpMap[tag-k]; ok {
			res += tmpMap[tag-k]
		}
		tmpMap[tag]++
	}
	return res
}

//前缀原理
