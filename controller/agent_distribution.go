package controller

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

func getSelfAgentUser(c *gin.Context) (*model.User, bool) {
	user, err := model.GetUserById(c.GetInt("id"), false)
	if err != nil {
		common.ApiError(c, err)
		return nil, false
	}
	if !model.IsAgentPortalVisible(user) {
		common.ApiErrorMsg(c, "当前账户未开通代理分销")
		return nil, false
	}
	return user, true
}

func GetSelfAgentDashboard(c *gin.Context) {
	user, ok := getSelfAgentUser(c)
	if !ok {
		return
	}
	pageInfo := common.GetPageQuery(c)
	stats, err := model.GetAgentDashboardStats(user.Id)
	if err != nil {
		common.ApiError(c, err)
		return
	}
	invitees, total, err := model.GetAgentInvitees(user.Id, pageInfo, strings.TrimSpace(c.Query("keyword")))
	if err != nil {
		common.ApiError(c, err)
		return
	}
	firstRate, repeatRate := model.GetUserEffectiveAgentRates(user)
	pageInfo.SetTotal(int(total))
	pageInfo.SetItems(invitees)
	common.ApiSuccess(c, gin.H{
		"agent_enabled":               user.AgentEnabled,
		"agent_portal_visible":        model.IsAgentPortalVisible(user),
		"agent_use_default_rates":     user.AgentUseDefaultRates,
		"agent_first_topup_rate":      user.AgentFirstTopupRate,
		"agent_repeat_topup_rate":     user.AgentRepeatTopupRate,
		"effective_first_topup_rate":  firstRate,
		"effective_repeat_topup_rate": repeatRate,
		"agent_commission_balance":    user.AgentCommissionBalance,
		"agent_commission_total":      user.AgentCommissionTotal,
		"agent_commission_withdrawn":  user.AgentCommissionWithdrawn,
		"aff_code":                    strings.TrimSpace(user.AffCode),
		"stats":                       stats,
		"invitees":                    pageInfo,
	})
}

func GetSelfAgentRecords(c *gin.Context) {
	user, ok := getSelfAgentUser(c)
	if !ok {
		return
	}
	pageInfo := common.GetPageQuery(c)
	items, total, err := model.GetAgentCommissionRecords(user.Id, pageInfo, strings.TrimSpace(c.Query("keyword")))
	if err != nil {
		common.ApiError(c, err)
		return
	}
	pageInfo.SetTotal(int(total))
	pageInfo.SetItems(items)
	common.ApiSuccess(c, pageInfo)
}

func GetSelfAgentTopUps(c *gin.Context) {
	user, ok := getSelfAgentUser(c)
	if !ok {
		return
	}
	pageInfo := common.GetPageQuery(c)
	items, total, err := model.GetAgentTopUps(user.Id, pageInfo, strings.TrimSpace(c.Query("keyword")))
	if err != nil {
		common.ApiError(c, err)
		return
	}
	pageInfo.SetTotal(int(total))
	pageInfo.SetItems(items)
	common.ApiSuccess(c, pageInfo)
}

func GetSelfAgentInsights(c *gin.Context) {
	user, ok := getSelfAgentUser(c)
	if !ok {
		return
	}
	days, err := strconv.Atoi(strings.TrimSpace(c.DefaultQuery("days", "7")))
	if err != nil {
		days = 7
	}
	insights, err := model.GetAgentInsights(user.Id, days)
	if err != nil {
		common.ApiError(c, err)
		return
	}
	common.ApiSuccess(c, insights)
}

type adminAgentWithdrawRequest struct {
	Amount float64 `json:"amount"`
	Remark string  `json:"remark"`
}

func AdminWithdrawAgentCommission(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil || userId <= 0 {
		common.ApiErrorMsg(c, "无效的用户 ID")
		return
	}
	var request adminAgentWithdrawRequest
	if err := common.DecodeJson(c.Request.Body, &request); err != nil || request.Amount <= 0 {
		common.ApiErrorMsg(c, "请输入正确的提现金额")
		return
	}
	amountCents := decimal.NewFromFloat(request.Amount).Mul(decimal.NewFromInt(100)).Round(0).IntPart()
	updated, err := model.WithdrawAgentCommission(userId, amountCents)
	if err != nil {
		if errors.Is(err, model.ErrAgentCommissionBalanceInsufficient) {
			common.ApiErrorMsg(c, "可提现佣金不足")
			return
		}
		common.ApiError(c, err)
		return
	}
	message := fmt.Sprintf("管理员处理代理佣金提现，扣减 ¥%.2f", request.Amount)
	if remark := strings.TrimSpace(request.Remark); remark != "" {
		message += "，备注：" + remark
	}
	model.RecordLog(userId, model.LogTypeManage, message)
	common.ApiSuccess(c, gin.H{
		"agent_commission_balance":   updated.AgentCommissionBalance,
		"agent_commission_total":     updated.AgentCommissionTotal,
		"agent_commission_withdrawn": updated.AgentCommissionWithdrawn,
	})
}
