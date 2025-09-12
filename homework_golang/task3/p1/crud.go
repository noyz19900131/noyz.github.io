/*
题目1：基本CRUD操作
假设有一个名为 students 的表，包含字段 id （主键，自增）、 name （学生姓名，字符串类型）、 age （学生年龄，整数类型）、 grade （学生年级，字符串类型）。
要求 ：
编写SQL语句向 students 表中插入一条新记录，学生姓名为 "张三"，年龄为 20，年级为 "三年级"。
编写SQL语句查询 students 表中所有年龄大于 18 岁的学生信息。
编写SQL语句将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"。
编写SQL语句删除 students 表中年龄小于 15 岁的学生记录。
*/

package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Students struct {
	ID    uint `gorm:"primaryKey;autoIncrement"`
	Name  string
	Age   uint
	Grade string
}

func main() {
	dsn := "root:kss@tcp(127.0.0.1:3306)/kss_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	//if err = db.AutoMigrate(&Students{}); err != nil {
	//	panic("failed to create table")
	//}

	/* 向 students 表中插入一条新记录，学生姓名为 "张三"，年龄为 20，年级为 "三年级" */
	//student := Students{Name: "张三", Age: 20, Grade: "三年级"}
	//result := db.Create(&student)
	//if result.Error != nil {
	//	panic("failed to insert data")
	//}

	/* 查询 students 表中所有年龄大于 18 岁的学生信息 */
	//var student Students
	//db.Where("age > ?", 18).Find(&student)
	//fmt.Println(student)

	/* 将 students 表中姓名为 "张三" 的学生年级更新为 "四年级" */
	//db.Model(&Students{}).Where("name = ?", "张三").Update("grade", "四年级")

	/* 删除 students 表中年龄小于 15 岁的学生记录 */
	db.Where("age < ?", 21).Delete(&Students{})
}
