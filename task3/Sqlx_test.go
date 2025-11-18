package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

/*
题目1：使用SQL扩展库进行查询
假设你已经使用Sqlx连接到一个数据库，并且有一个 employees 表，包含字段 id 、 name 、 department 、 salary 。
要求 ：
编写Go代码，使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息，并将结果映射到一个自定义的 Employee 结构体切片中。
编写Go代码，使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中。
*/

// 1. 基本结构和初始化

// Employee 结构体，使用 sqlx 标签映射数据库字段
type Employee struct {
	ID         int    `db:"id"`
	Name       string `db:"name"`
	Department string `db:"department"`
	Salary     int    `db:"salary"`
}

var db *sqlx.DB

func initDB() error {
	var err error
	// 链接数据库
	db, err = sqlx.Connect("mysql", "username:password@tcp(localhost:3306)/company?parseTime=true")
	if err != nil {
		return err
	}
	return nil
}

// 2. 查询技术部所有员工

// GetHighestPaidEmployee 查询工资最高的员工
func GetHighestPaidEmployee() ([]Employee, error) {
	var employee []Employee

	query := `SELECT id,name,department,salary FROM employees ORDER BY salary DESC limit 1 `

	err := db.Get(&employee, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("没有找到员工记录")
		}
		return nil, fmt.Errorf("查询最高工资员工失败: %v", err)
	}
	if len(employee) == 0 {
		return nil, fmt.Errorf("没有找到员工记录")
	}
	return employee, nil
}

func main() {
	if err := initDB(); err != nil {
		log.Fatal("数据库连接失败", err)
	}
	defer db.Close()
	// 查询技术部所有员工
	techEmployees, err := GetHighestPaidEmployee("技术部")
	if err != nil {
		log.Printf("查询技术部员工失败：%v", err)
	} else {
		for i, emp := range techEmployees {
			fmt.Printf("%d. ID: %d, 姓名: %s, 部门: %s, 工资: %d\n",
				i+1, emp.ID, emp.Name, emp.Department, emp.Salary)
		}
	}
	// 查询工资最高的员工
	fmt.Println("\n=== 查询员工最高工资 ===")
	highestPaid, err := GetHighestPaidEmployee()
	if err != nil {
		log.Printf("查询最高工资员工失败：v%", err)
	} else {
		for i, emp := range highestPaid {
			fmt.Printf("%d. ID: %d, 姓名: %s, 部门: %s, 工资: %d\n",
				i+1, emp.ID, emp.Name, emp.Department, emp.Salary)
		}
	}
}
