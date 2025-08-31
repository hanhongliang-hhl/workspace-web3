package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql" // 导入MySQL驱动并自动注册
	"github.com/jmoiron/sqlx"
)

func main() {
	//创建数据库连接
	dsn := "root:123456@tcp(127.0.0.1:3306)/rcs?charset=utf8mb4&parseTime=True&loc=Local"
	db := sqlx.MustOpen("mysql", dsn)
	defer db.Close() // 记得关闭数据库连接

	/* 题目1：使用SQL扩展库进行查询
	假设你已经使用Sqlx连接到一个数据库，并且有一个 employees 表，包含字段 id 、 name 、 department 、 salary 。
	要求 ：
	编写Go代码，使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息，并将结果映射到一个自定义的 Employee 结构体切片中。
	编写Go代码，使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中。*/
	db.Exec("CREATE TABLE IF NOT EXISTS employees (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255), department VARCHAR(255), salary DECIMAL(10, 2))")
	type Employee struct {
		Id         int     `db:"id"`
		Name       string  `db:"name"`
		Department string  `db:"department"`
		Salary     float64 `db:"salary"`
	}

	db.Exec("INSERT INTO employees (name, department, salary) VALUES (?, ?, ?)", "张三", "技术部", 5000.00)
	db.Exec("INSERT INTO employees (name, department, salary) VALUES (?, ?, ?)", "张三2", "科技部", 7000.00)
	db.Exec("INSERT INTO employees (name, department, salary) VALUES (?, ?, ?)", "张三3", "业务部", 8000.00)
	db.Exec("INSERT INTO employees (name, department, salary) VALUES (?, ?, ?)", "张三4", "技术部", 1000.00)

	var employees []Employee
	db.Select(&employees, "SELECT * FROM employees WHERE department = ?", "技术部")
	for _, employee := range employees {
		fmt.Println(employee.Name, employee.Department, employee.Salary)
	}

	var employee Employee
	db.Select(&employee, "SELECT * FROM employees ORDER BY salary DESC LIMIT 1")
	fmt.Println(employee.Name, employee.Department, employee.Salary)

	/*题目2：实现类型安全映射
	  假设有一个 books 表，包含字段 id 、 title 、 author 、 price 。
	  要求 ：
	  定义一个 Book 结构体，包含与 books 表对应的字段。
	  编写Go代码，使用Sqlx执行一个复杂的查询，例如查询价格大于 50 元的书籍，并将结果映射到 Book 结构体切片中，确保类型安全。 */
	db.Exec("CREATE TABLE IF NOT EXISTS books (id INT AUTO_INCREMENT PRIMARY KEY, title VARCHAR(255), author VARCHAR(255), price DECIMAL(10, 2))")
	type Book struct {
		Id     int     `db:"id"`
		Title  string  `db:"title"`
		Author string  `db:"author"`
		Price  float64 `db:"price"`
	}
	db.Exec("INSERT INTO books (title, author, price) VALUES (?, ?, ?)", "《Go语言实战》", "小王子", 50.00)
	db.Exec("INSERT INTO books (title, author, price) VALUES (?, ?, ?)", "《Go语言实战》", "小王子1", 50.00)
	var books []Book
	db.Select(&books, "SELECT * FROM books WHERE price > ?", 50.00)
	for _, book := range books {
		fmt.Println(book.Title, book.Author, book.Price)
	}

}
