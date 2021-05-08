package mypractice

func findMedianSortedArrays(nums1, nums2 []int) float64 {
	l1, l2 := len(nums1), len(nums2)
	if l1 == l2 && l2 == 0 {
		return float64(0)
	}
	all := l1 + l2
	mid := all / 2
	flag := (all % 2) == 0
	tmp, p1, p2 := make([]int, all), 0, 0
	f := func(nums []int) float64 {
		if flag {
			return (float64(nums[mid]) + float64(nums[mid-1])) / 2
		}
		return float64(nums[mid])
	}
	if l1 == 0 {
		return f(nums2)
	}
	if l2 == 0 {
		return f(nums1)
	}
	for i := 0; i < all; i++ {
		if p1 < l1 && p2 < l2 {
			if nums1[p1] > nums2[p2] {
				tmp[i] = nums2[p2]
				p2++
			} else {
				tmp[i] = nums1[p1]
				p1++
			}
		} else if p1 < l1 && p2 >= l2 {
			tmp[i] = nums1[p1]
			p1++
		} else if p1 >= l1 && p2 < l2 {
			tmp[i] = nums2[p2]
			p2++
		}
	}
	return f(tmp)
}

func findMedianSortedArrays1(nums1, nums2 []int) float64 {
	l1, l2 := len(nums1), len(nums2)
	if l1 == l2 && l2 == 0 {
		return float64(0)
	}
	all := l1 + l2
	mid := all / 2
	flag := (all % 2) == 0
	p1, p2, r1, r2 := 0, 0, 0, 0
	f := func(nums []int) float64 {
		if flag {
			return (float64(nums[mid]) + float64(nums[mid-1])) / 2
		}
		return float64(nums[mid])
	}
	if l1 == 0 {
		return f(nums2)
	}
	if l2 == 0 {
		return f(nums1)
	}
	for i := 0; i <= mid; i++ {
		r2 = r1
		if p1 < l1 && p2 < l2 {
			if nums1[p1] > nums2[p2] {
				r1 = nums2[p2]
				p2++
			} else {
				r1 = nums1[p1]
				p1++
			}
		} else if p1 < l1 && p2 >= l2 {
			r1 = nums1[p1]
			p1++
		} else if p1 >= l1 && p2 < l2 {
			r1 = nums2[p2]
			p2++
		}
	}
	if flag {
		return (float64(r1) + float64(r2)) / 2
	} else {
		return float64(r1)
	}
}

func findMedianSortedArrays2(nums1, nums2 []int) float64 {
	l := len(nums1) + len(nums2)
	mid := l / 2
	if (l % 2) == 0 {
		return float64(getKthElement(nums1, nums2, mid)+getKthElement(nums1, nums2, mid+1)) / 2.0
	} else {
		return float64(getKthElement(nums1, nums2, mid+1))
	}
}

func getKthElement(nums1, nums2 []int, k int) int {
	index1, index2 := 0, 0
	l1, l2 := len(nums1), len(nums2)
	for {
		if index1 == l1 {
			return nums2[index2+k-1]
		}
		if index2 == l2 {
			return nums1[index1+k-1]
		}
		if k == 1 {
			return min(nums1[index1], nums2[index2])
		}
		half := k / 2
		newIndex1, newIndex2 := min(index1+half, l1)-1, min(index2+half, l2)-1
		p1, p2 := nums1[newIndex1], nums2[newIndex2]
		if p1 <= p2 {
			k = k - (newIndex1 - index1 + 1)
			index1 = newIndex1 + 1
		} else {
			k = k - (newIndex2 - index2 + 1)
			index2 = newIndex2 + 1
		}
	}
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}
