/*
使用组合的方式创建一个 Person 结构体，包含 Name 和 Age 字段，
再创建一个 Employee 结构体，组合 Person 结构体并添加 EmployeeID 字段。
为 Employee 结构体实现一个 PrintInfo() 方法，输出员工的信息。

考察点 ：组合的使用、方法接收者。
*/
package main

import "fmt"

func main() {
	e := Employee{Person{"Jack", 22}, 220821}
	e.PrintInfo()
}

type Person struct {
	Name string
	Age  uint32
}

type Employee struct {
	Person
	EmployeeID uint32
}

func (e *Employee) PrintInfo() {
	fmt.Println("Employee Info")
	fmt.Println("Name:", e.Name)
	fmt.Println("Age:", e.Age)
	fmt.Println("ID:", e.EmployeeID)
}
