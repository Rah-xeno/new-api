package service

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/model"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupChannelSelectTest(t *testing.T) *gin.Context {
	t.Helper()

	gin.SetMode(gin.TestMode)
	common.SetDatabaseTypes(common.DatabaseTypeSQLite, common.DatabaseTypeSQLite)
	previousMemoryCacheEnabled := common.MemoryCacheEnabled
	previousDB := model.DB
	common.MemoryCacheEnabled = true

	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", strings.ReplaceAll(t.Name(), "/", "_"))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.Channel{}, &model.Ability{}))
	model.DB = db

	t.Cleanup(func() {
		common.MemoryCacheEnabled = previousMemoryCacheEnabled
		model.DB = previousDB
		sqlDB, dbErr := db.DB()
		if dbErr == nil {
			_ = sqlDB.Close()
		}
	})

	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	common.SetContextKey(ctx, constant.ContextKeyUserGroup, "default")
	return ctx
}

func addChannelSelectAbility(t *testing.T, group string, channelID int) {
	t.Helper()

	priority := int64(0)
	channel := &model.Channel{
		Id:       channelID,
		Name:     group,
		Key:      "test-key",
		Status:   common.ChannelStatusEnabled,
		Models:   "test-model",
		Group:    group,
		Priority: &priority,
	}
	require.NoError(t, model.DB.Create(channel).Error)
	require.NoError(t, model.DB.Create(&model.Ability{
		Group:     group,
		Model:     "test-model",
		ChannelId: channelID,
		Enabled:   true,
		Priority:  &priority,
	}).Error)
	model.InitChannelCache()
}

func TestChannelSelectionUsesBackupWhenPrimaryHasNoChannel(t *testing.T) {
	ctx := setupChannelSelectTest(t)
	addChannelSelectAbility(t, "backup", 2)

	param := &RetryParam{
		Ctx:         ctx,
		TokenGroup:  "primary",
		BackupGroup: "backup",
		ModelName:   "test-model",
		Retry:       common.GetPointer(0),
	}

	channel, selectedGroup, err := CacheGetRandomSatisfiedChannel(param)
	require.NoError(t, err)
	require.NotNil(t, channel)
	assert.Equal(t, 2, channel.Id)
	assert.Equal(t, "backup", selectedGroup)
	assert.Equal(t, "backup", param.CurrentGroup())
	assert.Equal(t, "backup", common.GetContextKeyString(ctx, constant.ContextKeyUsingGroup))
}

func TestBackupSelectionStateIsScopedToEachRequest(t *testing.T) {
	ctx := setupChannelSelectTest(t)
	addChannelSelectAbility(t, "primary", 1)
	addChannelSelectAbility(t, "backup", 2)

	firstRequest := &RetryParam{
		Ctx:         ctx,
		TokenGroup:  "primary",
		BackupGroup: "backup",
		ModelName:   "test-model",
		Retry:       common.GetPointer(0),
	}
	channel, selectedGroup, err := CacheGetRandomSatisfiedChannel(firstRequest)
	require.NoError(t, err)
	require.NotNil(t, channel)
	assert.Equal(t, 1, channel.Id)
	assert.Equal(t, "primary", selectedGroup)
	require.True(t, firstRequest.ActivateBackup())

	channel, selectedGroup, err = CacheGetRandomSatisfiedChannel(firstRequest)
	require.NoError(t, err)
	require.NotNil(t, channel)
	assert.Equal(t, 2, channel.Id)
	assert.Equal(t, "backup", selectedGroup)

	secondRequest := &RetryParam{
		Ctx:         ctx,
		TokenGroup:  "primary",
		BackupGroup: "backup",
		ModelName:   "test-model",
		Retry:       common.GetPointer(0),
	}
	channel, selectedGroup, err = CacheGetRandomSatisfiedChannel(secondRequest)
	require.NoError(t, err)
	require.NotNil(t, channel)
	assert.Equal(t, 1, channel.Id)
	assert.Equal(t, "primary", selectedGroup)
}
