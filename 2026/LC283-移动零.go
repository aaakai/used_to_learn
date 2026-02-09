package _026

func moveZeroes(nums []int) {
	fast, slow := 0, 0
	for fast < len(nums) {
		if nums[fast] != 0 {
			nums[fast], nums[slow] = nums[slow], nums[fast]
			slow++
		}
		fast++
	}
}
