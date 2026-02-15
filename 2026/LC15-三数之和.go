package _026

import "sort"

func threeSum(nums []int) [][]int {
	res := make([][]int, 0)
	sort.Ints(nums)
	lenNums := len(nums)

	for i := 0; i < lenNums; i++ {
		//跳过第一个相同的元素
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}
		right := lenNums - 1
		tag := 0 - nums[i]
		for left := i + 1; left < right; left++ {
			//跳过第一个相同的元素
			if left > i+1 && nums[left] == nums[left-1] {
				continue
			}
			for left < right && nums[left]+nums[right] > tag {
				right--
			}
			if left >= right {
				break
			}
			if nums[left]+nums[right] == tag {
				res = append(res, []int{nums[i], nums[left], nums[right]})
			}
		}
	}
	return res
}
