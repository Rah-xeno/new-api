package controller

import (
	"net/http"
	"strconv"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/dev"
	"github.com/QuantumNous/new-api/model"

	"github.com/gin-gonic/gin"
)

func GetInvitePlans(c *gin.Context) {
	plans, err := dev.GetAllInvitePlans(model.DB)
	if err != nil {
		common.ApiError(c, err)
		return
	}
	if plans == nil {
		plans = []dev.InvitePlan{}
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": plans})
}

func CreateInvitePlan(c *gin.Context) {
	var req struct {
		Plan dev.InvitePlan `json:"plan"`
	}
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
	if err := dev.CreateInvitePlan(model.DB, &req.Plan); err != nil {
		common.ApiError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": req.Plan})
}

func UpdateInvitePlan(c *gin.Context) {
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
	if planData, ok := body["plan"].(map[string]interface{}); ok {
		body = planData
	}
	if err := dev.UpdateInvitePlan(model.DB, id, body); err != nil {
		common.ApiError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func DeleteInvitePlan(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		common.ApiError(c, err)
		return
	}
	if err := dev.DeleteInvitePlan(model.DB, id); err != nil {
		common.ApiError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func GetInviteRewardRecords(c *gin.Context) {
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

	records, total, err := dev.GetRewardRecords(model.DB, planId, page, pageSize)
	if err != nil {
		common.ApiError(c, err)
		return
	}
	if records == nil {
		records = []dev.InviteRewardRecord{}
	}

	// Enrich with plan name and usernames
	type RewardRecordWithNames struct {
		dev.InviteRewardRecord
		PlanName    string `json:"plan_name"`
		InviterName string `json:"inviter_name"`
		InviteeName string `json:"invitee_name"`
	}

	enriched := make([]RewardRecordWithNames, len(records))

	// Collect plan names
	planNames := make(map[int]string)
	for _, r := range records {
		if _, ok := planNames[r.PlanId]; !ok {
			var name string
			if err := model.DB.Table("invite_plans").Where("id = ?", r.PlanId).Select("name").Scan(&name).Error; err == nil {
				planNames[r.PlanId] = name
			}
		}
	}

	// Collect user IDs
	userIds := make(map[int]bool)
	for _, r := range records {
		userIds[r.InviterId] = true
		userIds[r.InviteeId] = true
	}

	userNames := make(map[int]string)
	for id := range userIds {
		var username string
		if err := model.DB.Table("users").Where("id = ?", id).Select("username").Scan(&username).Error; err == nil {
			userNames[id] = username
		}
	}

	for i, r := range records {
		enriched[i] = RewardRecordWithNames{
			InviteRewardRecord: r,
			PlanName:           planNames[r.PlanId],
			InviterName:        userNames[r.InviterId],
			InviteeName:        userNames[r.InviteeId],
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"items": enriched,
			"total": total,
		},
	})
}
