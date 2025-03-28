# 权限服务设计

## 1. 概述

### 1.1 目的
基于 Casbin 的权限服务设计，提供统一的权限管理接口，支持 gRPC 和 HTTP 调用。

### 1.2 架构图
```
[业务系统] ---> [RPC/HTTP API] ---> [Casbin Engine] ---> [Storage]
                      |                    |
                      └-----> [Cache] <----┘
```

## 2. 接口设计

### 2.1 gRPC 接口定义
```protobuf
syntax = "proto3";

package casbin.v1;

service CasbinService {
    // 权限检查
    rpc Enforce(EnforceRequest) returns () {}
    rpc BatchEnforce(BatchEnforceRequest) returns () {}
    
    // 策略管理
    rpc AddPolicy(PolicyRequest) returns () {}
    rpc RemovePolicy(PolicyRequest) returns () {}
    rpc GetPolicy(GetPolicyRequest) returns () {}
    
    // 角色管理
    rpc AddRoleForUser(RoleRequest) returns () {}
    rpc DeleteRoleForUser(RoleRequest) returns () {}
    rpc GetRolesForUser(GetRoleRequest) returns () {}
    
    // 模型管理接口
    rpc GetCurrentModel(Empty) returns () {}
    rpc UpdateModel(UpdateModelRequest) returns () {}
    rpc ReloadModel(Empty) returns () {}
    rpc GetPredefinedModelTypes(Empty) returns () {}
}

message EnforceRequest {
    string subject = 1;
    string object = 2;
    string action = 3;
    string domain = 4;  // 可选，用于多租户
}
```

### 2.2 HTTP 接口定义
```
# 权限检查
POST /api/v1/enforce
Request:
{
    "subject": "alice",
    "object": "data1",
    "action": "read",
    "domain": "domain1"
}
Response:
{
    "allowed": true
}

# 批量权限检查
POST /api/v1/enforce/batch
Request:
{
    "requests": [
        {
            "subject": "alice",
            "object": "data1",
            "action": "read"
        }
    ]
}

# 策略管理
GET  /api/v1/policy?domain=domain1
POST /api/v1/policy
DELETE /api/v1/policy

# 角色管理
POST /api/v1/role
DELETE /api/v1/role
GET  /api/v1/role/:user

# 模型管理
GET /api/v1/model              # 获取当前使用的模型
PUT /api/v1/model              # 更新模型
POST /api/v1/model/reload      # 重新加载模型
```

## 3. 核心实现

### 3.1 预定义模型
```go
// 预定义模型类型
const (
    ModelACL           = "acl"           // 访问控制列表模型
    ModelRBAC          = "rbac"          // 基于角色的访问控制
    ModelRBACDomain    = "rbac_domain"   // 带域的RBAC
    ModelRBACResource  = "rbac_resource" // 基于资源的RBAC
    ModelABAC          = "abac"          // 基于属性的访问控制
)

// ModelManager 模型管理器
type ModelManager struct {
    // 获取预定义模型
    GetPredefinedModel(modelType string) (string, error)
    // 验证模型内容
    ValidateModel(modelText string) error
    // 获取所有预定义模型类型
    GetPredefinedModelTypes() []string
}
```

### 3.2 服务实现
```go
// CasbinService 权限服务
type CasbinService struct {
    enforcer *casbin.Enforcer
    cache    *cache.Cache
    config   *Config
    modelMgr *ModelManager
}

// NewCasbinService 创建权限服务
func NewCasbinService(config *Config) (*CasbinService, error) {}

// GetPredefinedModelTypes 获取所有预定义模型类型
func (s *CasbinService) GetPredefinedModelTypes(ctx context.Context, _ *pb.Empty) (*pb.ModelTypesResponse, error) {}

// Enforce 检查权限
func (s *CasbinService) Enforce(ctx context.Context, req *pb.EnforceRequest) (*pb.EnforceResponse, error) {}

// BatchEnforce 批量检查权限
func (s *CasbinService) BatchEnforce(ctx context.Context, req *pb.BatchEnforceRequest) (*pb.BatchEnforceResponse, error) {}

// GetCurrentModel 获取当前使用的模型
func (s *CasbinService) GetCurrentModel(ctx context.Context, _ *pb.Empty) (*pb.ModelResponse, error) {}

// UpdateModel 更新模型
func (s *CasbinService) UpdateModel(ctx context.Context, req *pb.UpdateModelRequest) (*pb.BoolResponse, error) {}

// ReloadModel 重新加载模型
func (s *CasbinService) ReloadModel(ctx context.Context, _ *pb.Empty) (*pb.BoolResponse, error) {}
```

### 3.3 配置定义
```go
// Config 服务配置
type Config struct {
    // 服务配置
    Server struct {
        HTTPPort int    `yaml:"http_port"`
        GRPCPort int    `yaml:"grpc_port"`
        Mode     string `yaml:"mode"` // debug or release
    } `yaml:"server"`
    
    // Casbin配置
    Casbin struct {
        // 模型配置
        ModelType string `yaml:"model_type"`  // 预定义模型类型，优先级最高
        ModelPath string `yaml:"model_path"`  // 自定义模型文件路径
        ModelText string `yaml:"model_text"`  // 直接配置模型内容
        
        // 适配器配置
        Adapter struct {
            Type       string `yaml:"type"` // mysql, postgres etc.
            Connection string `yaml:"connection"`
            TableName  string `yaml:"table_name"` // 自定义表名，默认为 casbin_rule
        } `yaml:"adapter"`
    } `yaml:"casbin"`
    
    // 缓存配置
    Cache struct {
        TTL             time.Duration `yaml:"ttl"`
        CleanupInterval time.Duration `yaml:"cleanup_interval"`
    } `yaml:"cache"`
}
```

### 3.4 数据库表结构
```sql
-- Casbin规则表
CREATE TABLE casbin_rule (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    ptype VARCHAR(10) NOT NULL,
    v0 VARCHAR(256), -- subject
    v1 VARCHAR(256), -- object
    v2 VARCHAR(256), -- action
    v3 VARCHAR(256), -- domain
    v4 VARCHAR(256),
    v5 VARCHAR(256),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_ptype (ptype),
    INDEX idx_v0_v1_v2_v3 (v0, v1, v2, v3)
);
```