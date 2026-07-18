package controller

import (
	"net/http"
	"strconv"
	"strings"

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

// --- User-facing invite API ---

type selfInviteBenefit struct {
	Id                    int     `json:"id"`
	Name                  string  `json:"name"`
	Remark                string  `json:"remark"`
	TriggerType           string  `json:"trigger_type"`
	TriggerTopupQuota     int64   `json:"trigger_topup_quota"`
	RewardType            string  `json:"reward_type"`
	RewardPercent         float64 `json:"reward_percent"`
	RewardMaxQuota        int64   `json:"reward_max_quota"`
	MaxInviteesPerInviter int     `json:"max_invitees_per_inviter"`
	EffectiveDays         int     `json:"effective_days"`
	ExpiresAt             int64   `json:"expires_at"`
}

type selfInviteeItem struct {
	Id             int    `json:"id"`
	Username       string `json:"username"`
	DisplayName    string `json:"display_name"`
	Status         int    `json:"status"`
	CreatedAt      int64  `json:"created_at"`
	Rewarded       bool   `json:"rewarded"`
	RewardQuota    int64  `json:"reward_quota"`
	RewardPlanName string `json:"reward_plan_name,omitempty"`
}

func GetSelfInviteDashboard(c *gin.Context) {
	userId := c.GetInt("id")
	user, err := model.GetUserById(userId, true)
	if err != nil {
		common.ApiError(c, err)
		return
	}

	// Auto-generate aff code if empty
	if strings.TrimSpace(user.AffCode) == "" {
		user.AffCode = common.GetRandomString(4)
		if err = user.Update(false); err != nil {
			common.ApiError(c, err)
			return
		}
	}

	now := common.GetTimestamp()

	// Get enabled plans as benefits list
	plans, err := dev.GetEnabledInvitePlans(model.DB)
	if err != nil {
		common.ApiError(c, err)
		return
	}
	benefits := make([]selfInviteBenefit, 0, len(plans))
	for _, p := range plans {
		if p.ExpiresAt > 0 && p.ExpiresAt < now {
			continue
		}
		benefits = append(benefits, selfInviteBenefit{
			Id:                    p.Id,
			Name:                  p.Name,
			Remark:                p.Remark,
			TriggerType:           p.TriggerType,
			TriggerTopupQuota:     p.TriggerTopupQuota,
			RewardType:            p.RewardType,
			RewardPercent:         p.RewardPercent,
			RewardMaxQuota:        p.RewardMaxQuota,
			MaxInviteesPerInviter: p.MaxInviteesPerInviter,
			EffectiveDays:         p.EffectiveDays,
			ExpiresAt:             p.ExpiresAt,
		})
	}

	// Count invitees
	var inviteeTotal int64
	if err := model.DB.Model(&model.User{}).Where("inviter_id = ?", userId).Count(&inviteeTotal).Error; err != nil {
		common.ApiError(c, err)
		return
	}

	// Count rewarded (distinct invitees with reward records)
	var rewardedTotal int64
	if err := model.DB.Model(&dev.InviteRewardRecord{}).
		Where("inviter_id = ?", userId).
		Distinct("invitee_id").
		Count(&rewardedTotal).Error; err != nil {
		common.ApiError(c, err)
		return
	}

	// Paginated invitee list
	page, _ := strconv.Atoi(c.DefaultQuery("p", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	type userLite struct {
		Id          int
		Username    string
		DisplayName string
		Status      int
		CreatedAt   int64
	}
	var invitees []userLite
	model.DB.Model(&model.User{}).
		Select("id, username, display_name, status, created_at").
		Where("inviter_id = ?", userId).
		Order("id desc").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&invitees)

	// Collect reward records for these invitees
	inviteeIDs := make([]int, len(invitees))
	for i, inv := range invitees {
		inviteeIDs[i] = inv.Id
	}

	type rewardInfo struct {
		InviteeId int
		PlanId    int
		Quota     int64
	}
	var rewards []rewardInfo
	if len(inviteeIDs) > 0 {
		model.DB.Model(&dev.InviteRewardRecord{}).
			Select("invitee_id, plan_id, reward_quota as quota").
			Where("inviter_id = ? AND invitee_id IN ?", userId, inviteeIDs).
			Order("id desc").
			Find(&rewards)
	}
	rewardMap := make(map[int]rewardInfo)
	for _, r := range rewards {
		if _, exists := rewardMap[r.InviteeId]; !exists {
			rewardMap[r.InviteeId] = r
		}
	}

	// Collect plan names for rewards
	planIDs := make(map[int]bool)
	for _, r := range rewards {
		planIDs[r.PlanId] = true
	}
	planNames := make(map[int]string)
	for pid := range planIDs {
		var name string
		if err := model.DB.Table("invite_plans").Where("id = ?", pid).Select("name").Scan(&name).Error; err == nil {
			planNames[pid] = name
		}
	}

	items := make([]selfInviteeItem, 0, len(invitees))
	for _, inv := range invitees {
		displayName := strings.TrimSpace(inv.DisplayName)
		if displayName == "" {
			displayName = strings.TrimSpace(inv.Username)
		}
		item := selfInviteeItem{
			Id:          inv.Id,
			Username:    inv.Username,
			DisplayName: displayName,
			Status:      inv.Status,
			CreatedAt:   inv.CreatedAt,
		}
		if rw, ok := rewardMap[inv.Id]; ok {
			item.Rewarded = true
			item.RewardQuota = rw.Quota
			item.RewardPlanName = planNames[rw.PlanId]
		}
		items = append(items, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"aff_code":          user.AffCode,
			"invite_total":      inviteeTotal,
			"rewarded_total":    rewardedTotal,
			"aff_quota":         user.AffQuota,
			"aff_history_quota": user.AffHistoryQuota,
			"benefits":          benefits,
			"invitees": gin.H{
				"page":      page,
				"page_size": pageSize,
				"total":     inviteeTotal,
				"items":     items,
			},
		},
	})
}

func GetSelfInviteLogs(c *gin.Context) {
	userId := c.GetInt("id")
	page, _ := strconv.Atoi(c.DefaultQuery("p", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Get invite reward records for this user as inviter
	var total int64
	model.DB.Model(&dev.InviteRewardRecord{}).Where("inviter_id = ?", userId).Count(&total)

	type logItem struct {
		dev.InviteRewardRecord
		PlanName    string `json:"plan_name"`
		InviteeName string `json:"invitee_name"`
	}
	var records []dev.InviteRewardRecord
	model.DB.Where("inviter_id = ?", userId).
		Order("id desc").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&records)

	// Enrich
	planIDs := make(map[int]bool)
	inviteeIDs := make(map[int]bool)
	for _, r := range records {
		planIDs[r.PlanId] = true
		inviteeIDs[r.InviteeId] = true
	}
	planNames := make(map[int]string)
	for pid := range planIDs {
		var name string
		model.DB.Table("invite_plans").Where("id = ?", pid).Select("name").Scan(&name)
		planNames[pid] = name
	}
	userNames := make(map[int]string)
	for uid := range inviteeIDs {
		var name string
		model.DB.Table("users").Where("id = ?", uid).Select("username").Scan(&name)
		userNames[uid] = name
	}

	items := make([]logItem, 0, len(records))
	for _, r := range records {
		items = append(items, logItem{
			InviteRewardRecord: r,
			PlanName:           planNames[r.PlanId],
			InviteeName:        userNames[r.InviteeId],
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"items": items,
			"total": total,
		},
	})
}
