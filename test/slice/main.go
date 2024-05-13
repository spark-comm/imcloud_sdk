package main

import (
	"fmt"
	"open_im_sdk/pkg/utils"
)

func main() {
	slice := []int64{2, 42}

	// 找出不连续的元素并补齐
	fixedSlice, missing := utils.FixNotConnected(slice)

	fmt.Println("修复后的切片：", fixedSlice)
	fmt.Println("缺失的元素：", missing)
}
