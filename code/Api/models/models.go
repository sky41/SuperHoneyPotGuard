package models

import (
	"time"
)

type User struct {
	ID            int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Username      string    `json:"username" gorm:"uniqueIndex;not null;size:50"`
	Password      string    `json:"-" gorm:"not null;size:255"`
	Email         *string   `json:"email" gorm:"uniqueIndex;size:100"`
	Phone         *string   `json:"phone" gorm:"size:20"`
	RealName      *string   `json:"realName" gorm:"column:real_name;size:50"`
	Status        int       `json:"status" gorm:"default:1;comment:0-禁用,1-启用"`
	LastLoginTime *time.Time `json:"lastLoginTime" gorm:"column:last_login_time"`
	LastLoginIP   *string   `json:"lastLoginIp" gorm:"column:last_login_ip;size:50"`
	CreatedAt     time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
	CreatedBy     *int      `json:"created_by"`
	UpdatedBy     *int      `json:"updated_by"`
	Roles         []Role    `json:"roles" gorm:"many2many:user_roles;"`
}

type Role struct {
	ID          int         `json:"id" gorm:"primaryKey;autoIncrement"`
	RoleName    string      `json:"roleName" gorm:"column:role_name;uniqueIndex;not null;size:50"`
	RoleCode    string      `json:"roleCode" gorm:"column:role_code;uniqueIndex;not null;size:50"`
	Description *string     `json:"description" gorm:"size:200"`
	Status      int         `json:"status" gorm:"default:1;comment:0-禁用,1-启用"`
	CreatedAt   time.Time   `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time   `json:"updatedAt" gorm:"autoUpdateTime"`
	CreatedBy   *int        `json:"created_by"`
	UpdatedBy   *int        `json:"updated_by"`
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;"`
	Users       []User      `json:"-" gorm:"many2many:user_roles;"`
}

type Permission struct {
	ID             int          `json:"id" gorm:"primaryKey;autoIncrement"`
	PermissionName string       `json:"permissionName" gorm:"column:permission_name;not null;size:50"`
	PermissionCode string       `json:"permissionCode" gorm:"column:permission_code;uniqueIndex;not null;size:100"`
	PermissionType string       `json:"permissionType" gorm:"column:permission_type;not null;size:20;comment:menu-菜单,button-按钮,api-接口"`
	ParentID       int          `json:"parentId" gorm:"column:parent_id;default:0"`
	Path           *string      `json:"path" gorm:"size:200"`
	Component      *string      `json:"component" gorm:"size:200"`
	Icon           *string      `json:"icon" gorm:"size:50"`
	SortOrder      int          `json:"sortOrder" gorm:"column:sort_order;default:0"`
	Description    *string      `json:"description" gorm:"size:200"`
	Status         int          `json:"status" gorm:"default:1;comment:0-禁用,1-启用"`
	CreatedAt      time.Time    `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt      time.Time    `json:"updatedAt" gorm:"autoUpdateTime"`
	Children       []Permission `json:"children" gorm:"-"`
	Roles          []Role       `json:"-" gorm:"many2many:role_permissions;"`
}

type UserRole struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    int       `json:"userId" gorm:"column:user_id;not null;index"`
	RoleID    int       `json:"roleId" gorm:"column:role_id;not null;index"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	CreatedBy *int      `json:"created_by"`
}

type RolePermission struct {
	ID           int       `json:"id" gorm:"primaryKey;autoIncrement"`
	RoleID       int       `json:"roleId" gorm:"column:role_id;not null;index"`
	PermissionID int       `json:"permissionId" gorm:"column:permission_id;not null;index"`
	CreatedAt    time.Time `json:"createdAt" gorm:"autoCreateTime"`
	CreatedBy    *int      `json:"created_by"`
}

// 提交目的：优化operation_logs表索引结构，提升INSERT性能
// 提交内容：移除不必要的索引，保留核心查询索引，添加复合索引
// 提交时间：2026-01-19

type OperationLog struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID      *int      `json:"userId" gorm:"column:user_id;index:idx_user_created"`
	Username    *string   `json:"username" gorm:"size:50"`
	Operation   string    `json:"operation" gorm:"not null;size:100"`
	Method      *string   `json:"method" gorm:"size:10"`
	URL         *string   `json:"url" gorm:"size:500"`
	IP          *string   `json:"ip" gorm:"size:50"`
	Location    *string   `json:"location" gorm:"size:100"`
	Params      *string   `json:"params" gorm:"type:text"`
	Result      *string   `json:"result" gorm:"type:text"`
	Status      int       `json:"status" gorm:"default:1;comment:0-失败,1-成功"`
	ErrorMsg    *string   `json:"errorMsg" gorm:"column:error_msg;size:500"`
	ExecuteTime int       `json:"executeTime" gorm:"column:execute_time;comment:执行时间(ms)"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime;index:idx_user_created"`
}

type RegisterRequest struct {
	Username string  `json:"username" binding:"required,min=3,max=50"`
	Password string  `json:"password" binding:"required,min=6"`
	Email    *string `json:"email" binding:"omitempty,email"`
	Phone    *string `json:"phone"`
	RealName *string `json:"realName"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateUserRequest struct {
	Username string  `json:"username" binding:"required,min=3,max=50"`
	Password string  `json:"password" binding:"required,min=6"`
	Email    *string `json:"email" binding:"omitempty,email"`
	Phone    *string `json:"phone"`
	RealName *string `json:"realName"`
	RoleIDs  []int   `json:"roleIds"`
	Status   *int    `json:"status"`
}

type UpdateUserRequest struct {
	Email    *string `json:"email" binding:"omitempty,email"`
	Phone    *string `json:"phone"`
	RealName *string `json:"realName"`
	RoleIDs  *[]int  `json:"roleIds"`
	Status   *int    `json:"status"`
}

type UpdateUserStatusRequest struct {
	Status int `json:"status" binding:"required,oneof=0 1"`
}

type ResetPasswordRequest struct {
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}

type CreateRoleRequest struct {
	RoleName      string  `json:"roleName" binding:"required"`
	RoleCode      string  `json:"roleCode" binding:"required"`
	Description   *string `json:"description"`
	PermissionIDs []int   `json:"permissionIds"`
	Status        *int    `json:"status"`
}

type UpdateRoleRequest struct {
	RoleName      *string `json:"roleName" binding:"required"`
	Description   *string `json:"description"`
	PermissionIDs *[]int  `json:"permissionIds"`
	Status        *int    `json:"status"`
}

type CreatePermissionRequest struct {
	PermissionName string  `json:"permissionName" binding:"required"`
	PermissionCode string  `json:"permissionCode" binding:"required"`
	PermissionType string `json:"permissionType" binding:"required,oneof=menu button api"`
	ParentID       *int    `json:"parentId"`
	Path           *string `json:"path"`
	Component      *string `json:"component"`
	Icon           *string `json:"icon"`
	SortOrder      *int    `json:"sortOrder"`
	Description    *string `json:"description"`
	Status         *int    `json:"status"`
}

type UpdatePermissionRequest struct {
	PermissionName *string `json:"permissionName" binding:"required"`
	PermissionType *string `json:"permissionType" binding:"required,oneof=menu button api"`
	ParentID       *int    `json:"parentId"`
	Path           *string `json:"path"`
	Component      *string `json:"component"`
	Icon           *string `json:"icon"`
	SortOrder      *int    `json:"sortOrder"`
	Description    *string `json:"description"`
	Status         *int    `json:"status"`
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PaginatedResponse struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

type Claims struct {
	UserID      int      `json:"userId"`
	Username    string   `json:"username"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}
