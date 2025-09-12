/*
实现类型安全映射
假设有一个 books 表，包含字段 id 、 title 、 author 、 price 。
要求 ：
定义一个 Book 结构体，包含与 books 表对应的字段。
编写Go代码，使用Sqlx执行一个复杂的查询，例如查询价格大于 50 元的书籍，并将结果映射到 Book 结构体切片中，确保类型安全。
*/
package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Book struct {
	ID     uint `gorm:"unique"`
	Title  string
	Author string
	Price  float64
}

func main() {
	dsn := "root:kss@tcp(127.0.0.1:3306)/kss_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	//if err = db.AutoMigrate(&Book{}); err != nil {
	//	panic("failed to create table books")
	//}

	//books := []Book{
	//	{Title: "西游记", Author: "吴承恩", Price: 50.62},
	//	{Title: "三国演义", Author: "罗贯中", Price: 38.66},
	//	{Title: "水浒传", Author: "施耐庵", Price: 42.31},
	//	{Title: "红楼梦", Author: "曹雪芹", Price: 56.58},
	//}
	//db.Create(&books)

	var bks []Book
	db.Raw("SELECT * FROM books WHERE price >= ?", 50).Scan(&bks)
	fmt.Println(bks)
}
