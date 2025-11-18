package main

import "fmt"

/*题目 ：使用组合的方式创建一个 Person 结构体，包含 Name 和 Age 字段，
再创建一个 Employee 结构体，组合 Person 结构体并添加 EmployeeID 字段。为 Employee 结构体实现一个 PrintInfo() 方法，输出员工的信息。
考察点 ：组合的使用、方法接收者。*/

type Person struct {
	Name string
	Age  int
}

type Employee struct {
	person     Person
	EmployeeID int
}

func (emp *Employee) PrintInfo() {
	fmt.Println("=== 员工信息 ===")
	fmt.Printf("ID: %d\n", emp.EmployeeID)
	fmt.Printf("姓名: %s\n", emp.person.Name)
	fmt.Printf("年龄: %d\n", emp.person.Age)
	fmt.Println("===============")
}

func main() {
	emp := Employee{
		person: Person{
			Name: "张三",
			Age:  18,
		},
		EmployeeID: 9527,
	}
	emp.PrintInfo()
}
