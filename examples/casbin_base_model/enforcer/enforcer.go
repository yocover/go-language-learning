package enforcer

import (
	"fmt"
	"sync"
	"time"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Config 定义权限配置
type Config struct {
	// 数据库配置
	DBType       string
	DBConnection string

	// 模型配置文件路径
	ModelPath string

	// 是否自动加载策略（用于动态更新策略）
	AutoLoad bool
	// 自动加载间隔（秒）
	AutoLoadInterval int
}

// Enforcer 封装 casbin enforcer
type Enforcer struct {
	enforcer *casbin.Enforcer
	adapter  *gormadapter.Adapter
	config   *Config
	mu       sync.RWMutex
}

// NewEnforcer 创建一个新的 enforcer 实例
func NewEnforcer(config *Config) (*Enforcer, error) {
	// 连接数据库
	db, err := gorm.Open(mysql.Open(config.DBConnection), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect database failed: %v", err)
	}

	// 创建 adapter
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, fmt.Errorf("create adapter failed: %v", err)
	}

	// 创建 enforcer
	e, err := casbin.NewEnforcer(config.ModelPath, adapter)
	if err != nil {
		return nil, fmt.Errorf("create enforcer failed: %v", err)
	}

	// 加载策略
	if err := e.LoadPolicy(); err != nil {
		return nil, fmt.Errorf("load policy failed: %v", err)
	}

	enforcer := &Enforcer{
		enforcer: e,
		adapter:  adapter,
		config:   config,
	}

	// 如果配置了自动加载，启动自动加载协程
	if config.AutoLoad {
		go enforcer.autoLoad()
	}

	return enforcer, nil
}

// 检查权限
func (e *Enforcer) Enforce(rvals ...interface{}) (bool, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.enforcer.Enforce(rvals...)
}

// AddPolicy 添加策略
func (e *Enforcer) AddPolicy(params ...interface{}) (bool, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.enforcer.AddPolicy(params...)
}

// RemovePolicy 删除策略
func (e *Enforcer) RemovePolicy(params ...interface{}) (bool, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.enforcer.RemovePolicy(params...)
}

// AddGroupingPolicy 添加角色继承关系
func (e *Enforcer) AddGroupingPolicy(params ...interface{}) (bool, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.enforcer.AddGroupingPolicy(params...)
}

// RemoveGroupingPolicy 删除角色继承关系
func (e *Enforcer) RemoveGroupingPolicy(params ...interface{}) (bool, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.enforcer.RemoveGroupingPolicy(params...)
}

// GetAllSubjects 获取所有主体
func (e *Enforcer) GetAllSubjects() []string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.enforcer.GetAllSubjects()
}

// GetAllObjects 获取所有对象
func (e *Enforcer) GetAllObjects() []string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.enforcer.GetAllObjects()
}

// GetAllActions 获取所有操作
func (e *Enforcer) GetAllActions() []string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.enforcer.GetAllActions()
}

// GetAllRoles 获取所有角色
func (e *Enforcer) GetAllRoles() []string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.enforcer.GetAllRoles()
}

// LoadPolicy 重新加载策略
func (e *Enforcer) LoadPolicy() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.enforcer.LoadPolicy()
}

// SavePolicy 保存策略到存储
func (e *Enforcer) SavePolicy() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.enforcer.SavePolicy()
}

// autoLoad 自动加载策略
func (e *Enforcer) autoLoad() {
	if e.config.AutoLoadInterval <= 0 {
		e.config.AutoLoadInterval = 10 // 默认10秒
	}

	ticker := time.NewTicker(time.Duration(e.config.AutoLoadInterval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if err := e.LoadPolicy(); err != nil {
			fmt.Printf("Auto load policy failed: %v\n", err)
		}
	}
}
