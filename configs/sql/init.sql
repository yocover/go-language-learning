-- 创建账户表
CREATE TABLE IF NOT EXISTS accounts (
    id INT PRIMARY KEY,
    name VARCHAR(50),
    balance DECIMAL(10,2)
);

-- 插入初始数据
INSERT INTO accounts (id, name, balance) VALUES
    (1, 'Alice', 1000.00),
    (2, 'Bob', 1000.00)
ON DUPLICATE KEY UPDATE
    name = VALUES(name),
    balance = VALUES(balance); 