package main

import (
	"context"
	"errors"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	// "gorm.io/gorm/logger"
)

type Student struct {
	// gorm.Model
	Id    uint   `gorm:"primaryKey;autoIncrement"` // 主键且自增
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Grade string `json:"grade"`
}

func (Student) TableName() string {
	return "students"
}

func main() {
	/* 题目1：基本CRUD操作
	假设有一个名为 students 的表，包含字段 id （主键，自增）、 name （学生姓名，字符串类型）、 age （学生年龄，整数类型）、 grade （学生年级，字符串类型）。
	要求 ：
	编写SQL语句向 students 表中插入一条新记录，学生姓名为 "张三"，年龄为 20，年级为 "三年级"。
	编写SQL语句查询 students 表中所有年龄大于 18 岁的学生信息。
	编写SQL语句将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"。
	编写SQL语句删除 students 表中年龄小于 15 岁的学生记录。 */

	//创建数据库连接
	var db *gorm.DB
	dsn := "root:123456@tcp(127.0.0.1:3306)/rcs?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info),
		// SkipDefaultTransaction: true,//禁用默认事务（单条ddl）
	})
	if err != nil {
		panic("数据库连接失败：" + err.Error())
	}

	//假设有一个名为 students 的表，包含字段 id （主键，自增）、 name （学生姓名，字符串类型）、 age （学生年龄，整数类型）、 grade （学生年级，字符串类型）。
	db.AutoMigrate(&Student{})

	//编写SQL语句向 students 表中插入一条新记录，学生姓名为 "张三"，年龄为 20，年级为 "三年级"。
	db.Create(&Student{Name: "张三", Age: 20, Grade: "三年级"})
	db.Create(&Student{Name: "李四", Age: 21, Grade: "三年级"})
	db.Create(&Student{Name: "王五", Age: 22, Grade: "三年级"})

	//编写SQL语句查询 students 表中所有年龄大于 18 岁的学生信息。
	var students []Student
	db.Where("age > ?", 18).Find(&students)

	//编写SQL语句将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"。
	ctx := context.Background()
	gorm.G[Student](db).Where("name = ?", "张三").Update(ctx, "grade", "四年级")
	//db.Model(&Student{}).Where("name = ?", "张三").Update("grade", "四年级")

	//编写SQL语句删除 students 表中年龄小于 15 岁的学生记录。 */
	//db.Where("age < ?", 15).Delete(&Student{})
	db.Debug().Delete(&Student{}, "age < ?", 21)

	//--------------------------------------------------------------------------------------------------------------

	/* 题目2：事务语句
	假设有两个表： accounts 表（包含字段 id 主键， balance 账户余额）和 transactions 表（包含字段 id 主键， from_account_id 转出账户ID，
	to_account_id 转入账户ID， amount 转账金额）。
	*/

	type Account struct {
		Id      uint    `gorm:"primaryKey;autoIncrement"`
		Balance float64 `json:"balance"`
	}

	type Transaction struct {
		Id            uint    `gorm:"primaryKey;autoIncrement"`
		FromAccountId uint    `json:"from_account_id"`
		ToAccountId   uint    `json:"to_account_id"`
		Amount        float64 `json:"amount"`
	}

	/* 要求 ：
	编写一个事务，实现从账户 A 向账户 B 转账 100 元的操作。在事务中，需要先检查账户 A 的余额是否足够，如果足够则从账户 A 扣除 100 元，
	向账户 B 增加 100 元，并在 transactions 表中记录该笔转账信息。如果余额不足，则回滚事务。 */
	db.AutoMigrate(&Account{}, &Transaction{})

	/* db.Exec("DELETE FROM accounts")
	db.Exec("ALTER TABLE accounts AUTO_INCREMENT = 1")

	db.Exec("DELETE FROM transactions")
	db.Exec("ALTER TABLE transactions AUTO_INCREMENT = 1")

	db.Create(&Account{Balance: 1000}) //账户 A
	db.Create(&Account{Balance: 500})  //账户 B */

	// db.Commit()
	// db.Rollback()

	//语句执行失败，自动提交或回滚事务
	err1 := db.Transaction(func(tx *gorm.DB) error {
		var account Account
		db.Where("id = ?", 1).Find(&account)
		var account2 Account
		db.Where("id = ?", 2).Find(&account2)
		if account.Balance < 100 {
			return errors.New(strconv.Itoa(int(account.Id)) +"余额不足")
		}
		db.Model(&Account{}).Where("id = ?", 1).Update("balance", account.Balance-100)
		db.Model(&Account{}).Where("id = ?", 2).Update("balance", account2.Balance+100)
	
		db.Create(&Transaction{FromAccountId: 1, ToAccountId: 2, Amount: 100})
		return nil
	})

	if err1 != nil {
		panic("余额不足，回滚事务"+err1.Error())
	}

}
