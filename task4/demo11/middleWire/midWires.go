package middleWire

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// 前置打印中间件（练习）
func LogPrintMiddlewire() gin.HandlerFunc {

	return func(c *gin.Context) {
		fmt.Println("前置打印中间件")
		c.Next()
		// c.Abort()
		fmt.Println("c.Next()后置打印中间件") //c.Next()后执行
	}
}
