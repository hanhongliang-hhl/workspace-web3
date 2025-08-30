package main

import (
	"fmt"
)

func main() {
	/* 输入： [0,0,1,1,1,2,2,3,3,4] 
	输出：[0,1,2,3,4] 长度为5
	解释：函数应该返回新的长度 2 ，并且原数组 nums 的前两个元素被修改为 1, 2 。不需要考虑数组中超出新长度后面的元素。
	*/
	nums := []int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}
	
	// 使用双指针法
    slow := 0  // 慢指针指向不重复元素的位置
    
    for fast := 1; fast < len(nums); fast++ {
        // 如果快指针指向的元素与慢指针不同
        if nums[fast] != nums[slow] {
            slow++  // 慢指针前移
            nums[slow] = nums[fast]  // 将不重复元素放到慢指针位置
        }
    }
	fmt.Println(nums[:slow+1])
}
