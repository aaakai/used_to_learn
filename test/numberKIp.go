package test

/*
给你一个非负整数数组 nums ，你最初位于数组的 第一个下标 。数组中的每个元素代表你在该位置可以跳跃的最大长度。

判断你是否能够到达最后一个下标，如果可以，返回 true ；否则，返回 false 。

示例 1：

输入：nums = [2,3,1,1,4]
输出：true
解释：可以先跳 1 步，从下标 0 到达下标 1, 然后再从下标 1 跳 3 步到达最后一个下标。
示例 2：

输入：nums = [3,2,1,0,4]
输出：false
解释：无论怎样，总会到达下标为 3 的位置。但该下标的最大跳跃长度是 0 ， 所以永远不可能到达最后一个下标。

提示：

1 <= nums.length <= 104
0 <= nums[i] <= 105

case [2,3,1,1,4] [3,2,1,0,4]
*/
func CanJump(nums []int) bool {
	if len(nums) == 0 {
		return false
	}
	n := nums[0]
	for i := 1; i < len(nums); i++ {
		if n == 0 {
			return false
		} else if n-1 >= nums[i] {
			n--
			continue
		} else {
			n = nums[i]
		}
	}
	return true
}

/*
给定一个长度为 n 的 0 索引整数数组 nums。初始位置为 nums[0]。

每个元素 nums[i] 表示从索引 i 向前跳转的最大长度。换句话说，如果你在 nums[i] 处，你可以跳转到任意 nums[i + j] 处:

0 <= j <= nums[i]
i + j < n
返回到达 nums[n - 1] 的最小跳跃次数。生成的测试用例可以到达 nums[n - 1]。

示例 1:

输入: nums = [2,3,1,1,4]
输出: 2
解释: 跳到最后一个位置的最小跳跃数是 2。

	从下标为 0 跳到下标为 1 的位置，跳 1 步，然后跳 3 步到达数组的最后一个位置。

示例 2:

输入: nums = [2,3,0,1,4]
输出: 2

提示:

1 <= nums.length <= 104
0 <= nums[i] <= 1000
题目保证可以到达 nums[n-1]

case [2,3,1,1,4]  [2,3,0,1,4]
*/
func CanJump2(nums []int) int {
	jump := 0
	end := 0
	pot := 0
	for i := 0; i < len(nums)-1; i++ {
		pot = max(pot, i+nums[i])
		if end == i {
			end = pot
			jump++
		}
	}
	return jump
}
