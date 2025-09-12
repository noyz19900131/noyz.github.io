/*
题目2：事务语句
假设有两个表： accounts 表（包含字段 id 主键， balance 账户余额）和 transactions 表
（包含字段 id 主键， from_account_id 转出账户ID， to_account_id 转入账户ID， amount 转账金额）。
要求 ：
编写一个事务，实现从账户 A 向账户 B 转账 100 元的操作。
在事务中，需要先检查账户 A 的余额是否足够，如果足够则从账户 A 扣除 100 元，向账户 B 增加 100 元，
并在 transactions 表中记录该笔转账信息。如果余额不足，则回滚事务。
*/
package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Accounts struct {
	ID      uint `gorm:"primaryKey"`
	Balance float64
}

type Transactions struct {
	ID              uint `gorm:"primaryKey"`
	From_account_id uint
	To_account_id   uint
	Amount          float64
}

func main() {
	dsn := "root:kss@tcp(127.0.0.1:3306)/kss_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	//建表
	// if err := db.AutoMigrate(&Accounts{}); err != nil {
	// 	panic("failed to create table accounts")
	// }
	// if err = db.AutoMigrate(&Transactions{}); err != nil {
	// 	panic("failed to create table transactions")
	// }

	//插入数据
	// accounts := []Accounts{{Balance: 101}, {Balance: 200}}
	// db.Create(&accounts)

	//账户A给账户B转账
	var a, b Accounts
	tx := db.Begin()
	tx.SavePoint("sp")
	tx.Where("id = ?", 1).Find(&a)
	tx.Where("id = ?", 2).Find(&b)
	if a.Balance < 100 {
		tx.RollbackTo("sp")
	} else {
		tx.Model(&Accounts{}).Where("id = ?", 1).Update("balance", a.Balance-100)
		tx.Model(&Accounts{}).Where("id = ?", 2).Update("balance", b.Balance+100)
		tx.Create(&Transactions{From_account_id: a.ID, To_account_id: b.ID, Amount: 100})
	}
	tx.Commit()
}
