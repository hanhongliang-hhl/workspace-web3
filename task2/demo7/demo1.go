package main

import (
	"fmt"

)

type Shape interface {
	Area() int64
	Perimeter(int)
}

type Rectangle struct{
	name string
	age int
}
type Circle struct{

}

func (r Rectangle) Area() int64 {
	return 0
}

func (r *Rectangle) Perimeter(age int) {
	r.age = age
}

func (c Circle) Area() int64 {
	return 0
}

func (r Circle) Perimeter(age int) {

}

type Person struct {
	Name string
	Age int
}

type Employee struct {
	Person  Person
	EmployeeID int
}

func (e Employee) PrintInfo() {
	fmt.Println(e.Person.Name, e.Person.Age, e.EmployeeID)
}

func main() {
	/* 题目 ：定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。然后创建 Rectangle 和 Circle 结构体，实现 Shape 接口。
	在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。*/
	var sha Shape
	sha = Circle{}
	value := sha.Area()
	fmt.Println(value)

	sha = &Rectangle{
		name: "长方形",
		age: 10,
	}
	sha.Perimeter(20)
	fmt.Println(sha) 


/*题目 ：使用组合的方式创建一个 Person 结构体，包含 Name 和 Age 字段，再创建一个 Employee 结构体，
组合 Person 结构体并添加 EmployeeID 字段。为 Employee 结构体实现一个 PrintInfo() 方法，输出员工的信息。 */

	 per := Person{Name: "张三", Age: 18}
	 emp := Employee{Person: per, EmployeeID: 1001}
	 emp.PrintInfo()

	fmt.Println("main 结束")

}
