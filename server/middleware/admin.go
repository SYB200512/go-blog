package middleware

import (
	"server/model/appTypes"
	"server/model/response"
	"server/utils"

	"github.com/gin-gonic/gin"
)

// 管理员权限认证
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID := utils.GetRoleID(c)
		if roleID != appTypes.Admin {
			response.NoAuth("Access denied ,admin permission required", c)
			c.Abort()
			return
		}
	}
}
