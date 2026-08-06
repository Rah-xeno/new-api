package model

import (
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGlobalSystemPromptOptionFallsBackAndPersistsOverride(t *testing.T) {
	require.NoError(t, DB.AutoMigrate(&Option{}))
	require.NoError(t, DB.Where("key = ?", GlobalSystemPromptAppendOptionKey).Delete(&Option{}).Error)

	previousRedisEnabled := common.RedisEnabled
	common.OptionMapRWMutex.Lock()
	previousOptionMap := common.OptionMap
	common.OptionMap = make(map[string]string)
	common.OptionMapRWMutex.Unlock()
	common.RedisEnabled = false
	t.Cleanup(func() {
		require.NoError(t, DB.Where("key = ?", GlobalSystemPromptAppendOptionKey).Delete(&Option{}).Error)
		common.RedisEnabled = previousRedisEnabled
		common.OptionMapRWMutex.Lock()
		common.OptionMap = previousOptionMap
		common.OptionMapRWMutex.Unlock()
	})

	assert.Equal(t, constant.GlobalSystemPromptAppend, GetGlobalSystemPromptAppend())

	const configuredPrompt = "configured global prompt"
	require.NoError(t, UpdateOption(GlobalSystemPromptAppendOptionKey, configuredPrompt))

	var saved Option
	require.NoError(t, DB.First(&saved, "key = ?", GlobalSystemPromptAppendOptionKey).Error)
	assert.Equal(t, configuredPrompt, saved.Value)
	assert.Equal(t, configuredPrompt, GetGlobalSystemPromptAppend())
}
