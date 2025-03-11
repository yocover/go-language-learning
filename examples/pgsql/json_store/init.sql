-- 如果表已存在，先删除
DROP TABLE IF EXISTS dm_excel_json_store;

-- 创建扩展（如果需要）
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 创建文件数据表
CREATE TABLE dm_excel_json_store (
    id SERIAL PRIMARY KEY,                         -- 自增主键
    file_id VARCHAR(255) NOT NULL,                 -- 文件唯一标识符
    file_name VARCHAR(255) NOT NULL,               -- 文件名称
    json_content TEXT NOT NULL,                    -- JSON内容，使用TEXT类型存储大型JSON字符串
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,  -- 创建时间
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,  -- 更新时间
    
    -- 添加唯一约束
    CONSTRAINT dm_excel_json_store_file_id_unique UNIQUE (file_id)
);

-- 创建索引
CREATE INDEX idx_dm_excel_json_store_created_at ON dm_excel_json_store(created_at);
CREATE INDEX idx_dm_excel_json_store_file_name ON dm_excel_json_store(file_name);

-- 创建更新时间触发器函数
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 创建触发器
DROP TRIGGER IF EXISTS update_dm_excel_json_store_updated_at ON dm_excel_json_store;
CREATE TRIGGER update_dm_excel_json_store_updated_at
    BEFORE UPDATE ON dm_excel_json_store
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 添加注释
COMMENT ON TABLE dm_excel_json_store IS 'Excel JSON数据存储表';
COMMENT ON COLUMN dm_excel_json_store.id IS '自增主键';
COMMENT ON COLUMN dm_excel_json_store.file_id IS '文件唯一标识符';
COMMENT ON COLUMN dm_excel_json_store.file_name IS '文件名称';
COMMENT ON COLUMN dm_excel_json_store.json_content IS 'JSON内容（TEXT格式）';
COMMENT ON COLUMN dm_excel_json_store.created_at IS '创建时间';
COMMENT ON COLUMN dm_excel_json_store.updated_at IS '更新时间';

-- 验证表是否创建成功
SELECT EXISTS (
    SELECT FROM information_schema.tables 
    WHERE table_schema = 'public' 
    AND table_name = 'dm_excel_json_store'
); 