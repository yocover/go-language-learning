# Casbin 权限服务设计文档

## 1. 概述

### 1.1 目的
本文档描述了基于 Casbin 的权限服务设计，提供统一的权限管理接口，支持 gRPC 和 HTTP 调用。

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
    
    // 角色管理
    rpc AddRoleForUser(RoleRequest) returns (BoolResponse) {}
    rpc DeleteRoleForUser(RoleRequest) returns (BoolResponse) {}
    rpc GetRolesForUser(GetRoleRequest) returns (GetRoleResponse) {}
    
    // 新增模型管理接口
    rpc GetCurrentModel(Empty) returns (ModelResponse) {}
    rpc UpdateModel(UpdateModelRequest) returns (BoolResponse) {}
    rpc ReloadModel(Empty) returns (BoolResponse) {}
    
    // 新增模型管理接口
    rpc GetPredefinedModelTypes(Empty) returns (ModelTypesResponse) {}
}

message EnforceRequest {
    string subject = 1;
    string object = 2;
    string action = 3;
    string domain = 4;  // 可选，用于多租户
}

message EnforceResponse {
    bool allowed = 1;
}

message BatchEnforceRequest {
    repeated EnforceRequest requests = 1;
}

message BatchEnforceResponse {
    repeated bool results = 1;
}

message PolicyRequest {
    string subject = 1;
    string object = 2;
    string action = 3;
    string domain = 4;
}

message GetPolicyRequest {
    string domain = 1;  // 可选，用于过滤特定域的策略
}

message GetPolicyResponse {
    repeated Policy policies = 1;
}

message Policy {
    string subject = 1;
    string object = 2;
    string action = 3;
    string domain = 4;
}

message RoleRequest {
    string user = 1;
    string role = 2;
    string domain = 3;
}

message GetRoleRequest {
    string user = 1;
    string domain = 2;
}

message GetRoleResponse {
    repeated string roles = 1;
}

message BoolResponse {
    bool result = 1;
}

message Empty {}

message ModelResponse {
    string model_text = 1;    // 当前使用的模型内容
    string model_path = 2;    // 当前使用的模型文件路径
}

message UpdateModelRequest {
    oneof model {
        string model_text = 1;  // 直接更新模型内容
        string model_path = 2;  // 使用新的模型文件
    }
}

