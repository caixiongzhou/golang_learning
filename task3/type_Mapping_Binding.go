package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

/*
题目2：实现类型安全映射
假设有一个 books 表，包含字段 id 、 title 、 author 、 price 。
要求 ：
定义一个 Book 结构体，包含与 books 表对应的字段。
编写Go代码，使用Sqlx执行一个复杂的查询，例如查询价格大于 50 元的书籍，并将结果映射到 Book 结构体切片中，确保类型安全。
*/

// Book 结构体，使用 sqlx 标签精确映射数据库字段
type Book struct {
	ID     int    `db:"id" json:"id"`
	Title  string `db:"title" json:"title"`
	Author string `db:"author" json:"author"`
	Price  int    `db:"price" json:"price"`
}

var db *sqlx.DB

// initDB 初始化数据库连接
func initDB() error {
	var err error
	// 连接数据库
	db, err := sqlx.Connect("mysql", "username:password@tcp(localhost:3306)/bookstore?parseTime=true")
	if err != nil {
		return fmt.Errorf("数据库连接失败：%v", err)
	}
	// 设置连接池参数
	db.SetMaxOpenConns(25) // 最大打开连接数
	db.SetMaxIdleConns(5)  // 最大空闲连接数
	return nil
}

// 2. 类型安全的查询方法
// GetExpensiveBooks 查询价格大于指定值的书籍
func GetExpensiveBooks(minPrice float64) ([]Book, error) {
	var books []Book

	query := `SELECT id,title,author,price FROM books WHERE price >? ORDER BY price DESC `
	err := db.Select(&books, query, minPrice)
	if err != nil {
		return nil, fmt.Errorf("查询高价书籍失败:%v", err)
	}
	return books, nil
}

// GetBooksByAuthor 查询指定作者的所有书籍
func GetBooksByAuthor(authorName string) ([]Book, error) {
	var books []Book
	query := `SELECT id,title,author,price FROM books WHERE author = ? ORDER BY price DESC`
	err := db.Select(&books, query, authorName)
	if err != nil {
		return nil, fmt.Errorf("查询作者书籍失败：%v", err)
	}
	return books, nil
}

// 3. 使用命名查询增强类型安全

// BookQuery 查询参数结构体
type BookQuery struct {
	MinPrice  float64 `db:"min_price"`
	MaxPrice  float64 `db:"max_price"`
	Author    string  `db:"author"`
	TitleLike string  `db:"title_like"`
}

// QueryBooks 复杂的多条件查询
func QueryBooks(params BookQuery) ([]Book, error) {
	var book []Book
	query := `SELECT id, title, author, price FROM books 
              WHERE 1=1 
 			  {{if .MinPrice}} AND price >= :min_price {{end}}
 			  {{if .MaxPrice}} AND price <= :max_price {{end}}
 			  {{if .Author}} AND author = :author {{end}}
 			  {{if .TitleLike}} AND title LIKE :title_like {{end}}
 			  ORDER BY price DESC`

	// 使用 sqlx.Named 处理命名参数
	namedQuery, args, err := sqlx.Named(query, params)
	if err != nil {
		return nil, fmt.Errorf("构建命名查询失败：%v", err)
	}

	// 重新绑定查询以适应具体数据库
	query = db.Rebind(namedQuery)
	err = db.Select(&book, query, args...)
	if err != nil {
		return nil, fmt.Errorf("执行查询失败：%v", err)
	}
	return book, err
}

// 使用事务的复杂操作

// UpdateBookPrice 更新书籍价格（带事务）
func UpdateBookPrice(bookID int, newPrice float64) error {
	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("开始事务失败: %v", err)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()
	// 检查书籍是否存在
	var existingBook Book
	err = tx.Get(&existingBook, "SELECT FORM books WHERE id = ?", bookID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("书籍不存在:%v", err)
	}

	// 更新价格                        ?
	tx.Exec("UPDATE books SET price = ? WHERE id = ? ", newPrice, bookID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("更新价格失败：%v", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return fmt.Errorf("提交事务失败：%v", err)
	}
	return nil
}

func main() {
	// 初始化数据库连接
	if err := initDB(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 1. 查询价格大于50元的书籍
	fmt.Println("=== 价格大于50元的书籍 ===")
	expensiveBooks, err := GetExpensiveBooks(50)
	if err != nil {
		log.Printf("查询失败：%v", err)
	} else {
		for i, book := range expensiveBooks {
			fmt.Printf("%d. 《%s》- %s ￥%.2f\n", i+1, book.Title, book.Author, book.Price)
		}
	}

	// 2. 复杂查询示例
	fmt.Println("\n=== 复杂查询：价格在30-80元之间的书籍 ===")
	queryParams := BookQuery{
		MinPrice: 30.0,
		MaxPrice: 80.0,
	}
	complexResult, err := QueryBooks(queryParams)
	if err != nil {
		for i, book := range complexResult {
			fmt.Printf("%d. 《%s》- %s ￥%.2f\n", i+1, book.Title, book.Author, book.Price) // i+1 的意思是将索引值转换为人类可读的序号,i也可以，从0开始而已
		}
	}
}
