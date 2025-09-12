/*
题目 ：定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。
然后创建 Rectangle 和 Circle 结构体，实现 Shape 接口。
在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。

考察点 ：接口的定义与实现、面向对象编程风格。
*/
package main

import (
	"fmt"
	"math"
)

func main() {
	rectangle := Rectangle{2.5, 2}
	circle := Circle{1.1}
	fmt.Println(rectangle.Area(), rectangle.Perimeter(), circle.Area(), circle.Perimeter())
}

type Shape interface {
	Area()
	Perimeter()
}

type Rectangle struct {
	length float64
	width  float64
}

type Circle struct {
	radius float64
}

func (r *Rectangle) Area() float64 {
	return r.length * r.width
}

func (r *Rectangle) Perimeter() float64 {
	return 2 * (r.length + r.width)
}

func (c *Circle) Area() float64 {
	return math.Pow(c.radius, 2) * math.Pi
}

func (c *Circle) Perimeter() float64 {
	return 2 * c.radius * math.Pi
}