message ModelTypesResponse {
    repeated string types = 1;
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

// 预定义模型内容
var PredefinedModels = map[string]string{
    ModelACL: `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && r.obj == p.obj && r.act == p.act
`,
    ModelRBAC: `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`,
    ModelRBACDomain: `
[request_definition]
r = sub, dom, obj, act

[policy_definition]
p = sub, dom, obj, act

[role_definition]
g = _, _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub, r.dom) && r.dom == p.dom && r.obj == p.obj && r.act == p.act
`,
    ModelRBACResource: `
[request_definition]
r = sub, dom, obj, act

[policy_definition]
p = sub, dom, obj, act

[role_definition]
g = _, _, _
g2 = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub, r.dom) && r.dom == p.dom && g2(r.obj, p.obj) && r.act == p.act
`,
    ModelABAC: `
[request_definition]
r = sub, obj, act, ctx

[policy_definition]
p = sub_rule, obj_rule, act, eft

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = eval(p.sub_rule) && eval(p.obj_rule) && r.act == p.act
`,
}

// ModelManager 模型管理器
type ModelManager struct {
    // 获取预定义模型
    GetPredefinedModel(modelType string) (string, error)
    // 验证模型内容
    ValidateModel(modelText string) error
    // 获取所有预定义模型类型
    GetPredefinedModelTypes() []string
}

// 实现 ModelManager
type modelManager struct{}

func NewModelManager() *modelManager {
    return &modelManager{}
}

func (m *modelManager) GetPredefinedModel(modelType string) (string, error) {
    if model, ok := PredefinedModels[modelType]; ok {
        return model, nil
    }
    return "", fmt.Errorf("model type %s not found", modelType)
}

func (m *modelManager) ValidateModel(modelText string) error {
    _, err := casbin.NewModel(modelText)
    return err
}

func (m *modelManager) GetPredefinedModelTypes() []string {
    types := make([]string, 0, len(PredefinedModels))
    for t := range PredefinedModels {
        types = append(types, t)
    }
    return types
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
func NewCasbinService(config *Config) (*CasbinService, error) {
    // 初始化模型管理器
    modelMgr := NewModelManager()
    
    // 初始化适配器
    a, err := gormadapter.NewAdapter(config.Casbin.Adapter.Type, 
        config.Casbin.Adapter.Connection,
        gormadapter.WithTableName(config.Casbin.Adapter.TableName))
    if err != nil {
        return nil, err
    }
    
    // 初始化 Enforcer
    var e *casbin.Enforcer
    
    // 根据配置选择模型
    switch {
    case config.Casbin.ModelType != "":
        // 使用预定义模型
        modelText, err := modelMgr.GetPredefinedModel(config.Casbin.ModelType)
        if err != nil {
            return nil, err
        }
        e, err = casbin.NewEnforcer(casbin.NewModel(modelText), a)
    case config.Casbin.ModelText != "":
        // 使用直接配置的模型内容
        if err := modelMgr.ValidateModel(config.Casbin.ModelText); err != nil {
            return nil, err
        }
        e, err = casbin.NewEnforcer(casbin.NewModel(config.Casbin.ModelText), a)
    case config.Casbin.ModelPath != "":
        // 使用模型文件
        e, err = casbin.NewEnforcer(config.Casbin.ModelPath, a)
    default:
        // 默认使用 RBAC with Domain 模型
        modelText, _ := modelMgr.GetPredefinedModel(ModelRBACDomain)
        e, err = casbin.NewEnforcer(casbin.NewModel(modelText), a)
    }
    
    if err != nil {
        return nil, err
    }
    
    // 启用自动加载和保存
    e.EnableAutoSave(true)
    
    // 初始化缓存
    c := cache.New(config.CacheTTL, config.CacheCleanupInterval)
    
    return &CasbinService{
        enforcer: e,
        cache:    c,
        config:   config,
        modelMgr: modelMgr,
    }, nil
}

// GetPredefinedModelTypes 获取所有预定义模型类型
func (s *CasbinService) GetPredefinedModelTypes(ctx context.Context, _ *pb.Empty) (*pb.ModelTypesResponse, error) {
    types := s.modelMgr.GetPredefinedModelTypes()
    return &pb.ModelTypesResponse{Types: types}, nil
}

// Enforce 检查权限
func (s *CasbinService) Enforce(ctx context.Context, req *pb.EnforceRequest) (*pb.EnforceResponse, error) {
    // 构造缓存键
    cacheKey := fmt.Sprintf("%s:%s:%s:%s", req.Subject, req.Object, req.Action, req.Domain)
    
    // 检查缓存
    if result, found := s.cache.Get(cacheKey); found {
        return &pb.EnforceResponse{Allowed: result.(bool)}, nil
    }
    
    // 执行权限检查
    allowed, err := s.enforcer.Enforce(req.Subject, req.Object, req.Action, req.Domain)
    if err != nil {
        return nil, err
    }
    
    // 设置缓存
    s.cache.Set(cacheKey, allowed, cache.DefaultExpiration)
    
    return &pb.EnforceResponse{Allowed: allowed}, nil
}

// BatchEnforce 批量检查权限
func (s *CasbinService) BatchEnforce(ctx context.Context, req *pb.BatchEnforceRequest) (*pb.BatchEnforceResponse, error) {
    results := make([]bool, len(req.Requests))
    
    for i, r := range req.Requests {
        resp, err := s.Enforce(ctx, r)
        if err != nil {
            return nil, err
        }
        results[i] = resp.Allowed
    }
    
    return &pb.BatchEnforceResponse{Results: results}, nil
}

// GetCurrentModel 获取当前使用的模型
func (s *CasbinService) GetCurrentModel(ctx context.Context, _ *pb.Empty) (*pb.ModelResponse, error) {
    return &pb.ModelResponse{
        ModelText: s.enforcer.GetModel().ToText(),
        ModelPath: s.config.Casbin.ModelPath,
    }, nil
}

// UpdateModel 更新模型
func (s *CasbinService) UpdateModel(ctx context.Context, req *pb.UpdateModelRequest) (*pb.BoolResponse, error) {
    var newModel model.Model
    var err error
    
    switch {
    case req.GetModelText() != "":
        newModel = casbin.NewModel(req.GetModelText())
    case req.GetModelPath() != "":
        newModel, err = model.NewModelFromFile(req.GetModelPath())
        if err != nil {
            return &pb.BoolResponse{Result: false}, err
        }
    default:
        return &pb.BoolResponse{Result: false}, errors.New("either model_text or model_path must be provided")
    }
    
    // 更新模型
    err = s.enforcer.SetModel(newModel)
    if err != nil {
        return &pb.BoolResponse{Result: false}, err
    }
    
    // 清除缓存
    s.cache.Flush()
    
    return &pb.BoolResponse{Result: true}, nil
}

// ReloadModel 重新加载模型
func (s *CasbinService) ReloadModel(ctx context.Context, _ *pb.Empty) (*pb.BoolResponse, error) {
    err := s.enforcer.LoadModel()
    if err != nil {
        return &pb.BoolResponse{Result: false}, err
    }
    
    // 清除缓存
    s.cache.Flush()
    
    return &pb.BoolResponse{Result: true}, nil
}
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

## 4. 使用示例

### 4.1 gRPC 调用示例
```go
func main() {
    // 创建 gRPC 连接
    conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()
    
    // 创建客户端
    client := pb.NewCasbinServiceClient(conn)
    
    // 检查权限
    resp, err := client.Enforce(context.Background(), &pb.EnforceRequest{
        Subject: "alice",
        Object:  "data1",
        Action:  "read",
        Domain:  "domain1",
    })
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Allowed: %v\n", resp.Allowed)
}
```

### 4.2 HTTP 调用示例
```go
func main() {
    // 检查权限
    resp, err := http.Post(
        "http://localhost:8080/api/v1/enforce",
        "application/json",
        strings.NewReader(`{
            "subject": "alice",
            "object": "data1",
            "action": "read",
            "domain": "domain1"
        }`),
    )
    if err != nil {
        log.Fatal(err)
    }
    
    var result struct {
        Allowed bool `json:"allowed"`
    }
    json.NewDecoder(resp.Body).Decode(&result)
    fmt.Printf("Allowed: %v\n", result.Allowed)
}
```

## 5. 部署指南

### 5.1 配置文件
```yaml
server:
  http_port: 8080
  grpc_port: 9090
  mode: release

casbin:
  # 使用预定义模型（优先级最高）
  model_type: "rbac_domain"  # 可选值: acl, rbac, rbac_domain, rbac_resource, abac
  
  # 或者使用模型文件（当 model_type 未指定时使用）
  model_path: "config/rbac_model.conf"
  
  # 或者直接配置模型内容（当 model_type 和 model_path 都未指定时使用）
  model_text: |
    [request_definition]
    r = sub, obj, act
    
    [policy_definition]
    p = sub, obj, act
    
    [role_definition]
    g = _, _
    
    [policy_effect]
    e = some(where (p.eft == allow))
    
    [matchers]
    m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
  
  adapter:
    type: mysql
    connection: "root:password@tcp(localhost:3306)/casbin"
    table_name: "my_casbin_rules"  # 自定义表名

cache:
  ttl: 5m
  cleanup_interval: 10m
```