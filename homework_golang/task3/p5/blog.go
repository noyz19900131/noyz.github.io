/*
题目1：模型定义
假设你要开发一个博客系统，有以下几个实体： User （用户）、 Post （文章）、 Comment （评论）。
要求 ：
使用Gorm定义 User 、 Post 和 Comment 模型，
其中 User 与 Post 是一对多关系（一个用户可以发布多篇文章），
Post 与 Comment 也是一对多关系（一篇文章可以有多个评论）。
编写Go代码，使用Gorm创建这些模型对应的数据库表。

题目2：关联查询
基于上述博客系统的模型定义。
要求 ：
编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
编写Go代码，使用Gorm查询评论数量最多的文章信息。

题目3：钩子函数
继续使用博客系统的模型。
要求 ：
为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。
*/

package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string
	PostCount int
	Posts     []Post
}

type Post struct {
	gorm.Model
	Title         string
	Content       string
	CommentCount  int
	CommentStatus string
	UserID        uint
	Comments      []Comment
}

type Comment struct {
	gorm.Model
	Content string
	PostID  uint
}

//func (p *Post) AfterCreate(tx *gorm.DB) error {
//	var user User
//	if err := tx.First(&user, p.UserID).Error; err != nil {
//		return err
//	}
//
//	return tx.Model(&user).Update("post_count", user.PostCount+1).Error
//}

func (c *Comment) AfterDelete(tx *gorm.DB) error {
	// 1. 查询当前文章的剩余评论数量
	var commentCount int64
	if err := tx.Model(&Comment{}).Where("post_id = ?", c.PostID).Count(&commentCount).Error; err != nil {
		return fmt.Errorf("查询评论数量失败: %w", err)
	}

	// 2. 准备更新的数据
	updateData := map[string]interface{}{
		"comment_count": commentCount,
	}

	// 3. 如果评论数量为0，更新评论状态为"无评论"
	if commentCount == 0 {
		updateData["comment_status"] = "无评论"
	} else {
		updateData["comment_status"] = "有评论"
	}

	// 4. 更新文章信息
	if err := tx.Model(&Post{}).Where("id = ?", c.PostID).Updates(updateData).Error; err != nil {
		return fmt.Errorf("更新文章状态失败: %w", err)
	}

	return nil
}

func main() {
	dsn := "root:kss@tcp(127.0.0.1:3306)/kss_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if err = db.AutoMigrate(&User{}, &Post{}, &Comment{}); err != nil {
		panic("failed to create table")
	}

	//users := []User{
	//	{Username: "Jack"},
	//	{Username: "Mike"},
	//}
	//db.Create(&users)
	//
	//posts := []Post{
	//	{Title: "Jack_title", Content: "Jack_Content", UserID: users[0].ID},
	//	{Title: "Mike_title", Content: "Mike_Content", UserID: users[1].ID},
	//}
	//db.Create(&posts)

	// comments := []Comment{
	// 	{Content: "1", PostID: posts[0].ID},
	// 	{Content: "2", PostID: posts[0].ID},
	// 	{Content: "3", PostID: posts[1].ID},
	// 	{Content: "4", PostID: posts[1].ID},
	// 	{Content: "5", PostID: posts[1].ID},
	// 	{Content: "6", PostID: posts[2].ID},
	// 	{Content: "7", PostID: posts[3].ID},
	// 	{Content: "8", PostID: posts[3].ID},
	// 	{Content: "9", PostID: posts[3].ID},
	// 	{Content: "10", PostID: posts[3].ID},
	// 	{Content: "11", PostID: posts[4].ID},
	// 	{Content: "12", PostID: posts[5].ID},
	// }
	// db.Create(&comments)

	// var user User
	// db.Preload("Posts").Preload("Posts.Comments").Find(&user, 3)
	// fmt.Println(user.ID, user.Username, user.Email)
	// for _, post := range user.Posts {
	// 	fmt.Println(post.Title)
	// 	for _, comment := range post.Comments {
	// 		fmt.Println(comment.Content)
	// 	}
	// }

	//var post Post
	//err = db.Model(&Post{}).
	//	Select("posts.*, COUNT(comments.id) as comment_count").
	//	Joins("LEFT JOIN comments ON posts.id = comments.post_id").
	//	Group("posts.id").
	//	Order("comment_count DESC").
	//	First(&post).Error
	//if err != nil {
	//	panic("failed to find article")
	//} else {
	//	fmt.Println(post.Title)
	//}

	//post := Post{Title: "Mike_title2", Content: "Mike_content2", UserID: 2}
	//db.Create(&post)
	//var updatedUser [2]User
	//db.Find(&updatedUser[0], 1)
	//fmt.Println(updatedUser[0].PostCount)
	//db.Find(&updatedUser[1], 2)
	//fmt.Println(updatedUser[1].PostCount)

	// 创建测试数据
	user := User{Username: "测试用户"}
	db.Create(&user)

	post := Post{
		Title:         "测试文章",
		Content:       "这是一篇用于测试的文章",
		UserID:        user.ID,
		CommentCount:  0,
		CommentStatus: "无评论",
	}
	db.Create(&post)

	// 添加两条评论
	comment1 := Comment{Content: "第一条评论", PostID: post.ID}
	comment2 := Comment{Content: "第二条评论", PostID: post.ID}
	db.Create(&comment1)
	db.Create(&comment2)

	// 手动更新文章的初始评论状态（实际应用中可在添加评论时通过钩子处理）
	db.Model(&post).Updates(map[string]interface{}{
		"comment_count":  2,
		"comment_status": "有评论",
	})

	fmt.Println("删除第一条评论后:")
	db.Delete(&comment1)
	var postAfterFirstDelete Post
	db.First(&postAfterFirstDelete, post.ID)
	fmt.Printf("文章评论数量: %d, 评论状态: %s\n",
		postAfterFirstDelete.CommentCount, postAfterFirstDelete.CommentStatus)

	fmt.Println("删除第二条评论后:")
	db.Delete(&comment2)
	var postAfterSecondDelete Post
	db.First(&postAfterSecondDelete, post.ID)
	fmt.Printf("文章评论数量: %d, 评论状态: %s\n",
		postAfterSecondDelete.CommentCount, postAfterSecondDelete.CommentStatus)

}
