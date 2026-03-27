// Package handler 提供HTTP请求处理器，负责接收HTTP请求、调用服务层、返回响应。
// 本包包含用户相关的HTTP处理器。
package handler

import (
	"final-whistle/backend/internal/middleware"
	"final-whistle/backend/internal/service"
	"final-whistle/backend/internal/utils"
	"github.com/gin-gonic/gin"
)

// UserHandler 处理用户相关的HTTP请求。
type UserHandler struct {
	service service.UserService
}

// NewUserHandler 创建并返回一个新的UserHandler实例。
func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// GetProfile 处理获取当前用户资料摘要的请求。
// 路径: GET /me/profile (受保护路由)
// 响应: UserProfileSummaryDTO 或错误响应
func (h *UserHandler) GetProfile(c *gin.Context) {
	user, ok := middleware.CurrentUser(c)
	if !ok {
		utils.UnauthorizedResponse(c, "Authentication required")
		return
	}

	result, err := h.service.GetProfileSummary(user.ID, middleware.CurrentLocale(c))
	if err != nil {
		if err == service.ErrNotFound {
			utils.NotFoundResponse(c, "User not found")
			return
		}
		utils.InternalErrorResponse(c, "Failed to load profile", nil)
		return
	}

	utils.OKResponse(c, result)
}

// GetCheckInHistory 处理获取当前用户签到历史的请求。
// 路径: GET /me/checkins (受保护路由)
// 查询参数:
//   - page: 页码，默认1
//   - pageSize: 每页大小，默认20，最大50
//
// 响应: UserCheckInHistoryResponseDTO 或错误响应
func (h *UserHandler) GetCheckInHistory(c *gin.Context) {
	user, ok := middleware.CurrentUser(c)
	if !ok {
		utils.UnauthorizedResponse(c, "Authentication required")
		return
	}

	// 解析分页参数，使用默认值
	page := parseInt(c.Query("page"), 1)
	pageSize := parseInt(c.Query("pageSize"), 20)

	result, err := h.service.GetCheckInHistory(user.ID, page, pageSize, middleware.CurrentLocale(c))
	if err != nil {
		utils.InternalErrorResponse(c, "Failed to load check-in history", nil)
		return
	}

	utils.OKResponse(c, result)
}
