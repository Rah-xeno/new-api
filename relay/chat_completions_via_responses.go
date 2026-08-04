package relay

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/relay/channel"
	openaichannel "github.com/QuantumNous/new-api/relay/channel/openai"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	relayconstant "github.com/QuantumNous/new-api/relay/constant"
	"github.com/QuantumNous/new-api/service"
	"github.com/QuantumNous/new-api/setting"
	"github.com/QuantumNous/new-api/types"

	"github.com/gin-gonic/gin"
)

func globalSystemPromptForGroup(group string) string {
	if setting.IsGroupExemptFromSystemPrompt(group) {
		return ""
	}
	return strings.TrimSpace(constant.GlobalSystemPromptAppend)
}

func applySystemPromptIfNeeded(c *gin.Context, info *relaycommon.RelayInfo, request *dto.GeneralOpenAIRequest) {
	if info == nil || request == nil {
		return
	}

	globalPrompt := globalSystemPromptForGroup(info.UsingGroup)
	channelPrompt := strings.TrimSpace(info.ChannelSetting.SystemPrompt)
	if globalPrompt == "" && channelPrompt == "" {
		return
	}

	systemRole := request.GetSystemRoleName()
	systemMessageIndex := -1
	for i, message := range request.Messages {
		if message.Role == systemRole {
			systemMessageIndex = i
			break
		}
	}

	prefixParts := make([]string, 0, 2)
	if globalPrompt != "" {
		prefixParts = append(prefixParts, globalPrompt)
	}
	if channelPrompt != "" && (systemMessageIndex == -1 || info.ChannelSetting.SystemPromptOverride) {
		prefixParts = append(prefixParts, channelPrompt)
	}
	if len(prefixParts) == 0 {
		return
	}
	prefixText := strings.Join(prefixParts, "\n")

	if systemMessageIndex == -1 {
		systemMessage := dto.Message{
			Role:    systemRole,
			Content: prefixText,
		}
		request.Messages = append([]dto.Message{systemMessage}, request.Messages...)
		return
	}

	common.SetContextKey(c, constant.ContextKeySystemPromptOverride, true)
	message := &request.Messages[systemMessageIndex]
	if message.IsStringContent() {
		existing := message.StringContent()
		if strings.TrimSpace(existing) == "" {
			message.SetStringContent(prefixText)
		} else {
			message.SetStringContent(prefixText + "\n" + existing)
		}
		return
	}
	contents := message.ParseContent()
	contents = append([]dto.MediaContent{{
		Type: dto.ContentTypeText,
		Text: prefixText,
	}}, contents...)
	message.Content = contents
}

func applySystemPromptToResponsesRequest(c *gin.Context, info *relaycommon.RelayInfo, request *dto.OpenAIResponsesRequest) {
	if info == nil || request == nil {
		return
	}

	globalPrompt := globalSystemPromptForGroup(info.UsingGroup)
	channelPrompt := strings.TrimSpace(info.ChannelSetting.SystemPrompt)
	hasExisting := len(request.Instructions) > 0 && string(request.Instructions) != "null"

	prefixParts := make([]string, 0, 2)
	if globalPrompt != "" {
		prefixParts = append(prefixParts, globalPrompt)
	}
	if channelPrompt != "" && (!hasExisting || info.ChannelSetting.SystemPromptOverride) {
		prefixParts = append(prefixParts, channelPrompt)
	}
	if len(prefixParts) == 0 {
		return
	}
	prefixText := strings.Join(prefixParts, "\n")

	if !hasExisting {
		if raw, err := common.Marshal(prefixText); err == nil {
			request.Instructions = raw
		}
		return
	}

	common.SetContextKey(c, constant.ContextKeySystemPromptOverride, true)
	var existingString string
	if err := common.Unmarshal(request.Instructions, &existingString); err == nil {
		raw, marshalErr := common.Marshal(prefixText + "\n" + existingString)
		if marshalErr == nil {
			request.Instructions = raw
		}
		return
	}

	var existingItems []json.RawMessage
	if err := common.Unmarshal(request.Instructions, &existingItems); err != nil {
		return
	}
	prefixItem, err := common.Marshal(prefixText)
	if err != nil {
		return
	}
	mergedItems := append([]json.RawMessage{prefixItem}, existingItems...)
	if raw, marshalErr := common.Marshal(mergedItems); marshalErr == nil {
		request.Instructions = raw
	}
}

