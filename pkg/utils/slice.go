package utils

import "fmt"

func DeleteSlice(a []int, elem int) []int {
	j := 0
	for _, v := range a {
		if v != elem {
			a[j] = v
			j++
		}
	}
	return a[:j]
}

// 数组去重
func RemoveDuplicate[T string | int](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func ToStringSlice[T int32 | int | int64](sliceList []T) []string {
	var list []string
	for _, item := range sliceList {
		item := item
		list = append(list, fmt.Sprintf("%v", item))
	}
	return list
}
