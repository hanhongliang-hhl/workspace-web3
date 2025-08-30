package main

import (
	"fmt"
)

/*
	题目 ：编写一个Go程序，定义一个函数，该函数接收一个整数指针作为参数，在函数内部将该指针指向的值增加10，然后在主函数中调用该函数并输出修改后的值。

考察点 ：指针的使用、值传递与引用传递的区别。
*/
func zhizhen(a *int) *int {
	*a += 10
	return a
}

/*
题目 ：实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2。
考察点 ：指针运算、切片操作。
*/
func qiepian(a *[]int) {
	for i := range *a { //*a表示切片值，a表示切片指针
		(*a)[i]*=2
	}
}

func main() {
	a := 10
	newVar := zhizhen(&a)
	fmt.Println(*newVar)

	newVar1 := make([]int, 5)
	newVar1 = []int{1, 2, 3, 4, 5}
	fmt.Println(newVar1, len(newVar1), cap(newVar1))
	newVar1 = append(newVar1, 6)
	fmt.Println(newVar1, len(newVar1), cap(newVar1))
	qiepian(&newVar1)
	fmt.Println(newVar1, len(newVar1), cap(newVar1))

}
