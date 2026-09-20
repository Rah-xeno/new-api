package relay

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/dto"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/setting"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplySystemPromptIfNeededPrependsFixedPrompt(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	info := &relaycommon.RelayInfo{
		ChannelMeta: &relaycommon.ChannelMeta{
			ChannelSetting: dto.ChannelSettings{SystemPrompt: "channel prompt"},
		},
	}
	request := &dto.GeneralOpenAIRequest{
		Messages: []dto.Message{
			{Role: "system", Content: "user system prompt"},
			{Role: "user", Content: "hello"},
		},
	}

	applySystemPromptIfNeeded(c, info, request)

	require.Len(t, request.Messages, 2)
	assert.Equal(t, "system", request.Messages[0].Role)
	assert.Equal(
		t,
		strings.TrimSpace(constant.GlobalSystemPromptAppend)+"\nuser system prompt",
		request.Messages[0].StringContent(),
	)
	assert.NotContains(t, request.Messages[0].StringContent(), "channel prompt")
}

func TestApplySystemPromptToResponsesRequestPrependsFixedPrompt(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	info := &relaycommon.RelayInfo{ChannelMeta: &relaycommon.ChannelMeta{}}
	existing, err := common.Marshal("user instructions")
	require.NoError(t, err)
	request := &dto.OpenAIResponsesRequest{Instructions: existing}

	applySystemPromptToResponsesRequest(c, info, request)

	var instructions string
	require.NoError(t, common.Unmarshal(request.Instructions, &instructions))
	assert.Equal(
		t,
		strings.TrimSpace(constant.GlobalSystemPromptAppend)+"\nuser instructions",
		instructions,
	)
}

func TestApplySystemPromptToResponsesRequestPreservesInstructionArray(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	info := &relaycommon.RelayInfo{ChannelMeta: &relaycommon.ChannelMeta{}}
	request := &dto.OpenAIResponsesRequest{
		Instructions: json.RawMessage(`[{"role":"developer","content":"existing"}]`),
	}

	applySystemPromptToResponsesRequest(c, info, request)

	var instructions []json.RawMessage
	require.NoError(t, common.Unmarshal(request.Instructions, &instructions))
	require.Len(t, instructions, 2)
	var prefix string
	require.NoError(t, common.Unmarshal(instructions[0], &prefix))
	assert.Equal(t, strings.TrimSpace(constant.GlobalSystemPromptAppend), prefix)
	assert.JSONEq(t, `{"role":"developer","content":"existing"}`, string(instructions[1]))
}

func TestApplySystemPromptIfNeededHonorsGroupExemption(t *testing.T) {
	require.NoError(t, setting.UpdateSystemPromptExemptGroupsByJSONString(`{"exempt":true}`))
	t.Cleanup(func() {
		require.NoError(t, setting.UpdateSystemPromptExemptGroupsByJSONString(`{}`))
	})

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	info := &relaycommon.RelayInfo{
		UsingGroup: "exempt",
		ChannelMeta: &relaycommon.ChannelMeta{
			ChannelSetting: dto.ChannelSettings{SystemPrompt: "channel prompt"},
		},
	}
	request := &dto.GeneralOpenAIRequest{
		Messages: []dto.Message{{Role: "user", Content: "hello"}},
	}

	applySystemPromptIfNeeded(c, info, request)

	require.Len(t, request.Messages, 2)
	assert.Equal(t, "system", request.Messages[0].Role)
	assert.Equal(t, "channel prompt", request.Messages[0].StringContent())
	assert.NotContains(t, request.Messages[0].StringContent(), strings.TrimSpace(constant.GlobalSystemPromptAppend))
}
