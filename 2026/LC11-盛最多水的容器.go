package _026

func maxArea(height []int) int {
	res := 0
	if len(height) < 1 {
		return 0
	}
	left, right := 0, len(height)-1
	for right > left {
		area := min(height[left], height[right]) * (right - left)
		res = max(res, area)
		if height[left] < height[right] {
			left++
		} else {
			right--
		}
	}
	return res
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
