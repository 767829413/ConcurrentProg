package myexec

import "fmt"

func findMedianSortedArrays(nums1 []int, nums2 []int) float64 {
	l1 := len(nums1)
	l2 := len(nums2)
	if l1 == l2 && l2 == 0 {
		return float64(0)
	}
	all := l1 + l2
	mid := all / 2
	flag := (all % 2) == 0
	tmp, p1, p2 := make([]int, all), 0, 0
	if l1 == 0 {
		if flag {
			return (float64(nums2[mid]) + float64(nums2[mid-1])) / 2
		}
		return float64(nums2[mid])
	}
	if l2 == 0 {
		if flag {
			return (float64(nums1[mid]) + float64(nums1[mid-1])) / 2
		}
		return float64(nums1[mid])
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
	if flag {
		return (float64(tmp[mid]) + float64(tmp[mid-1])) / 2
	}
	return float64(tmp[mid])
}

func findMedianSortedArrays1(nums1 []int, nums2 []int) float64 {
	l1 := len(nums1)
	l2 := len(nums2)
	if l1 == l2 && l2 == 0 {
		return float64(0)
	}
	all := l1 + l2
	mid := all / 2
	flag := (all % 2) == 0
	p1, p2 := 0, 0
	if l1 == 0 {
		if flag {
			return (float64(nums2[mid]) + float64(nums2[mid-1])) / 2
		}
		return float64(nums2[mid])
	}
	if l2 == 0 {
		if flag {
			return (float64(nums1[mid]) + float64(nums1[mid-1])) / 2
		}
		return float64(nums1[mid])
	}

	for i := 0; i < mid; i++ {
		if p1 < l1 && p2 < l2 {
			if nums1[p1] > nums2[p2] {
				p2++
			} else {
				p1++
			}
		} else if p1 < l1 && p2 >= l2 {
			p1++
		} else if p1 >= l1 && p2 < l2 {
			p2++
		}
	}
	fmt.Println(p1, p2)
	return 0.0
}
