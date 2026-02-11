package _026

func trap(height []int) int {
	res := 0
	left, right := 0, len(height)-1
	lMax, rMax := 0, 0
	for left < right {
		lMax = maxI(lMax, height[left])
		rMax = maxI(rMax, height[right])
		if height[left] < height[right] {
			res += lMax - height[left]
			left++
		} else {
			res += rMax - height[right]
			right--
		}
	}
	return res
}

func maxI(a, b int) int {
	if a > b {
		return a
	}
	return b
}
