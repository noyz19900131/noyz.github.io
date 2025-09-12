package data

import (
	"blog/cfg"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null" json:"username" url:"username" form:"username"`
	Password string `gorm:"not null" json:"password" url:"password" form:"password"`
	Email    string `gorm:"unique;not null" json:"email" url:"email" form:"email"`
}

type Post struct {
	gorm.Model
	Title   string `gorm:"not null" json:"title" url:"title" form:"title"`
	Content string `gorm:"not null" json:"content" url:"content" form:"content"`
	UserID  uint   `gorm:"not 0" json:"user_id" url:"user_id" form:"user_id"`
}

type Comment struct {
	gorm.Model
	Content string `gorm:"not null" json:"content" url:"content" form:"content"`
	UserID  uint   `gorm:"not 0" json:"user_id" url:"user_id" form:"user_id"`
	PostID  uint   `gorm:"not 0" json:"post_id" url:"post_id" form:"post_id"`
}

// 数据库连接和建表
func ConnectDatabase() *gorm.DB {
	//连接数据库
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local",
		cfg.CFG.Db.User,
		cfg.CFG.Db.Password,
		cfg.CFG.Db.Host,
		cfg.CFG.Db.Port,
		cfg.CFG.Db.Dbname,
		cfg.CFG.Db.Charset,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil
	}

	//自动迁移（无表则建表，有表则更新结构）
	if err = db.AutoMigrate(&User{}, &Post{}, &Comment{}); err != nil {
		panic(err.Error())
		return nil
	}

	return db
}
