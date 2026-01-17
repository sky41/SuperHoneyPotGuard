package middleware

import (
	"net/http"
	"superhoneypotguard/database"
	"superhoneypotguard/utils"

	"github.com/gin-gonic/gin"
)

func PermissionMiddleware(permissionCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := GetCurrentUser(c)

		var count int64
		database.DB.Raw(`
			SELECT COUNT(*)
			FROM users u
			INNER JOIN user_roles ur ON u.id = ur.user_id
			INNER JOIN roles r ON ur.role_id = r.id
			INNER JOIN role_permissions rp ON r.id = rp.role_id
			INNER JOIN permissions p ON rp.permission_id = p.id
			WHERE u.id = ? AND p.permission_code = ? AND p.status = 1 AND r.status = 1 AND u.status = 1
		`, user.UserID, permissionCode).Scan(&count)

		if count == 0 {
			utils.ErrorResponse(c, http.StatusForbidden, "权限不足")
			c.Abort()
			return
		}

		c.Next()
	}
}
