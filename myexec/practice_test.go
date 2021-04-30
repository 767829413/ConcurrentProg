package myexec

import "testing"

func TestFindMedianSortedArrays(t *testing.T) {
	arr1 := []int{1,3,5}
	arr2 := []int{2,4,6,8,10}
	//t.Log(findMedianSortedArrays(arr1, arr2))
	t.Log(findMedianSortedArrays1(arr1, arr2))
}
