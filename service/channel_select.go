package service

import (
	"errors"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/logger"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/setting"
	"github.com/gin-gonic/gin"
)

type RetryParam struct {
	Ctx          *gin.Context
	TokenGroup   string
	BackupGroup  string
	ModelName    string
	RequestPath  string
	Retry        *int
	resetNextTry bool
	backupActive bool
}

func (p *RetryParam) CurrentGroup() string {
	if p.backupActive {
		return strings.TrimSpace(p.BackupGroup)
	}
	return strings.TrimSpace(p.TokenGroup)
}

func (p *RetryParam) CanUseBackup() bool {
	backupGroup := strings.TrimSpace(p.BackupGroup)
	return !p.backupActive && backupGroup != "" && backupGroup != strings.TrimSpace(p.TokenGroup)
}

func (p *RetryParam) ActivateBackup() bool {
	if !p.CanUseBackup() {
		return false
	}
	p.backupActive = true
	p.SetRetry(0)
	return true
}

func (p *RetryParam) GetRetry() int {
	if p.Retry == nil {
		return 0
	}
	return *p.Retry
}

func (p *RetryParam) SetRetry(retry int) {
	p.Retry = &retry
}

func (p *RetryParam) IncreaseRetry() {
	if p.resetNextTry {
		p.resetNextTry = false
		return
	}
	if p.Retry == nil {
		p.Retry = new(int)
	}
	*p.Retry++
}

func (p *RetryParam) ResetRetryNextTry() {
	p.resetNextTry = true
}

// CacheGetRandomSatisfiedChannel tries the primary group first, then switches
// to the configured backup group when the primary has no available channel.
func CacheGetRandomSatisfiedChannel(param *RetryParam) (*model.Channel, string, error) {
	for {
		var channel *model.Channel
		var err error
		currentGroup := param.CurrentGroup()
		selectGroup := currentGroup
		userGroup := common.GetContextKeyString(param.Ctx, constant.ContextKeyUserGroup)

		if currentGroup == "auto" {
			if len(setting.GetAutoGroups()) == 0 {
				return nil, selectGroup, errors.New("auto groups is not enabled")
			}
			autoGroups := GetUserAutoGroup(userGroup)
			startGroupIndex := 0

			if lastGroupIndex, exists := common.GetContextKey(param.Ctx, constant.ContextKeyAutoGroupIndex); exists {
				if idx, ok := lastGroupIndex.(int); ok {
					startGroupIndex = idx
				}
			}

			for i := startGroupIndex; i < len(autoGroups); i++ {
				autoGroup := autoGroups[i]
				priorityRetry := param.GetRetry()
				if i > startGroupIndex {
					priorityRetry = 0
				}
				logger.LogDebug(param.Ctx, "Auto selecting group: %s, priorityRetry: %d", autoGroup, priorityRetry)

				channel, _ = model.GetRandomSatisfiedChannel(autoGroup, param.ModelName, priorityRetry, param.RequestPath)
				if channel == nil {
					logger.LogDebug(param.Ctx, "No available channel in group %s for model %s at priorityRetry %d, trying next group", autoGroup, param.ModelName, priorityRetry)
					common.SetContextKey(param.Ctx, constant.ContextKeyAutoGroupIndex, i+1)
					param.SetRetry(0)
					continue
				}
				selectGroup = autoGroup
				logger.LogDebug(param.Ctx, "Auto selected group: %s", autoGroup)

				common.SetContextKey(param.Ctx, constant.ContextKeyAutoGroupIndex, i)
				break
			}
		} else {
			channel, err = model.GetRandomSatisfiedChannel(currentGroup, param.ModelName, param.GetRetry(), param.RequestPath)
		}

		if err != nil {
			return nil, selectGroup, err
		}
		if channel != nil {
			common.SetContextKey(param.Ctx, constant.ContextKeyUsingGroup, selectGroup)
			common.SetContextKey(param.Ctx, constant.ContextKeyAutoGroup, selectGroup)
			return channel, selectGroup, nil
		}
		if !param.ActivateBackup() {
			return nil, selectGroup, nil
		}
		logger.LogDebug(param.Ctx, "No available channel in primary group %s for model %s, trying backup group %s", currentGroup, param.ModelName, param.CurrentGroup())
	}
}
