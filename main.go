package main

import "fmt"

func bubbleSort(arr []int) {
	n := len(arr)

	// 外层循环：需要比较 n-1 轮
	for i := 0; i < n-1; i++ {
		// 内层循环：每轮比较相邻元素
		// n-i-1：每轮结束后，末尾已经有序，不需要再比较
		for j := 0; j < n-i-1; j++ {
			// 如果前一个大于后一个，交换
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
}
func main() {
	arr := []int{64, 34, 25, 12, 22, 11, 90}
	fmt.Println("排序前:", arr)

	bubbleSort(arr)
	fmt.Println("排序后:", arr)
}
