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
    rpc Enforce(EnforceRequest) returns (EnforceResponse) {}
    rpc BatchEnforce(BatchEnforceRequest) returns (BatchEnforceResponse) {}
    
    // 策略管理
    rpc AddPolicy(PolicyRequest) returns (BoolResponse) {}
    rpc RemovePolicy(PolicyRequest) returns (BoolResponse) {}
    rpc GetPolicy(GetPolicyRequest) returns (GetPolicyResponse) {}
    rpc HasPolicy(PolicyRequest) returns (BoolResponse) {}
    rpc UpdatePolicy(UpdatePolicyRequest) returns (BoolResponse) {}
    
    // 角色管理
    rpc AddRoleForUser(RoleRequest) returns (BoolResponse) {}
    rpc DeleteRoleForUser(RoleRequest) returns (BoolResponse) {}
    rpc GetRolesForUser(GetRoleRequest) returns (GetRolesResponse) {}
    rpc HasRoleForUser(RoleRequest) returns (BoolResponse) {}
    rpc GetUsersForRole(GetUsersRequest) returns (GetUsersResponse) {}
    
    // 权限查询
    rpc GetImplicitPermissionsForUser(GetPermissionsRequest) returns (GetPermissionsResponse) {}
    rpc GetPermissionsForUser(GetPermissionsRequest) returns (GetPermissionsResponse) {}
    rpc HasPermissionForUser(PermissionRequest) returns (BoolResponse) {}
    
    // 域管理（多租户）
    rpc GetAllDomains(Empty) returns (GetDomainsResponse) {}
    rpc AddPolicyInDomain(DomainPolicyRequest) returns (BoolResponse) {}
    rpc RemovePolicyInDomain(DomainPolicyRequest) returns (BoolResponse) {}
    
    // 模型管理
    rpc GetCurrentModel(Empty) returns (ModelResponse) {}
    rpc ReloadPolicy(Empty) returns (BoolResponse) {}
}

// 请求和响应消息定义
message EnforceRequest {
    string subject = 1;
    string object = 2;
    string action = 3;
    string domain = 4;  // 可选，用于多租户
    map<string, string> context = 5;  // 可选，用于 ABAC
}

message PolicyRequest {
    string subject = 1;
    string object = 2;
    string action = 3;
    string domain = 4;
    repeated string values = 5;  // 扩展字段
}

message UpdatePolicyRequest {
    PolicyRequest old_rule = 1;
    PolicyRequest new_rule = 2;
}

message GetPermissionsRequest {
    string user = 1;
    string domain = 2;
}

message GetPermissionsResponse {
    repeated PolicyRule permissions = 1;
}

message PolicyRule {
    string subject = 1;
    string object = 2;
    string action = 3;
    string domain = 4;
    repeated string values = 5;
}

// ... 其他消息定义 ...
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

### 2.2 多模型支持

在实际项目中，通常不需要动态切换权限模型，而是采用以下方案：

1. **固定模型方案**
   - 在项目初期确定权限模型（如 RBAC）
   - 通过扩展现有模型来满足新需求
   - 保持模型稳定，避免运行时切换

2. **多模型并存方案**
   ```go
   // PermissionService 权限服务
   type PermissionService struct {
       rbacEnforcer     *casbin.Enforcer  // 基础 RBAC 模型
       deptEnforcer     *casbin.Enforcer  // 部门权限模型
       resourceEnforcer *casbin.Enforcer  // 资源权限模型
   }

   // NewPermissionService 创建权限服务
   func NewPermissionService(adapter persist.Adapter) (*PermissionService, error) {
       // 创建基础 RBAC Enforcer
       rbacEnforcer, err := casbin.NewEnforcer("rbac_model.conf", adapter)
       if err != nil {
           return nil, err
       }

       // 创建部门权限 Enforcer
       deptEnforcer, err := casbin.NewEnforcer("dept_model.conf", adapter)
       if err != nil {
           return nil, err
       }

       // 创建资源权限 Enforcer
       resourceEnforcer, err := casbin.NewEnforcer("resource_model.conf", adapter)
       if err != nil {
           return nil, err
       }

       return &PermissionService{
           rbacEnforcer:     rbacEnforcer,
           deptEnforcer:     deptEnforcer,
           resourceEnforcer: resourceEnforcer,
       }, nil
   }

   // CheckPermission 检查权限
   func (s *PermissionService) CheckPermission(ctx context.Context, req *PermissionRequest) (bool, error) {
       switch req.PermissionType {
       case "rbac":
           return s.rbacEnforcer.Enforce(req.Subject, req.Object, req.Action)
       case "dept":
           return s.deptEnforcer.Enforce(req.Subject, req.Object, req.Action, req.Department)
       case "resource":
           return s.resourceEnforcer.Enforce(req.Subject, req.Object, req.Action, req.ResourceType)
       default:
           return false, fmt.Errorf("unsupported permission type: %s", req.PermissionType)
       }
   }
   ```

3. **使用示例**
   ```go
   // 初始化权限服务
   permissionService, err := NewPermissionService(adapter)
   if err != nil {
       log.Fatalf("Failed to create permission service: %v", err)
   }

   // 检查 RBAC 权限
   allowed, err := permissionService.CheckPermission(ctx, &PermissionRequest{
       PermissionType: "rbac",
       Subject: "alice",
       Object: "data1",
       Action: "read",
   })

   // 检查部门权限
   allowed, err = permissionService.CheckPermission(ctx, &PermissionRequest{
       PermissionType: "dept",
       Subject: "alice",
       Object: "data1",
       Action: "read",
       Department: "dept1",
   })
   ```

4. **方案优势**
   - **稳定性**：避免运行时切换模型带来的风险
   - **清晰性**：不同类型的权限检查逻辑分离
   - **可维护性**：每个模型独立维护，互不影响
   - **性能优化**：可以针对不同模型进行独立的缓存优化
   - **平滑升级**：可以逐步迁移或替换某个模型

### 2.4 缓存设计

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