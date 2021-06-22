package sshmysql

import (
	"fmt"
	"testing"
)

func TestGetDeployRecordListBySpaceIds(t *testing.T) {
	ss := []string{"da778b0964024d118d6c79f967a98ede", "8b60c1b73d2743d0ba8fb3541625eedf", "fb4c0f80fe884033ac8be35a7916ec51"}
	re, err := GetDeployRecordListBySpaceIds(ss)
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range re {
		fmt.Println(v)
	}
}

func TestGetDeployOpRecordList(t *testing.T) {
	re, err := GetDeployOpRecordList()
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range re {
		fmt.Println(v)
	}
}
