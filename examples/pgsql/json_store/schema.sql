-- 创建扩展（如果需要）
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 创建文件数据表
CREATE TABLE IF NOT EXISTS dm_excel_json_store (
    id SERIAL PRIMARY KEY,                         -- 自增主键
    file_id VARCHAR(255) NOT NULL UNIQUE,          -- 文件唯一标识符
    file_name VARCHAR(255) NOT NULL,               -- 文件名称
    json_content TEXT NOT NULL,                    -- JSON内容，使用TEXT类型存储大型JSON字符串
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,  -- 创建时间
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,  -- 更新时间
    
    -- 添加索引以提高查询性能
    CONSTRAINT file_json_store_file_id_unique UNIQUE (file_id)
);

-- 创建索引（只保留基本索引）
CREATE INDEX IF NOT EXISTS idx_dm_excel_json_store_created_at ON dm_excel_json_store(created_at);
CREATE INDEX IF NOT EXISTS idx_dm_excel_json_store_file_name ON dm_excel_json_store(file_name);

-- 创建更新时间触发器
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_dm_excel_json_store_updated_at
    BEFORE UPDATE ON dm_excel_json_store
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column(); 