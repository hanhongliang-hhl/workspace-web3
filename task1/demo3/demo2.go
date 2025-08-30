package main

import (
	"fmt"
)

func main() {
	/* 输入：digits = [1,2,3,4]
	输出：[4,3,2,2]
	*/
	digits := []int{1, 2, 3, 4}
	for i := 0; i < len(digits); i++ {
		for j := i+1; j < len(digits); j++ {
			if digits[i] < digits[j] {
				temp := digits[i]
				digits[i] = digits[j]
				digits[j] = temp
			}
		}
	}
	digits[len(digits)-1] += 1

	digits = append(digits, 5)
	fmt.Println(digits)
}
