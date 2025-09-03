package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Comment struct {
	Id      uint   `gorm:"primaryKey;autoIncrement"`
	Content string `json:"content"`
	PostId  uint   `json:"post_id"`
}

func (Comment) TableName() string { return "comment" }

// 为 Comment 模型添加 AfterDelete 钩子函数
func (comment *Comment) AfterDelete(tx *gorm.DB) error {
	// 检查该文章是否还有其他评论
	var count int64
	tx.Model(&Comment{}).Where("post_id = ?", comment.PostId).Count(&count)

	// 如果没有评论了，更新文章状态为"无评论"
	if count == 0 {
		tx.Model(&Post{}).Where("id = ?", comment.PostId).Update("comment_status", "无评论")
	}

	return nil
}

type Post struct {
	Id            uint      `gorm:"primaryKey;autoIncrement"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	UserId        uint      `json:"user_id"`
	Comments      []Comment `gorm:"foreignKey:PostId;references:Id"`
	CommentStatus string    `json:"comment_status" gorm:"default:有评论"` // 添加评论状态字段
	User          User
}

func (Post) TableName() string { return "post" }

// 为 Post 模型添加 AfterCreate 钩子函数
func (post *Post) AfterCreate(tx *gorm.DB) error {
	// 在文章创建后，增加对应用户的文章数量统计
	if post.UserId > 0 {
		tx.Debug().Model(&User{}).Where("id = ?", post.UserId).UpdateColumn("post_count", gorm.Expr("post_count + ?", 1))
	}
	return nil
}

type User struct {
	Id        uint   `gorm:"primaryKey;autoIncrement"`
	Name      string `json:"name"`
	Posts     []Post `gorm:"foreignKey:UserId;references:Id"`
	PostCount int    `json:"post_count" gorm:"default:0"` // 添加文章数量统计字段
}

func (User) TableName() string { return "user" }

func main() {
	//创建数据库连接
	var db *gorm.DB
	dsn := "root:123456@tcp(127.0.0.1:3306)/rcs?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		// SkipDefaultTransaction: true,//禁用默认事务（单条ddl）
	})
	if err != nil {
		panic("数据库连接失败：" + err.Error())
	}

	/* 题目1：模型定义
	   假设你要开发一个博客系统，有以下几个实体： User （用户）、 Post （文章）、 Comment （评论）。
	   要求 ：
	   使用Gorm定义 User 、 Post 和 Comment 模型，其中 User 与 Post 是一对多关系（一个用户可以发布多篇文章）， Post 与 Comment 也是一对多关系（一篇文章可以有多个评论）。
	   编写Go代码，使用Gorm创建这些模型对应的数据库表。*/
	db.AutoMigrate(&User{}, &Post{}, &Comment{})
	/* db.AutoMigrate(&User{}, &Post{}, &Comment{})
	   db.Create(&User{Name: "张三"})
	   db.Create(&Post{Title: "文章1", Content: "内容1", UserId: 1})
	   db.Create(&Post{Title: "文章2", Content: "内容2", UserId: 1})
	   db.Create(&Comment{Content: "评论1", PostId: 1})
	   db.Create(&Comment{Content: "评论1", PostId: 2})
	   db.Create(&Comment{Content: "评论2", PostId: 2}) */

	/* user := User{Name: "张三",
	 Posts: []Post{Post{Title: "文章1", Content: "内容1",Comments: []Comment{Comment{Content: "评论1"}, Comment{Content: "评论2"}}},
	 			   Post{Title: "文章2", Content: "内容2",comments: []Comment{Comment{Content: "评论1"}, Comment{Content: "评论2"}}},
				}}} */

	/*题目2：关联查询
	  基于上述博客系统的模型定义。
	  要求 ：
	  编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
	  编写Go代码，使用Gorm查询评论数量最多的文章信息。*/
	var user User
	db.Debug().Where("id = ?", 1).Preload("Posts").Preload("Posts.Comments").Find(&user)
	fmt.Println("使用Gorm查询某个用户发布的所有文章及其对应的评论信息,", user)

	var post1 Post
	db.Where("Id=?", 1).Preload("User").Preload("Comments").Find(&post1)
	fmt.Println("查询指定post结构体文章的文章信息及用户信息", post1)

	var post Post
	db.Debug().Raw(`SELECT p.*, COUNT(c.id) as comment_count FROM post p 
LEFT JOIN comment c ON p.id = c.post_id 
GROUP BY p.id 
ORDER BY comment_count DESC 
LIMIT 1`).Scan(&post)

	fmt.Println("使用Gorm查询评论数量最多的文章信息,", post)

	/*题目3：钩子函数
	  继续使用博客系统的模型。
	  要求 ：
	  为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
	  为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。 */

	db.Debug().Create(&Post{Title: "文章3", Content: "内容3", UserId: 1})

	db.Delete(&Comment{}, 1) // 删除评论时触发

}
