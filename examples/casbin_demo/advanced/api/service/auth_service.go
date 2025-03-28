package service

import (
	"fmt"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/go-language-learning/examples/casbin_demo/advanced/api/models"
	"gorm.io/gorm"
)

// AuthService 认证服务
type AuthService struct {
	db       *gorm.DB
	enforcer *casbin.Enforcer
}

// NewAuthService 创建认证服务
func NewAuthService(db *gorm.DB, enforcer *casbin.Enforcer) *AuthService {
	return &AuthService{
		db:       db,
		enforcer: enforcer,
	}
}

// AddUserToGroup 将用户添加到用户组
func (s *AuthService) AddUserToGroup(userID uint, groupID uint) error {
	// 首先检查用户和用户组是否存在
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return fmt.Errorf("用户不存在: %w", err)
	}

	var group models.UserGroup
	if err := s.db.First(&group, groupID).Error; err != nil {
		return fmt.Errorf("用户组不存在: %w", err)
	}

	// 添加用户到用户组
	member := models.UserGroupMember{
		UserID:  userID,
		GroupID: groupID,
	}
	if err := s.db.Create(&member).Error; err != nil {
		return fmt.Errorf("添加用户到用户组失败: %w", err)
	}

	// 同步到 Casbin 策略
	_, err := s.enforcer.AddGroupingPolicy(fmt.Sprintf("user:%d", userID), fmt.Sprintf("group:%d", groupID))
	if err != nil {
		return fmt.Errorf("同步 Casbin 策略失败: %w", err)
	}

	return nil
}

// AssignRoleToUser 给用户分配角色
func (s *AuthService) AssignRoleToUser(userID uint, roleID uint) error {
	// 检查用户和角色是否存在
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return fmt.Errorf("用户不存在: %w", err)
	}

	var role models.Role
	if err := s.db.First(&role, roleID).Error; err != nil {
		return fmt.Errorf("角色不存在: %w", err)
	}

	// 添加用户角色关系
	userRole := models.UserRole{
		UserID: userID,
		RoleID: roleID,
	}
	if err := s.db.Create(&userRole).Error; err != nil {
		return fmt.Errorf("分配角色失败: %w", err)
	}

	// 同步到 Casbin 策略
	_, err := s.enforcer.AddRoleForUser(fmt.Sprintf("user:%d", userID), role.Name)
	if err != nil {
		return fmt.Errorf("同步 Casbin 策略失败: %w", err)
	}

	return nil
}

// AddGroupPolicy 添加用户组权限策略
func (s *AuthService) AddGroupPolicy(groupID uint, domain, obj, act string) error {
	// 检查用户组是否存在
	var group models.UserGroup
	if err := s.db.First(&group, groupID).Error; err != nil {
		return fmt.Errorf("用户组不存在: %w", err)
	}

	// 添加策略
	_, err := s.enforcer.AddPolicy(fmt.Sprintf("group:%d", groupID), domain, obj, act, "allow")
	if err != nil {
		return fmt.Errorf("添加策略失败: %w", err)
	}

	return nil
}

// CheckPermission 检查权限
func (s *AuthService) CheckPermission(userID uint, domain, obj, act string) (bool, error) {
	return s.enforcer.Enforce(fmt.Sprintf("user:%d", userID), domain, obj, act)
}

// GetUserGroups 获取用户所属的所有用户组
func (s *AuthService) GetUserGroups(userID uint) ([]models.UserGroup, error) {
	var groups []models.UserGroup
	err := s.db.Model(&models.UserGroup{}).
		Joins("JOIN user_group_members ON user_groups.id = user_group_members.group_id").
		Where("user_group_members.user_id = ?", userID).
		Find(&groups).Error
	if err != nil {
		return nil, fmt.Errorf("获取用户组失败: %w", err)
	}
	return groups, nil
}

// GetUserRoles 获取用户的所有角色
func (s *AuthService) GetUserRoles(userID uint) ([]models.Role, error) {
	var roles []models.Role
	err := s.db.Model(&models.Role{}).
		Joins("JOIN user_roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ?", userID).
		Find(&roles).Error
	if err != nil {
		return nil, fmt.Errorf("获取角色失败: %w", err)
	}
	return roles, nil
}

// GetGroupPolicies 获取用户组的所有权限策略
func (s *AuthService) GetGroupPolicies(groupID uint) ([][]string, error) {
	policies := s.enforcer.GetFilteredPolicy(0, fmt.Sprintf("group:%d", groupID))
	return policies, nil
}

// RemoveUserFromGroup 从用户组中移除用户
func (s *AuthService) RemoveUserFromGroup(userID uint, groupID uint) error {
	// 从数据库中移除
	err := s.db.Where("user_id = ? AND group_id = ?", userID, groupID).
		Delete(&models.UserGroupMember{}).Error
	if err != nil {
		return fmt.Errorf("从用户组移除用户失败: %w", err)
	}

	// 从 Casbin 策略中移除
	_, err = s.enforcer.RemoveGroupingPolicy(
		fmt.Sprintf("user:%d", userID),
		fmt.Sprintf("group:%d", groupID),
	)
	if err != nil {
		return fmt.Errorf("同步 Casbin 策略失败: %w", err)
	}

	return nil
}

