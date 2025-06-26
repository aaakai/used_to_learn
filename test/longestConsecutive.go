package test

import "sort"

/*
给定一个未排序的整数数组 nums ，找出数字连续的最长序列（不要求序列元素在原数组中连续）的长度。

请你设计并实现时间复杂度为 O(n) 的算法解决此问题。

示例 1：

输入：nums = [100,4,200,1,3,2]
输出：4
解释：最长数字连续序列是 [1, 2, 3, 4]。它的长度为 4。
示例 2：

输入：nums = [0,3,7,2,5,8,4,6,0,1]
输出：9
示例 3：

输入：nums = [1,0,1,2]
输出：3

提示：

0 <= nums.length <= 105
-109 <= nums[i] <= 109
*/

func LongestConsecutive(nums []int) int {
	n := len(nums)
	if n <= 1 {
		return n
	}

	sort.Ints(nums)
	ans := 1
	length := nums[0]
	for i := 1; i < n; i++ {
		if nums[i] == nums[i-1]+1 {
			length++
		} else {
			if nums[i] == nums[i-1] {
				continue
			}
			length = 1
		}
		ans = max(ans, length)
	}
	return ans
}

/*func LongestConsecutive(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	mapNums := make(map[int]bool)
	for _, num := range nums {
		mapNums[num] = true
	}
	maxLen := 0
	for num := range mapNums {
		if !mapNums[num-1] {
			curNum := num
			curLen := 1
			for mapNums[curNum+1] {
				curNum++
				curLen++
			}
			if curLen > maxLen {
				maxLen = curLen
			}
		}
	}
	return maxLen
}*/
