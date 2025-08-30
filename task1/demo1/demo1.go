package main

import "fmt"

func main() {
	/* 给定一个非空整数数组，除了某个元素只出现一次以外，其余每个元素均出现两次。找出那个只出现了一次的元素。
	可以使用 for 循环遍历数组，结合 if 条件判断和 map 数据结构来解决，例如通过 map 记录每个元素出现的次数，
	然后再遍历 map 找到出现次数为1的元素。 */

	var arr = [5]int{4, 1, 2, 1, 2}
	var map1 = make(map[int]int, len(arr))
	for i := 0; i < len(arr); i++ {
		//把数组的值当中map的key，如果key不存在把value置为1存入map，value+1,存在就value+1
		value, exists := map1[arr[i]]
		if exists {
			map1[arr[i]] = value + 1
		} else {
			map1[arr[i]] = 1
		}
	}
	fmt.Println(map1)
	//变量map1的key和value
	for k, v := range map1 {
		if v == 1 {
			fmt.Println(k)
		}
	}

}