func chatCompletionsViaResponses(c *gin.Context, info *relaycommon.RelayInfo, adaptor channel.Adaptor, request *dto.GeneralOpenAIRequest) (*dto.Usage, *types.NewAPIError) {
	chatJSON, err := common.Marshal(request)
	if err != nil {
		return nil, types.NewError(err, types.ErrorCodeConvertRequestFailed, types.ErrOptionWithSkipRetry())
	}

	chatJSON, err = relaycommon.RemoveDisabledFields(chatJSON, info.ChannelOtherSettings, info.ChannelSetting.PassThroughBodyEnabled)
	if err != nil {
		return nil, types.NewError(err, types.ErrorCodeConvertRequestFailed, types.ErrOptionWithSkipRetry())
	}

	if len(info.ParamOverride) > 0 {
		chatJSON, err = relaycommon.ApplyParamOverrideWithRelayInfo(chatJSON, info)
		if err != nil {
			return nil, newAPIErrorFromParamOverride(err)
		}
	}

	var overriddenChatReq dto.GeneralOpenAIRequest
	if err := common.Unmarshal(chatJSON, &overriddenChatReq); err != nil {
		return nil, types.NewError(err, types.ErrorCodeChannelParamOverrideInvalid, types.ErrOptionWithSkipRetry())
	}

	result, err := service.ConvertRequestVia(c, info, &overriddenChatReq, types.RelayFormatOpenAI, types.RelayFormatOpenAIResponses)
	if err != nil {
		return nil, types.NewErrorWithStatusCode(err, types.ErrorCodeInvalidRequest, http.StatusBadRequest, types.ErrOptionWithSkipRetry())
	}
	responsesReq, ok := result.Value.(*dto.OpenAIResponsesRequest)
	if !ok {
		return nil, types.NewError(fmt.Errorf("expected OpenAI responses request, got %T", result.Value), types.ErrorCodeConvertRequestFailed, types.ErrOptionWithSkipRetry())
	}

	savedRelayMode := info.RelayMode
	savedRequestURLPath := info.RequestURLPath
	defer func() {
		info.RelayMode = savedRelayMode
		info.RequestURLPath = savedRequestURLPath
	}()

	info.RelayMode = relayconstant.RelayModeResponses
	info.RequestURLPath = "/v1/responses"

	convertedRequest, err := adaptor.ConvertOpenAIResponsesRequest(c, info, *responsesReq)
	if err != nil {
		return nil, types.NewError(err, types.ErrorCodeConvertRequestFailed, types.ErrOptionWithSkipRetry())
	}
	relaycommon.AppendRequestConversionFromRequest(info, convertedRequest)

	jsonData, err := common.Marshal(convertedRequest)
	if err != nil {
		return nil, types.NewError(err, types.ErrorCodeConvertRequestFailed, types.ErrOptionWithSkipRetry())
	}

	jsonData, err = relaycommon.RemoveDisabledFields(jsonData, info.ChannelOtherSettings, info.ChannelSetting.PassThroughBodyEnabled)
	if err != nil {
		return nil, types.NewError(err, types.ErrorCodeConvertRequestFailed, types.ErrOptionWithSkipRetry())
	}

	body, size, closer, err := relaycommon.NewOutboundJSONBody(jsonData)
	if err != nil {
		return nil, types.NewError(err, types.ErrorCodeConvertRequestFailed, types.ErrOptionWithSkipRetry())
	}
	defer closer.Close()
	jsonData = nil
	info.UpstreamRequestBodySize = size
	var requestBody io.Reader = body

	var httpResp *http.Response
	resp, err := adaptor.DoRequest(c, info, requestBody)
	if err != nil {
		return nil, types.NewOpenAIError(err, types.ErrorCodeDoRequestFailed, http.StatusInternalServerError)
	}
	if resp == nil {
		return nil, types.NewOpenAIError(nil, types.ErrorCodeBadResponse, http.StatusInternalServerError)
	}

	statusCodeMappingStr := c.GetString("status_code_mapping")

	httpResp = resp.(*http.Response)
	clientStream := info.IsStream
	upstreamStream := isResponsesEventStreamContentType(httpResp.Header.Get("Content-Type"))
	info.IsStream = clientStream || upstreamStream
	if httpResp.StatusCode != http.StatusOK {
		newApiErr := service.RelayErrorHandler(c.Request.Context(), httpResp, false)
		service.ResetStatusCode(newApiErr, statusCodeMappingStr)
		return nil, newApiErr
	}

	if upstreamStream && clientStream {
		usage, newApiErr := openaichannel.OaiResponsesToChatStreamHandler(c, info, httpResp)
		if newApiErr != nil {
			service.ResetStatusCode(newApiErr, statusCodeMappingStr)
			return nil, newApiErr
		}
		return usage, nil
	}
	if upstreamStream {
		info.IsStream = false
		usage, newApiErr := openaichannel.OaiResponsesToChatBufferedStreamHandler(c, info, httpResp)
		if newApiErr != nil {
			service.ResetStatusCode(newApiErr, statusCodeMappingStr)
			return nil, newApiErr
		}
		return usage, nil
	}

	usage, newApiErr := openaichannel.OaiResponsesToChatHandler(c, info, httpResp)
	if newApiErr != nil {
		service.ResetStatusCode(newApiErr, statusCodeMappingStr)
		return nil, newApiErr
	}
	return usage, nil
}

func isResponsesEventStreamContentType(contentType string) bool {
	return strings.Contains(strings.ToLower(contentType), "text/event-stream")
}
