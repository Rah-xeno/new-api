package dev

import (
	"net/http"
	"strconv"

	"github.com/QuantumNous/new-api/common"

	"github.com/gin-gonic/gin"
)

// --- Admin API ---

type createInvitePlanReq struct {
	Plan InvitePlan `json:"plan"`
}

func CreateInvitePlanHandler(c *gin.Context) {
	var req createInvitePlanReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ApiError(c, err)
		return
	}
	if req.Plan.Name == "" {
		common.ApiErrorI18n(c, "plan name is required")
		return
	}
	req.Plan.CreatedAt = common.GetTimestamp()
	req.Plan.UpdatedAt = req.Plan.CreatedAt
	if err := CreateInvitePlan(&req.Plan); err != nil {
		common.ApiError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": req.Plan})
}

func UpdateInvitePlanHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		common.ApiError(c, err)
		return
	}
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		common.ApiError(c, err)
		return
	}
	// Extract plan fields if nested
	if planData, ok := body["plan"].(map[string]interface{}); ok {
		body = planData
	}
	if err := UpdateInvitePlan(id, body); err != nil {
		common.ApiError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func DeleteInvitePlanHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		common.ApiError(c, err)
		return
	}
	if err := DeleteInvitePlan(id); err != nil {
		common.ApiError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func GetInvitePlansHandler(c *gin.Context) {
	plans, err := GetAllInvitePlans()
	if err != nil {
		common.ApiError(c, err)
		return
	}
	if plans == nil {
		plans = []InvitePlan{}
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": plans})
}

func GetInviteRewardRecordsHandler(c *gin.Context) {
	planIdStr := c.Query("plan_id")
	planId, err := strconv.Atoi(planIdStr)
	if err != nil {
		common.ApiError(c, err)
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("p", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	records, total, err := GetRewardRecords(planId, page, pageSize)
	if err != nil {
		common.ApiError(c, err)
		return
	}
	if records == nil {
		records = []InviteRewardRecord{}
	}

	// Fetch inviter names
	inviterIds := make(map[int]bool)
	for _, r := range records {
		if r.InviterName == "" {
			inviterIds[r.InviterId] = true
		}
	}
	for id := range inviterIds {
		var username string
		if err := DB.Table("users").Where("id = ?", id).Select("username").Scan(&username).Error; err == nil {
			for i := range records {
				if records[i].InviterId == id {
					records[i].InviterName = username
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"items": records,
			"total": total,
		},
	})
}

// RegisterInvitePlanRoutes registers admin API routes for invite plans.
// adminAuth is the admin authentication middleware.
func RegisterInvitePlanRoutes(adminGroup *gin.RouterGroup, adminAuth gin.HandlerFunc) {
	group := adminGroup.Group("/invite-plan/admin")
	group.Use(adminAuth)
	{
		group.GET("/plans", GetInvitePlansHandler)
		group.POST("/plans", CreateInvitePlanHandler)
		group.PUT("/plans/:id", UpdateInvitePlanHandler)
		group.PATCH("/plans/:id", UpdateInvitePlanHandler)
		group.DELETE("/plans/:id", DeleteInvitePlanHandler)
		group.GET("/reward-records", GetInviteRewardRecordsHandler)
	}
}