// RemoveRoleFromUser 移除用户的角色
func (s *AuthService) RemoveRoleFromUser(userID uint, roleID uint) error {
	// 从数据库中移除
	err := s.db.Where("user_id = ? AND role_id = ?", userID, roleID).
		Delete(&models.UserRole{}).Error
	if err != nil {
		return fmt.Errorf("移除用户角色失败: %w", err)
	}

	// 获取角色名称
	var role models.Role
	if err := s.db.First(&role, roleID).Error; err != nil {
		return fmt.Errorf("角色不存在: %w", err)
	}

	// 从 Casbin 策略中移除
	_, err = s.enforcer.DeleteRoleForUser(fmt.Sprintf("user:%d", userID), role.Name)
	if err != nil {
		return fmt.Errorf("同步 Casbin 策略失败: %w", err)
	}

	return nil
}

// UpdateUserRole 更新用户角色
func (s *AuthService) UpdateUserRole(userID uint, oldRoleID uint, newRoleID uint) error {
	// 开启事务
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 删除旧角色
		if err := s.RemoveRoleFromUser(userID, oldRoleID); err != nil {
			return err
		}
		// 添加新角色
		return s.AssignRoleToUser(userID, newRoleID)
	})
}

// UpdateUserGroup 更新用户组
func (s *AuthService) UpdateUserGroup(userID uint, oldGroupID uint, newGroupID uint) error {
	// 开启事务
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 删除旧用户组
		if err := s.RemoveUserFromGroup(userID, oldGroupID); err != nil {
			return err
		}
		// 添加新用户组
		return s.AddUserToGroup(userID, newGroupID)
	})
}

// UpdateGroupPolicy 更新用户组权限策略
func (s *AuthService) UpdateGroupPolicy(groupID uint, oldPolicy, newPolicy []string) error {
	// 检查用户组是否存在
	var group models.UserGroup
	if err := s.db.First(&group, groupID).Error; err != nil {
		return fmt.Errorf("用户组不存在: %w", err)
	}

	// 更新策略
	_, err := s.enforcer.UpdatePolicy(
		append([]string{fmt.Sprintf("group:%d", groupID)}, oldPolicy...),
		append([]string{fmt.Sprintf("group:%d", groupID)}, newPolicy...),
	)
	if err != nil {
		return fmt.Errorf("更新策略失败: %w", err)
	}

	return nil
}

// RemoveGroupPolicy 删除用户组权限策略
func (s *AuthService) RemoveGroupPolicy(groupID uint, domain, obj, act string) error {
	// 检查用户组是否存在
	var group models.UserGroup
	if err := s.db.First(&group, groupID).Error; err != nil {
		return fmt.Errorf("用户组不存在: %w", err)
	}

	// 删除策略
	_, err := s.enforcer.RemovePolicy(fmt.Sprintf("group:%d", groupID), domain, obj, act)
	if err != nil {
		return fmt.Errorf("删除策略失败: %w", err)
	}

	return nil
}

// GetAllPolicies 获取所有权限策略
func (s *AuthService) GetAllPolicies() [][]string {
	return s.enforcer.GetPolicy()
}

// HasPolicy 检查是否存在特定策略
func (s *AuthService) HasPolicy(sub, dom, obj, act string) bool {
	return s.enforcer.HasPolicy(sub, dom, obj, act)
}

// GetImplicitPermissionsForUser 获取用户的所有权限（包括继承的权限）
func (s *AuthService) GetImplicitPermissionsForUser(userID uint) ([][]string, error) {
	return s.enforcer.GetImplicitPermissionsForUser(fmt.Sprintf("user:%d", userID))
}

// CheckBasicPermission 基本权限检查（不包含继承权限）
func (s *AuthService) CheckBasicPermission(userID uint, domain, obj, act string) (bool, error) {
	// 直接检查用户是否有指定权限
	return s.enforcer.HasPermissionForUser(
		fmt.Sprintf("user:%d", userID),
		domain,
		obj,
		act,
	)
}

// CheckInheritedPermission 继承权限检查（包含角色和用户组继承的权限）
func (s *AuthService) CheckInheritedPermission(userID uint, domain, obj, act string) (bool, error) {
	// 使用 Enforce 方法会检查所有继承的权限
	return s.enforcer.Enforce(
		fmt.Sprintf("user:%d", userID),
		domain,
		obj,
		act,
	)
}

// ValidatePolicy 策略验证
func (s *AuthService) ValidatePolicy(policy []string) (bool, error) {
	if len(policy) < 3 {
		return false, fmt.Errorf("策略格式错误：至少需要主体、对象和操作")
	}

	// 检查策略格式
	switch len(policy) {
	case 3: // [obj, act, eft]
		return true, nil
	case 4: // [sub, obj, act, eft]
		return true, nil
	case 5: // [sub, dom, obj, act, eft]
		return true, nil
	default:
		return false, fmt.Errorf("策略格式错误：参数数量不正确")
	}
}

