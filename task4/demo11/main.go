package main

import (
	"log"
	"net/http"
	"os"

	"web3Demo/task4/demo11/jwtToken"
	"web3Demo/task4/demo11/middleWire"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	gorm.Model
	Username string    `gorm:"unique;not null"`
	Password string    `gorm:"not null"`
	Email    string    `gorm:"unique;not null"`
	Posts    []Post    `gorm:"foreignKey:UserID;references:ID"`
	Comments []Comment `gorm:"foreignKey:UserID;references:ID"`
}

type Post struct {
	gorm.Model
	Title    string `validate:"required" gorm:"not null"`
	Content  string `validate:"required" gorm:"not null"`
	UserID   uint
	User     User
	Comments []Comment `gorm:"foreignKey:PostID;references:ID"`
}

type Comment struct {
	gorm.Model
	Content string `gorm:"not null"`
	UserID  uint
	User    User
	PostID  uint
	Post    Post
}

// HashPassword 对密码进行加密
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword 验证密码
func CheckPassword(password, hash string) bool {
	log.Println("验证密码")
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func main() {

	/* 6.错误处理与日志记录
	   对可能出现的错误进行统一处理，如数据库连接错误、用户认证失败、文章或评论不存在等，返回合适的 HTTP 状态码和错误信息。
	   使用日志库记录系统的运行信息和错误信息，方便后续的调试和维护。 */
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("无法打开日志文件:", err)
	}
	defer file.Close()
	log.SetOutput(file)

	log.Println("demo1开始")
	//创建数据库连接
	dsn := "root:123456@tcp(127.0.0.1:3306)/task4?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Println("数据库连接错误：" + err.Error())
		panic("数据库连接错误：" + err.Error())
	}

	/* 2.设计数据库表结构，至少包含以下几个表：
	   users 表：存储用户信息，包括 id 、 username 、 password 、 email 等字段。
	   posts 表：存储博客文章信息，包括 id 、 title 、 content 、 user_id （关联 users 表的 id ）、 created_at 、 updated_at 等字段。
	   comments 表：存储文章评论信息，包括 id 、 content 、 user_id （关联 users 表的 id ）、 post_id （关联 posts 表的 id ）、 created_at 等字段。
	   使用 GORM 定义对应的 Go 模型结构体。 */
	// 自动迁移模型
	db.AutoMigrate(&User{}, &Post{}, &Comment{})
	log.Println("数据库表结构创建成功")

	ginEngin := gin.Default() //创建Gin引擎

	//软删除错误数据
	//db.Delete(&Post{}, 2)

	/* 3.用户认证与授权
	   实现用户注册和登录功能，用户注册时需要对密码进行加密存储，登录时验证用户输入的用户名和密码。
	   使用 JWT（JSON Web Token）实现用户认证和授权，用户登录成功后返回一个 JWT，后续的需要认证的接口需要验证该 JWT 的有效性。 */
	userRegister := ginEngin.Group("/userRegister") //定义用户注册和登录路由组
	//用户注册
	userRegister.POST("/register", func(c *gin.Context) {
		var user User
		err := c.ShouldBindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if user.Username == "" || user.Password == "" || user.Email == "" {
			c.JSON(http.StatusBadRequest, "用户名、密码和邮箱为必填项")
			return
		}
		var userQuery User
		if db.Where("username=?", user.Username).Find(&userQuery).RowsAffected > 0 {
			log.Println("用户已存在")
			c.JSON(http.StatusBadRequest, gin.H{"error": "用户已存在"})
			return
		}
		password, errPass := HashPassword(user.Password)
		if errPass != nil {
			log.Println("用户注册时密码加密失败", errPass)
		}
		user.Password = password
		result := db.Create(&user)
		if result.Error != nil {
			log.Println("用户注册失败", result.Error.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "用户注册失败" + result.Error.Error()})
			return
		}
		log.Println("用户注册成功")
		c.JSON(http.StatusOK, gin.H{"message": "用户注册成功"})
	})

	//用户登录
	userRegister.POST("/login", func(c *gin.Context) {
		var user User
		err := c.ShouldBindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if user.Username == "" || user.Password == "" {
			c.JSON(http.StatusBadRequest, "用户名和密码为必填项")
			return
		}
		var userQuery User
		if db.Where("username=?", user.Username).Find(&userQuery).RowsAffected == 0 {
			log.Println("用户不存在")
			c.JSON(http.StatusBadRequest, gin.H{"error": "用户不存在"})
			return
		}
		if !CheckPassword(user.Password, userQuery.Password) {
			log.Println("密码错误")
			c.JSON(http.StatusBadRequest, gin.H{"error": "密码错误"})
			return
		}

		token, err := jwtToken.GenerateJWT(user.ID, user.Username)
		if err != nil {
			log.Println("生成JWT失败", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "生成JWT失败" + err.Error()})
			return
		}

		log.Println("用户登录成功")
		c.JSON(http.StatusOK, gin.H{"message": "用户登录成功", "token": token})
	})

	/* 4.文章管理功能
	   实现文章的创建功能，只有已认证的用户才能创建文章，创建文章时需要提供文章的标题和内容。
	   实现文章的读取功能，支持获取所有文章列表和单个文章的详细信息。
	   实现文章的更新功能，只有文章的作者才能更新自己的文章。
	   实现文章的删除功能，只有文章的作者才能删除自己的文章。 */
	postsManage := ginEngin.Group("/postsManage") //定义文章管理路由组
	// postsManage.Use(middleWire.AuthMiddlewir())
	//文章创建
	postsManage.POST("/create", middleWire.AuthMiddlewir(), func(c *gin.Context) {
		var post Post
		err := c.ShouldBindJSON(&post)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if post.Content == "" || post.Title == "" {
			c.JSON(http.StatusBadRequest, "文章的标题和内容为必输项")
			return
		}

		db.Create(&post)
		log.Println("文章创建成功")
		c.JSON(http.StatusOK, gin.H{"message": "文章创建成功"})
	})
	//文章读取
	postsManage.POST("/read", func(c *gin.Context) {
		var post Post
		err := c.ShouldBindJSON(&post)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if post.UserID != 0 {
			var postList []Post
			db.Unscoped().Where("user_id=?", post.UserID).Preload("User").Find(&postList)
			log.Println("文章多笔读取成功")
			c.JSON(http.StatusOK, gin.H{"message": "文章读取成功", "多笔查询": postList})
			return
		}
		if post.ID != 0 {
			var postRet Post
			db.Unscoped().Where("id=?", post.ID).Preload("User").Find(&postRet)
			log.Println("文章单笔读取成功")
			c.JSON(http.StatusOK, gin.H{"message": "文章读取成功", "单笔查询": postRet})
			return
		}

	})
	//文章更新
	postsManage.POST("/update", func(ctx *gin.Context) {
		var post Post
		err := ctx.ShouldBind(&post)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var postQuery Post
		postResult := db.Where("id=?", post.ID).Find(&postQuery)
		if postResult.RowsAffected == 0 {
			log.Println("文章不存在")
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "文章不存在"})
			return
		}
		if postResult.RowsAffected > 0 && postQuery.UserID != post.UserID {
			log.Println("只有文章作者才能更新文章")
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "只有文章作者才能更新文章"})
			return
		}
		db.Model(&post).Where("id=?", post.ID).Updates(Post{Title: post.Title, Content: post.Content})
		// db.Model(&post).Where("id=?", post.ID).Updates(gin.H{"Title": post.Title,"Content": post.Content})
		ctx.JSON(http.StatusOK, gin.H{"message": "文章更新成功"})
	})
	//文章删除
	postsManage.POST("/delete", func(ctx *gin.Context) {
		var post Post
		err := ctx.ShouldBind(&post)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var postQuery Post
		postResult := db.Where("id=?", post.ID).Find(&postQuery)
		if postResult.RowsAffected == 0 {
			log.Println("文章不存在")
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "文章不存在"})
			return
		}
		if postResult.RowsAffected > 0 && postQuery.UserID != post.UserID {
			log.Println("只有文章作者才能删除文章")
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "只有文章作者才能删除文章"})
			return
		}

		db.Delete(&post, post.ID) //软删除
		ctx.JSON(http.StatusOK, gin.H{"message": "文章删除成功"})
	})

	/* 5.评论功能
	实现评论的创建功能，已认证的用户可以对文章发表评论。
	实现评论的读取功能，支持获取某篇文章的所有评论列表。 */
	commentsManage := ginEngin.Group("/commentsManage")
	// commentsManage.Use(middleWire.AuthMiddlewir())
	commentsManage.POST("/create", middleWire.AuthMiddlewir(), middleWire.LogPrintMiddlewire(), func(ctx *gin.Context) {
		var comment Comment
		err := ctx.ShouldBindJSON(&comment)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result := db.Create(&comment)
		if result.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "评论创建成功"})
	})
	commentsManage.POST("/read", middleWire.LogPrintMiddlewire(), func(ctx *gin.Context) {
		var comment Comment
		err := ctx.ShouldBindJSON(&comment)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var commentList []Comment
		result := db.Preload("User").Preload("Post").Unscoped().Where("post_id=?", comment.PostID).Find(&commentList)
		if result.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
			return
		}
		if len(commentList) == 0 {
			log.Println("没有查询到数据")
			ctx.JSON(http.StatusOK, gin.H{"message": "没有查询到数据"})
			return
		}
		log.Println("评论读取成功")
		ctx.JSON(http.StatusOK, gin.H{"message": "评论读取成功", "多笔查询": commentList})
	})

	log.Println("启动服务成功")
	ginEngin.Run(":8080") //	启动gin服务
}
