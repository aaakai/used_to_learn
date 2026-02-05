package _026

func removeDuplicates(nums []int) int {
	numLen := len(nums)
	if numLen <= 1 {
		return numLen
	}
	slow := 0
	for i := 1; i < numLen; i++ {
		if nums[slow] != nums[i] {
			nums[slow] = nums[i]
			slow++
		}
	}
	return slow
}
