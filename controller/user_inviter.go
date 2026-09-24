package controller

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/i18n"
	"github.com/QuantumNous/new-api/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AdminBindUserInviter(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil || userId <= 0 {
		common.ApiErrorI18n(c, i18n.MsgInvalidParams)
		return
	}
	var request struct {
		AffCode string `json:"aff_code"`
	}
	if err := common.DecodeJson(c.Request.Body, &request); err != nil {
		common.ApiErrorI18n(c, i18n.MsgInvalidParams)
		return
	}
	user, err := model.GetUserById(userId, false)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			common.ApiErrorI18n(c, i18n.MsgUserNotExists)
		} else {
			common.ApiError(c, err)
		}
		return
	}
	if !canManageTargetRole(c.GetInt("role"), user.Role) {
		common.ApiErrorI18n(c, i18n.MsgUserNoPermissionHigherLevel)
		return
	}
	updated, err := model.BindUserInviter(userId, request.AffCode)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrUserInviterAlreadyBound):
			common.ApiErrorI18n(c, i18n.MsgUserInviterAlreadyBound)
		case errors.Is(err, model.ErrInvalidInviteCode):
			common.ApiErrorI18n(c, i18n.MsgUserInvalidInviteCode)
		case errors.Is(err, model.ErrUserInviterSelf):
			common.ApiErrorI18n(c, i18n.MsgUserInviterSelf)
		default:
			common.ApiError(c, err)
		}
		return
	}
	for _, id := range []int{userId, updated.InviterId} {
		if err := model.InvalidateUserCache(id); err != nil {
			common.SysError(fmt.Sprintf("failed to invalidate user cache for user %d: %s", id, err.Error()))
		}
	}
	recordManageAuditFor(c, userId, "user.inviter_bind", map[string]interface{}{
		"username":      user.Username,
		"id":            userId,
		"inviter_id":    updated.InviterId,
		"referral_mode": updated.ReferralMode,
	})
	common.ApiSuccess(c, gin.H{
		"id":            userId,
		"inviter_id":    updated.InviterId,
		"referral_mode": updated.ReferralMode,
	})
}
