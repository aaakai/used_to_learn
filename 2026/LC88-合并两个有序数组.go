package _026

func merge(nums1 []int, m int, nums2 []int, n int) {
	n1Tag := m - 1
	n2Tag := n - 1
	index := m + n - 1
	for index >= 0 {
		res := 0
		if n1Tag < 0 {
			res = nums2[n2Tag]
			n2Tag--
		} else if n2Tag < 0 {
			res = nums1[n1Tag]
			n1Tag--
		} else if nums1[n1Tag] > nums2[n2Tag] {
			res = nums1[n1Tag]
			n1Tag--
		} else {
			res = nums2[n2Tag]
			n2Tag--
		}
		nums1[index] = res
		index--
	}
}
