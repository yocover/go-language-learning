package mysql

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// Account 表示银行账户
type Account struct {
	ID      int
	Name    string
	Balance float64
}

// TransferRequest 表示转账请求
type TransferRequest struct {
	FromID     int
	ToID       int
	Amount     float64
	SavePoint  bool    // 是否使用保存点
	ErrorAfter float64 // 在转账多少金额后模拟错误
}

// TestTransaction 演示MySQL事务的使用
func TestTransaction() {
	fmt.Println("=== 测试MySQL事务 ===")

	// 连接数据库
	dsn := "dbadmin:DbUs3r@2024!@tcp(localhost:3306)/myapp_db?parseTime=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}
	defer db.Close()

	// 测试连接
	err = db.Ping()
	if err != nil {
		log.Fatal("无法连接到数据库:", err)
	}
	fmt.Println("成功连接到数据库")

	// 1. 创建测试表
	createTable(db)

	// 2. 插入测试数据
	initTestData(db)

	// 3. 测试正常转账（事务成功）
	fmt.Println("\n1. 测试正常转账:")
	req1 := TransferRequest{
		FromID: 1,
		ToID:   2,
		Amount: 100,
	}
	transferMoney(db, req1)
	printAccounts(db)

	// 4. 测试转账失败（事务回滚）
	fmt.Println("\n2. 测试转账失败（余额不足）:")
	req2 := TransferRequest{
		FromID: 1,
		ToID:   2,
		Amount: 10000, // 超出余额
	}
	transferMoney(db, req2)
	printAccounts(db)

	// 5. 测试保存点
	fmt.Println("\n3. 测试保存点:")
	req3 := TransferRequest{
		FromID:     1,
		ToID:       2,
		Amount:     300,
		SavePoint:  true,
		ErrorAfter: 100,
	}
	transferMoneyWithSavepoint(db, req3)
	printAccounts(db)
}

// createTable 创建测试表
func createTable(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS accounts (
			id INT PRIMARY KEY,
			name VARCHAR(50),
			balance DECIMAL(10,2)
		)
	`)
	if err != nil {
		log.Fatal("创建表失败:", err)
	}
}

// initTestData 初始化测试数据
func initTestData(db *sql.DB) {
	// 清空表
	db.Exec("DELETE FROM accounts")

	// 插入测试数据
	accounts := []Account{
		{ID: 1, Name: "Alice", Balance: 1000},
		{ID: 2, Name: "Bob", Balance: 1000},
	}

	for _, acc := range accounts {
		_, err := db.Exec("INSERT INTO accounts (id, name, balance) VALUES (?, ?, ?)",
			acc.ID, acc.Name, acc.Balance)
		if err != nil {
			log.Fatal("插入数据失败:", err)
		}
	}
}

// transferMoney 执行转账操作
func transferMoney(db *sql.DB, req TransferRequest) error {
	// 开始事务
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("开始事务失败: %v", err)
	}

	// 检查转出账户余额
	var balance float64
	err = tx.QueryRow("SELECT balance FROM accounts WHERE id = ? FOR UPDATE",
		req.FromID).Scan(&balance)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("查询转出账户失败: %v", err)
	}

	// 检查余额是否足够
	if balance < req.Amount {
		tx.Rollback()
		return fmt.Errorf("余额不足")
	}

	// 更新转出账户
	_, err = tx.Exec("UPDATE accounts SET balance = balance - ? WHERE id = ?",
		req.Amount, req.FromID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("更新转出账户失败: %v", err)
	}

	// 更新转入账户
	_, err = tx.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?",
		req.Amount, req.ToID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("更新转入账户失败: %v", err)
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return fmt.Errorf("提交事务失败: %v", err)
	}

	fmt.Printf("转账成功: 从账户%d转出%.2f到账户%d\n", req.FromID, req.Amount, req.ToID)
	return nil
}

// transferMoneyWithSavepoint 使用保存点的转账操作
func transferMoneyWithSavepoint(db *sql.DB, req TransferRequest) error {
	// 开始事务
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("开始事务失败: %v", err)
	}

	// 创建保存点
	_, err = tx.Exec("SAVEPOINT before_transfer")
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("创建保存点失败: %v", err)
	}

	// 第一次转账
	_, err = tx.Exec("UPDATE accounts SET balance = balance - ? WHERE id = ?",
		req.ErrorAfter, req.FromID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("第一次转账失败: %v", err)
	}

	_, err = tx.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?",
		req.ErrorAfter, req.ToID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("第一次转账失败: %v", err)
	}

	// 模拟错误
	if req.ErrorAfter > 0 {
		// 回滚到保存点
		_, err = tx.Exec("ROLLBACK TO SAVEPOINT before_transfer")
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("回滚到保存点失败: %v", err)
		}
		fmt.Println("发生错误，回滚到保存点")
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return fmt.Errorf("提交事务失败: %v", err)
	}

	fmt.Println("带保存点的转账完成")
	return nil
}

// printAccounts 打印所有账户信息
func printAccounts(db *sql.DB) {
	rows, err := db.Query("SELECT id, name, balance FROM accounts")
	if err != nil {
		log.Fatal("查询账户失败:", err)
	}
	defer rows.Close()

	fmt.Println("\n当前账户状态:")
	fmt.Println("ID\t姓名\t余额")
	fmt.Println("----------------------")
	for rows.Next() {
		var acc Account
		if err := rows.Scan(&acc.ID, &acc.Name, &acc.Balance); err != nil {
			log.Fatal("读取数据失败:", err)
		}
		fmt.Printf("%d\t%s\t%.2f\n", acc.ID, acc.Name, acc.Balance)
	}
}

/* MySQL 事务的重要概念：

1. ACID 特性：
   - 原子性（Atomicity）：事务是不可分割的工作单位，要么全部完成，要么全部失败
   - 一致性（Consistency）：事务必须使数据库从一个一致性状态变到另一个一致性状态
   - 隔离性（Isolation）：多个事务并发执行时，一个事务的执行不应影响其他事务
   - 持久性（Durability）：事务一旦提交，其结果就是永久性的

2. 事务的基本操作：
   - BEGIN/START TRANSACTION：开始事务
   - COMMIT：提交事务
   - ROLLBACK：回滚事务
   - SAVEPOINT：创建保存点
   - ROLLBACK TO SAVEPOINT：回滚到保存点

3. 事务隔离级别：
   - READ UNCOMMITTED：读未提交
   - READ COMMITTED：读已提交
   - REPEATABLE READ：可重复读（MySQL默认级别）
   - SERIALIZABLE：串行化

4. 注意事项：
   - 事务应该尽可能短小
   - 避免长时间持有锁
   - 合理使用事务隔离级别
   - 注意死锁问题
   - 在适当的时候使用保存点
*/
