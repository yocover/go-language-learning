package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// 清空表
func clearTables(db *sql.DB) error {
	// 由于有外键约束，需要先删除 orders 表的数据
	_, err := db.Exec("DELETE FROM orders")
	if err != nil {
		return err
	}

	// 然后删除 users 表的数据
	_, err = db.Exec("DELETE FROM users")
	if err != nil {
		return err
	}

	log.Println("表数据已清空")
	return nil
}

// 创建表
func createTables(db *sql.DB) error {
	// 创建用户表
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(50) NOT NULL UNIQUE,
			email VARCHAR(100) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	// 创建订单表
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS orders (
			id INT AUTO_INCREMENT PRIMARY KEY,
			user_id INT NOT NULL,
			product_name VARCHAR(100) NOT NULL,
			amount DECIMAL(10,2) NOT NULL,
			order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`)
	return err
}

// 插入测试数据
func insertTestData(db *sql.DB) error {
	// 插入用户数据并获取用户ID
	var userIDs = make(map[string]int64)

	// 插入 Alice
	result, err := db.Exec(`INSERT INTO users (username, email) VALUES ('alice', 'alice@example.com')`)
	if err != nil {
		return err
	}
	userIDs["alice"], _ = result.LastInsertId()

	// 插入 Bob
	result, err = db.Exec(`INSERT INTO users (username, email) VALUES ('bob', 'bob@example.com')`)
	if err != nil {
		return err
	}
	userIDs["bob"], _ = result.LastInsertId()

	// 插入 Charlie
	result, err = db.Exec(`INSERT INTO users (username, email) VALUES ('charlie', 'charlie@example.com')`)
	if err != nil {
		return err
	}
	userIDs["charlie"], _ = result.LastInsertId()

	log.Printf("插入了 3 个用户\n")

	// 插入订单数据
	_, err = db.Exec(`
		INSERT INTO orders (user_id, product_name, amount) VALUES 
		(?, 'iPhone', 999.99),
		(?, 'AirPods', 199.99),
		(?, 'MacBook', 1299.99),
		(?, 'iPad', 499.99)
	`, userIDs["alice"], userIDs["alice"], userIDs["bob"], userIDs["charlie"])
	if err != nil {
		return err
	}

	log.Printf("插入了 4 个订单\n")
	return nil
}

// 查询并显示数据
func queryData(db *sql.DB) error {
	// 查询用户及其订单
	rows, err := db.Query(`
		SELECT u.username, u.email, o.product_name, o.amount 
		FROM users u 
		LEFT JOIN orders o ON u.id = o.user_id
		ORDER BY u.username, o.order_date
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	log.Println("\n用户订单信息：")
	log.Println("----------------------------------------")
	for rows.Next() {
		var username, email, productName sql.NullString
		var amount sql.NullFloat64
		err := rows.Scan(&username, &email, &productName, &amount)
		if err != nil {
			return err
		}

		if productName.Valid {
			log.Printf("用户: %s (%s) - 订购: %s, 金额: %.2f\n",
				username.String, email.String, productName.String, amount.Float64)
		} else {
			log.Printf("用户: %s (%s) - 暂无订单\n",
				username.String, email.String)
		}
	}
	log.Println("----------------------------------------")

	return rows.Err()
}

func main() {
	// 使用 root 用户连接
	db, err := sql.Open("mysql", "root:MyR00t@2024!@tcp(localhost:3306)/myapp_db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 测试数据库连接
	err = db.Ping()
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	log.Println("数据库连接成功！")

	// 创建表
	if err := createTables(db); err != nil {
		log.Fatal("创建表失败:", err)
	}
	log.Println("表创建成功！")

	// 清空表
	if err := clearTables(db); err != nil {
		log.Fatal("清空表失败:", err)
	}

	// 插入测试数据
	if err := insertTestData(db); err != nil {
		log.Fatal("插入测试数据失败:", err)
	}
	log.Println("测试数据插入成功！")

	// 查询并显示数据
	if err := queryData(db); err != nil {
		log.Fatal("查询数据失败:", err)
	}
}
