package _026

func twoSum(nums []int, target int) []int {
	tagMap := map[int]int{}
	for i := 0; i < len(nums); i++ {
		if _, ok := tagMap[target-nums[i]]; ok {
			return []int{i, tagMap[target-nums[i]]}
		} else {
			tagMap[nums[i]] = i
		}
	}
	return []int{}
}
