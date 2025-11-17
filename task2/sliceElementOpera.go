package main

import "fmt"

/*
题目 ：实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2。
考察点 ：指针运算、切片操作。
*/
func main() {
	var slice1 = make([]int, 0)
	slice1 = append(slice1, 1, 2, 3, 4, 5)
	fmt.Println("slice1：", slice1) // 修正：不能用 + 连接字符串和切片

	fmt.Println("操作后的slice1：", sliceElementOpera(slice1))
}

func sliceElementOpera(sliceTest []int) []int {

	for i := range sliceTest {
		sliceTest[i] *= 2
	}
	fmt.Println("方法里的sliceTest：", sliceTest)
	return sliceTest
}
