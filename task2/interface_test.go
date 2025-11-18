package main

import (
	"fmt"
	"math"
)

/*
✅面向对象
题目 ：定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。然后创建 Rectangle 和 Circle 结构体，
实现 Shape 接口。在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。
考察点 ：接口的定义与实现、面向对象编程风格。
*/
type Shape interface {
	Area()
	Perimeter()
}

type Rectangle struct {
	Length float32
	Width  float32
}
type Circle struct {
	Radius float32
}

func (rect *Rectangle) Area() float32 {

	return rect.Length * rect.Width
}

func (rect *Rectangle) Perimeter() float32 {

	return 2 * (rect.Length + rect.Width)
}

func (circle *Circle) Area() float32 {

	return math.Pi * circle.Radius * circle.Radius // 使用乘法而不是 math.Pow
}

func (circle *Circle) Perimeter() float32 {

	return 2 * math.Pi * circle.Radius
}

func main() {
	RectangleCase := Rectangle{
		Length: 10.00,
		Width:  5.00,
	}
	CircleCase := Circle{
		Radius: 2.00,
	}

	fmt.Printf("长方形的面积是：%2f cm\n²,", RectangleCase.Area())
	fmt.Printf("长方形的周长是：%2f  cm\n", RectangleCase.Perimeter())
	fmt.Printf("圆形的面积是：%2f cm²\n", CircleCase.Area())
	fmt.Printf("圆形的周长是：%2f cm\n", CircleCase.Perimeter())

}