// GetPermissionsByRole 获取角色的所有权限
func (s *AuthService) GetPermissionsByRole(roleName string) [][]string {
	return s.enforcer.GetPermissionsForUser(roleName)
}

// GetPermissionsByGroup 获取用户组的所有权限
func (s *AuthService) GetPermissionsByGroup(groupID uint) [][]string {
	return s.enforcer.GetPermissionsForUser(fmt.Sprintf("group:%d", groupID))
}

// CheckUserInGroup 检查用户是否在指定用户组中（包括继承关系）
func (s *AuthService) CheckUserInGroup(userID uint, groupID uint) (bool, error) {
	return s.enforcer.HasGroupingPolicy(
		fmt.Sprintf("user:%d", userID),
		fmt.Sprintf("group:%d", groupID),
	)
}

// GetAllUserPermissions 获取用户的所有权限（包括直接权限、角色权限和用户组权限）
func (s *AuthService) GetAllUserPermissions(userID uint) (map[string][]string, error) {
	permissions := make(map[string][]string)

	// 1. 获取直接权限
	directPerms, err := s.CheckBasicPermission(userID, "*", "*", "*")
	if err != nil {
		return nil, fmt.Errorf("获取直接权限失败: %w", err)
	}
	if directPerms {
		permissions["direct"] = []string{"*", "*", "*"}
	}

	// 2. 获取角色权限
	roles, err := s.GetUserRoles(userID)
	if err != nil {
		return nil, fmt.Errorf("获取角色失败: %w", err)
	}
	for _, role := range roles {
		rolePerms := s.GetPermissionsByRole(role.Name)
		permissions["role:"+role.Name] = make([]string, 0)
		for _, perm := range rolePerms {
			permissions["role:"+role.Name] = append(permissions["role:"+role.Name], perm...)
		}
	}

	// 3. 获取用户组权限
	groups, err := s.GetUserGroups(userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户组失败: %w", err)
	}
	for _, group := range groups {
		groupPerms := s.GetPermissionsByGroup(group.ID)
		permissions["group:"+group.Name] = make([]string, 0)
		for _, perm := range groupPerms {
			permissions["group:"+group.Name] = append(permissions["group:"+group.Name], perm...)
		}
	}

	return permissions, nil
}

// CheckPolicyConflicts 检查权限策略冲突
func (s *AuthService) CheckPolicyConflicts(sub, dom, obj, act string) (bool, [][]string, error) {
	// 获取所有匹配的策略
	policies := s.enforcer.GetFilteredPolicy(0, sub)
	conflicts := make([][]string, 0)

	for _, policy := range policies {
		// 检查是否有冲突的策略（相同对象不同操作）
		if len(policy) >= 4 && policy[1] == dom && policy[2] == obj && policy[3] != act {
			conflicts = append(conflicts, policy)
		}
	}

	return len(conflicts) > 0, conflicts, nil
}

// CheckDataPermission 检查数据权限
func (s *AuthService) CheckDataPermission(userID uint, domain, obj, act string, projectID string) (bool, error) {
	return s.enforcer.Enforce(
		fmt.Sprintf("user:%d", userID),
		domain,
		obj,
		act,
		projectID,
	)
}

// AddDepartmentPolicy 添加部门数据权限策略
func (s *AuthService) AddDepartmentPolicy(deptID string, parentDeptID string) error {
	_, err := s.enforcer.AddNamedGroupingPolicy("g4", deptID, parentDeptID)
	return err
}

// AddProjectPolicy 添加项目数据权限策略
func (s *AuthService) AddProjectPolicy(sub string, domain string, obj string, act string, projectID string) error {
	_, err := s.enforcer.AddPolicy(sub, domain, obj, act, projectID, "allow")
	return err
}

// GetUserDepartmentProjects 获取用户部门可访问的所有项目
func (s *AuthService) GetUserDepartmentProjects(userID uint) ([]string, error) {
	return s.enforcer.GetImplicitPermissionsForUser(fmt.Sprintf("user:%d", userID))
}

// ShareDocumentWithUser 分享文档给指定用户
func (s *AuthService) ShareDocumentWithUser(ownerID uint, targetUserID uint, docID string, permissions []string) error {
	policies := make([][]string, len(permissions))
	for i, perm := range permissions {
		policies[i] = []string{
			fmt.Sprintf("user:%d", targetUserID),
			"platform",
			fmt.Sprintf("/api/documents/%s", docID),
			perm,
			"*",
			"allow",
		}
	}
	_, err := s.enforcer.AddPolicies(policies)
	return err
}

// GetSharedDocuments 获取分享给用户的文档
func (s *AuthService) GetSharedDocuments(userID uint) ([]string, error) {
	permissions, err := s.enforcer.GetImplicitPermissionsForUser(fmt.Sprintf("user:%d", userID))
	if err != nil {
		return nil, err
	}

	docs := make([]string, 0)
	for _, perm := range permissions {
		if strings.HasPrefix(perm[2], "/api/documents/") {
			docID := strings.TrimPrefix(perm[2], "/api/documents/")
			if !contains(docs, docID) {
				docs = append(docs, docID)
			}
		}
	}
	return docs, nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
