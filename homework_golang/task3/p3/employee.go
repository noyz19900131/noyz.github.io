/*
使用SQL扩展库进行查询
假设你已经使用Sqlx连接到一个数据库，并且有一个 employees 表，包含字段 id 、 name 、 department 、 salary 。
要求 ：
编写Go代码，使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息，并将结果映射到一个自定义的 Employee 结构体切片中。
编写Go代码，使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中。
*/
package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Employees struct {
	ID         uint `gorm:"unique"`
	Name       string
	Department string
	Salary     float64
}

func main() {
	dsn := "root:kss@tcp(127.0.0.1:3306)/kss_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// if err := db.AutoMigrate(&Employees{}); err != nil {
	// 	panic("failed to create table employees")
	// }

	// emps := []Employees{
	// 	{Name: "Jack", Department: "技术部", Salary: 8800},
	// 	{Name: "Mark", Department: "人事部", Salary: 9000},
	// 	{Name: "Harry", Department: "财务部", Salary: 8900},
	// }
	// db.Create(&emps)

	var emp Employees
	db.Raw("SELECT * FROM employees WHERE department = ?", "技术部").Scan(&emp)
	fmt.Println(emp)

	db.Raw("SELECT * FROM employees ORDER BY salary DESC LIMIT 0,1").Scan(&emp)
	fmt.Println(emp)
}
