package _026

func removeElement(nums []int, val int) int {
	if len(nums) <= 1 {
		return len(nums)
	}
	left := 0
	for right := 0; right < len(nums); right++ {
		if nums[right] != val {
			nums[left] = nums[right]
			left++
		}
	}
	return left
}
