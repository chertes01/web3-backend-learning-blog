package main

import (
	"fmt"
	"log"
	"myproject/config"
	"myproject/models"

	"github.com/jmoiron/sqlx"
)

func createTable(db *sqlx.DB) error {
	// 执行创建 employees 表的语句
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS employees (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		department VARCHAR(100),
		salary DECIMAL(10, 2) NOT NULL DEFAULT 0.00	
	);`)
	if err != nil {
		// 如果失败，返回一个带上下文的错误
		return fmt.Errorf("创建 'employees' 表失败: %w", err)
	}
	log.Println("'employees' 表创建成功 (或已存在)")

	// 执行创建 books 表的语句
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS books (
		id INT AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		author VARCHAR(100) NOT NULL,
		price DECIMAL(10, 2) NOT NULL DEFAULT 0.00
	);`)
	if err != nil {
		// 如果失败，返回一个带上下文的错误
		return fmt.Errorf("创建 'books' 表失败: %w", err)
	}
	log.Println("'books' 表创建成功 (或已存在)")

	// 所有表都创建成功
	return nil
}

func updateEmployeeSalary(db *sqlx.DB, employees []models.Employee) error {
	if len(employees) == 0 {
		log.Println("没有员工记录需要更新。")
		return nil
	}
	//准备更新语句
	sqlStr := "insert into employees set id=?, name=?, department=?, salary=?"
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		return fmt.Errorf("准备更新语句失败: %v", err)
	}
	defer stmt.Close()

	for _, emp := range employees {
		_, err := stmt.Exec(emp.ID, emp.Name, emp.Department, emp.Salary)
		if err != nil {
			return fmt.Errorf("更新员工 ID %d 失败: %v", emp.ID, err)
		}
		log.Printf("成功更新员工 ID %d", emp.ID)
	}
	return nil
}

// 查询 employees 表中所有部门为 "技术部" 的员工信息，并将结果映射到一个自定义的 Employee 结构体切片中。
func queryTechDepartmentEmployees(db *sqlx.DB) ([]models.Employee, error) {
	var employees []models.Employee
	err := db.Select(&employees, "select id,name,department,salary from employees where department=?", "技术部")
	if err != nil {
		return nil, fmt.Errorf("查询技术部员工失败: %v", err)
	}
	// 输出查询结果

	return employees, nil
}

//查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中

func queryHighestSalaryEmployee(db *sqlx.DB) (*models.Employee, error) {
	var employee models.Employee
	err := db.Get(&employee, "select id,name,department,salary from employees order by salary desc limit 1")
	if err != nil {
		return nil, fmt.Errorf("查询工资最高的员工失败: %v", err)
	}
	// 输出查询结果
	fmt.Println("工资最高的员工信息:")
	fmt.Printf("员工ID:%d,姓名:%s,部门:%s,薪水:%.2f\n", employee.ID, employee.Name, employee.Department, employee.Salary)
	return &employee, nil
}

// 定义一个 Book 结构体，包含与 books 表对应的字段
var book []models.Books

func updateBooks(db *sqlx.DB, books []models.Books) error {
	if len(books) == 0 {
		log.Println("没有书籍记录需要插入。")
		return nil
	}

	// 准备插入语句
	sqlStr := "INSERT INTO books (id, title, author, price) VALUES (?, ?, ?, ?)"
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		return fmt.Errorf("准备插入语句失败: %v", err)
	}
	defer stmt.Close()

	for _, book := range books {
		_, err := stmt.Exec(book.ID, book.Title, book.Author, book.Price)
		if err != nil {
			return fmt.Errorf("插入书籍 ID %d 失败: %v", book.ID, err)
		}
		log.Printf("成功插入书籍 ID %d：《%s》", book.ID, book.Title)
	}
	return nil
}

// 查询价格大于 50 元的书籍，并将结果映射到 Book 结构体切片中，确保类型安全
func queryExpensiveBooks(db *sqlx.DB) ([]models.Books, error) {
	err := db.Select(&book, "SELECT id, title, author, price FROM books WHERE price > ?", 70.0)
	if err != nil {
		return nil, fmt.Errorf("查询价格大于 50 元的书籍失败: %v", err)
	}
	return book, nil
}

func main() {
	config.CreatDB()

	err := config.InitDB()
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	// //创建数据表
	// err = createTable(config.DB)
	// if err != nil {
	// 	log.Fatalf("创建数据表失败: %v", err)
	// }

	// //更新员工信息
	// employeesToUpdate := []models.Employee{
	// 	{ID: 1, Name: "David", Department: "技术部", Salary: 9453.12},
	// 	{ID: 2, Name: "Eve", Department: "财务部", Salary: 6487.33},
	// 	{ID: 3, Name: "Grace", Department: "销售部", Salary: 7850.67},
	// 	{ID: 4, Name: "Helen", Department: "技术部", Salary: 8120.40},
	// 	{ID: 5, Name: "Bob", Department: "人力资源部", Salary: 7312.55},
	// }
	// err = updateEmployeeSalary(config.DB, employeesToUpdate)
	// if err != nil {
	// 	log.Fatalf("更新员工失败: %v", err)
	// }
	// log.Println("所有员工更新完成。")
	//查询技术部员工信息
	_, err = queryTechDepartmentEmployees(config.DB)
	if err != nil {
		log.Fatalf("查询技术部员工失败: %v", err)
	}
	//查询工资最高的员工信息
	_, err = queryHighestSalaryEmployee(config.DB)
	if err != nil {
		log.Fatalf("查询工资最高的员工失败: %v", err)
	}

	// booksToUpdate := []models.Books{
	// 	{ID: 1, Title: "The Silent Forest", Author: "Tom", Price: 68.50},
	// 	{ID: 2, Title: "Go Programming", Author: "Linda", Price: 89.90},
	// 	{ID: 3, Title: "Quantum Dreams", Author: "Zhang Wei", Price: 77.30},
	// 	{ID: 4, Title: "The Last Kingdom", Author: "John", Price: 59.80},
	// 	{ID: 5, Title: "Deep Learning Simplified", Author: "Sophia", Price: 95.60},
	// }
	// err = updateBooks(config.DB, booksToUpdate)
	// if err != nil {
	// 	log.Fatalf("更新失败: %v", err)
	// }
	// log.Println("更新完成。")
	//查询价格大于 50 元的书籍
	book, err := queryExpensiveBooks(config.DB)
	if err != nil {
		log.Fatalf("查询价格大于 50 元的书籍失败: %v", err)
	}
	// 输出查询结果
	fmt.Println("价格大于 50 元的书籍:")
	for _, b := range book {
		fmt.Printf("书籍ID:%d,标题:%s,作者:%s,价格:%.2f\n", b.ID, b.Title, b.Author, b.Price)
	}

}
